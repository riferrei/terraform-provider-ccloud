data "ccloud_environment" "default_env" {
  name = "default"
}

resource "ccloud_cluster" "new_cluster" {
  environment_id = data.ccloud_environment.default_env.id
  name = "cluster-created-by-terraform"
  cloud_provider = "aws"
  cloud_region = "us-east-1"
}

resource "ccloud_apikey" "api_key" {
  environment_id = data.ccloud_environment.default_env.id
  cluster_id = ccloud_cluster.new_cluster.id
}

output "_bootstrap_server" {
  value = ccloud_cluster.new_cluster.cluster_endpoint
}

output "api_key" {
  value = ccloud_apikey.api_key.key
}

output "secret" {
  value = ccloud_apikey.api_key.secret
}
