package workspace

import (
	"encoding/json"
	"net/http"

	"github.com/crikke/cms/pkg/contentmanagement/api/handlers"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/go-chi/chi/v5"
)

func NewWorkspaceRoute(app app.App) http.Handler {
	r := chi.NewRouter()

	wsHandler := handlers.WorkspaceHandler{App: app}

	r.Post("/", createWorkspace(app))
	r.Route("/{id}", func(r chi.Router) {
		r.Use(wsHandler.WorkspaceParamContext)

		r.Put("/", updateWorkspace(app))
		r.Get("/", getWorkspace(app))
	})
	return r
}

// getWorkspace 		godoc
// @Summary 		Get workspace
// @Description 	Get workspace by id
// @Tags 			workspace
// @Produces 		json
// @Param			id			path	string	true 	"uuid formatted ID." format(uuid)
// @Success			200			{object}	workspace.Workspace
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{id} [get]
func getWorkspace(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws := handlers.WithWorkspace(r.Context())

		data, err := json.Marshal(&ws)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

type CreateWorkspaceRequest struct {
	Name        string
	Description string
}

// createWorkspace 		godoc
// @Summary 		Create workspace
// @Description 	Create a new workspace
// @Tags 			workspace
// @Consumes 		json
// @Produces 		json
// @Param			workspace	body CreateWorkspaceRequest true "workspace body"
// @Success			201			{object}	workspace.Workspace
// @Header						201			{string}	Location
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace [post]
func createWorkspace(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// updateWorkspace 		godoc
// @Summary 		Update workspace
// @Description 	Update a new workspace
// @Tags 			workspace
// @Consumes 		json
// @Produces 		json
// @Param			id			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			workspace	body workspace.Workspace true "workspace body"
// @Success			200			{object}	workspace.Workspace
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{id} [put]
func updateWorkspace(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
