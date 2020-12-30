terraform {
    required_providers {
        cmccloud = {
            source  = "github.com/cmccloud/cmccloud"
			version = "0.1.0"
        }
    }
}

provider "cmccloud" {
    api_endpoint = "https://api.cloud.cmctelecom.vn"
    api_key = "vTMSG7F9mFKnNRYIz8eA9N9XrHJ4zP"
}


data "cmccloud_images" "all" {}

# Returns all images
output "all_images" {
  value = data.cmccloud_images.all
}