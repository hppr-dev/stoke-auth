# High availability

Stoke can run with multiple replicas behind a load balancer so that no single process is a bottleneck or single point of failure. This requires a **shared database** and optional **federated keystore** so that tokens issued by any replica are valid everywhere.

## Requirements

- **Shared database:** Use Postgres or MySQL for the Stoke database. All replicas must connect to the same database for users, groups, and claims. SQLite is not suitable for multi-replica deployments.
- **No key persistence in HA:** When high availability is enabled, signing keys are not stored in the database. Each replica keeps its own keys in memory. Configure the cluster so that each replica’s public keys are merged and exposed to clients.

## Enabling HA and federated keys

1. **Configure the cluster** in your Stoke config (e.g. `config.yaml` or the Helm-generated ConfigMap):

   ```yaml
   cluster:
     enabled: true
     discovery: static
     static_peers:
       - https://stoke-0:8080
       - https://stoke-1:8080
     refresh_sec: 30   # optional; default 30
     instance_id: "stoke-0"   # optional; unique per replica so key ids (kid) stay distinct in merged JWKS
   ```

   - `static_peers` must list the base URLs of **all** replicas (including this one if you want this instance to merge its own keys with peers). Use the URL that other replicas and clients use to reach each instance (e.g. service URL or ingress).
   - Each replica’s `/api/pkeys` will return a **merged** JWKS (this instance’s keys plus keys fetched from each peer). Tokens issued by any replica can then be verified by any replica and by resource servers that use `/api/pkeys`.

2. **Replicas and join/leave:** You can scale replicas up or down. Update `static_peers` when you add or remove replicas so that the list matches the current set. Discovery is refreshed periodically (`refresh_sec`); after a restart or config reload, the new list is used.

3. **Helm:** Set `server.replicaCount` to the desired number of replicas and configure `cluster` in the chart values (or via the generated config) as above. See [helm/README.md](../helm/README.md#high-availability).

## Behaviour

- **Issuance:** Any replica can issue tokens (login, renew). Tokens are signed with that replica’s in-memory key; the `kid` in the token identifies the key.
- **Verification:** Each replica merges its own public keys with those fetched from every peer. That merged set is served at `GET /api/pkeys` and used for token verification (e.g. middleware and token handlers). So a token issued by replica A is valid when verified by replica B or by a resource server that uses the federated JWKS.
- **Database:** All replicas read and write the same users, groups, and claims. Key storage is not used when `cluster.enabled` is true.

## Configuration reference

| Field           | Description |
|----------------|-------------|
| `cluster.enabled` | Set to `true` to enable HA: key persistence is disabled and `/api/pkeys` returns merged JWKS. |
| `cluster.discovery` | Discovery mechanism. Use `static` (default); `k8s` may be supported later. |
| `cluster.static_peers` | List of peer base URLs (e.g. `https://host:8080`) for merging keys. |
| `cluster.refresh_sec` | Seconds between refreshing the merged key set from peers; default 30. |
| `cluster.instance_id` | Optional unique id for this replica (e.g. `stoke-0`, `stoke1`). When set, signing key ids are prefixed (e.g. `stoke-0-p-0`) so merged JWKS from multiple replicas keeps all keys distinct. |

See the main [Configuration](../README.md#configuration) section and [values.yaml](../helm/values.yaml) for how to supply this in your deployment.
