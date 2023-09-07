resource "edison_ehscluster" "uwm-ehs" {
  region = "us-east-1"
  profile = "medium"
  release = "fenwood"
  tag = "uwm"
  dicom_endpoint = "test"
}

output "uwm-ehs" {
  value = edison_ehscluster.uwm-ehs
}