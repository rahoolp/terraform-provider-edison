resource "edison_eastore" "uwm-ea" {
  partition_space_tb = 12
}

output "uwm-ea" {
  value = edison_eastore.uwm-ea
}