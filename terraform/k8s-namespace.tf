resource "kubernetes_namespace" "lanchonete_ns" {
  metadata {
    name = "tc4-customer"
  }
}