resource "kubectl_manifest" "mongodb_deployment" {
  depends_on = [
    kubernetes_namespace.lanchonete_ns,
    kubectl_manifest.mongodb_service,
    kubectl_manifest.secrets
  ]

  override_namespace = "tc4-customer"
  wait               = true
  wait_for_rollout   = true

  timeouts {
    create = "10m"
  }

  yaml_body = <<YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb
  namespace: tc4-customer
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mongodb
  template:
    metadata:
      labels:
        app: mongodb
    spec:
      containers:
      - name: mongodb
        image: mongo:7.0
        ports:
        - containerPort: 27017
          name: mongodb
        env:
        - name: MONGO_INITDB_ROOT_USERNAME
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: MONGO_USER
        - name: MONGO_INITDB_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: MONGO_PASSWORD
        - name: MONGO_INITDB_DATABASE
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: MONGO_DB_NAME
        volumeMounts:
        - name: mongodb-data
          mountPath: /data/db
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        livenessProbe:
          exec:
            command:
            - mongosh
            - --eval
            - "db.adminCommand('ping')"
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 10
        readinessProbe:
          exec:
            command:
            - mongosh
            - --eval
            - "db.adminCommand('ping')"
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 10
      volumes:
      - name: mongodb-data
        emptyDir: {}

YAML
}

resource "kubectl_manifest" "mongodb_service" {
  depends_on = [kubernetes_namespace.lanchonete_ns]
  yaml_body  = <<YAML
apiVersion: v1
kind: Service
metadata:
  name: mongodb-service
  namespace: tc4-customer
spec:
  selector:
    app: mongodb
  ports:
  - protocol: TCP
    port: 27017
    targetPort: 27017
  clusterIP: None

YAML
}
