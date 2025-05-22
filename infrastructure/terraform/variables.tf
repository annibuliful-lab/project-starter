variable "api_key" {
  type = string
}

variable "region" {
  type    = string
  default = "sgp"
}

variable "k8s" {
  type = object({
    label           = string,
    node_quantity   = number
    node_pool_label = string,
    node_pool_plan  = string
    version         = string
  })

  default = {
    label           = "starter-labs-k8s"
    node_pool_label = "starter-labs-k8s-node-pool"
    node_pool_plan  = "vc2-2c-4gb"
    node_quantity   = 1
    version         = "v1.33.0+1" // https://api.vultr.com/v2/kubernetes/versions
  }
}

variable "s3" {
  type = object({
    cluster_id = number,
    label      = string
  })


  default = {
    cluster_id = 4 // singapore https://api.vultr.com/v2/object-storage/clusters
    label      = "livid-file"
  }
}
