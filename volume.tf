terraform {
    required_providers {
        cmccloud = {
            source  = "github.com/cmccloud/cmccloud"
			version = "0.1.0"
        }
    }
}

provider "cmccloud" {
    api_endpoint = "https://api.cloud.cmctelecom.vn/ver2"
    api_key = "vTMSG7F9mFKnNRYIz8eA9N9XrHJ4zP"
}


resource "cmccloud_volume" "vol1-ssd" {
    name = "vol1.ssd"
    size = 30
    type = "ssd"
    server_id = "040e6918-409e-4346-8827-cd328e113f91"
}