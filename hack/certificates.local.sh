# #!/bin/bash

SERVICE_NAME=core
NAMESPACE=default


# IP=$(ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}')
# echo IP:$IP

# docker exec -it k3d-dev-serverlb sh -c "echo '${ip} ${SERVICE_NAME}.${NAMESPACE}.svc' >> /etc/hosts"
# docker exec -it k3d-dev-agent-0 sh -c "echo '${ip} ${SERVICE_NAME}.${NAMESPACE}.svc' >> /etc/hosts"
# docker exec -it k3d-dev-server-0 sh -c "echo '${ip} ${SERVICE_NAME}.${NAMESPACE}.svc' >> /etc/hosts"

# openssl genrsa -out ./certificates/$SERVICE_NAME.key 2048
# openssl req -new -nodes -x509 -days 365 -key ./certificates/$SERVICE_NAME.key -out ./certificates/$SERVICE_NAME.pem -subj "/CN=${SERVICE_NAME}.${NAMESPACE}.svc" -addext "subjectAltName = DNS:${SERVICE_NAME}.${NAMESPACE}.svc"

CA_BUNDLE=$(cat ./certificates/$SERVICE_NAME.pem | base64 | tr -d "\n")
echo CA_BUNDLE:$CA_BUNDLE

cat <<EOF | kubectl apply -f -
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: $SERVICE_NAME
webhooks:
  - admissionReviewVersions:
      - v1beta1
      - v1
    clientConfig:
      caBundle: ${CA_BUNDLE}
      url: https://${SERVICE_NAME}.${NAMESPACE}.svc:8080/v1/validate
    failurePolicy: Fail
    matchPolicy: Exact
    name: core.atop.io
    rules:
      - apiGroups:
          - ""
        apiVersions:
          - v1
        operations:
          - CREATE
          - UPDATE
        resources:
          - pods
          - deployments
        scope: "*"
    objectSelector:
      matchExpressions:
        - key: atop.io/sidecar
          operator: In
          values:
            - enable
    sideEffects: None
    timeoutSeconds: 30

EOF

cat <<EOF | kubectl apply -f -
apiVersion: admissionregistration.k8s.io/v1
kind: MutatingWebhookConfiguration
metadata:
  name: $SERVICE_NAME
webhooks:
  - admissionReviewVersions:
      - v1beta1
      - v1
    clientConfig:
      caBundle: ${CA_BUNDLE}
      url: https://${SERVICE_NAME}.${NAMESPACE}.svc:8080/v1/inject
    failurePolicy: Ignore
    matchPolicy: Exact
    name: core.atop.io
    namespaceSelector: {}
    objectSelector:
      matchExpressions:
        - key: atop.io/sidecar
          operator: In
          values:
            - enable
    reinvocationPolicy: Never
    rules:
      - apiGroups:
          - ""
        apiVersions:
          - v1
        operations:
          - CREATE
          - UPDATE
        resources:
          - pods
          - deployments
        scope: "*"
    sideEffects: None
    timeoutSeconds: 30
EOF
