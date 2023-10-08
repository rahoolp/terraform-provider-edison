resource "edison_ehscluster" "tenant-ehs" {
  region = "us-east-1"
  profile = "medium"
  release = "fenwood"
  tag = "uwm"
}

output "tenant-ehs" {
  value = edison_ehscluster.tenant-ehs
}