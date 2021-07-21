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

type workshopResourceType struct {
}

func (w workshopResourceType) GetSchema(ctx context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
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
			"presenters": {
				Required: true,
				Attributes: schema.MapNestedAttributes(map[string]schema.Attribute{
					"title": {
						Optional: true,
						Type:     types.StringType,
					},
					"employer": {
						Optional: true,
						Type:     types.StringType,
					},
					"pronouns": {
						Optional: true,
						Type:     types.StringType,
					},
				}, schema.MapNestedAttributesOptions{}),
			},
			"meeting_info": {
				Required: true,
				Attributes: schema.SingleNestedAttributes(map[string]schema.Attribute{
					"url": {
						Type:     types.StringType,
						Required: true,
					},
					"password": {
						Type:      types.StringType,
						Optional:  true,
						Sensitive: true,
					},
				}),
			},
		},
	}, nil
}

type workshopData struct {
	ID              types.String                     `tfsdk:"id"`
	Title           string                           `tfsdk:"title"`
	Description     string                           `tfsdk:"description"`
	DurationMinutes int64                            `tfsdk:"duration_minutes"`
	Presenters      map[string]workshopPresenterData `tfsdk:"presenters"`
	MeetingInfo     workshopMeetingInfoData          `tfsdk:"meeting_info"`
}

type workshopPresenterData struct {
	Title    *string `tfsdk:"title"`
	Employer *string `tfsdk:"employer"`
	Pronouns *string `tfsdk:"pronouns"`
}

type workshopMeetingInfoData struct {
	URL      string  `tfsdk:"url"`
	Password *string `tfsdk:"password"`
}

func (w workshopResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
	prov, ok := p.(*provider)
	if !ok {
		// TODO: return an error
	}
	return workshopResource{client: prov.client}, nil
}

type workshopResource struct {
	client *hashitalks.Client
}

func (w workshopResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan workshopData
	err := req.Plan.Get(ctx, &plan)
	if err != nil {
		// TODO: return error
	}

	presenters := map[string]hashitalks.WorkshopPresenter{}
	for name, info := range plan.Presenters {
		presenters[name] = hashitalks.WorkshopPresenter{
			Title:    info.Title,
			Employer: info.Employer,
			Pronouns: info.Pronouns,
		}
	}

	workshop, err := w.client.Workshops.Create(ctx, hashitalks.Workshop{
		Title:           plan.Title,
		Description:     plan.Description,
		DurationMinutes: plan.DurationMinutes,
		Presenters:      presenters,
		MeetingInfo: hashitalks.WorkshopMeetingInfo{
			URL:      plan.MeetingInfo.URL,
			Password: plan.MeetingInfo.Password,
		},
	})
	if err != nil {
		// TODO: return error
	}
	plan.ID = types.String{Value: workshop.ID}

	err = resp.State.Set(ctx, &plan)
	if err != nil {
		// TODO: return error
	}
}

func (w workshopResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	workshop, err := w.client.Workshops.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, hashitalks.ErrWorkshopNotFound) {
		// TODO: return error
	} else if errors.Is(err, hashitalks.ErrWorkshopNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	presenters := map[string]workshopPresenterData{}
	for name, info := range workshop.Presenters {
		presenters[name] = workshopPresenterData{
			Title:    info.Title,
			Employer: info.Employer,
			Pronouns: info.Pronouns,
		}
	}

	err = resp.State.Set(ctx, &workshopData{
		ID:              types.String{Value: workshop.ID},
		Title:           workshop.Title,
		Description:     workshop.Description,
		DurationMinutes: workshop.DurationMinutes,
		Presenters:      presenters,
		MeetingInfo: workshopMeetingInfoData{
			URL:      workshop.MeetingInfo.URL,
			Password: workshop.MeetingInfo.Password,
		},
	})
	if err != nil {
		// TODO: return error
	}
}

func (w workshopResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	var plan workshopData
	err = req.Plan.Get(ctx, &plan)
	if err != nil {
		// TODO: return error
	}

	presenters := map[string]hashitalks.WorkshopPresenter{}
	for name, info := range plan.Presenters {
		presenters[name] = hashitalks.WorkshopPresenter{
			Title:    info.Title,
			Employer: info.Employer,
			Pronouns: info.Pronouns,
		}
	}

	_, err = w.client.Workshops.Update(ctx, hashitalks.Workshop{
		ID:              id.(types.String).Value,
		Title:           plan.Title,
		Description:     plan.Description,
		DurationMinutes: plan.DurationMinutes,
		Presenters:      presenters,
		MeetingInfo: hashitalks.WorkshopMeetingInfo{
			URL:      plan.MeetingInfo.URL,
			Password: plan.MeetingInfo.Password,
		},
	})
	if err != nil {
		// TODO: return error
	}

	plan.ID = id.(types.String)

	err = resp.State.Set(ctx, &plan)
	if err != nil {
		// TODO: return error
	}
}

func (w workshopResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		// TODO: return error
	}
	err = w.client.Workshops.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, hashitalks.ErrWorkshopNotFound) {
		// TODO: return error
	}
	resp.State.RemoveResource(ctx)
}
