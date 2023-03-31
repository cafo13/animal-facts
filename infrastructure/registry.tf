resource "google_artifact_registry_repository" "repository" {
  location      = "europe-west4"
  repository_id = "animal-facts"
  description   = "Storing images of animal facts services"
  format        = "DOCKER"
}

resource "google_service_account" "registry" {
  account_id   = "registry"
  display_name = "Registry Service Account"
  depends_on = [
    google_project_iam_member.owner,
    google_project_service.registry,
  ]
}

resource "google_iam_workload_identity_pool" "pool" {
  workload_identity_pool_id = "github-action"
  display_name              = "github-action"
}

resource "google_iam_workload_identity_pool_provider" "pool_oidc" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-action"
  display_name                       = "github-action"

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_project_iam_member" "registry" {
  project = google_project.project.id
  role    = "roles/iam.workloadIdentityUser"
  member  = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.pool.name}/attribute.repository/cafo13/animal-facts"
}

resource "google_project_iam_member" "registry_admin" {
  project = google_project.project.id
  role    = "roles/artifactregistry.admin"
  member  = google_service_account.registry.member
}
