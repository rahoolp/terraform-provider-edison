package provider

import (
	"context"
	"os"

	edison "github.com/rahoolp/terraform-provider-edison/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	client *edison.Client
}

func (p *provider) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_endpoint": {
				Type:     types.StringType,
				Optional: true,
			},
			"token": {
				Type:     types.StringType,
				Optional: true, //so that we can allow for enviornment variables as well
			},
		},
	}, nil
}

type providerData struct {
	Endpoint types.String `tfsdk:"api_endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	err := req.Config.Get(ctx, &config)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error parsing plan",
			Detail:   "An unexpected error was encountered parsing the plan. This is always a bug in the provider.\n\nDetails: " + err.Error(),
		})
		return
	}
	if config.Endpoint.Unknown {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity:  tfprotov6.DiagnosticSeverityError,
			Summary:   "Can't interpolate into provider block",
			Detail:    "Interpolating that value into the provider block doesn't give the provider enough information to run. Try hard-coding the value, instead.",
			Attribute: tftypes.NewAttributePath().WithAttributeName("api_endpoint"),
		})
		return
	}
	if config.Token.Unknown {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity:  tfprotov6.DiagnosticSeverityError,
			Summary:   "Can't interpolate into provider block",
			Detail:    "Interpolating that value into the provider block doesn't give the provider enough information to run. Try hard-coding the value, instead.",
			Attribute: tftypes.NewAttributePath().WithAttributeName("token"),
		})
		return
	}
	if config.Endpoint.Null {
		config.Endpoint.Value = os.Getenv("EDISON_API_ENDPOINT")
	}
	if config.Token.Null {
		config.Token.Value = os.Getenv("EDISON_TOKEN")
	}
	if config.Endpoint.Value == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity:  tfprotov6.DiagnosticSeverityError,
			Summary:   "Invalid provider config",
			Detail:    "api_endpoint must be set.",
			Attribute: tftypes.NewAttributePath().WithAttributeName("api_endpoint"),
		})
		return
	}
	if config.Token.Value == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity:  tfprotov6.DiagnosticSeverityError,
			Summary:   "Invalid provider config",
			Detail:    "token must be set.",
			Attribute: tftypes.NewAttributePath().WithAttributeName("token"),
		})
		return
	}
	client, err := edison.NewClient(config.Endpoint.Value, config.Token.Value)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error creating client",
			Detail:   "An unexpected error was encountered creating an API client.\n\nDetails: " + err.Error(),
		})
		return
	}
	p.client = client
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.ResourceType{
		"edison_eastore":    eastoreResourceType{},
		"edison_ehscluster": ehsclusterResourceType{},
		"edison_aw":         awResourceType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.DataSourceType{}, nil
}
