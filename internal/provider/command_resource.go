package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"terraform-provider-bingo/internal/pkg/cmp"
)

type commandResourceType struct{}

func (its commandResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "CMP指令",

		Attributes: map[string]tfsdk.Attribute{
			"host_type": {
				Type:                types.StringType,
				Required:            true,
				MarkdownDescription: "宿主机类型，1:虚拟机,2:物理机",
			},
			"content": {
				Type:                types.StringType,
				Required:            true,
				MarkdownDescription: "命令内容",
			},
			"instance_ids": {
				Type:                types.StringType,
				Required:            true,
				MarkdownDescription: "实例编号，多个用逗号分割",
			},
			"record_id": {
				Type:                types.StringType,
				Computed:            true,
				MarkdownDescription: "记录ID",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"task_id": {
				Type:                types.StringType,
				Computed:            true,
				MarkdownDescription: "任务ID",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"status": {
				Type:                types.StringType,
				Computed:            true,
				MarkdownDescription: "指令状态",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

func (its commandResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)
	return commandResource{
		provider: provider,
	}, diags
}

type commandResourceData struct {
	Id          types.String `tfsdk:"id"`
	RecordId    types.String `tfsdk:"record_id"`
	TaskId      types.String `tfsdk:"task_id"`
	Status      types.String `tfsdk:"status"`
	HostType    types.String `tfsdk:"host_type"`
	Content     types.String `tfsdk:"content"`
	InstanceIds types.String `tfsdk:"instance_ids"`
}

type commandResource struct {
	provider provider
}

func (its commandResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data commandResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// cmp, err := d.provider.client.CreateCMP(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cmp, got error: %s", err))
	//     return
	// }
	input := &cmp.CommandInput{}
	input.HostType = data.HostType.Value
	input.Name = "terraform-deploy-" + time.Now().Format("20060102150405")
	input.Description = "Created by `terraform-provider-bingo`"
	input.Content = data.Content.Value
	input.InstanceIds = data.InstanceIds.Value

	output, err := its.provider.cmpClient.CreateCommand(input)
	if err != nil {
		resp.Diagnostics.AddError("CMP", fmt.Sprintf("Unable to create command, got error: %s", err))
		return
	}

	data.Id = types.String{Value: output.RecordId}
	data.RecordId = types.String{Value: output.RecordId}
	data.TaskId = types.String{Value: output.TaskId}
	data.Status = types.String{Value: output.Status}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "Sent a command successfully", map[string]interface{}{
		"input":  input,
		"output": output,
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (its commandResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data commandResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// cmp, err := d.provider.client.ReadCMP(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cmp, got error: %s", err))
	//     return
	// }
	input := &cmp.DescribeCommandInput{
		ConStr: "deploy",
		SqlId:  "command.selectRecordById",
		Params: struct {
			Id string `json:"id"`
		}{Id: data.Id.Value},
	}
	output, err := its.provider.cmpClient.DescribeCommand(input)
	if err != nil {
		resp.Diagnostics.AddError("CMP", fmt.Sprintf("Unable to read command, got error: %s", err))
		return
	}

	data.Id = types.String{Value: output.Id}
	data.RecordId = types.String{Value: output.Id}
	data.TaskId = types.String{Value: output.TaskId}
	data.Status = types.String{Value: output.Status}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{cmp.CommandStatusNew, cmp.CommandStatusDeploying},
		Target:       []string{cmp.CommandStatusSuccess},
		Refresh:      refreshCommandStatus(its.provider.cmpClient, input, cmp.CommandStatusFailed),
		Timeout:      30 * time.Minute,
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("CMP", fmt.Sprintf("waiting for command (%s) : %s", output.Id, err))
		return
	}

	tflog.Trace(ctx, "Executed a command successfully", map[string]interface{}{
		"input":  input,
		"output": output,
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (its commandResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data commandResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// cmp, err := d.provider.client.UpdateCMP(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cmp, got error: %s", err))
	//     return
	// }

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (its commandResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data commandResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// cmp, err := d.provider.client.DeleteCMP(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cmp, got error: %s", err))
	//     return
	// }

	resp.State.RemoveResource(ctx)
}

func (its commandResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStateNotImplemented(ctx, "", resp)
}

func refreshCommandStatus(cmpClient *cmp.Client, input *cmp.DescribeCommandInput, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := cmpClient.DescribeCommand(input)
		if err != nil {
			return nil, "", err
		}
		if output.Status == failState {
			steps, err := cmpClient.DescribeCommandSteps(&cmp.DescribeCommandStepsInput{
				SqlId:  "task.listAllStepsForAgent",
				ConStr: "deploy",
				Params: struct {
					TaskId string `json:"taskId"`
				}{TaskId: output.TaskId},
				Page:     1,
				PageSize: 1,
			})
			if err != nil {
				return nil, "", err
			}
			if len(steps) != 1 {
				return nil, "", fmt.Errorf("command step not found (%s)", output.TaskId)
			}
			return output, output.Status, fmt.Errorf("failed to reach target state. Reason: %s", steps[0].StepLog)
		}
		return output, output.Status, nil
	}
}
