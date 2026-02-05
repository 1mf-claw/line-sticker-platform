package api

import (
	"database/sql"
	"fmt"
	"sync"

	"example.com/app/internal/ai"
	_ "modernc.org/sqlite"
)

type Store struct {
	mu sync.Mutex
	db *sql.DB
	ai ai.Pipeline
	aiSecrets map[string]AICredentialsRequest
}

func NewStore() *Store {
	db, err := sql.Open("sqlite", "file:data.db?_busy_timeout=5000&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}

	s := &Store{db: db, ai: ai.MockPipeline{}, aiSecrets: map[string]AICredentialsRequest{}}
	s.migrate()
	return s
}

func (s *Store) migrate() {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			title TEXT,
			theme TEXT,
			sticker_count INTEGER,
			status TEXT,
			character_id TEXT,
			ai_provider TEXT,
			ai_model TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS characters (
			id TEXT PRIMARY KEY,
			project_id TEXT,
			source_type TEXT,
			reference_image_url TEXT,
			status TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS drafts (
			id TEXT PRIMARY KEY,
			project_id TEXT,
			idx INTEGER,
			caption TEXT,
			image_prompt TEXT,
			status TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS stickers (
			id TEXT PRIMARY KEY,
			project_id TEXT,
			draft_id TEXT,
			image_url TEXT,
			transparent_url TEXT,
			status TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			type TEXT,
			status TEXT,
			progress INTEGER,
			error_message TEXT,
			project_id TEXT,
			target_id TEXT
		);`,
	}
	for _, stmt := range stmts {
		if _, err := s.db.Exec(stmt); err != nil {
			panic(err)
		}
	}

	// Add columns for legacy DBs
	s.ensureColumn("projects", "ai_provider", "TEXT")
	s.ensureColumn("projects", "ai_model", "TEXT")
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func (s *Store) CreateProject(title string, stickerCount int) *Project {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := newID("proj")
	_, _ = s.db.Exec(
		`INSERT INTO projects (id,title,theme,sticker_count,status,character_id,ai_provider,ai_model) VALUES (?,?,?,?,?,?,?,?)`,
		id, title, "", stickerCount, "DRAFT", "", "", "",
	)
	return &Project{ID: id, Title: title, StickerCount: stickerCount, Status: "DRAFT"}
}

func (s *Store) UpdateProjectTheme(projectID string, theme string) (*Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, _ := s.db.Exec(`UPDATE projects SET theme=? WHERE id=?`, theme, projectID)
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return nil, false
	}
	return s.GetProject(projectID)
}

func (s *Store) GetProject(projectID string) (*Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := s.db.QueryRow(`SELECT id,title,theme,sticker_count,status,character_id,ai_provider,ai_model FROM projects WHERE id=?`, projectID)
	p := &Project{}
	if err := row.Scan(&p.ID, &p.Title, &p.Theme, &p.StickerCount, &p.Status, &p.CharacterID, &p.AIProvider, &p.AIModel); err != nil {
		return nil, false
	}
	return p, true
}

func (s *Store) ListProjects() []*Project {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, _ := s.db.Query(`SELECT id,title,theme,sticker_count,status,character_id,ai_provider,ai_model FROM projects`)
	defer rows.Close()
	out := []*Project{}
	for rows.Next() {
		p := &Project{}
		_ = rows.Scan(&p.ID, &p.Title, &p.Theme, &p.StickerCount, &p.Status, &p.CharacterID, &p.AIProvider, &p.AIModel)
		out = append(out, p)
	}
	return out
}

func (s *Store) CreateCharacter(projectID string, req CharacterCreateRequest) (*Character, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.GetProject(projectID); !ok {
		return nil, false
	}
	id := newID("char")
	c := &Character{
		ID:                id,
		SourceType:        req.SourceType,
		ReferenceImageURL: req.ReferenceImageURL,
		Status:            "READY",
	}
	_, _ = s.db.Exec(`INSERT INTO characters (id,project_id,source_type,reference_image_url,status) VALUES (?,?,?,?,?)`,
		c.ID, projectID, c.SourceType, c.ReferenceImageURL, c.Status,
	)
	_, _ = s.db.Exec(`UPDATE projects SET character_id=? WHERE id=?`, c.ID, projectID)
	return c, true
}

func (s *Store) GenerateDrafts(projectID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.GetProject(projectID)
	if !ok {
		return nil, false
	}
	_, _ = s.db.Exec(`UPDATE projects SET status=? WHERE id=?`, "GENERATING_DRAFTS", projectID)

	charInput := s.getCharacterInput(projectID)
	ideas, _ := s.ai.GenerateDrafts(p.Theme, p.StickerCount, charInput)
	for i := 1; i <= p.StickerCount; i++ {
		id := newID("draft")
		caption := fmt.Sprintf("草稿 %d", i)
		prompt := fmt.Sprintf("主角依主題動作 %d", i)
		if i-1 < len(ideas) {
			caption = ideas[i-1].Caption
			prompt = ideas[i-1].ImagePrompt
		}
		_, _ = s.db.Exec(`INSERT INTO drafts (id,project_id,idx,caption,image_prompt,status) VALUES (?,?,?,?,?,?)`,
			id, projectID, i, caption, prompt, "DRAFT",
		)
	}
	_, _ = s.db.Exec(`UPDATE projects SET status=? WHERE id=?`, "DRAFT_READY", projectID)
	job := s.newJob("GENERATE_DRAFT", projectID, "")
	// draft generation finished
	s.setJobProgress(job.ID, 100, "SUCCESS")
	return job, true
}

func (s *Store) ListDrafts(projectID string) []*Draft {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, _ := s.db.Query(`SELECT id,project_id,idx,caption,image_prompt,status FROM drafts WHERE project_id=? ORDER BY idx`, projectID)
	defer rows.Close()
	out := []*Draft{}
	for rows.Next() {
		d := &Draft{}
		_ = rows.Scan(&d.ID, &d.ProjectID, &d.Index, &d.Caption, &d.ImagePrompt, &d.Status)
		out = append(out, d)
	}
	return out
}

func (s *Store) UpdateDraft(draftID string, req DraftUpdateRequest) (*Draft, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.db.Exec(`UPDATE drafts SET caption=COALESCE(NULLIF(?,''),caption), image_prompt=COALESCE(NULLIF(?,''),image_prompt) WHERE id=?`,
		req.Caption, req.ImagePrompt, draftID,
	)
	row := s.db.QueryRow(`SELECT id,project_id,idx,caption,image_prompt,status FROM drafts WHERE id=?`, draftID)
	d := &Draft{}
	if err := row.Scan(&d.ID, &d.ProjectID, &d.Index, &d.Caption, &d.ImagePrompt, &d.Status); err != nil {
		return nil, false
	}
	return d, true
}

func (s *Store) GenerateStickers(projectID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.db.Exec(`UPDATE projects SET status=? WHERE id=?`, "GENERATING_IMAGES", projectID)
	rows, _ := s.db.Query(`SELECT id,image_prompt FROM drafts WHERE project_id=?`, projectID)
	defer rows.Close()
	charInput := s.getCharacterInput(projectID)
	for rows.Next() {
		var draftID, prompt string
		_ = rows.Scan(&draftID, &prompt)
		imageURL, _ := s.ai.GenerateImage(prompt, charInput)
		id := newID("stk")
		_, _ = s.db.Exec(`INSERT INTO stickers (id,project_id,draft_id,image_url,transparent_url,status) VALUES (?,?,?,?,?,?)`,
			id, projectID, draftID, imageURL, "", "READY",
		)
	}
	_, _ = s.db.Exec(`UPDATE projects SET status=? WHERE id=?`, "IMAGES_READY", projectID)
	job := s.newJob("GENERATE_IMAGE", projectID, "")
	// sticker generation finished
	s.setJobProgress(job.ID, 100, "SUCCESS")
	return job, true
}

func (s *Store) ListStickers(projectID string) []*Sticker {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, _ := s.db.Query(`SELECT id,project_id,draft_id,image_url,transparent_url,status FROM stickers WHERE project_id=?`, projectID)
	defer rows.Close()
	out := []*Sticker{}
	for rows.Next() {
		st := &Sticker{}
		_ = rows.Scan(&st.ID, &st.ProjectID, &st.DraftID, &st.ImageURL, &st.TransparentURL, &st.Status)
		out = append(out, st)
	}
	return out
}

func (s *Store) RemoveBackground(projectID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, _ := s.db.Query(`SELECT id,image_url FROM stickers WHERE project_id=?`, projectID)
	defer rows.Close()
	for rows.Next() {
		var id, imageURL string
		_ = rows.Scan(&id, &imageURL)
		transparentURL, _ := s.ai.RemoveBackground(imageURL)
		_, _ = s.db.Exec(`UPDATE stickers SET transparent_url=? WHERE id=?`, transparentURL, id)
	}
	job := s.newJob("REMOVE_BG", projectID, "")
	// remove-bg finished
	s.setJobProgress(job.ID, 100, "SUCCESS")
	return job, true
}

func (s *Store) Export(projectID string) (*ExportResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.db.Exec(`UPDATE projects SET status=? WHERE id=?`, "DONE", projectID)
	return &ExportResponse{DownloadURL: "https://example.com/export.zip"}, true
}

func (s *Store) RegenerateSticker(stickerID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := s.db.QueryRow(`SELECT id,project_id,draft_id FROM stickers WHERE id=?`, stickerID)
	var id, projectID, draftID string
	if err := row.Scan(&id, &projectID, &draftID); err != nil {
		return nil, false
	}
	promptRow := s.db.QueryRow(`SELECT image_prompt FROM drafts WHERE id=?`, draftID)
	var prompt string
	_ = promptRow.Scan(&prompt)
	charInput := s.getCharacterInput(projectID)
	imageURL, _ := s.ai.GenerateImage(prompt, charInput)
	_, _ = s.db.Exec(`UPDATE stickers SET image_url=?, status=? WHERE id=?`, imageURL, "READY", stickerID)
	job := s.newJob("GENERATE_IMAGE", projectID, stickerID)
	// regenerate finished
	s.setJobProgress(job.ID, 100, "SUCCESS")
	return job, true
}

func (s *Store) GetJob(jobID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := s.db.QueryRow(`SELECT id,type,status,progress,error_message FROM jobs WHERE id=?`, jobID)
	j := &Job{}
	if err := row.Scan(&j.ID, &j.Type, &j.Status, &j.Progress, &j.ErrorMessage); err != nil {
		return nil, false
	}
	return j, true
}

func (s *Store) newJob(jobType string, projectID string, targetID string) *Job {
	id := newID("job")
	j := &Job{ID: id, Type: jobType, Status: "RUNNING", Progress: 0}
	_, _ = s.db.Exec(`INSERT INTO jobs (id,type,status,progress,error_message,project_id,target_id) VALUES (?,?,?,?,?,?,?)`,
		id, jobType, j.Status, j.Progress, "", projectID, targetID,
	)
	return j
}

func (s *Store) setJobProgress(jobID string, progress int, status string) {
	_, _ = s.db.Exec(`UPDATE jobs SET progress=?, status=? WHERE id=?`, progress, status, jobID)
}
