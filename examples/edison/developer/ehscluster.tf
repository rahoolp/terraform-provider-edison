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

resource "edison_ehscluster" "uwm-ehs" {
  region = var.region
  profile = var.profile
  release = var.release
  tag = var.tag
  depends_on = [edison_eastore.uwm-ea]
  dicom_endpoint = join("@", [join(":", [edison_eastore.uwm-ea.ip_address, edison_eastore.uwm-ea.ip_port]), edison_eastore.uwm-ea.aet])
}

output "uwm-ehs" {
  value = edison_ehscluster.uwm-ehs
}