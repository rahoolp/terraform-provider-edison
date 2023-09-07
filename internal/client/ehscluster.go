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
	ErrEHSClusterNotFound = errors.New("ehscluster not found")
)

type EHSClustersService struct {
	basePath string
	client   *Client
}

func newEHSClusterService(basePath string, client *Client) *EHSClustersService {
	return &EHSClustersService{
		basePath: basePath,
		client:   client,
	}
}

type EHSCluster struct {
	ID                string `json:"id,omitempty"`
	Region            string `json:"region"`
	Profile           string `json:"profile"`
	Release           string `json:"release"`
	Tag               string `json:"tag"`
	ClusterName       string `json:"cluster_name"`
	DicomEndPoint     string `json:"dicom_endpoint"`
	APIServerEndPoint string `json:"api_server_endpoint,omitempty"`
	VPC               string `json:"vpc,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

func (s EHSClustersService) buildURL(p string) string {
	return path.Join(s.basePath, p)
}

func (s EHSClustersService) Create(ctx context.Context, ehscluster EHSCluster) (EHSCluster, error) {
	b, err := json.Marshal(ehscluster)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error serialising ehscluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPost, s.buildURL("/"), buf)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return EHSCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return EHSCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return EHSCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrConflict,
		Field: "/id",
	}) {
		return EHSCluster{}, errors.New("EHS Cluster already exists")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/region",
	}) {
		return EHSCluster{}, errors.New("region must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/dicom_endpoint",
	}) {
		return EHSCluster{}, errors.New("dicom_endpoint must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/profile",
	}) {
		return EHSCluster{}, errors.New("profile must be set")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/release",
	}) {
		return EHSCluster{}, errors.New("release must be set")
	}

	if len(resp.Errors) > 0 {
		return EHSCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.EHSClusters) < 1 {
		return EHSCluster{}, errors.New("no EHS Cluster returned in response")
	}
	return resp.EHSClusters[0], nil
}

func (s EHSClustersService) Get(ctx context.Context, id string) (EHSCluster, error) {
	if id == "" {
		return EHSCluster{}, errors.New("id must be specified")
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, s.buildURL("/"+id), nil)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return EHSCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return EHSCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return EHSCluster{}, ErrEHSClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return EHSCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.EHSClusters) < 1 {
		return EHSCluster{}, errors.New("no EHS Cluster returned in response")
	}
	return resp.EHSClusters[0], nil
}

func (s EHSClustersService) Update(ctx context.Context, ehscluster EHSCluster) (EHSCluster, error) {
	if ehscluster.ID == "" {
		return EHSCluster{}, errors.New("id must be specified")
	}
	b, err := json.Marshal(ehscluster)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error serialising EHS Cluster: %w", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := s.client.NewRequest(ctx, http.MethodPut, s.buildURL("/"+ehscluster.ID), buf)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error constructing request: %w", err)
	}
	res, err := s.client.Do(req)
	if err != nil {
		return EHSCluster{}, fmt.Errorf("error making request: %w", err)
	}
	resp, err := responseFromBody(res)
	if err != nil {
		return EHSCluster{}, err
	}

	if resp.Errors.Contains(serverError) {
		return EHSCluster{}, errors.New("server error")
	}
	if resp.Errors.Contains(invalidFormatError) {
		return EHSCluster{}, errors.New("invalid format error returned")
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrNotFound,
		Param: "id",
	}) {
		return EHSCluster{}, ErrEHSClusterNotFound
	}
	if resp.Errors.Contains(RequestError{
		Slug:  requestErrMissing,
		Field: "/partition_space_tb",
	}) {
		return EHSCluster{}, errors.New("EHS Cluster partition_space_tb must be set")
	}
	if len(resp.Errors) > 0 {
		return EHSCluster{}, fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	if len(resp.EHSClusters) < 1 {
		return EHSCluster{}, errors.New("no EHS Cluster returned in response")
	}
	return resp.EHSClusters[0], nil
}

func (s EHSClustersService) Delete(ctx context.Context, id string) error {
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
		return ErrEHSClusterNotFound
	}
	if len(resp.Errors) > 0 {
		return fmt.Errorf("unexpected error in response: %+v", resp.Errors)
	}
	return nil
}
