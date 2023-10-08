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

type eastoreResourceType struct {
}

func (e eastoreResourceType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"partition_space_tb": {
				Type:     types.NumberType,
				Required: true,
			},
			"ip_address": {
				Type:     types.StringType,
				Computed: true,
			},
			"ip_port": {
				Type:     types.StringType,
				Computed: true,
			},
			"aet": {
				Type:     types.StringType,
				Computed: true,
			},
			"account_id": {
				Type:     types.StringType,
				Computed: true,
			},
			"service_ep": {
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

type eastoreData struct {
	ID               types.String `tfsdk:"id"`
	PartitionSpaceTB int64        `tfsdk:"partition_space_tb"`
	IPAddress        types.String `tfsdk:"ip_address"`
	IPPort           types.String `tfsdk:"ip_port"`
	AET              types.String `tfsdk:"aet"`
	AccountID        types.String `tfsdk:"account_id"`
	ServiceEP        types.String `tfsdk:"service_ep"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

func (s eastoreResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
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
	return eastoreResource{client: prov.client}, nil
}

type eastoreResource struct {
	client *edison.Client
}

func (e eastoreResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

	tflog.Info(ctx, "EA Store Create..")
	fmt.Println("EA Store Create")

	var eastr eastoreData
	err := req.Plan.Get(ctx, &eastr)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error parsing plan",
			Detail:   "An unexpected error was encountered parsing the plan. This is always a bug in the provider.\n\nDetails: " + err.Error(),
		})
		return
	}

	var id string = "1"
	var ipAddress string = "192.168.1.1"
	var ipPort string = "4242"
	var aet string = "AET1"
	var account_id string = "3091120"
	var service_ep string = "eadicom_ep"
	now := time.Now()
	var createdAt string = now.Format("2006-01-02 15:04:05")
	var updatedAt string = now.Format("2006-01-02 15:04:05")

	eastore, err := e.client.EAStores.Create(ctx, edison.EAStore{
		ID:               id,
		PartitionSpaceTB: eastr.PartitionSpaceTB,
		IPAddress:        ipAddress,
		IPPort:           ipPort,
		AET:              aet,
		AccountID:        account_id,
		ServiceEP:        service_ep,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "EA Store Create: "+err.Error())
	}

	eastr.ID = types.String{Value: eastore.ID}
	eastr.PartitionSpaceTB = eastore.PartitionSpaceTB
	eastr.IPAddress = types.String{Value: ipAddress}
	eastr.IPPort = types.String{Value: ipPort}
	eastr.AET = types.String{Value: aet}
	eastr.AccountID = types.String{Value: account_id}
	eastr.ServiceEP = types.String{Value: service_ep}
	eastr.CreatedAt = types.String{Value: createdAt}
	eastr.UpdatedAt = types.String{Value: updatedAt}

	err = resp.State.Set(ctx, &eastr)
	if err != nil {
		tflog.Info(ctx, "EA Store Create: "+err.Error())
	}
}

func (e eastoreResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

	tflog.Info(ctx, "EA Store Read..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "EA Store Read: "+err.Error())
	}

	eastr, err := e.client.EAStores.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrEAStoreNotFound) {
		tflog.Info(ctx, "EA Store Read: "+err.Error())
	} else if errors.Is(err, edison.ErrEAStoreNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	err = resp.State.Set(ctx, &eastoreData{
		ID:               types.String{Value: eastr.ID},
		PartitionSpaceTB: eastr.PartitionSpaceTB,
		IPAddress:        types.String{Value: eastr.IPAddress},
		IPPort:           types.String{Value: eastr.IPPort},
		AET:              types.String{Value: eastr.AET},
		AccountID:        types.String{Value: eastr.AccountID},
		ServiceEP:        types.String{Value: eastr.ServiceEP},
		CreatedAt:        types.String{Value: eastr.CreatedAt},
		UpdatedAt:        types.String{Value: eastr.UpdatedAt},
	})

	if err != nil {
		tflog.Info(ctx, "EA Store: "+err.Error())
	}
}

func (e eastoreResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

	tflog.Info(ctx, "EA Store Update..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "EA Store Update: "+err.Error())
	}

	var eastr eastoreData
	err = req.Plan.Get(ctx, &eastr)
	if err != nil {
		tflog.Info(ctx, "EA Store Update: "+err.Error())
	}

	now := time.Now()
	var updatedAt string = now.Format("2006-01-02 15:04:05")

	_, err = e.client.EAStores.Update(ctx, edison.EAStore{
		ID:               id.(types.String).Value,
		PartitionSpaceTB: eastr.PartitionSpaceTB,
		IPAddress:        eastr.IPAddress.Value,
		IPPort:           eastr.IPPort.Value,
		AET:              eastr.AET.Value,
		AccountID:        eastr.AccountID.Value,
		ServiceEP:        eastr.IPAddress.Value,
		CreatedAt:        eastr.CreatedAt.Value,
		UpdatedAt:        updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "EA Store Update: "+err.Error())
	}
	eastr.ID = id.(types.String)

	err = resp.State.Set(ctx, &eastr)
	if err != nil {
		tflog.Info(ctx, "EA Store Update: "+err.Error())
	}
}

func (e eastoreResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

	tflog.Info(ctx, "EA Store Delete..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "EA Store Delete: "+err.Error())
	}
	err = e.client.EAStores.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrEAStoreNotFound) {
		tflog.Info(ctx, "EA Store Delete: "+err.Error())
	}
	resp.State.RemoveResource(ctx)
}
