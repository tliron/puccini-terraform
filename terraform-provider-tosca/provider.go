package main

import (
	contextpkg "context"
	"os"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/tosca/parser"
)

var stderr = os.Stderr

//
// ToscaProvider
//

type ToscaProvider struct {
	parserContext *parser.Context
	urlContext    *urlpkg.Context
	configured    bool
	client        *hashicups.Client
}

func ProviderFactory() tfsdk.Provider {
	return &ToscaProvider{
		parserContext: parser.NewContext(),
		urlContext:    urlpkg.NewContext(),
	}
}

// tfsdk.Provider interface
func (self *ToscaProvider) GetSchema(context contextpkg.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"username": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}, nil
}

// tfsdk.Provider interface
func (self *ToscaProvider) Configure(context contextpkg.Context, request tfsdk.ConfigureProviderRequest, response *tfsdk.ConfigureProviderResponse) {
	helper := ProviderConfigureHelper{context, request, response}

	type ProviderConfig struct {
		Username types.String `tfsdk:"username"`
		Host     types.String `tfsdk:"host"`
		Password types.String `tfsdk:"password"`
	}

	var config ProviderConfig
	if !helper.Get(&config) {
		return
	}

	// TODO: fake, remove!

	// User must provide a user to the provider
	var username string
	if config.Username.Unknown {
		// Cannot connect to client with an unknown value
		helper.AddWarning("Unable to create client", "Cannot use unknown value as username")
		return
	}

	if config.Username.Null {
		username = os.Getenv("HASHICUPS_USERNAME")
	} else {
		username = config.Username.Value
	}

	if username == "" {
		// Error vs warning - empty value must stop execution
		helper.AddError("Unable to find username", "Username cannot be an empty string")
		return
	}

	// User must provide a password to the provider
	var password string
	if config.Password.Unknown {
		// Cannot connect to client with an unknown value
		helper.AddError("Unable to create client", "Cannot use unknown value as password")
		return
	}

	if config.Password.Null {
		password = os.Getenv("HASHICUPS_PASSWORD")
	} else {
		password = config.Password.Value
	}

	if password == "" {
		// Error vs warning - empty value must stop execution
		helper.AddError("Unable to find password", "password cannot be an empty string")
		return
	}

	// User must specify a host
	var host string
	if config.Host.Unknown {
		// Cannot connect to client with an unknown value
		helper.AddError("Unable to create client", "Cannot use unknown value as host")
		return
	}

	if config.Host.Null {
		host = os.Getenv("HASHICUPS_HOST")
	} else {
		host = config.Host.Value
	}

	if host == "" {
		// Error vs warning - empty value must stop execution
		helper.AddError("Unable to find host", "Host cannot be an empty string")
		return
	}

	self.configured = true
}

// tfsdk.Provider interface
func (self *ToscaProvider) GetResources(context contextpkg.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		PROVIDER_NAME + "_clout": new(CloutResourceType),
	}, nil
}

// tfsdk.Provider interface
func (self *ToscaProvider) GetDataSources(context contextpkg.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return nil, nil
}
