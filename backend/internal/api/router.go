package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
			writeJSON(w, http.StatusOK, []Project{})
		case http.MethodPost:
			writeJSON(w, http.StatusOK, Project{
				ID:           "proj_123",
				Status:       "DRAFT",
				StickerCount: 8,
			})
		default:
			writeStatus(w, http.StatusMethodNotAllowed)
		}
		return
	}

	// /projects/{projectId}
	if len(segments) == 2 && segments[0] == "projects" {
		projectID := segments[1]
		switch r.Method {
		case http.MethodGet, http.MethodPatch:
			writeJSON(w, http.StatusOK, Project{
				ID:           projectID,
				Title:        "LINE Sticker Project",
				Theme:        "Office Life",
				StickerCount: 8,
				Status:       "DRAFT_READY",
				CharacterID:  "char_456",
			})
		default:
			writeStatus(w, http.StatusMethodNotAllowed)
		}
		return
	}

	// /projects/{projectId}/character
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "character" {
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, Character{
				ID:                "char_456",
				SourceType:        "AI",
				ReferenceImageURL: "",
				Status:            "READY",
			})
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
			writeJSON(w, http.StatusOK, Job{
				ID:       "job_789",
				Type:     "GENERATE_DRAFT",
				Status:   "QUEUED",
				Progress: 0,
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/drafts
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "drafts" {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, []Draft{
				{
					ID:          "draft_1",
					ProjectID:   segments[1],
					Index:       1,
					Caption:     "早安",
					ImagePrompt: "橘貓揮手說早安",
					Status:      "DRAFT",
				},
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /drafts/{draftId}
	if len(segments) == 2 && segments[0] == "drafts" {
		if r.Method == http.MethodPatch {
			writeJSON(w, http.StatusOK, Draft{
				ID:     segments[1],
				Status: "DRAFT",
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/stickers:generate
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "stickers:generate" {
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, Job{
				ID:       "job_888",
				Type:     "GENERATE_IMAGE",
				Status:   "QUEUED",
				Progress: 0,
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /stickers/{stickerId}:regenerate
	if len(segments) == 2 && strings.HasPrefix(segments[0], "stickers") && strings.HasSuffix(segments[1], ":regenerate") {
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, Job{
				ID:       "job_777",
				Type:     "GENERATE_IMAGE",
				Status:   "QUEUED",
				Progress: 0,
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/stickers
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "stickers" {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, []Sticker{
				{
					ID:             "stk_1",
					ProjectID:      segments[1],
					DraftID:        "draft_1",
					ImageURL:       "https://example.com/sticker.png",
					TransparentURL: "https://example.com/sticker-transparent.png",
					Status:         "READY",
				},
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/stickers:remove-bg
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "stickers:remove-bg" {
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, Job{
				ID:       "job_999",
				Type:     "REMOVE_BG",
				Status:   "QUEUED",
				Progress: 0,
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /projects/{projectId}/export
	if len(segments) == 3 && segments[0] == "projects" && segments[2] == "export" {
		if r.Method == http.MethodPost {
			writeJSON(w, http.StatusOK, ExportResponse{
				DownloadURL: "https://example.com/export.zip",
			})
			return
		}
		writeStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// /jobs/{jobId}
	if len(segments) == 2 && segments[0] == "jobs" {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, Job{
				ID:       segments[1],
				Type:     "GENERATE_DRAFT",
				Status:   "RUNNING",
				Progress: 50,
			})
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
