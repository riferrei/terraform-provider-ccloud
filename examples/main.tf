variable "ccloud_username" {
  type = string
}

variable "ccloud_password" {
  type = string
}

provider "ccloud" {
  username = var.ccloud_username
  password = var.ccloud_password
}
