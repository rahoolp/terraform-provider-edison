package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	edison "github.com/rahoolp/terraform-provider-edison/internal/client"
)

type ehsclusterResourceType struct {
}

func (e ehsclusterResourceType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"region": {
				Type:     types.StringType,
				Required: true,
			},
			"profile": {
				Type:     types.StringType,
				Required: true,
			},
			"release": {
				Type:     types.StringType,
				Required: true,
			},
			"tag": {
				Type:     types.StringType,
				Required: true,
			},
			"dicom_endpoint": {
				Type:     types.StringType,
				Required: true,
			},
			"api_server_endpoint": {
				Type:     types.StringType,
				Computed: true,
			},
			"vpc": {
				Type:     types.StringType,
				Computed: true,
			},
			"created_at": {
				Type:     types.StringType,
				Computed: true,
			},
			"updated_at": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

type ehsclusterData struct {
	ID                types.String `tfsdk:"id"`
	Region            types.String `tfsdk:"region"`
	Profile           types.String `tfsdk:"profile"`
	Release           types.String `tfsdk:"release"`
	Tag               types.String `tfsdk:"tag"`
	DicomEndPoint     types.String `tfsdk:"dicom_endpoint"`
	APIServerEndPoint types.String `tfsdk:"api_server_endpoint"`
	VPC               types.String `tfsdk:"vpc"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

func (s ehsclusterResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
	prov, ok := p.(*provider)
	if !ok {
		return nil, []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error converting provider",
				Detail:   fmt.Sprintf("An unexpected error was encountered converting the provider. This is always a bug in the provider.\n\nType: %T", p),
			},
		}
	}
	return ehsclusterResource{client: prov.client}, nil
}

type ehsclusterResource struct {
	client *edison.Client
}

func (e ehsclusterResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

	tflog.Info(ctx, "EHS Cluster Create..")
	fmt.Println("EHS Cluster Create")

	var ehscluster ehsclusterData
	err := req.Plan.Get(ctx, &ehscluster)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error parsing plan",
			Detail:   "An unexpected error was encountered parsing the plan. This is always a bug in the provider.\n\nDetails: " + err.Error(),
		})
		return
	}

	var vpc string = "vpc-0c6aa52f85161d3cc"
	var apiSrvEP string = "https://5a4028bb2291be0fa29ab9717a8b9e92.gr7.us-east-1.eks.amazonaws.com/"
	now := time.Now()
	var createdAt string = now.Format("2006-01-02 15:04:05")
	var updatedAt string = now.Format("2006-01-02 15:04:05")

	ecluster, err := e.client.EHSClusters.Create(ctx, edison.EHSCluster{
		Region:            ehscluster.Region.Value,
		Profile:           ehscluster.Profile.Value,
		Release:           ehscluster.Release.Value,
		VPC:               vpc,
		APIServerEndPoint: apiSrvEP,
		DicomEndPoint:     ehscluster.DicomEndPoint.Value,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Create: "+err.Error())
	}

	ehscluster.ID = types.String{Value: ecluster.ID}
	ehscluster.CreatedAt = types.String{Value: ecluster.CreatedAt}
	ehscluster.UpdatedAt = types.String{Value: ecluster.UpdatedAt}
	ehscluster.APIServerEndPoint = types.String{Value: apiSrvEP}
	ehscluster.VPC = types.String{Value: vpc}

	err = resp.State.Set(ctx, &ehscluster)
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Create: "+err.Error())
	}
}

func (e ehsclusterResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

	tflog.Info(ctx, "EHS Cluster Read..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Read: "+err.Error())
	}

	ehscluster, err := e.client.EHSClusters.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrSpeakerNotFound) {
		tflog.Info(ctx, "EHS Cluster Read: "+err.Error())
	} else if errors.Is(err, edison.ErrSpeakerNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	err = resp.State.Set(ctx, &ehsclusterData{
		ID:                types.String{Value: ehscluster.ID},
		Profile:           types.String{Value: ehscluster.Profile},
		Region:            types.String{Value: ehscluster.Region},
		Release:           types.String{Value: ehscluster.Release},
		VPC:               types.String{Value: ehscluster.VPC},
		Tag:               types.String{Value: ehscluster.Tag},
		APIServerEndPoint: types.String{Value: ehscluster.APIServerEndPoint},
		DicomEndPoint:     types.String{Value: ehscluster.DicomEndPoint},
		CreatedAt:         types.String{Value: ehscluster.CreatedAt},
		UpdatedAt:         types.String{Value: ehscluster.UpdatedAt},
	})

	if err != nil {
		tflog.Info(ctx, "EHS Cluster: "+err.Error())
	}
}

func (e ehsclusterResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

	tflog.Info(ctx, "EHS Cluster Update..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Update: "+err.Error())
	}

	var ehscluster ehsclusterData
	err = req.Plan.Get(ctx, &ehscluster)
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Update: "+err.Error())
	}

	now := time.Now()
	var updatedAt string = now.Format("2006-01-02 15:04:05")

	_, err = e.client.EHSClusters.Update(ctx, edison.EHSCluster{
		ID:                id.(types.String).Value,
		Profile:           ehscluster.Profile.Value,
		Region:            ehscluster.Region.Value,
		Release:           ehscluster.Release.Value,
		VPC:               ehscluster.VPC.Value,
		Tag:               ehscluster.Tag.Value,
		APIServerEndPoint: ehscluster.APIServerEndPoint.Value,
		DicomEndPoint:     ehscluster.DicomEndPoint.Value,
		CreatedAt:         ehscluster.CreatedAt.Value,
		UpdatedAt:         updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Update: "+err.Error())
	}
	ehscluster.ID = id.(types.String)

	err = resp.State.Set(ctx, &ehscluster)
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Update: "+err.Error())
	}
}

func (e ehsclusterResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

	tflog.Info(ctx, "EHS Cluster Delete..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "EHS Cluster Delete: "+err.Error())
	}
	err = e.client.EHSClusters.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrEHSClusterNotFound) {
		tflog.Info(ctx, "EHS Cluster Delete: "+err.Error())
	}
	resp.State.RemoveResource(ctx)
}
