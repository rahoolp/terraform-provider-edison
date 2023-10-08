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

type avResourceType struct {
}

func (e avResourceType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"account_id": {
				Type:     types.StringType,
				Required: true,
			},
			"tenant_id": {
				Type:     types.StringType,
				Required: true,
			},
			"tenant_folder": {
				Type:     types.StringType,
				Computed: true,
			},
			"tenant_queue": {
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

type avData struct {
	ID           types.String `tfsdk:"id"`
	AccountID    types.String `tfsdk:"account_id"`
	TenantID     types.String `tfsdk:"tenant_id"`
	TenantFolder types.String `tfsdk:"tenant_folder"`
	TenantQueue  types.String `tfsdk:"tenant_queue"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func (s avResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
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

	return avResource{client: prov.client}, nil
}

type avResource struct {
	client *edison.Client
}

func (e avResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

	tflog.Info(ctx, "AV Create..")
	fmt.Println("AV Create")

	var av avData
	err := req.Plan.Get(ctx, &av)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error parsing plan",
			Detail:   "An unexpected error was encountered parsing the plan. This is always a bug in the provider.\n\nDetails: " + err.Error(),
		})
		return
	}

	//go func() {
	var tenantFolder string = "http://s3.amazonaws.com/av_bucket/" + av.TenantID.Value
	now := time.Now()
	var createdAt string = now.Format("2006-01-02 15:04:05")
	var updatedAt string = now.Format("2006-01-02 15:04:05")
	var tenantQueue string = "arn:aws:mq:us-east-1:" + av.TenantID.Value

	eav, err := e.client.AVs.Create(ctx, edison.AV{
		AccountID:    av.AccountID.Value,
		TenantID:     av.TenantID.Value,
		TenantFolder: tenantFolder,
		TenantQueue:  tenantQueue,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "AV Create: "+err.Error())
	}

	av.ID = types.String{Value: eav.ID}
	av.CreatedAt = types.String{Value: eav.CreatedAt}
	av.UpdatedAt = types.String{Value: eav.UpdatedAt}
	av.TenantFolder = types.String{Value: tenantFolder}
	av.TenantQueue = types.String{Value: tenantQueue}

	err = resp.State.Set(ctx, &av)
	if err != nil {
		tflog.Info(ctx, "AV Create: "+err.Error())
	}
	//}()
}

// func (e avResource) Create1(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
// 	tflog.Info(ctx, "AV Create..")

// 	var av avData
// 	err := req.Plan.Get(ctx, &av)
// 	if err != nil {
// 		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
// 			Severity: tfprotov6.DiagnosticSeverityError,
// 			Summary:  "Error parsing plan",
// 			Detail:   "An unexpected error was encountered parsing the plan. This is always a bug in the provider.\n\nDetails: " + err.Error(),
// 		})
// 		return
// 	}

// 	// Queue the resource creation task (asynchronous)
// 	go func() {
// 		var dnsEP string = "https://av-04.ehs.edison.gehealthcare.com/"
// 		// Implement resource creation logic here
// 		// This can involve API calls, external processes, etc.

// 		// Simulate resource creation (replace with your actual logic)
// 		//time.Sleep(5 * time.Second)

// 		now := time.Now()
// 		var createdAt string = now.Format("2006-01-02 15:04:05")
// 		var updatedAt string = now.Format("2006-01-02 15:04:05")

// 		// Once the resource is ready, set the state
// 		// This should be done carefully to avoid race conditions
// 		// You might need to use a synchronization mechanism
// 		err := resp.State.Set(ctx, &avData{
// 			ID:              av.ID, // Replace with the actual resource ID
// 			ConcurrentUsers: av.ConcurrentUsers,
// 			DicomEndPoint:   av.DicomEndPoint,
// 			EHSClusterID:    types.String{Value: av.EHSClusterID.Value},
// 			DNSEndPoint:     types.String{Value: dnsEP},
// 			CreatedAt:       types.String{Value: createdAt},
// 			UpdatedAt:       types.String{Value: updatedAt},
// 		})

// 		if err != nil {
// 			tflog.Info(ctx, "AV Create: "+err.Error())
// 		}
// 	}()

// 	// Return immediately, allowing Terraform to continue
// }

func (e avResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

	tflog.Info(ctx, "AV Read..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "AV Read: "+err.Error())
	}

	av, err := e.client.AVs.Get(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrEHSClusterNotFound) {
		tflog.Info(ctx, "AV Read: "+err.Error())
	} else if errors.Is(err, edison.ErrEHSClusterNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}

	err = resp.State.Set(ctx, &avData{
		ID:           types.String{Value: av.ID},
		TenantID:     types.String{Value: av.TenantID},
		AccountID:    types.String{Value: av.AccountID},
		TenantFolder: types.String{Value: av.TenantFolder},
		TenantQueue:  types.String{Value: av.TenantQueue},
		CreatedAt:    types.String{Value: av.CreatedAt},
		UpdatedAt:    types.String{Value: av.UpdatedAt},
	})

	if err != nil {
		if err != nil {
			tflog.Info(ctx, "AV: "+err.Error())
			// Return a RetryableError here if you want to retry this resource read
			//resp.State.RetryableError(ctx, "Failed to set resource state")
			return
		}
		tflog.Info(ctx, "AV: "+err.Error())
	}
}

func (e avResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

	tflog.Info(ctx, "AV Update..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "AV Update: "+err.Error())
	}

	var av avData
	err = req.Plan.Get(ctx, &av)
	if err != nil {
		tflog.Info(ctx, "AV Update: "+err.Error())
	}

	now := time.Now()
	var updatedAt string = now.Format("2006-01-02 15:04:05")

	_, err = e.client.AVs.Update(ctx, edison.AV{
		ID:           id.(types.String).Value,
		TenantID:     av.TenantID.Value,
		AccountID:    av.AccountID.Value,
		TenantFolder: av.TenantFolder.Value,
		TenantQueue:  av.TenantQueue.Value,
		CreatedAt:    av.CreatedAt.Value,
		UpdatedAt:    updatedAt,
	})
	if err != nil {
		tflog.Info(ctx, "AV Update: "+err.Error())
	}
	av.ID = id.(types.String)

	err = resp.State.Set(ctx, &av)
	if err != nil {
		tflog.Info(ctx, "AV Update: "+err.Error())
	}
}

func (e avResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

	tflog.Info(ctx, "AV Delete..")

	id, err := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"))
	if err != nil {
		tflog.Info(ctx, "AV Delete: "+err.Error())
	}
	err = e.client.AVs.Delete(ctx, id.(types.String).Value)
	if err != nil && !errors.Is(err, edison.ErrEHSClusterNotFound) {
		tflog.Info(ctx, "AV Delete: "+err.Error())
	}
	resp.State.RemoveResource(ctx)
}
