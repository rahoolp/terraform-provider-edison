package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type AW struct {
	ID              string `json:"id,omitempty"`
	ConcurrentUsers int    `json:"concurrent_users"`
	EHSClusterID    string `json:"ehs_cluster_id"`
	DicomEndPoint   string `json:"dicom_endpoint"`
	DNSEndPoint     string `json:"dns_endpoint,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

func (a API) handleGetAW(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetAW(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAWNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AWs: []AW{ap}})
}

func (a API) handlePostAW(w http.ResponseWriter, r *http.Request) {
	var ap AW
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
	err = a.Storer.CreateAW(ap)
	if err != nil {
		if err == ErrAWAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{AWs: []AW{ap}})
}

func (a API) handlePutAW(w http.ResponseWriter, r *http.Request) {
	var ap AW
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID = trout.RequestVars(r).Get("id")
	err = a.Storer.UpdateAW(ap)
	if err != nil {
		if err == ErrAWNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AWs: []AW{ap}})
}

func (a API) handleDeleteAW(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetAW(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAWNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteAW(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrAWNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{AWs: []AW{ap}})
}
