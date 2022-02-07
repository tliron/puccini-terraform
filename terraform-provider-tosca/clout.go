package main

import (
	contextpkg "context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/segmentio/ksuid"
	"github.com/tliron/kutil/format"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/clout/js"
)

//
// CloutResourceType
//

type CloutResourceType struct{}

// tfsdk.ResourceType interface
func (self *CloutResourceType) GetSchema(context contextpkg.Context) (tfsdk.Schema, diag.Diagnostics) {
	//  Terraform requires *at least one* attribute in the schema
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"service_template": {
				Type:        types.StringType,
				Required:    true,
				Description: "service template URL",
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"clout": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

// tfsdk.ResourceType interface
func (self *CloutResourceType) NewResource(context contextpkg.Context, provider tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &CloutResource{
		provider: provider.(*ToscaProvider),
	}, nil
}

//
// CloutResource
//

type CloutResource struct {
	provider *ToscaProvider
}

// tfsdk.Resource interface
func (self *CloutResource) Create(context contextpkg.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	helper := CreateResourceHelper{context, request, response}

	helper.SetAttribute("id", ksuid.New().String())

	var serviceTemplateUrl string
	if !helper.GetAttribute("service_template", &serviceTemplateUrl) {
		return
	}
	helper.SetAttribute("service_template", serviceTemplateUrl)

	url, err := urlpkg.NewValidURL(serviceTemplateUrl, nil, self.provider.urlContext)
	if err != nil {
		helper.AddError("bad URL", err.Error())
		return
	}

	_, serviceTemplate, problems, err := self.provider.parserContext.Parse(url, nil, nil, nil)
	if err != nil {
		helper.AddError("TOSCA parsing error", err.Error()+problems.String())
		return
	}

	clout, err := serviceTemplate.Compile(false)
	if err != nil {
		helper.AddError("TOSCA compilation", err.Error()+problems.String())
		return
	}

	js.Resolve(clout, problems, self.provider.urlContext, true, "yaml", true, false, false)
	if !problems.Empty() {
		helper.AddError("TOSCA resolution", problems.String())
		return
	}

	js.Coerce(clout, problems, self.provider.urlContext, true, "yaml", true, false, false)
	if !problems.Empty() {
		helper.AddError("TOSCA coercion", problems.String())
		return
	}

	cloutYaml, err := format.EncodeYAML(clout, " ", false)
	if err != nil {
		helper.AddError("Clout encoding", err.Error())
		return
	}

	helper.SetAttribute("clout", cloutYaml)
}

// tfsdk.Resource interface
func (self *CloutResource) Read(context contextpkg.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
}

// tfsdk.Resource interface
func (self *CloutResource) Update(context contextpkg.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
}

// tfsdk.Resource interface
func (self *CloutResource) Delete(context contextpkg.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
}

// tfsdk.Resource interface
func (self *CloutResource) ImportState(context contextpkg.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	helper := ImportResourceStateHelper{context, request, response}
	helper.PassthroughID("id")
}
