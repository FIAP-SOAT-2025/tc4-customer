resource "kubectl_manifest" "db_seed_job" {
  depends_on = [
    kubectl_manifest.mongodb_deployment,
    kubectl_manifest.secrets,
    kubectl_manifest.configmap
  ]
  yaml_body = <<YAML
apiVersion: batch/v1
kind: Job
metadata:
  name: db-seed-job
  namespace: tc4-customer
spec:
  template:
    spec:
      initContainers:
      - name: wait-for-mongodb
        image: busybox:1.36
        command: ['sh', '-c', 'until nc -z mongodb-service.tc4-customer.svc.cluster.local 27017; do echo waiting for mongodb; sleep 2; done;']
      containers:
      - name: tc4-customer-seed-db
        image: tlnob/tc4-customer:latest
        imagePullPolicy: IfNotPresent
        command: ["./customer-service", "seed"]
        envFrom:
        - configMapRef:
            name: api-configmap
        - secretRef:
            name: api-secrets
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
      restartPolicy: Never
  backoffLimit: 4

YAML
}