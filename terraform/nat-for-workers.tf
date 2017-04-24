variable "latest_ubuntu" {
  type = "string"
  default = "ubuntu-1404-trusty-v20161109"
}

provider "google" {
  project = "${var.projectid}"
  region = "${var.region}"
}

variable "projectid" { }

variable "region" { }

variable "zone" { }

variable "prefix" { }

variable "network" { }

variable "subnetwork" { }

variable "tcp_router_ip" { }

resource "google_compute_route" "nat-for-worker" {
  name        = "${var.prefix}nat-for-worker"
  dest_range  = "${var.tcp_router_ip}/32"
  network       = "${var.network}"
  next_hop_instance = "${google_compute_instance.nat-instance-private-with-nat-primary.name}"
  next_hop_instance_zone = "${var.zone}"
  priority    = 500
  tags = ["${var.prefix}worker-node"]
}


// NAT server (primary)
resource "google_compute_instance" "nat-instance-private-with-nat-primary" {
  name         = "${var.prefix}nat-instance-for-worker"
  machine_type = "n1-standard-1"
  zone         = "${var.zone}"

  tags = ["nat", "internal"]

  disk {
    image = "${var.latest_ubuntu}"
  }

  network_interface {
    subnetwork = "${var.subnetwork}"
    access_config {
      // Ephemeral IP
    }
  }

  can_ip_forward = true

  metadata_startup_script = <<EOT
  #!/bin/bash
  sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"
  iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
  EOT
}