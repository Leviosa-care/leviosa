resource "hcloud_server" "manager" {
  name               = "ubuntu-2gb-nbg1-2"
  image              = "ubuntu-24.04"
  server_type        = "cpx11"
  placement_group_id = 480516
}

resource "aws_s3_bucket" "bucket" {
  bucket = "staging-leviosa-assets"
}
