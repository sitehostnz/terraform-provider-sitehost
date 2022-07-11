variable "sitehost_client_id" {
  description = "SiteHost API v1.1 Client ID"
}

variable "sitehost_api_key" {
  description = "SiteHost API v1.1 API Key"
}

variable "location" {
  default = "AKLCITY"
}

variable "product_code" {
  default = "XENLIT"
}

variable "image" {
  default = "ubuntu-jammy-pvh.amd64"
}

variable "public_ssh_key" {
  description = "SSH Public Key Fingerprint"
  default     = "~/.ssh/id_rsa.pub"
}

variable "private_ssh_key" {
  description = "SSH Private Key Fingerprint"
  default     = "~/.ssh/id_rsa"
}