package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
				Optional: true,
			},
			"ip_port": {
				Type:     types.StringType,
				Optional: true,
			},
			"aet": {
				Type:     types.StringType,
				Optional: true,
			},
			"created_at": {
				Type:     types.StringType,
				Optional: true,
			},
			"updated_at": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

type eastoreData struct {
	ID               types.String `tfsdk:"id"`
	PartitionSpaceTB int64        `tfsdk:"partition_space_tb"`
	IPAddress        *string      `tfsdk:"ip_address"`
	IPPort           *string      `tfsdk:"ip_port"`
	AET              *string      `tfsdk:"aet"`
	CreatedAt        *string      `tfsdk:"created_at"`
	UpdatedAt        *string      `tfsdk:"updated_at"`
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

	eastore, err := e.client.EAStores.Create(ctx, edison.EAStore{
		PartitionSpaceTB: eastr.PartitionSpaceTB,
		IPAddress:        eastr.IPAddress,
		IPPort:           eastr.IPPort,
		AET:              eastr.AET,
		CreatedAt:        eastr.CreatedAt,
		UpdatedAt:        eastr.UpdatedAt,
	})
	if err != nil {
		// TODO: return error
	}
	eastr.ID = types.String{Value: eastore.ID}

	err = resp.State.Set(ctx, &eastr)
	if err != nil {
		// TODO: return error
	}
}

func (s eastoreResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	speaker, err := s.client.Speakers.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrSpeakerNotFound) {
		// TODO: return error
	} else if errors.Is(err, edison.ErrSpeakerNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	err = resp.State.Set(ctx, &speakerData{
		ID:       types.String{Value: speaker.ID},
		Name:     speaker.Name,
		Title:    speaker.Title,
		Employer: speaker.Employer,
		Pronouns: speaker.Pronouns,
		Photo:    speaker.Photo,
	})
	if err != nil {
		// TODO: return error
	}
}

func (e eastoreResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	var spkr speakerData
	err = req.Plan.Get(ctx, &spkr)
	if err != nil {
		// TODO: return error
	}

	_, err = e.client.Speakers.Update(ctx, edison.Speaker{
		ID:       id.(types.String).Value,
		Name:     spkr.Name,
		Title:    spkr.Title,
		Employer: spkr.Employer,
		Pronouns: spkr.Pronouns,
		Photo:    spkr.Photo,
	})
	if err != nil {
		// TODO: return error
	}
	spkr.ID = id.(types.String)

	err = resp.State.Set(ctx, &spkr)
	if err != nil {
		// TODO: return error
	}
}

func (e eastoreResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	err = e.client.Speakers.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrSpeakerNotFound) {
		// TODO: return error
	}
	resp.State.RemoveResource(ctx)
}
