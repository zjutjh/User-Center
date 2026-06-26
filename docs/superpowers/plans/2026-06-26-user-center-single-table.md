# User Center Single-Table Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the single-table user-center design with encrypted third-party password storage, plaintext internal RPC password lookup, and a cleaner go-zero project layout.

**Architecture:** Keep the go-zero API/RPC scaffold generated from `user.api` and `user.proto`. Put reusable third-party credential rules in `apps/user-rpc/internal/domain/credential`, keep database setup under `apps/user-rpc/internal/infra/mysql`, and leave API logic as HTTP/session/RPC orchestration only.

**Tech Stack:** Go, go-zero, goctl, gRPC/protobuf, GORM/gen, MySQL, standard-library AES-GCM encryption.

---

## File Structure

- Modify `apps/user-rpc/internal/config/config.go`: add credential encryption config.
- Modify `config.example.yaml`: document the credential encryption key.
- Modify `config.yaml`: add local development credential key if missing.
- Create `apps/user-rpc/internal/domain/credential/crypto.go`: AES-GCM encryption/decryption for reversible third-party secrets.
- Create `apps/user-rpc/internal/domain/credential/credential.go`: bind validation, DB update mapping, and password response decryption.
- Create `apps/user-rpc/internal/domain/credential/credential_test.go`: unit tests for encryption and mapping.
- Create `apps/user-rpc/internal/infra/mysql/mysql.go`: moved DB constructor from `apps/user-rpc/internal/model/mysql.go`.
- Delete `apps/user-rpc/internal/model/mysql.go`: after imports are updated.
- Modify `apps/user-rpc/internal/svc/service_context.go`: use `infra/mysql` and initialize credential service.
- Modify `apps/user-rpc/internal/logic/bind_logic.go`: delegate credential updates to the domain package.
- Modify `apps/user-rpc/internal/logic/get_user_password_logic.go`: decrypt third-party passwords before returning them.
- Modify `apps/user-rpc/internal/logic/get_user_info_logic.go`: avoid returning encrypted passwords in profile responses or decrypt only if the existing API contract requires those fields.
- Remove `apps/user-rpc/pb/pb`: clean accidental nested generated output after confirming no imports reference it.
- Run `make generate-api` and `make generate-rpc` only if `user.api` or `user.proto` changes during execution.

---

### Task 1: Add Credential Encryption Config

**Files:**
- Modify: `apps/user-rpc/internal/config/config.go`
- Modify: `config.example.yaml`
- Modify: `config.yaml`

- [ ] **Step 1: Add config type**

Edit `apps/user-rpc/internal/config/config.go`:

```go
package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zjutjh/User-Center/common/dbx"
	"github.com/zjutjh/User-Center/common/sessionx"
)

type CredentialConf struct {
	SecretKey string
}

type Config struct {
	zrpc.RpcServerConf
	Mysql      dbx.MysqlConf
	Session    sessionx.Config
	Credential CredentialConf
}
```

- [ ] **Step 2: Add example configuration**

Under `UserRPC` in `config.example.yaml`, add:

```yaml
  Credential:
    SecretKey: 0123456789abcdef0123456789abcdef
```

- [ ] **Step 3: Add local configuration**

Under `UserRPC` in `config.yaml`, add the same block if no `Credential` block exists:

```yaml
  Credential:
    SecretKey: 0123456789abcdef0123456789abcdef
```

- [ ] **Step 4: Verify config compiles**

Run: `go test ./apps/user-rpc/internal/config`

Expected: package compiles or reports `? ... [no test files]`.

- [ ] **Step 5: Commit**

```bash
git add apps/user-rpc/internal/config/config.go config.example.yaml config.yaml
git commit -m "config: add credential encryption settings"
```

---

### Task 2: Implement Credential Domain Package

**Files:**
- Create: `apps/user-rpc/internal/domain/credential/crypto.go`
- Create: `apps/user-rpc/internal/domain/credential/credential.go`
- Create: `apps/user-rpc/internal/domain/credential/credential_test.go`

- [ ] **Step 1: Write tests**

Create `apps/user-rpc/internal/domain/credential/credential_test.go`:

