variable "partition_space_tb" {
  description = "Size of EA tenant for DICOM data"
  type        = string
}

resource "edison_eastore" "tenant-ea" {
  partition_space_tb = var.partition_space_tb
}

output "tenant-ea" {
  value = edison_eastore.tenant-ea
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

variable "account_id" {
  description = "Customer AWS account id"
  type        = string
}

variable "tenant_id" {
  description = "Edison Customer unique id"
  type        = string
}

resource "edison_av" "tenant-av" {
  account_id = var.account_id
  tenant_id = var.tenant_id
}

variable "concurrent_users" {
  description = "Concurrent Users using AW"
  type        = string
}

resource "edison_aw" "tenant-aw" {
  concurrent_users = var.concurrent_users
  depends_on = [edison_eastore.tenant-ea, edison_ehscluster.tenant-ehs, edison_av.tenant-av]
  ehs_cluster_id = edison_ehscluster.tenant-ehs.id
  dicom_endpoint = join("@", [join(":", [edison_eastore.tenant-ea.ip_address, edison_eastore.tenant-ea.ip_port]), edison_eastore.tenant-ea.aet])
  ea_account_id = edison_eastore.tenant-ea.account_id 
  ea_service_ep = edison_eastore.tenant-ea.service_ep
}

# resource "edison_eastore" "vpcep_to_scpep" {
#   depends_on = [edison_aw.tenant-aw]
#   ea_account_id = edison_eastore.tenant-ea.account_id 
#   ea_service_ep = edison_eastore.tenant-ea.service_ep
# }


output "tenant-aw" {
  value = edison_aw.tenant-aw
}


output "tenant-av" {
  value = edison_av.tenant-av
}

