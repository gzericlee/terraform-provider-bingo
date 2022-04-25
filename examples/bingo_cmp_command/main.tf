terraform {
  required_providers {
    bingo = {
      source = "gzericlee/bingo"
    }
  }
}

provider "bingo" {

}

resource "bingo_cmp_command" "cmd" {
  host_type   	= "1"
  content     	= "pwd"
  instance_ids	= "##"
}

output "TASK_ID" {
  description = "任务编号"
  value = bingo_cmp_command.cmd.task_id
}