package cmp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/levigross/grequests"

	"terraform-provider-bingo/utils"
)

type CommandInput struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	HostType    string `json:"hostType"`
	InstanceIds string `json:"instanceIds"`
	Description string `json:"description"`
}

func (its CommandInput) String() string {
	return utils.Prettify(its)
}

type CommandOutput struct {
	RecordId   string    `json:"recordId"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	UserId     string    `json:"userId"`
	CreateTime time.Time `json:"createTime"`
	Status     string    `json:"status"`
	Machines   string    `json:"machines"`
	TaskId     string    `json:"taskId"`
}

func (its CommandOutput) String() string {
	return utils.Prettify(its)
}

type DescribeCommandInput struct {
	ConStr string `json:"conStr"`
	SqlId  string `json:"sqlId"`
	Params struct {
		Id string `json:"id"`
	} `json:"params"`
}

func (its DescribeCommandInput) String() string {
	return utils.Prettify(its)
}

type DescribeCommandOutput struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Content     string      `json:"content"`
	UserId      string      `json:"userId"`
	CreateTime  time.Time   `json:"createTime"`
	Status      string      `json:"status"`
	Machines    string      `json:"machines"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     time.Time   `json:"endTime"`
	TaskId      string      `json:"taskId"`
	Description interface{} `json:"description"`
}

func (its DescribeCommandOutput) String() string {
	return utils.Prettify(its)
}

type DescribeCommandStepsInput struct {
	SqlId  string `json:"sqlId"`
	ConStr string `json:"conStr"`
	Params struct {
		TaskId string `json:"taskId"`
	} `json:"params"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (its DescribeCommandStepsInput) String() string {
	return utils.Prettify(its)
}

type DescribeCommandStepsOutput struct {
	StepId       string    `json:"stepId"`
	MachineId    string    `json:"machineId"`
	StepContent  string    `json:"stepContent"`
	StepDesc     string    `json:"stepDesc"`
	StepStatus   string    `json:"stepStatus"`
	Progress     string    `json:"progress"`
	StepLog      string    `json:"stepLog"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	CreateTime   time.Time `json:"createTime"`
	MachineName  string    `json:"machineName"`
	MachineCode  string    `json:"machineCode"`
	InstanceCode string    `json:"instanceCode"`
	Agent        string    `json:"agent"`
}

func (its DescribeCommandStepsOutput) String() string {
	return utils.Prettify(its)
}

func (its *Client) CreateCommand(input *CommandInput) (*CommandOutput, error) {
	options := its.config.Options
	options.JSON = input
	resp, err := grequests.Post(fmt.Sprintf("%v/%v/api/command/sendCommand", its.config.Endpoint, its.config.MainApiContext), &options)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	content := resp.String()
	if !resp.Ok {
		err = fmt.Errorf("[CMP] Response code: [%v]，result: [%s]", resp.StatusCode, content)
		return nil, err
	}

	output := &CommandOutput{}
	err = json.Unmarshal([]byte(content), &output)

	return output, err
}

func (its *Client) DescribeCommand(input *DescribeCommandInput) (*DescribeCommandOutput, error) {
	options := its.config.Options
	options.JSON = input
	resp, err := grequests.Post(fmt.Sprintf("%v/%v/api/getEntity", its.config.Endpoint, its.config.MainApiContext), &options)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	content := resp.String()
	if !resp.Ok {
		err = fmt.Errorf("[CMP] Response code: [%v]，result: [%s]", resp.StatusCode, content)
		return nil, err
	}

	output := &DescribeCommandOutput{}
	err = json.Unmarshal([]byte(content), &output)

	return output, err
}

func (its *Client) DescribeCommandSteps(input *DescribeCommandStepsInput) ([]*DescribeCommandStepsOutput, error) {
	options := its.config.Options
	options.JSON = input

	resp, err := grequests.Post(fmt.Sprintf("%v/%v/api/queryPageList", its.config.Endpoint, its.config.MainApiContext), &options)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	content := resp.String()
	if !resp.Ok {
		err = fmt.Errorf("[CMP] Response code: [%v]，result: [%s]", resp.StatusCode, content)
		return nil, err
	}

	var steps []*DescribeCommandStepsOutput
	err = json.Unmarshal([]byte(content), &steps)

	return steps, err
}
