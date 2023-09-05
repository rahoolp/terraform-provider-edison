package edison

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
	ErrSpeakerNotFound = errors.New("speaker not found")
)

type SpeakersService struct {
	basePath string
	client   *Client
}

func newSpeakersService(basePath string, client *Client) *SpeakersService {
	return &SpeakersService{
		basePath: basePath,
		client:   client,
	}
}

type Speaker struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name"`
	Pronouns *string `json:"pronouns"`
	Employer *string `json:"employer"`
	Title    *string `json:"title"`
	Photo    *string `json:"photo"`
}

func (s SpeakersService) buildURL(p string) string {
	return path.Join(s.basePath, p)
}

func (s SpeakersService) Create(ctx context.Context, speaker Speaker) (Speaker, error) {
	b, err := json.Marshal(speaker)
	if err != nil {
		return Speaker{}, fmt.Errorf("error serialising speaker: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPost, s.buildURL("/"), buf)
	if err != nil {
		return Speaker{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Speaker{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Speaker{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Speaker{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return Speaker{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return Speaker{}, errors.New("speaker already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return Speaker{}, errors.New("name must be set")
	}
	if len(resp.Errors) > 0 {
		return Speaker{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Speakers) < 1 {
		return Speaker{}, errors.New("no speaker returned in response")
	}
	return resp.Speakers[0], nil
}

func (s SpeakersService) Get(ctx context.Context, id string) (Speaker, error) {
	if id == "" {
		return Speaker{}, errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return Speaker{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Speaker{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Speaker{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Speaker{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return Speaker{}, ErrSpeakerNotFound
	}
	if len(resp.Errors) > 0 {
		return Speaker{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Speakers) < 1 {
		return Speaker{}, errors.New("no speaker returned in response")
	}
	return resp.Speakers[0], nil
}

func (s SpeakersService) Update(ctx context.Context, speaker Speaker) (Speaker, error) {
	if speaker.ID == "" {
		return Speaker{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(speaker)
	if err != nil {
		return Speaker{}, fmt.Errorf("error serialising speaker: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPut, s.buildURL("/"+speaker.ID), buf)
	if err != nil {
		return Speaker{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Speaker{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Speaker{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Speaker{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return Speaker{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return Speaker{}, ErrSpeakerNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return Speaker{}, errors.New("speaker name must be set")
	}
	if len(resp.Errors) > 0 {
		return Speaker{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Speakers) < 1 {
		return Speaker{}, errors.New("no speaker returned in response")
	}
	return resp.Speakers[0], nil
}

func (s SpeakersService) Delete(ctx context.Context, id string) error {
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
		return ErrSpeakerNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
