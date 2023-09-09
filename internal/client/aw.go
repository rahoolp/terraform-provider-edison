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
	ErrAWNotFound = errors.New("aw not found")
)

type AWsService struct {
	basePath string
	client   *Client
}

func newAWService(basePath string, client *Client) *AWsService {
	return &AWsService{
		basePath: basePath,
		client:   client,
	}
}

type AW struct {
	ID              string `json:"id,omitempty"`
	ConcurrentUsers int    `json:"concurrent_users"`
	EHSClusterID    string `json:"ehs_cluster_id"`
	DicomEndPoint   string `json:"dicom_endpoint"`
	DNSEndPoint     string `json:"dns_endpoint,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

func (s AWsService) buildURL(p string) string {
	return path.Join(s.basePath, p)
}

func (s AWsService) Create(ctx context.Context, aw AW) (AW, error) {
	b, err := json.Marshal(aw)
	if err != nil {
		return AW{}, fmt.Errorf("error serialising aw: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPost, s.buildURL("/"), buf)
	if err != nil {
		return AW{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return AW{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AW{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AW{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return AW{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return AW{}, errors.New("AW already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/concurrent_users",
	}) {
		return AW{}, errors.New("concurrent users count must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/ehs_cluster_id",
	}) {
		return AW{}, errors.New("ehs cluster id must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/DicomEndPoint",
	}) {
		return AW{}, errors.New("Dicom End Point must be set")
	}

	if len(resp.Errors) > 0 {
		return AW{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AWs) < 1 {
		return AW{}, errors.New("no AW returned in response")
	}
	return resp.AWs[0], nil
}

func (s AWsService) Get(ctx context.Context, id string) (AW, error) {
	if id == "" {
		return AW{}, errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return AW{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return AW{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AW{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AW{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return AW{}, ErrAWNotFound
	}
	if len(resp.Errors) > 0 {
		return AW{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AWs) < 1 {
		return AW{}, errors.New("no AW returned in response")
	}
	return resp.AWs[0], nil
}

func (s AWsService) Update(ctx context.Context, aw AW) (AW, error) {
	if aw.ID == "" {
		return AW{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(aw)
	if err != nil {
		return AW{}, fmt.Errorf("error serialising AW: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPut, s.buildURL("/"+aw.ID), buf)
	if err != nil {
		return AW{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return AW{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return AW{}, err
	}

	if resp.Errors.Contains(serverError) {
		return AW{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return AW{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return AW{}, ErrAWNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/partition_space_tb",
	}) {
		return AW{}, errors.New("AW partition_space_tb must be set")
	}
	if len(resp.Errors) > 0 {
		return AW{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.AWs) < 1 {
		return AW{}, errors.New("no AW returned in response")
	}
	return resp.AWs[0], nil
}

func (s AWsService) Delete(ctx context.Context, id string) error {
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
		return ErrAWNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
