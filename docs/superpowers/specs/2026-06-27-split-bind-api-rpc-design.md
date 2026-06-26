# Split Bind API/RPC Design

## Goal

Replace the generic credential bind endpoint with platform-specific HTTP and RPC methods. Parameter-shape validation should live at the API/RPC request boundary, while logic files keep only business behavior.

## Decisions

- Expose three authenticated HTTP endpoints:
  - `POST /api/user/bind/oauth`
  - `POST /api/user/bind/zf`
  - `POST /api/user/bind/yxy`
- Expose three RPC methods:
  - `BindOauth(BindOauthRequest) returns (BindOauthResponse)`
  - `BindZf(BindZfRequest) returns (BindZfResponse)`
  - `BindYxy(BindYxyRequest) returns (BindYxyResponse)`
- Remove the public `BindType` switch contract from protobuf and generated client use.
- Keep API logic as orchestration: get current user id, normalize obvious whitespace/case where already established, call the matching RPC.
- Keep RPC logic as business behavior: confirm the user exists, encrypt or prepare credential fields, and update the user record.
- Move empty-field and invalid-shape checks into `.api` definitions where go-zero parsing can enforce them for HTTP callers.
- Remove simple parameter-invalid guards from logic files. Domain errors such as password length, user mismatch, user existence, and identity verification remain in domain/RPC logic.

## Request Shapes

HTTP request bodies:

```json
{"oauth_password":"..."}
{"zf_password":"..."}
{"device_id":"...","yxy_uid":"..."}
```

RPC request messages:

```proto
message BindOauthRequest {
    int64 user_id = 1;
    string oauth_password = 2;
}

message BindZfRequest {
    int64 user_id = 1;
    string zf_password = 2;
}

message BindYxyRequest {
    int64 user_id = 1;
    string device_id = 2;
    string yxy_uid = 3;
}
```

## Implementation Notes

- `apps/user-api/user.api` is the source for HTTP request types, routes, handlers, and swagger output.
- `apps/user-rpc/user.proto` is the source for RPC messages, server stubs, protobuf types, and the generated RPC client.
- The credential domain should expose platform-specific prepare methods so logic no longer depends on a generic `pb.BindRequest` or map-based update builder.
- Existing `GetUserPassword` behavior should continue returning trimmed/decrypted credential values.
- Existing generated files should be regenerated with `/Users/mgg/go/bin/goctl` to preserve goctl 1.10.1 output.

## Testing

- Add or update credential domain tests for platform-specific bind prepare methods.
- Update RPC integration bind tests to call `BindYxy`, `BindZf`, and `BindOauth`.
- Run focused package tests for credential logic and RPC compile paths.
- Run `go test ./...` and `go build ./...` before completion.
