resource "edison_eastore" "tenant-ea" {
  partition_space_tb = 15
}

output "tenant-ea" {
  value = edison_eastore.tenant-ea
}