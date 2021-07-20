package provider

import (
	"context"
	"errors"
	"math/big"

	hashitalks "github.com/hashicorp/terraform-provider-hashitalks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type talkResourceType struct {
}

func (s talkResourceType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"title": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Required: true,
			},
			"duration_minutes": {
				Type:     types.NumberType,
				Required: true,
			},
			"prerecorded": {
				Type:     types.BoolType,
				Required: true,
			},
			"speaker_ids": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
			},
			"recordings": {
				Computed: true,
				Type: types.MapType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"codec": types.StringType,
							"resolution": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"width":  types.NumberType,
									"height": types.NumberType,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

type talkData struct {
	ID              types.String `tfsdk:"id"`
	Title           string       `tfsdk:"title"`
	Description     string       `tfsdk:"description"`
	DurationMinutes int64        `tfsdk:"duration_minutes"`
	Prerecorded     bool         `tfsdk:"prerecorded"`
	SpeakerIDs      []string     `tfsdk:"speaker_ids"`
	Recordings      types.Map    `tfsdk:"recordings"`
}

type talkRecordingData struct {
	Codec      string                      `tfsdk:"codec"`
	Resolution talkRecordingResolutionData `tfsdk:"resolution"`
}

type talkRecordingResolutionData struct {
	Width  int64 `tfsdk:"width"`
	Height int64 `tfsdk:"height"`
}

func (s talkResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
	prov, ok := p.(*provider)
	if !ok {
		// TODO: return an error
	}
	return talkResource{client: prov.client}, nil
}

type talkResource struct {
	client *hashitalks.Client
}

func (s talkResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan talkData
	err := req.Plan.Get(ctx, &plan)
	if err != nil {
		// TODO: return error
	}

	talk, err := s.client.Talks.Create(ctx, hashitalks.Talk{
		Title:           plan.Title,
		Description:     plan.Description,
		DurationMinutes: plan.DurationMinutes,
		Prerecorded:     plan.Prerecorded,
		SpeakerIDs:      plan.SpeakerIDs,
	})
	if err != nil {
		// TODO: return error
	}
	plan.ID = types.String{Value: talk.ID}

	err = resp.State.Set(ctx, &plan)
	if err != nil {
		// TODO: return error
	}
}

func (s talkResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	talk, err := s.client.Talks.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, hashitalks.ErrTalkNotFound) {
		// TODO: return error
	} else if errors.Is(err, hashitalks.ErrTalkNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	recordings := types.Map{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"codec": types.StringType,
				"resolution": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"width":  types.NumberType,
						"height": types.NumberType,
					},
				},
			},
		},
		Elems: map[string]attr.Value{},
	}
	for user, recording := range talk.Recordings {
		recordings.Elems[user] = types.Object{
			AttrTypes: map[string]attr.Type{
				"codec": types.StringType,
				"resolution": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"width":  types.NumberType,
						"height": types.NumberType,
					},
				},
			},
			Attrs: map[string]attr.Value{
				"codec": types.String{Value: recording.Codec},
				"resolution": types.Object{
					AttrTypes: map[string]attr.Type{
						"width":  types.NumberType,
						"height": types.NumberType,
					},
					Attrs: map[string]attr.Value{
						"width":  types.Number{Value: big.NewFloat(float64(recording.Resolution.Width))},
						"height": types.Number{Value: big.NewFloat(float64(recording.Resolution.Height))},
					},
				},
			},
		}
	}
	err = resp.State.Set(ctx, &talkData{
		ID:              types.String{Value: talk.ID},
		Title:           talk.Title,
		Description:     talk.Description,
		DurationMinutes: talk.DurationMinutes,
		Prerecorded:     talk.Prerecorded,
		SpeakerIDs:      talk.SpeakerIDs,
		Recordings:      recordings,
	})
	if err != nil {
		// TODO: return error
	}
}

func (s talkResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	var plan talkData
	err = req.Plan.Get(ctx, &plan)
	if err != nil {
		// TODO: return error
	}

	_, err = s.client.Talks.Update(ctx, hashitalks.Talk{
		ID:              id.(types.String).Value,
		Title:           plan.Title,
		Description:     plan.Description,
		DurationMinutes: plan.DurationMinutes,
		Prerecorded:     plan.Prerecorded,
		SpeakerIDs:      plan.SpeakerIDs,
	})
	if err != nil {
		// TODO: return error
	}

	err = resp.State.Set(ctx, &plan)
	if err != nil {
		// TODO: return error
	}
}

func (s talkResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	err = s.client.Talks.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, hashitalks.ErrTalkNotFound) {
		// TODO: return error
	}
	resp.State.RemoveResource(ctx)
}
