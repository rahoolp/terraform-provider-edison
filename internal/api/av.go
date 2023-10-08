package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type AV struct {
	ID           string `json:"id,omitempty"`
	AccountID    string `json:"account_id"`
	TenantID     string `json:"tenant_id"`
	TenantFolder string `json:"tenant_folder,omitempty"`
	TenantQueue  string `json:"tenant_queue,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

func (a API) handleGetAV(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetAV(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAVNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AVs: []AV{ap}})
}

func (a API) handlePostAV(w http.ResponseWriter, r *http.Request) {
	var ap AV
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
	err = a.Storer.CreateAV(ap)
	if err != nil {
		if err == ErrAVAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{AVs: []AV{ap}})
}

func (a API) handlePutAV(w http.ResponseWriter, r *http.Request) {
	var ap AV
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID = trout.RequestVars(r).Get("id")
	err = a.Storer.UpdateAV(ap)
	if err != nil {
		if err == ErrAVNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AVs: []AV{ap}})
}

func (a API) handleDeleteAV(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetAV(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAVNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteAV(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAVNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AVs: []AV{ap}})
}
