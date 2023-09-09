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

type awResourceType struct {
}

func (e awResourceType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"concurrent_users": {
				Type:     types.NumberType,
				Required: true,
			},
			"ehs_cluster_id": {
				Type:     types.StringType,
				Computed: true,
			},
			"dicom_endpoint": {
				Type:     types.StringType,
				Required: true,
			},
			"dns_endpoint": {
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

type awData struct {
	ID              types.String `tfsdk:"id"`
	ConcurrentUsers int          `tfsdk:"concurrent_users"`
	DicomEndPoint   types.String `tfsdk:"dicom_endpoint"`
	DNSEndPoint     types.String `tfsdk:"dns_endpoint"`
	EHSClusterID    types.String `tfsdk:"ehs_cluster_id"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

func (s awResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
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
	return awResource{client: prov.client}, nil
}

type awResource struct {
	client *edison.Client
}

func (e awResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

	tflog.Info(ctx, "AW Create..")
	fmt.Println("AW Create")

	var aw awData
	err := req.Plan.Get(ctx, &aw)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error parsing plan",
			Detail:   "An unexpected error was encountered parsing the plan. This is always a bug in the provider.\n\nDetails: " + err.Error(),
		})
		return
	}

	var dnsEP string = "https://aw-04.ehs.edison.gehealthcare.com/"
	now := time.Now()
	var createdAt string = now.Format("2006-01-02 15:04:05")
	var updatedAt string = now.Format("2006-01-02 15:04:05")

	eaw, err := e.client.AWs.Create(ctx, edison.AW{
		ConcurrentUsers: aw.ConcurrentUsers,
		DicomEndPoint:   aw.DicomEndPoint.Value,
		EHSClusterID:    aw.EHSClusterID.Value,
		DNSEndPoint:     dnsEP,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "AW Create: "+err.Error())
	}

	aw.ID = types.String{Value: eaw.ID}
	aw.CreatedAt = types.String{Value: eaw.CreatedAt}
	aw.UpdatedAt = types.String{Value: eaw.UpdatedAt}
	aw.DNSEndPoint = types.String{Value: dnsEP}

	err = resp.State.Set(ctx, &aw)
	if err != nil {
		tflog.Info(ctx, "AW Create: "+err.Error())
	}
}

func (e awResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

	tflog.Info(ctx, "AW Read..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "AW Read: "+err.Error())
	}

	aw, err := e.client.AWs.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrEHSClusterNotFound) {
		tflog.Info(ctx, "AW Read: "+err.Error())
	} else if errors.Is(err, edison.ErrEHSClusterNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	err = resp.State.Set(ctx, &awData{
		ID:              types.String{Value: aw.ID},
		ConcurrentUsers: aw.ConcurrentUsers,
		DicomEndPoint:   types.String{Value: aw.DicomEndPoint},
		DNSEndPoint:     types.String{Value: aw.DNSEndPoint},
		CreatedAt:       types.String{Value: aw.CreatedAt},
		UpdatedAt:       types.String{Value: aw.UpdatedAt},
	})

	if err != nil {
		tflog.Info(ctx, "AW: "+err.Error())
	}
}

func (e awResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

	tflog.Info(ctx, "AW Update..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "AW Update: "+err.Error())
	}

	var aw awData
	err = req.Plan.Get(ctx, &aw)
	if err != nil {
		tflog.Info(ctx, "AW Update: "+err.Error())
	}

	now := time.Now()
	var updatedAt string = now.Format("2006-01-02 15:04:05")

	_, err = e.client.AWs.Update(ctx, edison.AW{
		ID:              id.(types.String).Value,
		ConcurrentUsers: aw.ConcurrentUsers,
		DicomEndPoint:   aw.DicomEndPoint.Value,
		DNSEndPoint:     aw.DNSEndPoint.Value,
		CreatedAt:       aw.CreatedAt.Value,
		UpdatedAt:       updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "AW Update: "+err.Error())
	}
	aw.ID = id.(types.String)

	err = resp.State.Set(ctx, &aw)
	if err != nil {
		tflog.Info(ctx, "AW Update: "+err.Error())
	}
}

func (e awResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

	tflog.Info(ctx, "AW Delete..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "AW Delete: "+err.Error())
	}
	err = e.client.AWs.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrEHSClusterNotFound) {
		tflog.Info(ctx, "AW Delete: "+err.Error())
	}
	resp.State.RemoveResource(ctx)
}
