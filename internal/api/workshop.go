package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type Workshop struct {
	ID              string                       `json:"id,omitempty"`
	Title           string                       `json:"title"`
	Description     string                       `json:"description"`
	DurationMinutes int64                        `json:"durationMinutes"`
	Presenters      map[string]WorkshopPresenter `json:"presenters"`
	MeetingInfo     WorkshopMeetingInfo          `json:"meetingInfo"`
}

type WorkshopPresenter struct {
	Title    *string `json:"title"`
	Employer *string `json:"employer"`
	Pronouns *string `json:"pronouns"`
}

type WorkshopMeetingInfo struct {
	URL      string  `json:"url"`
	Password *string `json:"password"`
}

func (a API) handleGetWorkshop(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetWorkshop(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrWorkshopNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Workshops: []Workshop{ap}})
}

func (a API) handlePostWorkshop(w http.ResponseWriter, r *http.Request) {
	var ap Workshop
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID, err = uuid.GenerateUUID()
	if err != nil {
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.CreateWorkshop(ap)
	if err != nil {
		if err == ErrWorkshopAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{Workshops: []Workshop{ap}})
}

func (a API) handlePutWorkshop(w http.ResponseWriter, r *http.Request) {
	var ap Workshop
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID = trout.RequestVars(r).Get("id")
	err = a.Storer.UpdateWorkshop(ap)
	if err != nil {
		if err == ErrWorkshopNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Workshops: []Workshop{ap}})
}

func (a API) handleDeleteWorkshop(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetWorkshop(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrWorkshopNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteWorkshop(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrWorkshopNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Workshops: []Workshop{ap}})
}
