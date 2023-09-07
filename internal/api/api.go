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

	router.Endpoint("/talks").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostTalk))
	router.Endpoint("/talks/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetTalk))
	router.Endpoint("/talks/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutTalk))
	router.Endpoint("/talks/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteTalk))

	router.Endpoint("/speakers").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostSpeaker))
	router.Endpoint("/speakers/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetSpeaker))
	router.Endpoint("/speakers/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutSpeaker))
	router.Endpoint("/speakers/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteSpeaker))

	router.Endpoint("/workshops").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostWorkshop))
	router.Endpoint("/workshops/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetWorkshop))
	router.Endpoint("/workshops/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutWorkshop))
	router.Endpoint("/workshops/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteWorkshop))

	router.Endpoint("/eastores").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostEAStore))
	router.Endpoint("/eastores/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetEAStore))
	router.Endpoint("/eastores/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutEAStore))
	router.Endpoint("/eastores/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteEAStore))

	router.Endpoint("/ehsclusters").Methods(http.MethodPost).Handler(http.HandlerFunc(a.handlePostEHSCluster))
	router.Endpoint("/ehsclusters/{id}").Methods(http.MethodGet).Handler(http.HandlerFunc(a.handleGetEHSCluster))
	router.Endpoint("/ehsclusters/{id}").Methods(http.MethodPut).Handler(http.HandlerFunc(a.handlePutEHSCluster))
	router.Endpoint("/ehsclusters/{id}").Methods(http.MethodDelete).Handler(http.HandlerFunc(a.handleDeleteEHSCluster))

	return api.NegotiateMiddleware(router)
}

func isAuthenticated(r *http.Request) bool {
	return r.Header.Get("Authentication") == "secrettoken"
}

type Response struct {
	Talks       []Talk             `json:"talks,omitempty"`
	Speakers    []Speaker          `json:"speakers,omitempty"`
	Workshops   []Workshop         `json:"workshops,omitempty"`
	Errors      []api.RequestError `json:"errors,omitempty"`
	Status      int                `json:"-"`
	EAStores    []EAStore          `json:"eastores,omitempty"`
	EHSClusters []EHSCluster       `json:"ehsclusters,omitempty"`
}
