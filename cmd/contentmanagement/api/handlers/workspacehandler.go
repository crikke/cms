package handlers

import (
	"context"
	"net/http"

	"github.com/crikke/cms/cmd/contentmanagement/app"
	"github.com/crikke/cms/pkg/workspace"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type key string

const workspacekey = key("workspace")
const WorkspaceQueryParam = "workspace"

type WorkspaceHandler struct {
	App app.App
}

func (h WorkspaceHandler) WorkspaceQueryContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		param := r.URL.Query().Get(WorkspaceQueryParam)

		if param == "" {
			http.Error(w, "workspace param: missing", http.StatusBadRequest)
			return
		}

		uid, err := uuid.Parse(param)

		if err != nil {
			http.Error(w, "workspace param: bad format", http.StatusBadRequest)
			return
		}

		ws, err := h.App.Queries.WorkspaceQueries.GetWorkspace.Handle(r.Context(), uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), workspacekey, ws)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h WorkspaceHandler) WorkspaceParamContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		param := chi.URLParam(r, WorkspaceQueryParam)

		if param == "" {
			http.Error(w, "workspace param: missing", http.StatusBadRequest)
			return
		}

		uid, err := uuid.Parse(param)

		if err != nil {
			http.Error(w, "workspace param: bad format", http.StatusBadRequest)
			return
		}

		ws, err := h.App.Queries.WorkspaceQueries.GetWorkspace.Handle(r.Context(), uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), workspacekey, ws)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithWorkspace(ctx context.Context) workspace.Workspace {
	if ws := ctx.Value(workspacekey); ws != nil {
		return ws.(workspace.Workspace)
	}
	return workspace.Workspace{}
}
