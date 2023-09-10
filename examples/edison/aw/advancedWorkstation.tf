variable "partition_space_tb" {
  description = "Size of EA tenant for DICOM data"
  type        = string
}

resource "edison_eastore" "tenant-ea" {
  partition_space_tb = var.partition_space_tb
}


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
}

output "tenant-ehs" {
  value = edison_ehscluster.tenant-ehs
}

variable "concurrent_users" {
  description = "Concurrent Users using AW"
  type        = string
}

resource "edison_aw" "tenant-aw" {
  concurrent_users = var.concurrent_users
  depends_on = [edison_eastore.tenant-ea, edison_ehscluster.tenant-ehs]
  ehs_cluster_id = edison_ehscluster.tenant-ehs.id
  dicom_endpoint = join("@", [join(":", [edison_eastore.tenant-ea.ip_address, edison_eastore.tenant-ea.ip_port]), edison_eastore.tenant-ea.aet])
}

output "tenant-aw" {
  value = edison_aw.tenant-aw
}


