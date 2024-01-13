# #!/bin/bash

SERVICE_NAME=core
NAMESPACE=default

# openssl genrsa -out ./certificates/$SERVICE_NAME.key 2048
# openssl req -new -nodes -x509 -days 365 -key ./certificates/$SERVICE_NAME.key -out ./certificates/$SERVICE_NAME.pem -subj "/CN=host.docker.internal" -addext "subjectAltName = DNS:host.docker.internal"

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
      url: https://host.docker.internal:8080/v1/validate
    failurePolicy: Fail
    matchPolicy: Exact
    name: core.atop.io
    namespaceSelector: {}
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
        - key: atop.io/enable
          operator: In
          values:
            - "true"
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
      url: https://host.docker.internal:8080/v1/inject
    failurePolicy: Ignore
    matchPolicy: Exact
    name: core.atop.io
    namespaceSelector: {}
    objectSelector:
      matchExpressions:
        - key: atop.io/enable
          operator: In
          values:
            - "true"
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
