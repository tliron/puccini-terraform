package main

import (
	contextpkg "context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

//
// ProviderConfigureHelper
//

type ProviderConfigureHelper struct {
	context  contextpkg.Context
	request  tfsdk.ConfigureProviderRequest
	response *tfsdk.ConfigureProviderResponse
}

func (self *ProviderConfigureHelper) Get(target interface{}) bool {
	diagnostics := self.request.Config.Get(self.context, target)
	self.response.Diagnostics.Append(diagnostics...)
	return !self.response.Diagnostics.HasError()
}

func (self *ProviderConfigureHelper) AddWarning(summary string, detail string) {
	self.response.Diagnostics.AddWarning(summary, detail)
}

func (self *ProviderConfigureHelper) AddError(summary string, detail string) {
	self.response.Diagnostics.AddError(summary, detail)
}

//
// CreateResourceHelper
//

type CreateResourceHelper struct {
	context  contextpkg.Context
	request  tfsdk.CreateResourceRequest
	response *tfsdk.CreateResourceResponse
}

func (self *CreateResourceHelper) AddWarning(summary string, detail string) {
	self.response.Diagnostics.AddWarning(summary, detail)
}

func (self *CreateResourceHelper) AddError(summary string, detail string) {
	self.response.Diagnostics.AddError(summary, detail)
}

func (self *CreateResourceHelper) GetAttribute(name string, target interface{}) bool {
	path := tftypes.NewAttributePath().WithAttributeName(name)
	diagnostics := self.request.Plan.GetAttribute(self.context, path, target)
	self.response.Diagnostics.Append(diagnostics...)
	return !self.response.Diagnostics.HasError()
}

func (self *CreateResourceHelper) SetAttribute(name string, value interface{}) bool {
	path := tftypes.NewAttributePath().WithAttributeName(name)
	diagnostics := self.response.State.SetAttribute(self.context, path, value)
	self.response.Diagnostics.Append(diagnostics...)
	return !self.response.Diagnostics.HasError()
}

//
// ImportResourceStateHelper
//

type ImportResourceStateHelper struct {
	context  contextpkg.Context
	request  tfsdk.ImportResourceStateRequest
	response *tfsdk.ImportResourceStateResponse
}

func (self *ImportResourceStateHelper) PassthroughID(name string) {
	path := tftypes.NewAttributePath().WithAttributeName(name)
	tfsdk.ResourceImportStatePassthroughID(self.context, path, self.request, self.response)
}
