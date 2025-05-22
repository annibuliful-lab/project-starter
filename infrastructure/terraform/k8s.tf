resource "vultr_kubernetes" "k8s" {
  region  = var.region
  label   = var.k8s.label
  version = var.k8s.version

  node_pools {
    node_quantity = var.k8s.node_quantity
    plan          = var.k8s.node_pool_plan
    label         = var.k8s.node_pool_label
  }
}
