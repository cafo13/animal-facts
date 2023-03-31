variable "project" {
  type = string
}

variable "user" {
  type = string
}

variable "region" {
  type = string
}

variable "firebase_location" {
  type = string
}

variable "billing_account" {
  type        = string
  description = "Billing account display name"
}

variable "app_version" {
  type        = string
  description = "The version of animal-facts to deploy"
}
