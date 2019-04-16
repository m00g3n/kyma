apiVersion: v1
kind: Secret
metadata:
  name: application-connector-certificate-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    kyma-project.io/installation: ""
type: Opaque
data:
  global.applicationConnectorCa: "__REMOTE_ENV_CA__"
  global.applicationConnectorCaKey: "__REMOTE_ENV_CA_KEY__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-certificate-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.tlsCrt: "__TLS_CERT__"
  global.tlsKey: "__TLS_KEY__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: installation-config-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.domainName: "__DOMAIN__"
  global.loadBalancerIP: "__EXTERNAL_PUBLIC_IP__"
  global.etcdBackup.containerName: "__ETCD_BACKUP_ABS_CONTAINER_NAME__"
  global.etcdBackup.enabled: "__ENABLE_ETCD_BACKUP__"
  nginx-ingress.controller.service.loadBalancerIP: "__REMOTE_ENV_IP__"
  cluster-users.users.adminGroup: "__ADMIN_GROUP__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: monitoring-config-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    component: monitoring
    kyma-project.io/installation: ""
data:
  global.alertTools.credentials.slack.apiurl: "__SLACK_API_URL_VALUE__"
  global.alertTools.credentials.slack.channel: "__SLACK_CHANNEL_VALUE__"
  global.alertTools.credentials.victorOps.routingkey: "__VICTOR_OPS_ROUTING_KEY_VALUE__"
  global.alertTools.credentials.victorOps.apikey: "__VICTOR_OPS_API_KEY_VALUE__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: istio-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    component: istio
    kyma-project.io/installation: ""
data:
  gateways.istio-ingressgateway.loadBalancerIP: "__EXTERNAL_PUBLIC_IP__"
  global.proxy.excludeIPRanges: "__PROXY_EXCLUDE_IP_RANGES__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: knative-serving-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    component: knative-serving
    kyma-project.io/installation: ""
data:
  knative-serving.domainName: "__DOMAIN__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: intallation-logging-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    component: logging
    kyma-project.io/installation: ""
data:
  global.logging.promtail.config.name: "__PROMTAIL_CONFIG_NAME__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: core-test-ui-acceptance-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    component: core
    kyma-project.io/installation: ""
data:
  test.acceptance.ui.logging.enabled: "__LOGGING_INSTALL_ENABLED__"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: assetstore-overrides
  namespace: kyma-installer
  labels:
    installer: overrides
    component: assetstore
    kyma-project.io/installation: ""
data:
  minio.persistence.enabled: "__MINIO_PERSISTENCE_ENABLED__"
  minio.gcsgateway.enabled: "__MINIO_GCS_GATEWAY_ENABLED__"
  minio.gcsgateway.gcsKeySecret: "__MINIO_GCS_GATEWAY_GCS_KEY_SECRET__"
  minio.gcsgateway.enabled: "__MINIO_GCS_GATEWAY_PROJECT_ID__"
  minio.defaultBucket.enabled: "__MINIO_PERSISTENCE_ENABLED__"