output "artifacts_registry_workload_identity_pool_provider" {
  value = google_iam_workload_identity_pool_provider.pool_oidc.name
}

output "artifacts_registry_service_account" {
  value = google_service_account.registry.member
}
