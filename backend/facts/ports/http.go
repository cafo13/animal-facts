package ports

import (
	"context"
	"net/http"

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
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var appFacts []query.Fact

	if user.Role == "trainer" {
		appFacts, err = h.app.Queries.AllFacts.Handle(r.Context(), query.AllFacts{})
	} else {
		appFacts, err = h.app.Queries.FactsForUser.Handle(r.Context(), query.FactsForUser{User: user})
	}

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	facts := appFactsToResponse(appFacts)
	factsResp := Facts{facts}

	render.Respond(w, r, factsResp)
}

func (h HttpServer) CreateFact(w http.ResponseWriter, r *http.Request) {
	postFact := PostFact{}
	if err := render.Decode(r, &postFact); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "attendee" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	cmd := command.ScheduleFact{
		FactUUID: uuid.New().String(),
		UserUUID:     user.UUID,
		UserName:     user.DisplayName,
		FactTime: postFact.Time,
		Notes:        postFact.Notes,
	}
	err = h.app.Commands.ScheduleFact.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.Header().Set("content-location", "/facts/"+cmd.FactUUID)
	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) CancelFact(w http.ResponseWriter, r *http.Request, factUUID string) {
	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.CancelFact.Handle(r.Context(), command.CancelFact{
		FactUUID: factUUID,
		User:         user,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RescheduleFact(w http.ResponseWriter, r *http.Request, factUUID string) {
	rescheduleFact := PostFact{}
	if err := render.Decode(r, &rescheduleFact); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.RescheduleFact.Handle(r.Context(), command.RescheduleFact{
		User:         user,
		FactUUID: factUUID,
		NewTime:      rescheduleFact.Time,
		NewNotes:     rescheduleFact.Notes,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RequestRescheduleFact(w http.ResponseWriter, r *http.Request, factUUID string) {
	rescheduleFact := PostFact{}
	if err := render.Decode(r, &rescheduleFact); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.RequestFactReschedule.Handle(r.Context(), command.RequestFactReschedule{
		User:         user,
		FactUUID: factUUID,
		NewTime:      rescheduleFact.Time,
		NewNotes:     rescheduleFact.Notes,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) ApproveRescheduleFact(w http.ResponseWriter, r *http.Request, factUUID string) {
	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.ApproveFactReschedule.Handle(r.Context(), command.ApproveFactReschedule{
		User:         user,
		FactUUID: factUUID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func (h HttpServer) RejectRescheduleFact(w http.ResponseWriter, r *http.Request, factUUID string) {
	user, err := newDomainUserFromAuthUser(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.RejectFactReschedule.Handle(r.Context(), command.RejectFactReschedule{
		User:         user,
		FactUUID: factUUID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
}

func appFactsToResponse(appFacts []query.Fact) []Fact {
	var facts []Fact
	for _, tm := range appFacts {
		t := Fact{
			CanBeCancelled:     tm.CanBeCancelled,
			MoveProposedBy:     tm.MoveProposedBy,
			MoveRequiresAccept: tm.CanBeCancelled,
			Notes:              tm.Notes,
			ProposedTime:       tm.ProposedTime,
			Time:               tm.Time,
			User:               tm.User,
			UserUuid:           tm.UserUUID,
			Uuid:               tm.UUID,
		}

		facts = append(facts, t)
	}

	return facts
}

func newDomainUserFromAuthUser(ctx context.Context) (fact.User, error) {
	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		return fact.User{}, err
	}

	userType, err := fact.NewUserTypeFromString(user.Role)
	if err != nil {
		return fact.User{}, err
	}

	return fact.NewUser(user.UUID, userType)
}
