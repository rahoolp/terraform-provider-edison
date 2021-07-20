package provider

import (
	"context"
	"os"

	hashitalks "github.com/hashicorp/terraform-provider-hashitalks/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	client *hashitalks.Client
}

func (p provider) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_endpoint": {
				Type:     types.StringType,
				Optional: true,
			},
			"token": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

type providerData struct {
	Endpoint types.String `tfsdk:"api_endpoint"`
	Token    types.String `tfsdk:"token"`
}

func (p provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	err := req.Config.Get(ctx, &config)
	if err != nil {
		// TODO: return error
	}
	if config.Endpoint.Unknown {
		// TODO: return error
	}
	if config.Token.Unknown {
		// TODO: return error
	}
	if config.Endpoint.Null {
		config.Endpoint.Value = os.Getenv("HASHITALKS_API_ENDPOINT")
	}
	if config.Token.Null {
		config.Token.Value = os.Getenv("HASHITALKS_TOKEN")
	}
	if config.Endpoint.Value == "" {
		// TODO: return error
	}
	if config.Token.Value == "" {
		// TODO: return error
	}
	client, err := hashitalks.NewClient(config.Endpoint.Value, config.Token.Value)
	if err != nil {
		// TODO: return error
	}
	p.client = client
}

func (p provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.DataSourceType{}, nil
}
