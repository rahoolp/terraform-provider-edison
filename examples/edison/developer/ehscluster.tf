variable "region" {
  description = "EHS cluster AWS region"
  type        = string
}
variable "profile" {
  description = "EHS cluster profile"
  type        = string
}
variable "release" {
  description = "EHS cluster release"
  type        = string
}
variable "tag" {
  description = "EHS cluster tag value for all resources"
  type        = string
}

resource "edison_ehscluster" "tenant-ehs" {
  region = var.region
  profile = var.profile
  release = var.release
  tag = var.tag
  depends_on = [edison_eastore.tenant-ea]
  dicom_endpoint = join("@", [join(":", [edison_eastore.tenant-ea.ip_address, edison_eastore.tenant-ea.ip_port]), edison_eastore.tenant-ea.aet])
}

output "tenant-ehs" {
  value = edison_ehscluster.tenant-ehs
}