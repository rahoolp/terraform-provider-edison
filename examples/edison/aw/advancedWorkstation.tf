resource "edison_eastore" "uwm-ea" {
  partition_space_tb = 10
}

resource "edison_ehscluster" "uwm-ehs" {
  region = "us-east-1"
  profile = "medium"
  release = "fenwood"
  tag = "uwm"
}

resource "edison_aw" "uwm-aw" {
  concurrent_users = 4
  depends_on = [edison_eastore.uwm-ea, edison_ehscluster.uwm-ehs]
  ehs_cluster_id = edison_ehscluster.uwm-ehs.id
  dicom_endpoint = join("@", [join(":", [edison_eastore.uwm-ea.ip_address, edison_eastore.uwm-ea.ip_port]), edison_eastore.uwm-ea.aet])
}

output "uwm-aw" {
  value = edison_aw.uwm-aw
}




# output "uwm-ea" {
#   value = edison_eastore.uwm-ea
# }



# output "uwm-ehs" {
#   value = edison_ehscluster.uwm-ehs
# }