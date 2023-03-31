variable "project" {
  type = string
}

variable "name" {
  type = string
}

variable "image" {
  type = string
}

variable "location" {
  type = string
}

variable "envs" {
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}
