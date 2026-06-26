# Split Bind API/RPC Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split generic credential binding into three platform-specific HTTP endpoints and RPC methods, and move simple request-shape validation out of logic files.

**Architecture:** `user.api` owns public HTTP request shapes and route generation. `user.proto` owns RPC method/message generation. API logic only authenticates the current user and calls the matching RPC, while RPC logic delegates platform-specific credential updates to the credential domain service.

**Tech Stack:** Go, go-zero/goctl 1.10.1, gRPC/protobuf, GORM repositories, existing `errorsx` normalization.

---

## File Structure

- Modify `apps/user-api/user.api`: add three bind request types and three authenticated routes.
- Regenerate `apps/user-api/internal/types/types.go`, handlers, routes, and swagger from `user.api`.
- Modify `apps/user-api/internal/logic/user/bind_oauth_logic.go`, `bind_zf_logic.go`, and `bind_yxy_logic.go`: call the corresponding RPC.
- Modify `apps/user-rpc/user.proto`: remove generic `BindType`/`BindRequest`/`BindResponse`, add split bind messages and methods using `Oauth`, `Zf`, and `Yxy` names.
- Regenerate `apps/user-rpc/pb`, `apps/user-rpc/internal/server`, and `apps/user-rpc/usercenterservice`.
- Replace `apps/user-rpc/internal/logic/bind_logic.go` with split bind logic files.
- Modify `apps/user-rpc/internal/dao/repo/user.go`: expose typed `UpdateOauthBindByID`, `UpdateZfBindByID`, and `UpdateYxyBindByID` methods instead of generic map updates.
- Modify `apps/user-rpc/internal/domain/credential/credential.go`: expose `PrepareOauthBind`, `PrepareZfBind`, and `PrepareYxyBind`.
- Modify tests in `apps/user-rpc/internal/domain/credential/credential_test.go` and `apps/user-rpc/test/usercenter_test.go`.
- Remove simple parameter-invalid checks from logic files where request definitions now own them.

## Tasks

- [ ] Write failing tests for split RPC bind methods.
- [ ] Update API and proto definitions.
- [ ] Regenerate go-zero API/RPC code with `/Users/mgg/go/bin/goctl`.
- [ ] Implement split API logic.
- [ ] Implement split RPC logic and credential domain builders.
- [ ] Remove obsolete generic bind references.
- [ ] Remove simple parameter-invalid guards from logic files.
- [ ] Run focused tests, then `go test ./...` and `go build ./...`.

## Verification Commands

```bash
go test ./apps/user-rpc/internal/domain/credential
go test ./apps/user-rpc/internal/logic
go test ./apps/user-api/internal/logic/user
go test ./...
go build ./...
```
