package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crikke/cms/pkg/contentmanagement/api/handlers"
	"github.com/crikke/cms/pkg/contentmanagement/app"
	"github.com/crikke/cms/pkg/contentmanagement/app/command"
	"github.com/crikke/cms/pkg/contentmanagement/app/query"
	"github.com/go-chi/chi/v5"
)

type key string

const tagKey = key("tag")

func NewWorkspaceRoute(app app.App) http.Handler {
	r := chi.NewRouter()

	wsHandler := handlers.WorkspaceHandler{App: app}

	r.Post("/", createWorkspace(app))
	r.Route("/{workspace}", func(r chi.Router) {
		r.Use(wsHandler.WorkspaceParamContext)

		r.Put("/", updateWorkspace(app))
		r.Get("/", getWorkspace(app))

		r.Route("/tags", func(r chi.Router) {
			r.Get("/", listTags(app))
			r.Post("/", createTag(app))
			r.Route("/{tag}", func(r chi.Router) {
				r.Use(tagContext)
				r.Get("/", getTag(app))
				r.Put("/", updateTag(app))
				r.Delete("/", deleteTag(app))
			})
		})
	})
	return r
}

func tagContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tag := chi.URLParam(r, "tag")

		if tag == "" {
			return
		}

		ctx := context.WithValue(r.Context(), tagKey, tag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withTag(ctx context.Context) string {

	tag := ctx.Value(tagKey)
	return tag.(string)
}

// getWorkspace 		godoc
// @Summary 		Get workspace
// @Description 	Get workspace by id
// @Tags 			workspace
// @Produces 		json
// @Param			workspace			path	string	true 	"uuid formatted ID." format(uuid)
// @Success			200			{object}	workspace.Workspace
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{workspace} [get]
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

// createWorkspace 		godoc
// @Summary 		Create workspace
// @Description 	Create a new workspace
// @Tags 			workspace
// @Consumes 		json
// @Produces 		json
// @Param			workspace	body command.CreateWorkspace true "workspace body"
// @Success			201			{object}	workspace.Workspace
// @Header						201			{string}	Location
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace [post]
func createWorkspace(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cmd := &command.CreateWorkspace{}
		err := json.NewDecoder(r.Body).Decode(cmd)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := app.Commands.WorkspaceCommands.CreateWorkspace.Handle(r.Context(), *cmd)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ws, err := app.Queries.WorkspaceQueries.GetWorkspace.Handle(r.Context(), id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(&ws)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Location", fmt.Sprintf("%s/%s", r.URL.String(), id.String()))
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

type TagBody struct {
	Name string
	ID   string
}

// updateWorkspace 		godoc
// @Summary 		Update workspace
// @Description 	Update a new workspace
// @Tags 			workspace
// @Consumes 		json
// @Produces 		json
// @Param			workspace			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			workspace	body workspace.Workspace true "workspace body"
// @Success			200			{object}	workspace.Workspace
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{workspace} [put]
func updateWorkspace(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws := handlers.WithWorkspace(r.Context())
		tag := &TagBody{}

		err := json.NewDecoder(r.Body).Decode(tag)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = app.Commands.WorkspaceCommands.UpdateTag.Handle(r.Context(), command.UpdateTag{
			WorkspaceId: ws.ID,
			Name:        tag.Name,
			Id:          tag.ID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(tag)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(data)
	}
}

// createTag 		godoc
// @Summary 		Create tag
// @Description 	Creates a tag in given workspace
// @Tags 			workspace
// @Consumes 		json
// @Produces 		json
// @Param			workspace			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			body			body 	TagBody true "Tag"
// @Success			201			{object}	workspace.Workspace
// @Header						201			{string}	Location
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{workspace}/tag [post]
func createTag(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ws := handlers.WithWorkspace(r.Context())
		tag := &TagBody{}

		err := json.NewDecoder(r.Body).Decode(tag)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = app.Commands.WorkspaceCommands.UpdateTag.Handle(r.Context(), command.UpdateTag{
			WorkspaceId: ws.ID,
			Name:        tag.Name,
			Id:          tag.ID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		url := r.URL.String()

		w.Header().Add("Location", fmt.Sprintf("%s/%s", url, tag.ID))
		w.WriteHeader(http.StatusCreated)
	}
}

// createTag 		godoc
// @Summary 		List all tags in workspace
// @Description 	List all tags in workspace
// @Tags 			workspace
// @Param			workspace			path	string	true 	"uuid formatted ID." format(uuid)
// @Produces 		json
// @Success			200			{object}	[]query.Tag
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{workspace}/tag [get]
func listTags(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workspace := handlers.WithWorkspace(r.Context())

		tags, err := app.Queries.WorkspaceQueries.ListTags.Handle(r.Context(), workspace.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}
}

// getTag 		godoc
// @Summary 		Get tag
// @Description 	Get tag by id
// @Tags 			workspace
// @Param			workspace			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			tag					path	string	true 	"name"
// @Produces 		json
// @Success			200			{object}	query.Tag
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{workspace}/tag/{tag} [get]
func getTag(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagId := withTag(r.Context())
		workspace := handlers.WithWorkspace(r.Context())
		tag, err := app.Queries.WorkspaceQueries.GetTag.Handle(r.Context(), query.GetTag{Id: tagId, WorkspaceId: workspace.ID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(&tag)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(data)
	}
}

// updateTag 		godoc
// @Summary 		Update tag
// @Description 	Update tag by id
// @Tags 			workspace
// @Produces 		json
// @Param			workspace			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			tag					path	string	true 	"name"
// @Param			body	body string true "Tag"
// @Success			200			{object}	query.Tag
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{workspace}/tag/{tag} [put]
func updateTag(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws := handlers.WithWorkspace(r.Context())
		tag := &TagBody{}

		err := json.NewDecoder(r.Body).Decode(tag)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = app.Commands.WorkspaceCommands.UpdateTag.Handle(r.Context(), command.UpdateTag{
			WorkspaceId: ws.ID,
			Name:        tag.Name,
			Id:          tag.ID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updated, err := app.Queries.WorkspaceQueries.GetTag.Handle(r.Context(), query.GetTag{Id: tag.ID, WorkspaceId: ws.ID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(&updated)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(data)
	}
}

// deleteTag 		godoc
// @Summary 		Delete tag
// @Description 	Delete tag by id
// @Tags 			workspace
// @Produces 		json
// @Param			workspace			path	string	true 	"uuid formatted ID." format(uuid)
// @Param			tag					path	string	true 	"tag id"
// @Failure			default		{object}	models.GenericError
// @Router			/contentmanagement/workspace/{workspace}/tag/{tag} [delete]
func deleteTag(app app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws := handlers.WithWorkspace(r.Context())
		tagId := withTag(r.Context())

		err := app.Commands.WorkspaceCommands.DeleteTag.Handle(r.Context(), command.DeleteTag{
			WorkspaceId: ws.ID,
			Id:          tagId,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
