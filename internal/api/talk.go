package api

import (
	"net/http"

	"darlinggo.co/api"
	"darlinggo.co/trout/v2"
	"github.com/hashicorp/go-uuid"
)

type Talk struct {
	ID              string                   `json:"id,omitempty"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	DurationMinutes int64                    `json:"durationMinutes"`
	Prerecorded     bool                     `json:"prerecorded"`
	SpeakerIDs      []string                 `json:"speakerIDs"`
	Recordings      map[string]TalkRecording `json:"recordings"`
}

type TalkRecording struct {
	Resolution TalkRecordingResolution `json:"resolution"`
	Codec      string                  `json:"codec"`
}

type TalkRecordingResolution struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

func (a API) handleGetTalk(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetTalk(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrTalkNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Talks: []Talk{ap}})
}

func (a API) handlePostTalk(w http.ResponseWriter, r *http.Request) {
	var ap Talk
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
	ap.Recordings = map[string]TalkRecording{}
	for _, id := range ap.SpeakerIDs {
		speaker, err := a.Storer.GetSpeaker(id)
		if err != nil {
			api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
			return
		}
		ap.Recordings[speaker.Name] = TalkRecording{
			Codec: "h264",
			Resolution: TalkRecordingResolution{
				Width:  3840,
				Height: 2160,
			},
		}
	}
	err = a.Storer.CreateTalk(ap)
	if err != nil {
		if err == ErrTalkAlreadyExists {
			api.Encode(w, r, http.StatusBadRequest, Response{Errors: []api.RequestError{{Field: "/id", Slug: api.RequestErrConflict}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusCreated, Response{Talks: []Talk{ap}})
}

func (a API) handlePutTalk(w http.ResponseWriter, r *http.Request) {
	var ap Talk
	err := api.Decode(r, &ap)
	if err != nil {
		api.Encode(w, r, http.StatusBadRequest, Response{Errors: api.InvalidFormatError})
		return
	}
	ap.ID = trout.RequestVars(r).Get("id")
	err = a.Storer.UpdateTalk(ap)
	if err != nil {
		if err == ErrTalkNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Talks: []Talk{ap}})
}

func (a API) handleDeleteTalk(w http.ResponseWriter, r *http.Request) {
	ap, err := a.Storer.GetTalk(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrTalkNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	err = a.Storer.DeleteTalk(trout.RequestVars(r).Get("id"))
	if err != nil {
		if err == ErrTalkNotFound {
			api.Encode(w, r, http.StatusNotFound, Response{Errors: []api.RequestError{{Param: "id", Slug: api.RequestErrNotFound}}})
			return
		}
		api.Encode(w, r, http.StatusInternalServerError, Response{Errors: api.ActOfGodError})
		return
	}
	api.Encode(w, r, http.StatusOK, Response{Talks: []Talk{ap}})
}
