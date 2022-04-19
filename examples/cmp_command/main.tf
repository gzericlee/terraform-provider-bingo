provider "bingo" {
  sso_endpoint      = "##"
  cmp_access_token	= "##"
  #cmp_client_secret	= "clientSecret1"
  cmp_endpoint		= "##"
}

resource "bingo_cmp_command" "send" {
  host_type   	= "##"
  content     	= "##"
  instance_ids	= "c0dea473-cfc0-49a7-830e-a7edc8f1125d"
}