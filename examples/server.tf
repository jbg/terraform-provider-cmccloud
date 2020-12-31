terraform {
    required_providers {
        cmccloud = {
            source  = "github.com/cmc-cloud/cmccloud"
			version = "0.1.0"
        }
    }
}

provider "cmccloud" {
    api_endpoint = "https://api.cloud.cmctelecom.vn/ver2"
    api_key = "vTMSG7F9mFKnNRYIz8eA9N9XrHJ4zP"
}


resource "cmccloud_server" "dat-terraform" {
    name = "dat.terraform2.com"
    cpu = 1
    ram = 1
    root = 20
    image_type = "image"
    image_id = "472c8730-bc28-42d6-8680-151402be1169"
    region = "hn"
	auto_backup = true
	enable_private_network = false
	backup_schedule = "daily-03:00"
}


resource "cmccloud_volume" "vol1-ssd" {    
    name    = "vol1.ssd"
    region  = "HN"
    size    = 60
	type    = "ssd"
	server_id = "040e6918-409e-4346-8827-cd328e113f91"
}

resource "cmccloud_vpc" "terra_vpc" {
    name = "terraform.vpc"
    description = "test via terraform"
    region = "hn"
    cidr = "10.10.0.0/16"
}

resource "cmccloud_floating_ip" "terra_floatip_1"{
    vpc_id = cmccloud_vpc.terra_vpc.id
}

resource "cmccloud_network" "private_network1" {
    name = "terraform.private.network1"
    description = "network create from terraform"
    gateway = "10.10.0.1"
    netmask = "255.255.255.0"
    firewall_id = "allow"
    vpc_id = cmccloud_vpc.terra_vpc.id
    server_ids = [ cmccloud_server.dat-terraform.id, "cf967e0a-2784-4245-8b7e-583b89d7fb66" ]
}

resource "cmccloud_firewall_vpc" "terra_firewall1" {
    name = "terraform.firewall1"
    description = "vpc firewall create from terraform"
    vpc_id = cmccloud_vpc.terra_vpc.id
    inbound_rule {
        protocol         = "tcp"
        port_range       = "22"
        cidrs = ["192.168.1.0/24", "192.168.2.0/24"]
    }
    inbound_rule {
        protocol         = "tcp"
        port_range       = "80"
        cidrs = ["192.168.5.0/24"]
    }
    inbound_rule {
        protocol         = "tcp"
        port_range       = "22"
        cidrs = ["192.168.3.0/24", "192.168.4.0/24"]
    }
    outbound_rule {
        protocol         = "tcp"
        port_range       = "22"
        cidrs  = ["192.168.1.0/24", "192.168.2.0/24"]
    }
}

resource "cmccloud_firewall_direct" "vm_firewall" { 
    server_id = "040e6918-409e-4346-8827-cd328e113f91" 
    ip_address = "203.171.21.11"
    inbound_rule {
        name = "allow ssh"
        protocol         = "tcp"
        port_range       = "22"
        port_type = "local"
        src  = "0.0.0.0/0"
        action = "allow"
    }
} 

resource "cmccloud_snapshot" "snapshot1" { 
    name = "test-snapshot"
    server_id = "040e6918-409e-4346-8827-cd328e113f91"
} 
