Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

This is an unofficial (therefore, unsupported) implementation of [Confluent Cloud](https://www.confluent.io/confluent-cloud) Terraform provider. It is intended to aid developers in the creation of immutable software architectures relying on Confluent Cloud for everything related to [Apache Kafka](https://kafka.apache.org).

**Important**: this project has no relationship with Confluent and will not be supported but any means. Use this code at your own risk and treat it as if it was your own code. Eventually, Confluent will release a official Terraform provider for Confluent Cloud, which will be registered within Terraform Registry.

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Building the Provider
----------------------

In order to use the provider you first have to build it. Then you must install the native executable generated either in the same folder where your .tf files reside or installing it as a plugin as explained [here](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin).

```console
make install
```

Once you have done this, you can run `terraform init` to initialize your project.

Examples
----------------------

Creating a Kafka cluster in an existing environment:

```
provider "ccloud" {
  username = "<YOUR_CCLOUD_USERNAME>"
  password = "<YOUR_CCLOUD_PASSWORD>"
}

data "ccloud_environment" "existing_env" {
  name = "your-env-name"
}

resource "ccloud_cluster" "new_cluster" {
  environment_id = data.ccloud_environment.existing_env.id
  name = "new-cluster"
  cloud_provider = "azure"
  cloud_region = "westus2"
}

output "bootstrap_server" {
  value = ccloud_cluster.new_cluster.cluster_endpoint
}
```

Creating an environment and then creating two Kafka clusters on it:

```
provider "ccloud" {
  username = "<YOUR_CCLOUD_USERNAME>"
  password = "<YOUR_CCLOUD_PASSWORD>"
}

resource "ccloud_environment" "new_env" {
  name = "new-env-name"
}

resource "ccloud_cluster" "cluster_1" {
  environment_id = data.ccloud_environment.new_env.id
  name = "cluster-1"
  cloud_provider = "azure"
  cloud_region = "westus2"
}

resource "ccloud_cluster" "cluster_2" {
  environment_id = data.ccloud_environment.new_env.id
  name = "cluster-2"
  cloud_provider = "azure"
  cloud_region = "westus2"
}
```

Creating an environment, a Kafka cluster, and an API Key for usage:

```
provider "ccloud" {
  username = "<YOUR_CCLOUD_USERNAME>"
  password = "<YOUR_CCLOUD_PASSWORD>"
}

resource "ccloud_environment" "new_env" {
  name = "new-env-name"
}

resource "ccloud_cluster" "new_cluster" {
  environment_id = data.ccloud_environment.new_env.id
  name = "new-cluster"
  cloud_provider = "azure"
  cloud_region = "westus2"
}

resource "ccloud_apikey" "new_apikey" {
  environment_id = ccloud_environment.new_env.id
  cluster_id = ccloud_cluster.new_cluster.id
}

output "bootstrap_server" {
  value = ccloud_cluster.new_cluster.cluster_endpoint
}

output "api_key" {
  value = ccloud_apikey.new_apikey.key
}

output "api_secret" {
  value = ccloud_apikey.new_apikey.secret
}
```