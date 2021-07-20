package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type Speaker struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name"`
	Pronouns *string `json:"pronouns"`
	Employer *string `json:"employer"`
	Title    *string `json:"title"`
	Photo    *string `json:"photo"`
}

func (a API) handleGetSpeaker(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetSpeaker(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrSpeakerNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Speakers: []Speaker{ap}})
}

func (a API) handlePostSpeaker(w http.ResponseWriter, r *http.Request) {
	var ap Speaker
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
	err = a.Storer.CreateSpeaker(ap)
	if err != nil {
		if err == ErrSpeakerAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{Speakers: []Speaker{ap}})
}

func (a API) handlePutSpeaker(w http.ResponseWriter, r *http.Request) {
	var ap Speaker
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID = trout.RequestVars(r).Get("id")
	err = a.Storer.UpdateSpeaker(ap)
	if err != nil {
		if err == ErrSpeakerNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Speakers: []Speaker{ap}})
}

func (a API) handleDeleteSpeaker(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetSpeaker(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrSpeakerNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteSpeaker(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrSpeakerNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Speakers: []Speaker{ap}})
}
