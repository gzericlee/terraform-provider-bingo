package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-bingo/internal/pkg/cmp"
	"terraform-provider-bingo/internal/pkg/sso"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	//
	// TODO: If appropriate, implement upstream provider SDK or HTTP client.
	cmpClient *cmp.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	SsoEndpoint     types.String `tfsdk:"sso_endpoint"`
	CmpEndpoint     types.String `tfsdk:"cmp_endpoint"`
	CmpAccessToken  types.String `tfsdk:"cmp_access_token"`
	CmpClientSecret types.String `tfsdk:"cmp_client_secret"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Configuration values are now available.
	// if data.Example.Null { /* ... */ }
	if data.CmpEndpoint.Null {
		resp.Diagnostics.AddError("CMP", "Endpoint is required")
	}
	if data.SsoEndpoint.Null {
		resp.Diagnostics.AddError("SSO", "Endpoint is required")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.

	if !data.CmpClientSecret.Null {
		ssoClient := sso.New(data.SsoEndpoint.Value, data.CmpClientSecret.Value)
		auth, err := ssoClient.GenerateAccessToken()
		if err != nil {
			resp.Diagnostics.AddError("SSO", fmt.Sprintf("Generate AccessToken failed: %s", err))
			return
		}
		data.CmpAccessToken = types.String{Value: auth.AccessToken}
		tflog.Trace(ctx, "Generate AT by clientSecret", map[string]interface{}{
			"input":  []string{data.SsoEndpoint.Value, data.CmpClientSecret.Value},
			"output": auth,
		})
	}

	p.cmpClient = cmp.New(data.CmpEndpoint.Value, data.CmpAccessToken.Value)

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"bingo_cmp_command": commandResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"sso_endpoint": {
				Type:                types.StringType,
				Required:            true,
				MarkdownDescription: "SSO地址，例如https://sso.bingosoft.net",
			},
			"cmp_endpoint": {
				Type:                types.StringType,
				Required:            true,
				MarkdownDescription: "CMP地址，例如https://cmp-dev.bingosoft.net",
			},
			"cmp_access_token": {
				Type:                types.StringType,
				Optional:            true,
				MarkdownDescription: "IAM颁发的用户凭证",
			},
			"cmp_client_secret": {
				Type:                types.StringType,
				Optional:            true,
				MarkdownDescription: "IAM颁发给CMP的客户端凭证",
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
