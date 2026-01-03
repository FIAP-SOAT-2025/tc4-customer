resource "kubectl_manifest" "secrets" {
  depends_on = [kubernetes_namespace.lanchonete_ns]
  yaml_body  = <<YAML
apiVersion: v1
kind: Secret
metadata:
  name: api-secrets
  namespace: tc4-customer
type: Opaque
data:
  MONGO_USER: ${base64encode(var.mongo_user)}
  MONGO_PASSWORD: ${base64encode(var.mongo_password)}
  MONGO_DB_NAME: ${base64encode(var.mongo_db_name)}
  MONGODB_URI: ${base64encode("mongodb://${urlencode(var.mongo_user)}:${urlencode(var.mongo_password)}@mongodb-service.tc4-customer.svc.cluster.local:27017/${var.mongo_db_name}?authSource=admin")}
  MONGODB_DATABASE: ${base64encode(var.mongo_db_name)}
YAML
}