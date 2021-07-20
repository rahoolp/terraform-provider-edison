package provider

import (
	"context"
	"errors"

	hashitalks "github.com/hashicorp/terraform-provider-hashitalks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type speakerResourceType struct {
}

func (s speakerResourceType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"title": {
				Type:     types.StringType,
				Optional: true,
			},
			"employer": {
				Type:     types.StringType,
				Optional: true,
			},
			"pronouns": {
				Type:     types.StringType,
				Optional: true,
			},
			"photo": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

type speakerData struct {
	ID       types.String `tfsdk:"id"`
	Name     string       `tfsdk:"name"`
	Title    *string      `tfsdk:"title"`
	Employer *string      `tfsdk:"employer"`
	Pronouns *string      `tfsdk:"pronouns"`
	Photo    *string      `tfsdk:"photo"`
}

func (s speakerResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
	prov, ok := p.(*provider)
	if !ok {
		// TODO: return an error
	}
	return speakerResource{client: prov.client}, nil
}

type speakerResource struct {
	client *hashitalks.Client
}

func (s speakerResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var spkr speakerData
	err := req.Plan.Get(ctx, &spkr)
	if err != nil {
		// TODO: return error
	}

	speaker, err := s.client.Speakers.Create(ctx, hashitalks.Speaker{
		Name:     spkr.Name,
		Title:    spkr.Title,
		Employer: spkr.Employer,
		Pronouns: spkr.Pronouns,
		Photo:    spkr.Photo,
	})
	if err != nil {
		// TODO: return error
	}
	spkr.ID = types.String{Value: speaker.ID}

	err = resp.State.Set(ctx, &spkr)
	if err != nil {
		// TODO: return error
	}
}

func (s speakerResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	speaker, err := s.client.Speakers.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, hashitalks.ErrSpeakerNotFound) {
		// TODO: return error
	} else if errors.Is(err, hashitalks.ErrSpeakerNotFound) {
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

func (s speakerResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var spkr speakerData
	err := req.Plan.Get(ctx, &spkr)
	if err != nil {
		// TODO: return error
	}

	_, err = s.client.Speakers.Update(ctx, hashitalks.Speaker{
		ID:       spkr.ID.Value, // TODO: fixme
		Name:     spkr.Name,
		Title:    spkr.Title,
		Employer: spkr.Employer,
		Pronouns: spkr.Pronouns,
		Photo:    spkr.Photo,
	})
	if err != nil {
		// TODO: return error
	}

	err = resp.State.Set(ctx, &spkr)
	if err != nil {
		// TODO: return error
	}
}

func (s speakerResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	err = s.client.Speakers.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, hashitalks.ErrSpeakerNotFound) {
		// TODO: return error
	}
	resp.State.RemoveResource(ctx)
}
