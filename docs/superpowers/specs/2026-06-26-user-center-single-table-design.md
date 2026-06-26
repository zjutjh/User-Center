# User Center Single-Table Design

## Goal

Continue the user-center migration on go-zero with a single user table. The API service remains the public HTTP boundary for login and profile operations. The RPC service remains the internal user domain service and provides password lookup for trusted internal callers.

## Decisions

- Keep the single `user` table for core user data and third-party platform secrets.
- Keep `mini_program_user` as a separate open-id mapping table because it models quick-login identity mapping rather than profile data.
- Return usable plaintext third-party passwords from the internal RPC `GetUserPassword` method.
- Store user-center login passwords as one-way hashes.
- Store third-party passwords with reversible encryption so RPC can decrypt them for internal services.
- Do not implement audit logging in this iteration.
- Preserve go-zero scaffold boundaries: API/RPC declarations are changed in `user.api` and `user.proto`, then regenerated with goctl.

## API Layer

`apps/user-api` is responsible for HTTP request parsing, public response formatting, session cookie handling, and RPC calls. It must not access the database directly and must not implement third-party password encryption or decryption.

Public API capabilities:

- Register a student account.
- Login by user-center password.
- Login through unified OAuth credentials.
- Login through mini-program open-id mapping.
- Read the authenticated user's profile.
- Reset password.
- Delete account.

## RPC Layer

`apps/user-rpc` owns user-domain behavior:

- Register users.
- Validate user-center password login.
- Validate third-party login credentials where required.
- Bind or update third-party platform credentials.
- Return user profile data.
- Return internal plaintext third-party passwords through `GetUserPassword`.

The `GetUserPassword` method is an internal RPC capability. It returns fields such as `student_id`, `device_id`, `yxy_uid`, `zf_password`, and `oauth_password`. Empty strings indicate credentials that are not bound.

## Data Model

The existing single-table shape remains:

- `user.student_id`: canonical student identifier.
- `user.password`: user-center login password hash.
- `user.phone_num`: phone number from Yi Xiao Yuan binding.
- `user.user_type`: user category.
- `user.email`: user email.
- `user.device_id`: Yi Xiao Yuan device id.
- `user.yxy_uid`: Yi Xiao Yuan uid.
- `user.zf_password`: encrypted Zheng Fang password.
- `user.oauth_password`: encrypted unified-auth password.

The table should not store plaintext third-party passwords after this iteration.

## Credential Domain

Create `apps/user-rpc/internal/domain/credential` as a small domain package for third-party credential rules. It is responsible for:

- Platform-specific bind validation.
- Mapping bind input to database update fields.
- Encrypting `zf_password` and `oauth_password` before persistence.
- Decrypting `zf_password` and `oauth_password` before internal RPC responses.
- Computing whether a credential is bound when needed by profile logic.

This keeps password rules out of individual go-zero logic files.

## Directory Design

Target RPC directory layout:

```text
apps/user-rpc/
  user.proto
  pb/
  usercenterservice/
  internal/
    server/
    logic/
    svc/
    config/
    dao/
      model/
      query/
      repo/
    domain/
      credential/
    infra/
      mysql/
    httpclient/
      oauth/
```

`apps/user-rpc/internal/model/mysql.go` should move to `apps/user-rpc/internal/infra/mysql/mysql.go` to avoid confusion with generated DAO models. The accidental nested `apps/user-rpc/pb/pb` directory should be removed after confirming it is not referenced.

## Generation Rules

- API scaffold changes must start from `apps/user-api/user.api`.
- RPC scaffold changes must start from `apps/user-rpc/user.proto`.
- API code is regenerated with `make generate-api`.
- RPC code is regenerated with `make generate-rpc`.
- DAO model/query code is regenerated with `make generate-model` only when DDL changes require it.
- Generated `handler`, `server`, `pb`, `usercenterservice`, and `types` code should not be hand-written.

## Testing

Verification should include:

- Unit tests for credential encryption/decryption and bind update mapping.
- RPC logic tests for binding encrypted credentials and returning decrypted passwords.
- Existing common, API, and RPC tests.
- `go test ./...`
- `go build ./...`

