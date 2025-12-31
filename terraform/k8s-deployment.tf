resource "kubectl_manifest" "deployment" {
  depends_on = [
    kubernetes_namespace.lanchonete_ns,
    kubectl_manifest.mongodb_statefulset,
    kubectl_manifest.db_seed_job,
    kubectl_manifest.secrets,
    kubectl_manifest.configmap
  ]
  yaml_body = <<YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tc4-customer-api
  namespace: tc4-customer
spec:
  replicas: 1
  selector:
    matchLabels:
      app: tc4-customer-api
  template:
    metadata:
      labels:
        app: tc4-customer-api
    spec:
      containers:
      - name: tc4-customer-api
        image: tlnob/fiap-tc4:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 15
          periodSeconds: 5
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        envFrom:
        - configMapRef:
            name: api-configmap
        - secretRef:
            name: api-secrets
        resources:
          requests:
            memory: "256Mi"
            cpu: "200m"
          limits:
            memory: "512Mi"
            cpu: "500m"

YAML
}