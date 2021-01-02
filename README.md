CMC Cloud Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Using The CMC Provider
---------------------

```mkdir test
cd test
wget https://raw.githubusercontent.com/cmc-cloud/terraform-provider-cmccloud/main/examples/server.tf server.tf
```

get your Cloud Api Key from https://portal.cloud.cmctelecom.vn/account-settings?to=settings
change api_key to your own api key in server.tf

```terraform init
```

Building The CMC CProvider
---------------------

Clone repository to: `$GOPATH/src/github.com/cmc-cloud/terraform-provider-cmccloud`

```sh
$ mkdir -p $GOPATH/src/github.com/cmc-cloud; cd $GOPATH/src/github.com/cmc-cloud
$ git clone git@github.com:cmc-cloud/terraform-provider-cmccloud
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/cmc-cloud/terraform-provider-cmccloud
$ make build
```
