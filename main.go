package main

import (
	"context"

	"github.com/hashicorp/terraform-provider-hashitalks/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	ctx := context.Background()

	tfsdk.Serve(ctx, provider.New, tfsdk.ServeOpts{
		Name: "registry.terraform.io/hashicorp/hashitalks",
	})
}
