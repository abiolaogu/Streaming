storage "file" {
  path = "/vault/data"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_cert_file = "/vault/tls/vault.crt"
  tls_key_file  = "/vault/tls/vault.key"
}

api_addr = "https://vault.streamverse.io:8200"
cluster_addr = "https://vault.streamverse.io:8201"

ui = true

disable_mlock = true

seal "awskms" {
  region     = "us-east-1"
  kms_key_id = "arn:aws:kms:us-east-1:ACCOUNT_ID:key/KEY_ID"
}

