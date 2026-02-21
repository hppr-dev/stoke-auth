/**
 * Fixture for generated-dbinit E2E test (US-012).
 * Generate password_hash and password_salt by running:
 *   echo 'e2edbinit' | go run ./cmd/ hash-password
 * (or use the built stoke-server binary) and paste the output below.
 */
export const E2E_DBINIT_PASSWORD = 'e2edbinit';
export const E2E_DBINIT_USERNAME = 'e2edbinit-user';

// Replace these with output from: echo 'e2edbinit' | go run ./cmd/ hash-password
let _passwordHash = '';
let _passwordSalt = '';

export function setGeneratedDbinitCredentials(hash: string, salt: string) {
  _passwordHash = hash;
  _passwordSalt = salt;
}

export function getGeneratedDbinitYaml(): string {
  if (!_passwordHash || !_passwordSalt) {
    throw new Error(
      'Generated dbinit fixture: run "echo e2edbinit | go run ./cmd/ hash-password" and call setGeneratedDbinitCredentials(hash, salt) before tests, or use precomputed values.'
    );
  }
  return `
users:
  - username: ${E2E_DBINIT_USERNAME}
    first_name: E2E
    last_name: Dbinit
    email: e2edbinit@test.example
    password_hash: ${_passwordHash}
    password_salt: ${_passwordSalt}
    groups:
      - E2E Dbinit Group

groups:
  - name: E2E Dbinit Group
    description: Group for generated dbinit E2E test
    claims:
      - E2E Dbinit Claim

claims:
  - name: E2E Dbinit Claim
    description: Claim for generated dbinit E2E test
    short_name: e2edb
    value: acc
`.trim();
}
