terraform {
  required_providers {
    bingo = {
      source = "gzericlee/bingo"
    }
  }
}

provider "bingo" {
  ssoEndpoint      = "https://sso.bingosoft.net"
  cmpAccessToken	= "YmluZ286YmluZ29fbWVtYmVyOjZtMnNIMEhSTg"
  #cmp_client_secret	= "clientSecret1"
  cmpEndpoint		= "https://cmp-dev.bingosoft.net"
}

resource "bingo_cmp_command" "cmd" {
  host_type   	= "1"
  content     	= "pwd"
  instance_ids	= "c0dea473-cfc0-49a7-830e-a7edc8f1125d"
}

output "TASK_ID" {
  description = "任务编号"
  value = bingo_cmp_command.cmd.task_id
}