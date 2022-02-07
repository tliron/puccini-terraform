package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/tliron/kutil/logging"
	_ "github.com/tliron/kutil/logging/simple"
)

const PROVIDER_NAME = "tosca"

func main() {
	logging.Configure(0, nil)

	tfsdk.Serve(context.Background(), ProviderFactory, tfsdk.ServeOpts{
		Name: PROVIDER_NAME,
	})
}
