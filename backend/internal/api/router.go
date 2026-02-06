package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

var store = NewStore()

func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/", apiHandler)
	return mux
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	// /projects
	if len(segments) == 1 && segments[0] == "projects" {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, store.ListProjects())
		case http.MethodPost:
			var req ProjectCreateRequest
			if !decodeJSON(w, r, &req) {
				return
			}
			p := store.CreateProject(req.Title, req.StickerCount)
			writeJSON(w, http.StatusOK, p)
		default:
			writeStatus(w, http.StatusMethodNotAllowed)
		}
		return
	}

	// /projects/{projectId}
	if len(segments) == 2 && segments[0] == "projects" {
		projectID := segments[1]
		switch r.Method {
		case http.MethodGet:
			if p, ok := store.GetProject(projectID); ok {
				writeJSON(w, http.StatusOK, p)
				return
			}
			writeStatus(w, http.StatusNotFound)
		case http.MethodPatch:
			var req ProjectUpdateRequest
			if !decodeJSON(w, r, &req) {
				return
			}
			if p, ok := store.UpdateProjectTheme(projectID, req.Theme); ok {
				writeJSON(w, http.StatusOK, p)
				return
			}
			writeStatus(w, http.StatusNotFound)
		default:
			writeStatus(w, http.StatusMethodNotAllowed)
		}
		return
	}

	// /projects/{projectId}/character
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "character" {
		if r.Method == http.MethodPost {
			var req CharacterCreateRequest
			if !decodeJSON(w, r, &req) {
				return
			}
			if c, ok := store.CreateCharacter(segments[1], req); ok {
				writeJSON(w, http.StatusOK, c)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /providers
	if len(segments) == 1 && segments[0] == "providers" {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, ListProviders())
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/ai-config
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "ai-config" {
		if r.Method == http.MethodPatch {
			var req AIConfigUpdateRequest
			if !decodeJSON(w, r, &req) {
				return
			}
			if p, ok := store.UpdateProjectAI(segments[1], req); ok {
				writeJSON(w, http.StatusOK, p)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/ai-pipeline
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "ai-pipeline" {
		if r.Method == http.MethodPatch {
			var req AIPipelineConfigRequest
			if !decodeJSON(w, r, &req) {
				return
			}
			if p, ok := store.UpdateProjectPipeline(segments[1], req); ok {
				writeJSON(w, http.StatusOK, p)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/ai-credentials
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "ai-credentials" {
		if r.Method == http.MethodPost {
			var req AICredentialsRequest
			if !decodeJSON(w, r, &req) {
				return
			}
			if ok := store.StoreAICredentials(segments[1], req); ok {
				writeStatus(w, http.StatusNoContent)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/ai-verify
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "ai-verify" {
		if r.Method == http.MethodPost {
			if err := store.VerifyAICredentials(segments[1]); err != nil {
				writeStatus(w, http.StatusBadRequest)
				return
			}
			writeStatus(w, http.StatusNoContent)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/theme:suggest
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "theme:suggest" {
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, ThemeSuggestResponse{
				Suggestions: []string{"Office Life", "Commute", "Weekend"},
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/drafts:generate
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "drafts:generate" {
		if r.Method == http.MethodPost {
			if job, ok := store.GenerateDrafts(segments[1]); ok {
				writeJSON(w, http.StatusOK, job)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/drafts
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "drafts" {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, store.ListDrafts(segments[1]))
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /drafts/{draftId}
	if len(segments) == 2 && segments[0] == "drafts" {
		if r.Method == http.MethodPatch {
			var req DraftUpdateRequest
			if !decodeJSON(w, r, &req) {
				return
			}
			if d, ok := store.UpdateDraft(segments[1], req); ok {
				writeJSON(w, http.StatusOK, d)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/stickers:generate
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "stickers:generate" {
		if r.Method == http.MethodPost {
			if job, ok := store.GenerateStickers(segments[1]); ok {
				writeJSON(w, http.StatusOK, job)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /stickers/{stickerId}:regenerate
	if len(segments) == 2 && segments[0] == "stickers" && strings.HasSuffix(segments[1], ":regenerate") {
		if r.Method == http.MethodPost {
			stickerID := strings.TrimSuffix(segments[1], ":regenerate")
			if job, ok := store.RegenerateSticker(stickerID); ok {
				writeJSON(w, http.StatusOK, job)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/stickers
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "stickers" {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, store.ListStickers(segments[1]))
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/stickers:remove-bg
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "stickers:remove-bg" {
		if r.Method == http.MethodPost {
			if job, ok := store.RemoveBackground(segments[1]); ok {
				writeJSON(w, http.StatusOK, job)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/export
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "export" {
		if r.Method == http.MethodPost {
			if res, ok := store.Export(segments[1]); ok {
				writeJSON(w, http.StatusOK, res)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /jobs/{jobId}
	if len(segments) == 2 && segments[0] == "jobs" {
		if r.Method == http.MethodGet {
			if j, ok := store.GetJob(segments[1]); ok {
				writeJSON(w, http.StatusOK, j)
				return
			}
			writeStatus(w, http.StatusNotFound)
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	writeStatus(w, http.StatusNotFound)
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeStatus(w, http.StatusBadRequest)
		return false
	}
	return true
}
