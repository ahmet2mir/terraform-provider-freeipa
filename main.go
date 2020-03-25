package main

import (
	"github.com/ahmet2mir/terraform-provider-freeipa/freeipa"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return freeipa.Provider()
		},
	})
}
