package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-bingo/internal/pkg/cmp"
)

func resourceCmpCommand() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "CMP指令",

		CreateContext: resourceCmpCommandCreate,
		ReadContext:   resourceCmpCommandRead,
		UpdateContext: resourceCmpCommandUpdate,
		DeleteContext: resourceCmpCommandDelete,

		Schema: map[string]*schema.Schema{
			"host_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "宿主机类型，1:虚拟机,2:物理机",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "命令内容",
			},
			"instance_ids": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例编号，多个用逗号分割",
			},
			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "记录ID",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "任务ID",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "指令状态",
			},
		},
	}
}

func resourceCmpCommandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bingoCloudClient)

	input := &cmp.CommandInput{}
	input.HostType = d.Get("host_type").(string)
	input.Name = "terraform-deploy-" + time.Now().Format("20060102150405")
	input.Description = "Created by `terraform-provider-bingo`"
	input.Content = d.Get("content").(string)
	input.InstanceIds = d.Get("instance_ids").(string)

	output, err := client.cmpClient.CreateCommand(input)
	if err != nil {
		return diag.Errorf(fmt.Sprintf("[CMP] Unable to create command, got error: %s", err))
	}

	d.SetId(output.RecordId)
	d.Set("record_id", output.RecordId)
	d.Set("task_id", output.TaskId)
	d.Set("status", output.Status)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Debug(ctx, "Sent a command successfully", map[string]interface{}{
		"input":  input,
		"output": output,
	})

	return nil
}

func resourceCmpCommandRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bingoCloudClient)

	input := &cmp.DescribeCommandInput{
		ConStr: "deploy",
		SqlId:  "command.selectRecordById",
		Params: struct {
			Id string `json:"id"`
		}{Id: d.Id()},
	}
	output, err := client.cmpClient.DescribeCommand(input)
	if err != nil {
		return diag.Errorf(fmt.Sprintf("[CMP] Unable to read command, got error: %s", err))
	}

	d.SetId(output.Id)
	d.Set("record_id", output.Id)
	d.Set("task_id", output.TaskId)
	d.Set("status", output.Status)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{cmp.CommandStatusNew, cmp.CommandStatusDeploying},
		Target:       []string{cmp.CommandStatusSuccess},
		Refresh:      refreshCommandStatus(client.cmpClient, input, cmp.CommandStatusFailed),
		Timeout:      30 * time.Minute,
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(fmt.Sprintf("[CMP] Waiting for command (%s) : %s", output.Id, err))
	}

	tflog.Debug(ctx, "[CMP] Executed a command successfully", map[string]interface{}{
		"input":  input,
		"output": output,
	})

	return nil
}

func resourceCmpCommandUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	tflog.Debug(ctx, "[CMP] Updated a command successfully", map[string]interface{}{})
	return nil
}

func resourceCmpCommandDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	tflog.Debug(ctx, "[CMP] Deleted a command successfully", map[string]interface{}{})
	return nil
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
