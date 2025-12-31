resource "kubernetes_service" "api_service" {
  metadata {
    name      = "api-service"
    namespace = "tc4-customer"
    annotations = {
      "service.beta.kubernetes.io/aws-load-balancer-type" = "nlb"
    }
  }
  spec {
    selector = {
      app = "tc4-customer-api"
    }
    port {
      protocol    = "TCP"
      port        = 80
      target_port = 8080
    }
    type = "LoadBalancer"
  }
  depends_on = [kubernetes_namespace.lanchonete_ns]
}