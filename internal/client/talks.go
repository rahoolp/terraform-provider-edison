package hashitalks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
)

var (
	ErrTalkNotFound = errors.New("talk not found")
)

type TalksService struct {
	basePath string
	client   *Client
}

func newTalksService(basePath string, client *Client) *TalksService {
	return &TalksService{
		basePath: basePath,
		client:   client,
	}
}

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

func (s TalksService) buildURL(p string) string {
	return path.Join(s.basePath, p)
}

func (s TalksService) Create(ctx context.Context, talk Talk) (Talk, error) {
	b, err := json.Marshal(talk)
	if err != nil {
		return Talk{}, fmt.Errorf("error serialising talk: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPost, s.buildURL("/"), buf)
	if err != nil {
		return Talk{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Talk{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Talk{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Talk{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return Talk{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return Talk{}, errors.New("talk already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return Talk{}, errors.New("name must be set")
	}
	if len(resp.Errors) > 0 {
		return Talk{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Talks) < 1 {
		return Talk{}, errors.New("no talk returned in response")
	}
	return resp.Talks[0], nil
}

func (s TalksService) Get(ctx context.Context, id string) (Talk, error) {
	if id == "" {
		return Talk{}, errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return Talk{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Talk{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Talk{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Talk{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return Talk{}, ErrTalkNotFound
	}
	if len(resp.Errors) > 0 {
		return Talk{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Talks) < 1 {
		return Talk{}, errors.New("no talk returned in response")
	}
	return resp.Talks[0], nil
}

func (s TalksService) Update(ctx context.Context, talk Talk) (Talk, error) {
	if talk.ID == "" {
		return Talk{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(talk)
	if err != nil {
		return Talk{}, fmt.Errorf("error serialising talk: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPut, s.buildURL("/"+talk.ID), buf)
	if err != nil {
		return Talk{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Talk{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Talk{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Talk{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return Talk{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return Talk{}, ErrTalkNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return Talk{}, errors.New("talk name must be set")
	}
	if len(resp.Errors) > 0 {
		return Talk{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Talks) < 1 {
		return Talk{}, errors.New("no talk returned in response")
	}
	return resp.Talks[0], nil
}

func (s TalksService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return err
	}

	if resp.Errors.Contains(serverError) {
		return errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return ErrTalkNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
