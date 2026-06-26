# 用户中心接入文档

本文档面向接入用户中心 HTTP API 的前端、后端服务或第三方应用。接口定义以 `apps/user-api/user.api` 为准。

## 基础约定

- 默认服务地址：按部署环境提供，例如本地开发常见为 `http://127.0.0.1:8080`。
- 请求格式：`Content-Type: application/json`。
- 响应格式：所有接口统一返回：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

- HTTP 状态码：业务成功和业务失败当前都返回 HTTP `200`，接入方必须以响应体中的 `code` 判断结果。
- 登录态：用户中心通过 Cookie 维护会话，默认 Cookie 名称为 `jh-passport`。调用登录态接口时，客户端需要携带该 Cookie。

## 登录态说明

成功调用以下接口会写入登录 Cookie：

- `POST /api/user/create/student`
- `POST /api/user/auth/password`

浏览器接入时，跨域场景需要前端请求启用凭证携带，例如：

```ts
fetch(url, {
  method: "POST",
  credentials: "include",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(payload)
})
```

服务端接入时，需要保存响应中的 `Set-Cookie`，后续请求通过 `Cookie` 请求头回传。

## 数据结构

### UserResp

```json
{
  "user": {
    "id": 1,
    "username": "302024000000",
    "student_id": "302024000000",
    "bind": {
      "zf": true,
      "yxy": false,
      "oauth": true
    },
    "user_type": "student",
    "email": "user@example.com",
    "phone_num": "",
    "create_time": "2026-06-27T00:00:00+08:00"
  }
}
```

### EmptyResp

空对象：

```json
{}
```

## 接口清单

### 创建学生账号

`POST /api/user/create/student`

创建学生账号。成功后会写入登录 Cookie，并返回用户 ID。

请求：

```json
{
  "student_id": "302024000000",
  "password": "secret123",
  "card_id": "330100200001010000",
  "email": "user@example.com"
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "user_id": 1
  }
}
```

字段说明：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `student_id` | 是 | 学号 |
| `password` | 是 | 用户中心密码 |
| `card_id` | 是 | 身份证号或身份校验字段 |
| `email` | 否 | 邮箱 |

### 账号密码登录

`POST /api/user/auth/password`

使用学号和密码登录。成功后会写入登录 Cookie，并返回用户信息。

请求：

```json
{
  "username": "302024000000",
  "password": "secret123"
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "user": {
      "id": 1,
      "username": "302024000000",
      "student_id": "302024000000",
      "bind": {
        "zf": false,
        "yxy": false,
        "oauth": false
      },
      "user_type": "student",
      "email": "",
      "phone_num": "",
      "create_time": "2026-06-27T00:00:00+08:00"
    }
  }
}
```

### 统一认证登录

`POST /api/user/auth/oauth`

接口已在 API 中定义，当前业务逻辑仍为预留实现。接入前需要确认服务端是否已补齐该登录方式。

请求：

```json
{
  "username": "302024000000",
  "password": "oauth-password"
}
```

### 小程序登录

`POST /api/user/auth/mini-program`

接口已在 API 中定义，当前业务逻辑仍为预留实现。接入前需要确认服务端是否已补齐该登录方式。

请求：

```json
{
  "app_type": "wechat",
  "code": "login-code"
}
```

### 获取当前用户信息

`POST /api/user/info`

需要登录 Cookie。返回当前登录用户信息。

请求体为空。

响应：同 `UserResp`。

### 重置密码

`POST /api/user/reset_password`

需要登录 Cookie。通过学号和身份信息校验后重置密码。

请求：

```json
{
  "stuid": "302024000000",
  "iid": "330100200001010000",
  "password": "new-secret123"
}
```

响应：`EmptyResp`。

字段说明：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `stuid` | 是 | 学号 |
| `iid` | 是 | 身份证号或身份校验字段 |
| `password` | 是 | 新密码 |

### 注销账号

`POST /api/user/delete`

需要登录 Cookie。通过学号和身份信息校验后注销当前账号。

请求：

```json
{
  "stuid": "302024000000",
  "iid": "330100200001010000"
}
```

响应：`EmptyResp`。

账号注销成功后，接入方应清理本地用户状态，并引导用户重新登录。

### 绑定统一认证密码

`POST /api/user/bind/oauth`

需要登录 Cookie。绑定统一认证密码。

请求：

```json
{
  "oauth_password": "oauth-password"
}
```

响应：`EmptyResp`。

### 绑定正方密码

`POST /api/user/bind/zf`

需要登录 Cookie。绑定正方教务系统密码。

请求：

```json
{
  "zf_password": "zf-password"
}
```

响应：`EmptyResp`。

### 绑定易校园信息

`POST /api/user/bind/yxy`

需要登录 Cookie。绑定易校园设备和用户标识。

请求：

```json
{
  "device_id": "device-id",
  "yxy_uid": "yxy-user-id"
}
```

响应：`EmptyResp`。

## 常见错误码

| code | msg | 说明 |
| --- | --- | --- |
| `0` | `ok` | 成功 |
| `10000` | `系统异常，请稍后重试` | 未分类系统错误 |
| `10001` | `下游服务调用失败` | 依赖服务异常 |
| `20000` | `用户未登录或登录已过期` | 未携带 Cookie、Cookie 过期或无效登录态 |
| `20003` | `参数非法` | JSON 格式错误、字段缺失或类型不匹配 |
| `40001` | `密码长度必须在6~20位之间` | 密码长度不符合要求 |
| `40002` | `账号或密码错误` | 登录凭据错误 |
| `40003` | `账号未激活` | 账号未激活 |
| `40004` | `用户不存在` | 用户不存在 |
| `40005` | `统一身份认证夜间不对外开放` | 统一认证当前不可用 |
| `40006` | `用户已经存在` | 注册时用户已存在 |
| `40007` | `密码需要修改` | 密码需要修改后继续使用 |
| `40008` | `无效的 Cookie` | Cookie 无法解析或被篡改 |

## 推荐接入流程

1. 登录或注册：调用 `POST /api/user/auth/password` 或 `POST /api/user/create/student`。
2. 保存登录态：浏览器依赖 Cookie 自动保存；服务端调用方保存 `Set-Cookie`。
3. 拉取用户信息：调用 `POST /api/user/info`，根据 `bind` 判断是否需要补充绑定。
4. 补充绑定：按需调用 `/api/user/bind/oauth`、`/api/user/bind/zf`、`/api/user/bind/yxy`。
5. 处理异常：所有请求都检查响应体 `code`；`20000` 或 `40008` 时重新登录。

## curl 示例

登录：

```bash
curl -i \
  -X POST "http://127.0.0.1:8080/api/user/auth/password" \
  -H "Content-Type: application/json" \
  -d '{"username":"302024000000","password":"secret123"}'
```

携带 Cookie 获取用户信息：

```bash
curl -i \
  -X POST "http://127.0.0.1:8080/api/user/info" \
  -H "Cookie: jh-passport=<cookie-value>"
```
