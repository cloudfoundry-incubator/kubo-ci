provider "google" {
  project = "${var.projectid}"
  region = "${var.region}"
  credentials = "${file("${path.cwd}/${var.service_account_key_path}")}"
}

variable "projectid" { }
variable "region" { }
variable "service_account_key_path" { }

variable "prefix" {
    type = "string"
    default = "cfcr"
}

variable "private_key_path" {}
variable "certificate_path" {}


resource "google_compute_instance_template" "default" {
  name_prefix  = "instance-template-"
  machine_type = "n1-standard-1"
  region       = "us-central1"

  disk {
    source_image = "debian-cloud/debian-8"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_region_instance_group_manager" "default" {
  name = "tuna-masters-instance-group-manager"

  base_instance_name         = "${var.prefix}"
  instance_template          = "${google_compute_instance_template.default.self_link}"
  region                     = "${var.region}"

  named_port {
    name = "kubernetes-master"
    port = 8443
  }
}


resource "google_compute_target_https_proxy" "default" {
  name             = "test-proxy"
  description      = "a description"
  url_map          = "${google_compute_url_map.default.self_link}"
  ssl_certificates = ["${google_compute_ssl_certificate.default.self_link}"]
}

resource "google_compute_ssl_certificate" "default" {
  name        = "kubo.sh-cert"
  description = "certificate "
  private_key = "${file(${var.private_key_path})}"
  certificate = "${file(${var.certificate_path})}"
}

resource "google_compute_url_map" "default" {
  name        = "url-map"
  description = "a description"

  default_service = "${google_compute_backend_service.default.self_link}"
}

resource "google_compute_backend_service" "default" {
  name        = "default-backend"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10

  health_checks = ["${google_compute_http_health_check.default.self_link}"]
}

resource "google_compute_http_health_check" "default" {
  name               = "test"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

