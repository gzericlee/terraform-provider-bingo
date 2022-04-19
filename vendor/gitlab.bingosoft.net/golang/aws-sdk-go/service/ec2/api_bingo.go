package ec2

import (
	"gitlab.bingosoft.net/golang/aws-sdk-go/aws/awsutil"
	"gitlab.bingosoft.net/golang/aws-sdk-go/aws/request"
)

const opDescribeServices = "DescribeServices"

func (c *EC2) DescribeServicesRequest(input *DescribeServicesInput) (req *request.Request, output *DescribeServicesOutput) {
	op := &request.Operation{
		Name:       opDescribeServices,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DescribeServicesInput{}
	}

	output = &DescribeServicesOutput{}
	req = c.newRequest(op, input, output)

	return
}

func (c *EC2) DescribeServices(input *DescribeServicesInput) (*DescribeServicesOutput, error) {
	req, out := c.DescribeServicesRequest(input)
	return out, req.Send()
}

type DescribeServicesInput struct {
	_ struct{} `type:"structure"`

	ServiceCodes []*string `min:"1" type:"list"`
}

func (s DescribeServicesInput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeServicesInput) GoString() string {
	return s.String()
}

type DescribeServicesOutput struct {
	_ struct{} `type:"structure"`

	DescribeServicesResult *struct{
		_ struct{} `type:"structure"`

		ServiceInfo *struct{

			Items []*struct {
				_ struct{} `type:"structure"`

				ServiceId *string `locationName:"ServiceId" type:"string"`

				ServiceName *string `locationName:"ServiceName" type:"string"`

				ServiceCode *string `locationName:"ServiceCode" type:"string"`

				ApiAddress *string `locationName:"ApiAddress" type:"string"`

				UiAddress *string `locationName:"UiAddress" type:"string"`

				Status *string `locationName:"Status" type:"string"`

			} `xml:"serviceInfo" type:"list"`

		} `locationName:"ServiceInfo" type:"structure"`

	} `locationName:"DescribeServicesResult" type:"structure"`

	ResponseMetadata *struct{

		RequestId *string `locationName:"RequestId" type:"string"`

	}`locationName:"ResponseMetadata" type:"structure"`
}

func (s DescribeServicesOutput) String() string {
	return awsutil.Prettify(s)
}

func (s DescribeServicesOutput) GoString() string {
	return s.String()
}






const opResizeVolume = "ResizeVolume"

func (c *EC2) ResizeVolumeRequest(input *ResizeVolumeInput) (req *request.Request, output *ResizeVolumeOutput) {
	op := &request.Operation{
		Name:       opResizeVolume,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &ResizeVolumeInput{}
	}

	output = &ResizeVolumeOutput{}
	req = c.newRequest(op, input, output)

	return
}

func (c *EC2) ResizeVolume(input *ResizeVolumeInput) (*ResizeVolumeOutput, error) {
	req, out := c.ResizeVolumeRequest(input)
	return out, req.Send()
}

type ResizeVolumeInput struct {
	_ struct{} `type:"structure"`

	VolumeId *string `locationName:"volumeId" type:"string"`

	Size *int64 `locationName:"size" type:"integer"`
}

func (s ResizeVolumeInput) String() string {
	return awsutil.Prettify(s)
}

func (s ResizeVolumeInput) GoString() string {
	return s.String()
}

type ResizeVolumeOutput struct {
	_ struct{} `type:"structure"`

	VolumeId *string `locationName:"volumeId" type:"string"`

	Size *int64 `locationName:"size" type:"integer"`
}

func (s ResizeVolumeOutput) String() string {
	return awsutil.Prettify(s)
}

func (s ResizeVolumeOutput) GoString() string {
	return s.String()
}