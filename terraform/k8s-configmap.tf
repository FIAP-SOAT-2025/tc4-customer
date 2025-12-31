resource "kubectl_manifest" "configmap" {
  depends_on = [kubernetes_namespace.lanchonete_ns]
  yaml_body  = <<YAML
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-configmap
  namespace: tc4-customer
data:
  APP_PORT: "8080"
  MONGODB_HOST: "mongodb-service.tc4-customer.svc.cluster.local"
  MONGODB_PORT: "27017"

YAML
}