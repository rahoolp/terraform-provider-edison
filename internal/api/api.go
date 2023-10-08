package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
)

type API struct {
	Storer *Storer
}

func (a API) Server(baseURL string) http.Handler {
	var router trout.Router
	router.SetPrefix(baseURL)

	router.Endpoint("/eastores").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostEAStore))
	router.Endpoint("/eastores/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetEAStore))
	router.Endpoint("/eastores/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutEAStore))
	router.Endpoint("/eastores/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteEAStore))

	router.Endpoint("/ehsclusters").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostEHSCluster))
	router.Endpoint("/ehsclusters/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetEHSCluster))
	router.Endpoint("/ehsclusters/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutEHSCluster))
	router.Endpoint("/ehsclusters/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteEHSCluster))

	router.Endpoint("/aws").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostAW))
	router.Endpoint("/aws/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetAW))
	router.Endpoint("/aws/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutAW))
	router.Endpoint("/aws/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteAW))

	router.Endpoint("/avs").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostAV))
	router.Endpoint("/avs/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetAV))
	router.Endpoint("/avs/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutAV))
	router.Endpoint("/avs/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteAV))

	return api.NegotiateMiddleware(router)
}

func isAuthenticated(r *http.Request) bool {
	return r.Header.Get("Authentication") == "secrettoken"
}

type Response struct {
	Errors      []api.RequestError `json:"errors,omitempty"`
	Status      int                `json:"-"`
	EAStores    []EAStore          `json:"eastores,omitempty"`
	EHSClusters []EHSCluster       `json:"ehsclusters,omitempty"`
	AWs         []AW               `json:"aws,omitempty"`
	AVs         []AV               `json:"avs,omitempty"`
}
