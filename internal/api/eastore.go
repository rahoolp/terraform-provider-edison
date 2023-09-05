package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type EAStore struct {
	ID               string `json:"id,omitempty"`
	PartitionSpaceTB int64  `json:"partition_space_tb"`
	IPAddress        string `json:"ip_address"`
	IPPort           string `json:"ip_port"`
	AET              string `json:"aet"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	//DeletedAt        sql.NullString `db:"deleted_at" json:"-"`
}

func (a API) handleGetEAStore(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetEAStore(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrEAStoreNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{EAStores: []EAStore{ap}})
}

func (a API) handlePostEAStore(w http.ResponseWriter, r *http.Request) {
	var ap EAStore
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
	err = a.Storer.CreateEAStore(ap)
	if err != nil {
		if err == ErrEAStoreAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{EAStores: []EAStore{ap}})
}

func (a API) handlePutEAStore(w http.ResponseWriter, r *http.Request) {
	var ap EAStore
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID = trout.RequestVars(r).Get("id")
	err = a.Storer.UpdateEAStore(ap)
	if err != nil {
		if err == ErrEAStoreNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{EAStores: []EAStore{ap}})
}

func (a API) handleDeleteEAStore(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetEAStore(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrEAStoreNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteEAStore(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrEAStoreNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{EAStores: []EAStore{ap}})
}
