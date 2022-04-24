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
	SSO_ENDPOINT      = "SSO_ENDPOINT"
	ACCESS_TOKEN      = "ACCESS_TOKEN"
	CMP_ENDPOINT      = "CMP_ENDPOINT"
	CMP_CLIENT_SECRET = "CMP_CLIENT_SECRET"
)

type bingoCloudClient struct {
	cmpClient *cmp.Client
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"sso_endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc(SSO_ENDPOINT, nil),
					Description: "SSO地址",
				},
				"access_token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc(ACCESS_TOKEN, nil),
					Description: "IAM颁发的用户凭证",
				},
				"cmp_endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc(CMP_ENDPOINT, nil),
					Description: "CMP地址",
				},
				"cmp_client_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(CMP_CLIENT_SECRET, nil),
					Description: "IAM颁发给CMP的客户端凭证",
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
		ssoEndpoint := r.Get("sso_endpoint").(string)
		cmpEndpoint := r.Get("cmp_endpoint").(string)
		accessToken := r.Get("access_token").(string)
		cmpClientSecret := r.Get("cmp_client_secret").(string)

		if cmpClientSecret != "" {
			ssoClient := sso.New(ssoEndpoint, cmpClientSecret)
			auth, err := ssoClient.GenerateAccessToken()
			if err != nil {
				return nil, diag.Errorf(fmt.Sprintf("[SSO] Generate AccessToken failed: %s", err))
			}
			accessToken = auth.AccessToken
			tflog.Trace(ctx, "Generate AT by clientSecret", map[string]interface{}{
				"input":  []string{ssoEndpoint, cmpClientSecret},
				"output": auth,
			})
		}

		return &bingoCloudClient{
			cmpClient: cmp.New(cmpEndpoint, accessToken),
		}, nil
	}
}
