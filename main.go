package main

import (
	"context"

	"github.com/rahoolp/terraform-provider-edison/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	ctx := context.Background()

	tfsdk.Serve(ctx, provider.New, tfsdk.ServeOpts{
		Name: "registry.terraform.io/hashicorp/edison",
	})
}
