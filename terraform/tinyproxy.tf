variable "latest_ubuntu" {
  type = "string"
  default = "ubuntu-1404-trusty-v20161109"
}

provider "google" {
  project = "${var.projectid}"
  region = "${var.region}"
  credentials = "${file("${path.cwd}/${var.service_account_key_path}")}"
}

variable "service_account_key_path" { }

variable "projectid" { }

variable "region" { }

variable "zone" { }

variable "prefix" { }

variable "subnetwork" { }

resource "google_compute_instance" "tinyproxy" {
  name         = "${var.prefix}-tinyproxy"
  machine_type = "n1-standard-1"
  zone         = "${var.zone}"

  tags = ["internal"]

  disk {
    image = "${var.latest_ubuntu}"
  }

  network_interface {
    subnetwork = "${var.subnetwork}"
    access_config {
      // Ephemeral IP
    }
  }


  metadata_startup_script = <<EOT
  #!/bin/bash

  apt-get update
  apt-get install -y tinyproxy
  sed -i 's#Allow 127.0.0.1#Allow 0.0.0.0/0#g' /etc/tinyproxy.conf
  service tinyproxy restart
  EOT
}

output "proxy_ip" {
  value = "${google_compute_instance.tinyproxy.network_interface.0.address}"
}
