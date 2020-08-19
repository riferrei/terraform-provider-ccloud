data "ccloud_environment" "default_env" {
  name = "default"
}

resource "ccloud_cluster" "new_cluster" {
  environment_id = data.ccloud_environment.default_env.id
  name = "cluster-created-by-terraform"
  cloud_provider = "aws"
  cloud_region = "us-east-1"
}

output "bootstrap_server" {
  value = ccloud_cluster.new_cluster.cluster_endpoint
}
