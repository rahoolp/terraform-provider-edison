package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type EHSCluster struct {
	ID      string `json:"id,omitempty"`
	Region  string `json:"region"`
	Profile string `json:"profile"`
	Release string `json:"release"`
	Tag     string `json:"tag"`
	//DicomEndPoint     string `json:"dicom_endpoint"`
	APIServerEndPoint string `json:"api_server_endpoint,omitempty"`
	VPC               string `json:"vpc,omitempty"`
	ClusterName       string `json:"cluster_name,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

func (a API) handleGetEHSCluster(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetEHSCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrEHSClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{EHSClusters: []EHSCluster{ap}})
}

func (a API) handlePostEHSCluster(w http.ResponseWriter, r *http.Request) {
	var ap EHSCluster
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
	err = a.Storer.CreateEHSCluster(ap)
	if err != nil {
		if err == ErrEHSClusterAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{EHSClusters: []EHSCluster{ap}})
}

func (a API) handlePutEHSCluster(w http.ResponseWriter, r *http.Request) {
	var ap EHSCluster
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID = trout.RequestVars(r).Get("id")
	err = a.Storer.UpdateEHSCluster(ap)
	if err != nil {
		if err == ErrEHSClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{EHSClusters: []EHSCluster{ap}})
}

func (a API) handleDeleteEHSCluster(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetEHSCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrEHSClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteEHSCluster(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrEHSClusterNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{EHSClusters: []EHSCluster{ap}})
}