```go
package credential

import (
	"testing"

	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
)

const testKey = "0123456789abcdef0123456789abcdef"

func TestEncryptDecrypt(t *testing.T) {
	codec, err := NewCodec(testKey)
	if err != nil {
		t.Fatalf("NewCodec() error = %v", err)
	}

	ciphertext, err := codec.Encrypt("secret-password")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if ciphertext == "" || ciphertext == "secret-password" {
		t.Fatalf("ciphertext = %q, want non-empty encrypted value", ciphertext)
	}

	plaintext, err := codec.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if plaintext != "secret-password" {
		t.Fatalf("plaintext = %q, want secret-password", plaintext)
	}
}

func TestBuildBindUpdatesEncryptsOAuthPassword(t *testing.T) {
	service, err := NewService(testKey)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	updates, err := service.BuildBindUpdates(&pb.BindRequest{
		UserId:        1,
		Type:          pb.BindType_BIND_TYPE_OAUTH,
		OauthPassword: "oauth-password",
	})
	if err != nil {
		t.Fatalf("BuildBindUpdates() error = %v", err)
	}

	encrypted, ok := updates["oauth_password"].(string)
	if !ok {
		t.Fatalf("oauth_password update type = %T, want string", updates["oauth_password"])
	}
	if encrypted == "" || encrypted == "oauth-password" {
		t.Fatalf("encrypted oauth password = %q, want encrypted value", encrypted)
	}

	plaintext, err := service.DecryptPassword(encrypted)
	if err != nil {
		t.Fatalf("DecryptPassword() error = %v", err)
	}
	if plaintext != "oauth-password" {
		t.Fatalf("plaintext = %q, want oauth-password", plaintext)
	}
}

func TestBuildBindUpdatesForYXY(t *testing.T) {
	service, err := NewService(testKey)
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	updates, err := service.BuildBindUpdates(&pb.BindRequest{
		UserId:   1,
		Type:     pb.BindType_BIND_TYPE_YXY,
		DeviceId: "device-1",
		YxyUid:   "uid-1",
	})
	if err != nil {
		t.Fatalf("BuildBindUpdates() error = %v", err)
	}
	if updates["device_id"] != "device-1" || updates["yxy_uid"] != "uid-1" {
		t.Fatalf("updates = %#v, want yxy fields", updates)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./apps/user-rpc/internal/domain/credential`

Expected: FAIL because `NewCodec`, `NewService`, and related methods are not defined.

- [ ] **Step 3: Add AES-GCM codec**

Create `apps/user-rpc/internal/domain/credential/crypto.go`:

```go
package credential

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type Codec struct {
	aead cipher.AEAD
}

func NewCodec(secretKey string) (*Codec, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm cipher: %w", err)
	}
	return &Codec{aead: aead}, nil
}

func (c *Codec) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	sealed := c.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (c *Codec) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	nonceSize := c.aead.NonceSize()
	if len(raw) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce := raw[:nonceSize]
	payload := raw[nonceSize:]
	plaintext, err := c.aead.Open(nil, nonce, payload, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt ciphertext: %w", err)
	}
	return string(plaintext), nil
}
```

- [ ] **Step 4: Add credential service**

Create `apps/user-rpc/internal/domain/credential/credential.go`:

```go
package credential

import (
	"strings"

	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type Service struct {
	codec *Codec
}

func NewService(secretKey string) (*Service, error) {
	codec, err := NewCodec(secretKey)
	if err != nil {
		return nil, err
	}
	return &Service{codec: codec}, nil
}

func (s *Service) BuildBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	switch in.Type {
	case pb.BindType_BIND_TYPE_YXY:
		return s.buildYXYUpdates(in)
	case pb.BindType_BIND_TYPE_ZF:
		return s.buildZFUpdates(in)
	case pb.BindType_BIND_TYPE_OAUTH:
		return s.buildOAuthUpdates(in)
	default:
		return nil, errorsx.ErrParameterInvalid
	}
}

func (s *Service) DecryptPassword(ciphertext string) (string, error) {
	return s.codec.Decrypt(ciphertext)
}

func (s *Service) buildYXYUpdates(in *pb.BindRequest) (map[string]any, error) {
	deviceID := strings.TrimSpace(in.DeviceId)
	yxyUID := strings.TrimSpace(in.YxyUid)
	if deviceID == "" || yxyUID == "" {
		return nil, errorsx.ErrParameterInvalid
	}
	return map[string]any{
		"device_id": deviceID,
		"yxy_uid":   yxyUID,
	}, nil
}

func (s *Service) buildZFUpdates(in *pb.BindRequest) (map[string]any, error) {
	password := strings.TrimSpace(in.ZfPassword)
	if password == "" {
		return nil, errorsx.ErrParameterInvalid
	}
	encrypted, err := s.codec.Encrypt(password)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"zf_password": encrypted,
	}, nil
}

func (s *Service) buildOAuthUpdates(in *pb.BindRequest) (map[string]any, error) {
	password := strings.TrimSpace(in.OauthPassword)
	if password == "" {
		return nil, errorsx.ErrParameterInvalid
	}
	encrypted, err := s.codec.Encrypt(password)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"oauth_password": encrypted,
	}, nil
}
```

