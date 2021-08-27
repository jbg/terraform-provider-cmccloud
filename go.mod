module terraform-provider-cmccloud

go 1.15

require (
	// run: go get -d github.com/cmc-cloud/gocmcapi@eb9c186 to get correct lastest version, eb9c186 = hash github commit
	github.com/cmc-cloud/gocmcapi v0.0.0-20201230035051-54fa9ab87e70
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/zclconf/go-cty v1.4.1 // indirect
	golang.org/x/tools v0.0.0-20201110201400-7099162a900a // indirect
)

// uncomment this line when build from code
// replace github.com/cmc-cloud/gocmcapi => C:\code\terraform\gocmcapi
