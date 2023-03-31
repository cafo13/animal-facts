locals {
  default_envs = [
    {
      name  = "GCP_PROJECT"
      value = var.project
    }
  ]
}

resource "google_cloud_run_service" "service" {
  name     = var.name
  location = var.location

  template {
    spec {
      containers {
        image = var.image
        ports {
          container_port = 80
        }

        dynamic "env" {
          for_each = concat(local.default_envs, var.envs)
          content {
            name  = env.value.name
            value = env.value.value
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "3"
      }
    }
  }

  autogenerate_revision_name = true
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth_policy" {
  location = google_cloud_run_service.service.location
  service  = google_cloud_run_service.service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
