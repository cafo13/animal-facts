module "cloud_run_animal_facts_facts_api" {
  source = "./cloud_run"

  project  = var.project
  location = var.region

  name  = "facts-api"
  image = "europe-west4-docker.pkg.dev/animal-facts-project/animal-facts/facts-api:${var.app_version}"

  envs = [
    {
      name  = "API_PORT"
      value = 80
    },
    {
      name  = "GIN_MODE"
      value = "release"
    }
  ]
}