- [ ] **Step 5: Run tests**

Run: `go test ./apps/user-rpc/internal/domain/credential`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add apps/user-rpc/internal/domain/credential
git commit -m "feat: add credential encryption domain"
```

---

### Task 3: Move MySQL Infrastructure Package

**Files:**
- Create: `apps/user-rpc/internal/infra/mysql/mysql.go`
- Delete: `apps/user-rpc/internal/model/mysql.go`
- Modify: `apps/user-rpc/internal/svc/service_context.go`

- [ ] **Step 1: Create new MySQL package**

Create `apps/user-rpc/internal/infra/mysql/mysql.go` with the content from `apps/user-rpc/internal/model/mysql.go`, but use package name `mysql`:

```go
package mysql

import (
	"fmt"

	"github.com/zjutjh/User-Center/common/dbx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(c dbx.MysqlConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

- [ ] **Step 2: Update service context import and fields**

Edit `apps/user-rpc/internal/svc/service_context.go`:

```go
package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/config"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/credential"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/infra/mysql"
	"github.com/zjutjh/User-Center/common/sessionx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Query      *query.Query
	Session    *sessionx.Manager
	Credential *credential.Service
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := mysql.NewDB(c.Mysql)
	logx.Must(err)

	credentialSvc, err := credential.NewService(c.Credential.SecretKey)
	logx.Must(err)

	return &ServiceContext{
		Config:     c,
		DB:         db,
		Query:      query.Use(db),
		Session:    sessionx.NewManager(c.Session),
		Credential: credentialSvc,
	}
}
```

- [ ] **Step 3: Remove old package file**

Delete `apps/user-rpc/internal/model/mysql.go`.

- [ ] **Step 4: Verify imports**

Run: `rg "internal/model|rpcmodel" apps/user-rpc`

Expected: no output.

- [ ] **Step 5: Run tests**

Run: `go test ./apps/user-rpc/internal/svc ./apps/user-rpc/internal/infra/mysql`

Expected: packages compile. If tests try to connect to MySQL through service context, run `go test ./apps/user-rpc/internal/infra/mysql` and defer service-context integration to the full test step.

- [ ] **Step 6: Commit**

```bash
git add apps/user-rpc/internal/svc/service_context.go apps/user-rpc/internal/infra/mysql apps/user-rpc/internal/model/mysql.go
git commit -m "refactor: move rpc mysql infrastructure"
```

---

### Task 4: Use Credential Domain in Bind Logic

**Files:**
- Modify: `apps/user-rpc/internal/logic/bind_logic.go`

- [ ] **Step 1: Replace local bind strategies**

Edit `apps/user-rpc/internal/logic/bind_logic.go` so `Bind` calls the domain service:

```go
updates, err := l.svcCtx.Credential.BuildBindUpdates(in)
if err != nil {
	return nil, err
}
```

Remove the local `bindStrategy` type, `bindStrategies` map, `buildBindUpdates`, `buildYXYBindUpdates`, `buildZFBindUpdates`, and `buildOAuthBindUpdates` functions. Remove the unused `strings` import.

- [ ] **Step 2: Convert encryption errors to service errors**

After calling `BuildBindUpdates`, handle unexpected encryption errors:

```go
updates, err := l.svcCtx.Credential.BuildBindUpdates(in)
if err != nil {
	if errors.Is(err, errorsx.ErrParameterInvalid) {
		return nil, err
	}
	l.Errorf("构建用户绑定信息失败: %v", err)
	return nil, errorsx.ErrUnknown
}
```

- [ ] **Step 3: Run targeted tests**

Run: `go test ./apps/user-rpc/internal/logic ./apps/user-rpc/internal/domain/credential`

Expected: PASS or no test files for logic package.

- [ ] **Step 4: Commit**

```bash
git add apps/user-rpc/internal/logic/bind_logic.go
git commit -m "refactor: route credential binding through domain"
```

---

### Task 5: Decrypt Passwords in Internal RPC Lookup

**Files:**
- Modify: `apps/user-rpc/internal/logic/get_user_password_logic.go`
- Modify: `apps/user-rpc/internal/logic/get_user_info_logic.go`

- [ ] **Step 1: Decrypt in GetUserPassword**

In `apps/user-rpc/internal/logic/get_user_password_logic.go`, decrypt fields before returning:

```go
zfPassword, err := l.svcCtx.Credential.DecryptPassword(user.ZfPassword)
if err != nil {
	l.Errorf("解密正方密码失败: %v", err)
	return nil, errorsx.ErrUnknown
}
oauthPassword, err := l.svcCtx.Credential.DecryptPassword(user.OauthPassword)
if err != nil {
	l.Errorf("解密统一认证密码失败: %v", err)
	return nil, errorsx.ErrUnknown
}

return &pb.GetUserPasswordResponse{
	StudentId:     user.StudentID,
	DeviceId:      user.DeviceID,
	YxyUid:        user.YxyUID,
	ZfPassword:    zfPassword,
	OauthPassword: oauthPassword,
}, nil
```

- [ ] **Step 2: Avoid encrypted password leakage in GetUserInfo**

Open `apps/user-rpc/internal/logic/get_user_info_logic.go`. If it returns `user.ZfPassword` or `user.OauthPassword`, replace those fields with empty strings unless the API contract explicitly requires decrypted third-party passwords in user profile responses:

```go
ZfPassword:    "",
OauthPassword: "",
```

If the existing API `UserInfo` does not expose password fields, prefer empty strings to avoid returning encrypted data through non-password RPCs.

- [ ] **Step 3: Run targeted tests**

Run: `go test ./apps/user-rpc/internal/logic ./apps/user-rpc/internal/domain/credential`

Expected: PASS or no test files for logic package.

- [ ] **Step 4: Commit**

```bash
git add apps/user-rpc/internal/logic/get_user_password_logic.go apps/user-rpc/internal/logic/get_user_info_logic.go
git commit -m "feat: decrypt internal password lookup"
```

---

### Task 6: Clean Generated Output Layout

**Files:**
- Delete: `apps/user-rpc/pb/pb`
- Verify: `apps/user-rpc/pb`

- [ ] **Step 1: Confirm nested pb is unused**

Run: `rg "apps/user-rpc/pb/pb|/pb/pb|pb/pb" .`

Expected: no imports reference nested `pb/pb`.

- [ ] **Step 2: Remove nested generated directory**

Delete `apps/user-rpc/pb/pb`.

- [ ] **Step 3: Regenerate RPC scaffold if proto generation path is wrong**

Run: `make generate-rpc`

Expected: generated files remain under `apps/user-rpc/pb`, `apps/user-rpc/internal/server`, and `apps/user-rpc/usercenterservice`, with no new `apps/user-rpc/pb/pb` directory.

- [ ] **Step 4: Commit**

```bash
git add apps/user-rpc/pb apps/user-rpc/internal/server apps/user-rpc/usercenterservice
git commit -m "chore: clean rpc generated layout"
```

---

### Task 7: Full Verification

**Files:**
- Modify if needed: files touched by earlier verification fixes.

- [ ] **Step 1: Format code**

Run:

```bash
goctl api format --dir apps/user-api
go fmt ./...
```

Expected: formatting completes without errors.

- [ ] **Step 2: Run all tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 3: Run full build**

Run: `go build ./...`

Expected: PASS.

- [ ] **Step 4: Inspect worktree**

Run: `git status --short`

Expected: only intentional changes remain.

- [ ] **Step 5: Commit verification fixes if any**

If formatting or verification changed files:

```bash
git add <changed-files>
git commit -m "test: verify user center credential flow"
```

---

## Self-Review

- Spec coverage: The plan covers single-table persistence, reversible third-party password storage, plaintext internal RPC lookup, API/RPC boundary preservation, directory cleanup, and goctl generation rules.
- Placeholder scan: No open-ended implementation placeholders remain; steps include exact files, commands, and code snippets.
- Type consistency: `CredentialConf`, `credential.Service`, `BuildBindUpdates`, and `DecryptPassword` are defined before use by `svc`, `bind_logic`, and `get_user_password_logic`.

