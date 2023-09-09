variable "partition_space_tb" {
  description = "Size of EA tenant for DICOM data"
  type        = string
}

resource "edison_eastore" "uwm-ea" {
  partition_space_tb = var.partition_space_tb
}

output "uwm-ea" {
  value = edison_eastore.uwm-ea
}