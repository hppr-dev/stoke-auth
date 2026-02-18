# Stoke Auth Helm Chart

This chart deploys the Stoke auth server: a single Deployment, optional PostgreSQL or MySQL subcharts for the database, and ConfigMap-based configuration.

## Prerequisites

* Helm 3
* A Kubernetes cluster
* Optional: a container registry and image pull secret if using a custom server image

## Install and upgrade

From the `helm/` directory:

* Install: `helm install <release> .`
* Upgrade: `helm upgrade <release> .`
* List releases: `helm list`
* Uninstall: `helm uninstall <release>`

## Values summary

The main [values.yaml](values.yaml) sections are:

* **server** — Image, pull policy, port, timeout, admin UI (enabled, allowedHosts), TLS (privateKey, publicCert).
* **postgresql** / **mysql** — Subchart configuration. Enable one with `enabled: true`; auth (database, username, password) and service account are configurable.
* **logging** — level, pretty output.
* **tokens** — algorithm, numBits, persistKeys, tokenRefreshLimit; expiration (key, token); claims (issuer, subject, audience, etc.); userInfo (claim key names).
* **users** — createStokeClaims; policy (allowSuperuserOverride, readOnlyMode, protectedUsers, protectedClaims, protectedGroups); **providers** (list of LDAP and/or OIDC provider configs).
* **dbInit** — Optional. When set, the chart generates a database init snippet and mounts it at `/etc/stoke/dbinit.d/` so the server applies it on start. Structure: users, groups, claims under dbInit.

The server config is rendered into a single ConfigMap (see [templates/server-config-map.yaml](templates/server-config-map.yaml)). Providers are configured inline via the `providers` value; there is no separate mount for a providers.d-style directory.

When dbInit is provided, the deployment uses a projected volume so both the main config and the generated dbinit file are available under `/etc/stoke/` (see [templates/stoke-server-deployment.yaml](templates/stoke-server-deployment.yaml)).
