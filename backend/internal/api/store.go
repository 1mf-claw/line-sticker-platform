package api

import (
	"fmt"
	"sync"
)

type Store struct {
	mu         sync.Mutex
	projects   map[string]*Project
	characters map[string]*Character
	drafts     map[string]*Draft
	stickers   map[string]*Sticker
	jobs       map[string]*Job
	counters   map[string]int
}

func NewStore() *Store {
	return &Store{
		projects:   map[string]*Project{},
		characters: map[string]*Character{},
		drafts:     map[string]*Draft{},
		stickers:   map[string]*Sticker{},
		jobs:       map[string]*Job{},
		counters:   map[string]int{},
	}
}

func (s *Store) nextID(prefix string) string {
	s.counters[prefix]++
	return fmt.Sprintf("%s_%d", prefix, s.counters[prefix])
}

func (s *Store) CreateProject(title string, stickerCount int) *Project {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID("proj")
	p := &Project{
		ID:           id,
		Title:        title,
		StickerCount: stickerCount,
		Status:       "DRAFT",
	}
	s.projects[id] = p
	return p
}

func (s *Store) UpdateProjectTheme(projectID string, theme string) (*Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.projects[projectID]
	if !ok {
		return nil, false
	}
	p.Theme = theme
	return p, true
}

func (s *Store) GetProject(projectID string) (*Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.projects[projectID]
	return p, ok
}

func (s *Store) ListProjects() []*Project {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Project, 0, len(s.projects))
	for _, p := range s.projects {
		out = append(out, p)
	}
	return out
}

func (s *Store) CreateCharacter(projectID string, req CharacterCreateRequest) (*Character, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.projects[projectID]
	if !ok {
		return nil, false
	}
	id := s.nextID("char")
	c := &Character{
		ID:                id,
		SourceType:        req.SourceType,
		ReferenceImageURL: req.ReferenceImageURL,
		Status:            "READY",
	}
	s.characters[id] = c
	p.CharacterID = id
	return c, true
}

func (s *Store) GenerateDrafts(projectID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.projects[projectID]
	if !ok {
		return nil, false
	}
	p.Status = "GENERATING_DRAFTS"
	for i := 1; i <= p.StickerCount; i++ {
		id := s.nextID("draft")
		s.drafts[id] = &Draft{
			ID:          id,
			ProjectID:   projectID,
			Index:       i,
			Caption:     fmt.Sprintf("草稿 %d", i),
			ImagePrompt: fmt.Sprintf("主角依主題動作 %d", i),
			Status:      "DRAFT",
		}
	}
	p.Status = "DRAFT_READY"
	job := s.newJob("GENERATE_DRAFT")
	return job, true
}

func (s *Store) ListDrafts(projectID string) []*Draft {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []*Draft{}
	for _, d := range s.drafts {
		if d.ProjectID == projectID {
			out = append(out, d)
		}
	}
	return out
}

func (s *Store) UpdateDraft(draftID string, req DraftUpdateRequest) (*Draft, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.drafts[draftID]
	if !ok {
		return nil, false
	}
	if req.Caption != "" {
		d.Caption = req.Caption
	}
	if req.ImagePrompt != "" {
		d.ImagePrompt = req.ImagePrompt
	}
	return d, true
}

func (s *Store) GenerateStickers(projectID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.projects[projectID]
	if !ok {
		return nil, false
	}
	p.Status = "GENERATING_IMAGES"
	for _, d := range s.drafts {
		if d.ProjectID != projectID {
			continue
		}
		id := s.nextID("stk")
		s.stickers[id] = &Sticker{
			ID:        id,
			ProjectID: projectID,
			DraftID:   d.ID,
			ImageURL:  "https://example.com/sticker.png",
			Status:    "READY",
		}
	}
	p.Status = "IMAGES_READY"
	job := s.newJob("GENERATE_IMAGE")
	return job, true
}

func (s *Store) ListStickers(projectID string) []*Sticker {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []*Sticker{}
	for _, st := range s.stickers {
		if st.ProjectID == projectID {
			out = append(out, st)
		}
	}
	return out
}

func (s *Store) RemoveBackground(projectID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, st := range s.stickers {
		if st.ProjectID == projectID {
			st.TransparentURL = "https://example.com/sticker-transparent.png"
		}
	}
	job := s.newJob("REMOVE_BG")
	return job, true
}

func (s *Store) Export(projectID string) (*ExportResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.projects[projectID]
	if !ok {
		return nil, false
	}
	p.Status = "DONE"
	return &ExportResponse{DownloadURL: "https://example.com/export.zip"}, true
}

func (s *Store) RegenerateSticker(stickerID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stickers[stickerID]; !ok {
		return nil, false
	}
	job := s.newJob("GENERATE_IMAGE")
	return job, true
}

func (s *Store) GetJob(jobID string) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[jobID]
	return j, ok
}

func (s *Store) newJob(jobType string) *Job {
	id := s.nextID("job")
	j := &Job{ID: id, Type: jobType, Status: "SUCCESS", Progress: 100}
	s.jobs[id] = j
	return j
}
