helm repo add openfga https://openfga.github.io/helm-charts
helm install openfga openfga/openfga \
  --set datastore.engine=postgres \
  --set datastore.uri="postgres://postgres:password@openfga-postgresql.default.svc.cluster.local:5432/postgres?sslmode=disable" \
  --set postgres.enabled=true \
  --set postgresql.auth.postgresPassword=password \
  --set postgresql.auth.database=postgres