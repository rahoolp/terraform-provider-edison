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
	ErrEAStoreNotFound = errors.New("eastore not found")
)

type EAStoresService struct {
	basePath string
	client   *Client
}

func newEAStoreService(basePath string, client *Client) *EAStoresService {
	return &EAStoresService{
		basePath: basePath,
		client:   client,
	}
}

type EAStore struct {
	ID               string `json:"id,omitempty"`
	PartitionSpaceTB int64  `json:"partition_space_tb"`
	IPAddress        string `json:"ip_address,omitempty"` //omitempty allows for null, aka null value
	IPPort           string `json:"ip_port,omitempty"`
	AET              string `json:"aet,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
	//DeletedAt        sql.NullString `db:"deleted_at" json:"-"`
}

func (s EAStoresService) buildURL(p string) string {
	return path.Join(s.basePath, p)
}

func (s EAStoresService) Create(ctx context.Context, eastore EAStore) (EAStore, error) {
	b, err := json.Marshal(eastore)
	if err != nil {
		return EAStore{}, fmt.Errorf("error serialising eastore: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPost, s.buildURL("/"), buf)
	if err != nil {
		return EAStore{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return EAStore{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return EAStore{}, err
	}

	if resp.Errors.Contains(serverError) {
		return EAStore{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return EAStore{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return EAStore{}, errors.New("EA Store already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/partition_space_tb",
	}) {
		return EAStore{}, errors.New("partition_space_tb must be set")
	}
	if len(resp.Errors) > 0 {
		return EAStore{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.EAStores) < 1 {
		return EAStore{}, errors.New("no EA Store returned in response")
	}
	return resp.EAStores[0], nil
}

func (s EAStoresService) Get(ctx context.Context, id string) (EAStore, error) {
	if id == "" {
		return EAStore{}, errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return EAStore{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return EAStore{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return EAStore{}, err
	}

	if resp.Errors.Contains(serverError) {
		return EAStore{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return EAStore{}, ErrEAStoreNotFound
	}
	if len(resp.Errors) > 0 {
		return EAStore{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.EAStores) < 1 {
		return EAStore{}, errors.New("no EA Store returned in response")
	}
	return resp.EAStores[0], nil
}

func (s EAStoresService) Update(ctx context.Context, eastore EAStore) (EAStore, error) {
	if eastore.ID == "" {
		return EAStore{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(eastore)
	if err != nil {
		return EAStore{}, fmt.Errorf("error serialising EA Store: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPut, s.buildURL("/"+eastore.ID), buf)
	if err != nil {
		return EAStore{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return EAStore{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return EAStore{}, err
	}

	if resp.Errors.Contains(serverError) {
		return EAStore{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return EAStore{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return EAStore{}, ErrEAStoreNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/partition_space_tb",
	}) {
		return EAStore{}, errors.New("EA Store partition_space_tb must be set")
	}
	if len(resp.Errors) > 0 {
		return EAStore{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.EAStores) < 1 {
		return EAStore{}, errors.New("no EA Store returned in response")
	}
	return resp.EAStores[0], nil
}

func (s EAStoresService) Delete(ctx context.Context, id string) error {
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
		return ErrEAStoreNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
