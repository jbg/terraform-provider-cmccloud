module terraform-provider-cmccloud

go 1.15

require (
	github.com/cmc-cloud/gocmcapi v0.0.0
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	golang.org/x/tools v0.0.0-20201110201400-7099162a900a
)

replace github.com/cmc-cloud/gocmcapi => C:\code\CMC\terraform\gocmcapi