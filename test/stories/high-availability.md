# High availability (federated keystore)

Stories for running Stoke with multiple replicas, shared database, and federated public keys so that tokens issued by any replica are valid everywhere.

## Context

- With `cluster.enabled: true`, each replica has its own in-memory signing keys and merges its public keys with all peers. `GET /api/pkeys` returns the merged JWKS. Tokens issued by one replica (signed with that replica’s key) must be verifiable by any replica and by resource servers using the federated JWKS.
- See [docs/high-availability.md](../../docs/high-availability.md) for configuration and behaviour.

## Stories and test coverage

### US-013: With HA enabled, each replica serves a merged JWKS at /api/pkeys

**Acceptance criteria:** When cluster is enabled and static peers are configured, `GET /api/pkeys` returns 200 and a JWK set whose `keys` array includes this replica’s public key(s) and the public keys fetched from each peer (merged, deduplicated by key id).

**Tested by:**

| Test | Location | How |
|------|----------|-----|
| Go | [internal/key/federated_issuer_test.go](../../internal/key/federated_issuer_test.go) | `TestFederatedTokenIssuer_PublicKeys_ReturnsMergedJWKS`: mock inner + one peer URL; assert merged JWKS has 2 keys |
| Go | [internal/cluster/federated_test.go](../../internal/cluster/federated_test.go) | `TestMergeJWKS_LocalOnly`, `TestMergeJWKS_OnePeer`, `TestMergeJWKS_DedupByKeyId`: merge behaviour and dedup |
| E2E | [test/e2e/specs/ha.spec.ts](../e2e/specs/ha.spec.ts) | `GET /api/pkeys returns 200 and valid JWKS` @US-013 |

---

### US-014: A token issued by one replica is accepted when verified with the federated JWKS

**Acceptance criteria:** A JWT signed by replica A’s key is successfully verified when the verifier uses the merged JWKS (from any replica’s `/api/pkeys` or from the same federated issuer). Resource servers and other replicas can therefore accept tokens issued by any replica.

**Tested by:**

| Test | Location | How |
|------|----------|-----|
| Go | [internal/key/federated_issuer_test.go](../../internal/key/federated_issuer_test.go) | `TestFederatedTokenIssuer_ParseClaims_VerifiesTokenFromPeer`: token signed by peer key B is verified by federated issuer with merged A+B keys |

---

## Test task: ensure HA works as intended

Run the following to validate high-availability behaviour:

1. **Unit tests (federated merge and verification)**  
   From repo root:
   ```bash
   go test ./internal/key/... ./internal/cluster/... -v -run "FederatedTokenIssuer|MergeJWKS"
   ```
   Ensures: merged JWKS from local + peers, dedup by kid, and token signed by a peer key verifies via merged set.

2. **E2E: /api/pkeys returns valid JWKS**  
   Start a Stoke server (with or without cluster enabled), then:
   ```bash
   cd test/e2e && npx playwright test -g "@US-013"
   ```
   Or from repo root: `task test e2e` with `STOKE_BASE_URL` set.  
   Ensures: public key endpoint returns 200 and a well-formed JWKS so clients and replicas can merge and verify.

3. **Optional: two-replica integration**  
   For a full two-replica run (shared DB, cluster enabled, static_peers pointing at each other), start two Stoke processes (e.g. different ports and configs), then: login against replica A, get token; `GET /api/pkeys` from replica B and confirm the merged set includes at least two keys; verify the token from A using the JWKS from B. This is not automated in this repo; see [docs/high-availability.md](../../docs/high-availability.md).
