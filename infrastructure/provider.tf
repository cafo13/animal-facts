provider "google" {
  project = var.project
  region  = var.region
}

provider "google-beta" {
  project     = var.project
  region      = var.region
  credentials = base64decode(google_service_account_key.firebase_key.private_key)
}
