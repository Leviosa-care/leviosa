storage "s3" {
  bucket     = "staging-leviosa-vault-storage"
  region     = "eu-central-1"
}

seal "awskms" {
  region     = "eu-central-1"
  kms_key_id = "arn:aws:kms:eu-central-1:970547355404:key/e4663dde-9acd-41d4-a31e-e206bb78bfe1"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}

api_addr = "http://0.0.0.0:8200"
cluster_addr = "http://0.0.0.0:8201"
ui = true

# Disable mlock for containerized environments
disable_mlock = true
