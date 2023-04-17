provider "google" {
  credentials = "{\"type\":\"service_account\"}"
  region      = "us-central1"
}

# Add example resources for DataprocCluster below

# resource "google_dataproc_cluster" "dataproc_cluster" {
# }

resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
}

resource "google_dataproc_cluster" "dataproc_cluster" {
  name     = "mycluster"
  region   = "us-central1"
  graceful_decommission_timeout = "120s"
  labels = {
    foo = "bar"
  }

  cluster_config {
    staging_bucket = "dataproc-staging-bucket"

    master_config {
      num_instances = 1
      machine_type  = "n1-standard-2"
      disk_config {
        boot_disk_type    = "pd-standard"
        boot_disk_size_gb = 700
      }
    }

    worker_config {
      num_instances    = 2
      machine_type     = "n1-standard-4"
      min_cpu_platform = "Intel Skylake"
      disk_config {
        boot_disk_size_gb = 200
        num_local_ssds    = 1
      }
    }

    preemptible_worker_config {
      num_instances = 1
      disk_config {
        boot_disk_size_gb = 300
      }
    }

    # Override or set some custom properties
    software_config {
      image_version = "2.0.35-debian10"
      override_properties = {
        "dataproc:dataproc.allow.zero.workers" = "true"
      }
    }

    gce_cluster_config {
      tags = ["foo", "bar"]
      # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
      service_account = google_service_account.default.email
      service_account_scopes = [
        "cloud-platform"
      ]
    }

    # You can define multiple initialization_action blocks
    initialization_action {
      script      = "gs://dataproc-initialization-actions/stackdriver/stackdriver.sh"
      timeout_sec = 500
    }
  }
}