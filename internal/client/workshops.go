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
	ErrWorkshopNotFound = errors.New("workshop not found")
)

type WorkshopsService struct {
	basePath string
	client   *Client
}

func newWorkshopsService(basePath string, client *Client) *WorkshopsService {
	return &WorkshopsService{
		basePath: basePath,
		client:   client,
	}
}

type Workshop struct {
	ID              string                       `json:"id,omitempty"`
	Title           string                       `json:"title"`
	Description     string                       `json:"description"`
	DurationMinutes int64                        `json:"durationMinutes"`
	Presenters      map[string]WorkshopPresenter `json:"presenters"`
	MeetingInfo     WorkshopMeetingInfo          `json:"meetingInfo"`
}

type WorkshopPresenter struct {
	Title    *string `json:"title"`
	Employer *string `json:"employer"`
	Pronouns *string `json:"pronouns"`
}

type WorkshopMeetingInfo struct {
	URL      string  `json:"url"`
	Password *string `json:"password"`
}

func (s WorkshopsService) buildURL(p string) string {
	return path.Join(s.basePath, p)
}

func (s WorkshopsService) Create(ctx context.Context, workshop Workshop) (Workshop, error) {
	b, err := json.Marshal(workshop)
	if err != nil {
		return Workshop{}, fmt.Errorf("error serialising workshop: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPost, s.buildURL("/"), buf)
	if err != nil {
		return Workshop{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Workshop{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Workshop{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Workshop{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return Workshop{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return Workshop{}, errors.New("workshop already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return Workshop{}, errors.New("name must be set")
	}
	if len(resp.Errors) > 0 {
		return Workshop{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Workshops) < 1 {
		return Workshop{}, errors.New("no workshop returned in response")
	}
	return resp.Workshops[0], nil
}

func (s WorkshopsService) Get(ctx context.Context, id string) (Workshop, error) {
	if id == "" {
		return Workshop{}, errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return Workshop{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Workshop{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Workshop{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Workshop{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return Workshop{}, ErrWorkshopNotFound
	}
	if len(resp.Errors) > 0 {
		return Workshop{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Workshops) < 1 {
		return Workshop{}, errors.New("no workshop returned in response")
	}
	return resp.Workshops[0], nil
}

func (s WorkshopsService) Update(ctx context.Context, workshop Workshop) (Workshop, error) {
	if workshop.ID == "" {
		return Workshop{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(workshop)
	if err != nil {
		return Workshop{}, fmt.Errorf("error serialising workshop: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPut, s.buildURL("/"+workshop.ID), buf)
	if err != nil {
		return Workshop{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return Workshop{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return Workshop{}, err
	}

	if resp.Errors.Contains(serverError) {
		return Workshop{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return Workshop{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return Workshop{}, ErrWorkshopNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/name",
	}) {
		return Workshop{}, errors.New("workshop name must be set")
	}
	if len(resp.Errors) > 0 {
		return Workshop{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.Workshops) < 1 {
		return Workshop{}, errors.New("no workshop returned in response")
	}
	return resp.Workshops[0], nil
}

func (s WorkshopsService) Delete(ctx context.Context, id string) error {
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
		return ErrWorkshopNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
