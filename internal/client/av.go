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
	ErrAVNotFound = errors.New("av not found")
)

type AVsService struct {
	basePath string
	client   *Client
}

func newAVService(basePath string, client *Client) *AVsService {
	return &AVsService{
		basePath: basePath,
		client:   client,
	}
}

type AV struct {
	ID           string `json:"id,omitempty"`
	AccountID    string `json:"account_id"`
	TenantID     string `json:"tenant_id"`
	TenantFolder string `json:"tenant_folder,omitempty"`
	TenantQueue  string `json:"tenant_queue,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

func (s AVsService) buildURL(p string) string {
	return path.Join(s.basePath, p)
}

func (s AVsService) Create(ctx context.Context, av AV) (AV, error) {
	b, err := json.Marshal(av)
	if err != nil {
		return AV{}, fmt.Errorf("error serialising av: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPost, s.buildURL("/"), buf)
	if err != nil {
		return AV{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return AV{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AV{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AV{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return AV{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return AV{}, errors.New("AV already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/tenant_id",
	}) {
		return AV{}, errors.New("tenant id must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/account_id",
	}) {
		return AV{}, errors.New("account id must be set")
	}

	if len(resp.Errors) > 0 {
		return AV{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AVs) < 1 {
		return AV{}, errors.New("no AV returned in response")
	}
	return resp.AVs[0], nil
}

func (s AVsService) Get(ctx context.Context, id string) (AV, error) {
	if id == "" {
		return AV{}, errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return AV{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return AV{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AV{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AV{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return AV{}, ErrAVNotFound
	}
	if len(resp.Errors) > 0 {
		return AV{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AVs) < 1 {
		return AV{}, errors.New("no AV returned in response")
	}
	return resp.AVs[0], nil
}

func (s AVsService) Update(ctx context.Context, av AV) (AV, error) {
	if av.ID == "" {
		return AV{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(av)
	if err != nil {
		return AV{}, fmt.Errorf("error serialising AV: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPut, s.buildURL("/"+av.ID), buf)
	if err != nil {
		return AV{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return AV{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AV{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AV{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return AV{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return AV{}, ErrAVNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/tenant_id",
	}) {
		return AV{}, errors.New("Tenant ID must be set")
	}
	if len(resp.Errors) > 0 {
		return AV{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AVs) < 1 {
		return AV{}, errors.New("no AV returned in response")
	}
	return resp.AVs[0], nil
}

func (s AVsService) Delete(ctx context.Context, id string) error {
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
		return ErrAVNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
