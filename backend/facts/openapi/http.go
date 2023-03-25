package openapi

import (
	"net/http"

	"github.com/cafo13/animal-facts/backend/common/auth"
	"github.com/cafo13/animal-facts/backend/common/server/httperr"
	"github.com/cafo13/animal-facts/backend/facts/app"
	"github.com/cafo13/animal-facts/backend/facts/app/command"
	"github.com/cafo13/animal-facts/backend/facts/app/query"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) GetFact(w http.ResponseWriter, r *http.Request) {
	appFact, err := h.app.Queries.RandomFact.Handle(r.Context(), query.RandomFact{})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	factResp := appFactToResponse(appFact)

	render.Respond(w, r, factResp)
}

func (h HttpServer) GetFactByID(w http.ResponseWriter, r *http.Request, factUUID uuid.UUID) {
	appFact, err := h.app.Queries.FactByID.Handle(r.Context(), query.FactByID{UUID: factUUID.String()})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	factResp := appFactToResponse(appFact)

	render.Respond(w, r, factResp)
}

func (h HttpServer) CreateFact(w http.ResponseWriter, r *http.Request) {
	postFact := Fact{}
	if err := render.Decode(r, &postFact); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "admin" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	cmd := command.CreateFact{
		FactUUID:   uuid.New(),
		FactText:   postFact.Text,
		FactSource: postFact.Source,
	}
	err = h.app.Commands.CreateFact.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.Header().Set("content-location", "/fact/"+cmd.FactUUID.String())
	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) UpdateFactByID(w http.ResponseWriter, r *http.Request, factUUID uuid.UUID) {
	updateFact := UpdateFact{}
	if err := render.Decode(r, &updateFact); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "admin" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	err = h.app.Commands.UpdateFact.Handle(r.Context(), command.UpdateFact{
		FactUUID:  factUUID,
		NewText:   updateFact.Text,
		NewSource: updateFact.Source,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) DeleteFactByID(w http.ResponseWriter, r *http.Request, factUUID uuid.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "admin" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	err = h.app.Commands.DeleteFact.Handle(r.Context(), command.DeleteFact{
		FactUUID: factUUID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func appFactToResponse(appFact query.Fact) Fact {
	return Fact{
		Uuid:   appFact.UUID,
		Text:   appFact.Text,
		Source: appFact.Source,
	}
}
