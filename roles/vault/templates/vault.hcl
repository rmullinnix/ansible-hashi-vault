backend "consul" {
  address = "127.0.0.1:{{ consul_port }}"
  path = "vault"
}

listener "tcp" {
 address = "127.0.0.1:{{ vault_port }}"
 tls_disable = {{ vault_tls_disable }} 
}

disable_mlock = {{ vault_mlock_disable }}
