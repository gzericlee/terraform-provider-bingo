package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-bingo/internal/pkg/cmp"
	"terraform-provider-bingo/internal/pkg/sso"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

const (
	IAM_ENDPOINT      = "IAM_ENDPOINT"
	CMP_ENDPOINT      = "CMP_ENDPOINT"
	IAM_CLIENT_ID     = "IAM_CLIENT_ID"
	IAM_CLIENT_SECRET = "IAM_CLIENT_SECRET"
)

type bingoCloudClient struct {
	cmpClient *cmp.Client
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"iam_endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc(IAM_ENDPOINT, nil),
					Description: "IAM地址",
				},
				"cmp_endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc(CMP_ENDPOINT, nil),
					Description: "CMP地址",
				},
				"iam_client_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc(IAM_CLIENT_ID, nil),
					Description: "IAM颁发的客户端ID",
				},
				"iam_client_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(IAM_CLIENT_SECRET, nil),
					Description: "IAM颁发的客户端Secret",
				},
			},

			DataSourcesMap: map[string]*schema.Resource{},

			ResourcesMap: map[string]*schema.Resource{
				"bingo_cmp_command": resourceCmpCommand(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, r *schema.ResourceData) (interface{}, diag.Diagnostics) {
		iamEndpoint := r.Get("iam_endpoint").(string)
		cmpEndpoint := r.Get("cmp_endpoint").(string)
		iamClientId := r.Get("iam_client_id").(string)
		iamClientSecret := r.Get("iam_client_secret").(string)

		ssoClient := sso.New(iamEndpoint, iamClientId, iamClientSecret)
		auth, err := ssoClient.GenerateAccessToken()
		if err != nil {
			return nil, diag.Errorf(fmt.Sprintf("[SSO] Generate AccessToken failed: %s", err))
		}

		tflog.Trace(ctx, "Generate AT by clientSecret", map[string]interface{}{
			"input":  []string{iamEndpoint, iamClientId, iamClientSecret},
			"output": auth,
		})

		return &bingoCloudClient{
			cmpClient: cmp.New(cmpEndpoint, auth.AccessToken),
		}, nil
	}
}
