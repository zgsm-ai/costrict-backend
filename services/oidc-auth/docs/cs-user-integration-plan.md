# oidc-auth 对接 cs-user 用户中心集成计划（形态 A）

> 状态：十项关键决策已落定（见 §9 决策记录），可进入实施。2026-08-21 修订：新增决策 9/10（邀请码机制与 GitHub star 机制整体移除），并同步更新 §1–§8 相关章节；决策 8 已被决策 9 取代。2026-08-21 二次修订：决策 9 改判为「**邀请码保留但减负，判新改用 cs-user `is_new_user`**」，§1–§8 邀请码相关表述随之回退，GitHub star 移除（决策 10）不变。

## 1. 背景与目标

### 1.1 现状

oidc-auth 以 Casdoor 为上游 IdP，自持 PostgreSQL `users` 表（`internal/repository/models.go` 的 `AuthUser`）。每次登录回调时：

- `CasdoorProvider.GetUserInfo`（`internal/providers/casdoor.go:182`）用 `DecodeJWTPayloadUnverified` **非验证解码** Casdoor JWT，提取 phone / name / universal_id / github_id / github_name / employee_number / email；
- `CasdoorProvider.Update`（`internal/providers/casdoor.go:93`）把身份属性覆盖写回本地 `users` 表，并按 `machine_code + vscode_version` 合并设备行。

即：身份事实的来源是 Casdoor，但信任边界的建立方式是「不校验签名直接读 payload」。本地库同时承担身份档案与设备会话两种职责。

### 1.2 目标（形态 A）

身份职责委托给 cs-user，oidc-auth 保留为**设备会话与插件业务服务**：

| 职责 | 归属 |
|---|---|
| Casdoor JWT 验签（JWKS）、claims 规范化、external_key 计算 | cs-user（`/api/internal/auth/parse-identity`） |
| 用户身份档案（identity）、跨服务账号唯一性 | cs-user（`/api/internal/users/get-or-create`） |
| GitHub 身份（github_id / github_name） | cs-user（`/api/internal/users/{subject_id}/auth-identities`） |
| 设备列表、token 哈希握手、同设备踢下线 | oidc-auth（本地保留） |
| UserCode / DeviceCode 生成与校验 | oidc-auth（本地保留） |
| 邀请码机制（生成 / 校验 / 判新） | **保留，判新改用 cs-user `is_new_user`**（决策 9 修订，见 §3.7） |
| GitHub star 同步（sync_star）与 `isStar` 字段 | **整体移除**（决策 10） |
| 插件登录轮询、web 登录回调 | oidc-auth（本地保留） |
| Casdoor 账号合并（原绑定流程的 MergeByCasdoor） | **退役，不再提供**（决策 7 修订：绑定流程下线，见 §2.4） |
| 身份信息本地缓存（登录时写入 + userinfo 软 TTL 刷新） | oidc-auth 本地（决策 6） |

收益：

1. 消除非验证解码的安全问题——身份信任边界移入 cs-user 的 JWKS 校验；
2. oidc-auth 与其他服务（server 等）共享同一账号体系，`external_key` 全局唯一去重；
3. 为 cs-user ADR Phase 3（接管 `/oidc-auth/api/v1/plugin/login` 与 `/token` 端点的 strangler 计划）铺路，届时 oidc-auth 的登录链路可整体退役；
4. **范围收窄（决策 10）**：oidc-auth 只保留 device token 轮询 + Login flow；GitHub star 本地账号体系能力随本次改造移除，本地库回归「设备会话 + 登录」单一职责。邀请码按决策 9 修订保留为 Login flow 的获客能力，仅判新逻辑外移给 cs-user。

**cs-user 侧零改动**（决策 1 选定多身份接口方案后确认，见 §9）。

## 2. 目标调用链

### 2.1 插件登录（`/oidc-auth/api/v1/plugin/login`）

现状（`internal/handler/login_handler.go`）：

```
Casdoor 回调 → GetUserByOauth (login_handler.go:216)
  → ExchangeToken(code)
  → GetUserInfo(accessToken)            ← 非验证解码 JWT
  → 追加 Device（TokenProvider="custom"）
→ handleInviterCodeValidation           ← 本地库查 github_id/phone/email 判新老用户
→ provider.Update(user)                 ← 本地 Upsert + 设备合并
→ generateTokenPair → 302 重定向（state=HashToken(accessToken)）
```

目标（三次 cs-user RPC，串行）：

```
Casdoor 回调 → GetUserByOauth
  → ExchangeToken(code)                       ← 不变
  → ① usercenter.ParseIdentity(accessToken)   ← 验签 + claims 规范化 + external_key
  → ② usercenter.GetOrCreate(claims)          ← 确保 cs-user 侧账号存在，取 subject_id + is_new_user
  → 若 inviterCode != "" 且 is_new_user        ← 邀请码判新改用 cs-user（决策 9 修订）
  → ③ usercenter.ListIdentities(subject_id)   ← 取 provider=="github" 行 → GithubID/GithubName
  → 由 profile + github 身份组装 AuthUser（ID=universal_id）
  → 追加 Device                               ← 不变
→ 本地 Upsert（设备会话写入）                   ← 语义不变，实现收敛
→ generateTokenPair → 302 重定向               ← 不变
```

三次 RPC **全部 fail-closed**（决策 2）：任一失败即登录失败（503），不做本地降级。

受影响入口共两处：`login_handler.go:164`（插件回调）、`web_handler.go:80`（web 回调）走三连 RPC；原 `manger_handler.go:135`（绑定回调）随绑定流程整体下线（§2.4），不再是改造点。

### 2.2 Web 登录（`/oidc-auth/api/v1/manager/login`）

`internal/handler/web_handler.go:55` 的 `webLoginCallbackHandler` 同步改造：`GetWebUserByOauth` 内的 `GetUserInfo` 调用替换为同一 `usercenter` 客户端三次调用，虚拟 web 设备（`MachineCode="web-"+uuid[:8]`）逻辑不变；web 侧邀请码处理（`inviter_code` query 参数、判新、写 `InviterID`）保留，判新同样改用 cs-user `is_new_user`（决策 9 修订，见 §3.7）。

### 2.3 cs-user 侧接口契约（已核实，零改动可用）

三个端点均在 `X-Internal-Token` 共享密钥保护之下（`cs-user/internal/app/app.go:128`）。租户上下文：oidc-auth **不发送** `X-Tenant-Id`，全部落 `default` 租户（决策 4）。

**① `POST /api/internal/auth/parse-identity`**

```json
// 请求
{ "token": "<原始 Casdoor JWT>" }
// 响应
{
  "external_key": "casdoor:<universal_id>",
  "profile": {
    "id": "...", "sub": "...", "universal_id": "...",
    "name": "...", "preferred_username": "...",
    "email": "...", "phone": "...", "picture": "...",
    "owner": "...", "provider": "...", "provider_user_id": "..."
  }
}
```

错误语义：400 缺 token；503 未配置 JWKS；401 验签失败。cs-user 内部经 `NormalizeClaimsMap` 规范化，external_key 按账号粒度 `casdoor:<universal_id>` 计算（2026-08-18 坍缩迁移后的格式）。

**② `POST /api/internal/users/get-or-create`**

请求体为 `models.JWTClaims` 形状（由 ① 的 profile 映射，`external_claims` 可不传）：

```json
{
  "id": "...", "universal_id": "...", "name": "...",
  "preferred_username": "...", "email": "...", "phone": "...",
  "picture": "...", "owner": "...", "provider": "...",
  "provider_user_id": "..."
}
```

响应：`{ "user": {...}, "is_new_user": bool }`（user 含 `subject_id`，供 ③ 使用；`is_new_user` 供邀请码判新，决策 9 修订，见 §3.7）。错误映射：claims 为空 / 无有效标识 → 400，其余 → 500；`ErrExplicitlyUnbound` 场景会拒绝静默重建（用户主动解绑过的身份不允许无感复活）。

**③ `GET /api/internal/users/{subject_id}/auth-identities`**

响应：`{ "identities": [...] }`，每行含 `provider / provider_user_id / display_name / external_subject / is_primary` 等（`cs-user/internal/models/models.go:67`）。取 `provider=="github"` 行：`provider_user_id` → GithubID，`display_name` → GithubName。

该方案成立的前提已核实：cs-user 的 normalize 会从 `properties.oauth_GitHub_*` 推导 provider=github 并收割 `provider_user_id`（`cs-user/internal/auth/normalize.go:66-87`）；`get-or-create` 内部经 `BindIdentityToUser` 将其落成 identity 行（`cs-user/internal/user/service.go:565`）；Casdoor 用户无论以何种方式登录，JWT 的 properties 均携带全部已关联身份，跨次登录经 Path 4 跨提供者重挂逐步累积。

### 2.4 绑定流程（决策 7 修订：整体下线）

原 `/manager/bind/account` + `/callback`（router.go:50-51）的绑定回调核心动作是 **Casdoor 账号合并**：`MergeByCasdoor`（manger_handler.go:212，定义于 service/manager.go:34，全仓库唯一调用点）在 IdP 层把两个 Casdoor 账号合并为一个（被并方 universal_id 消失），外加本地行合并与 `MergeUserQuota`。2026-08-18 修订：**该流程不做迁移，随本次改造整体下线**，理由三条（均已核实）：

1. **可见范围内无调用方**：costrict-web 全仓库（portal / server 等）对 `bind/account`、`manager/bind`、`bindType` 零引用；server 侧另有独立绑定链路（`POST /api/auth/bind/start` → `bindAuthCallback` → 必要时 `confirm-merge`，server/cmd/api/main.go:641-642，带测试），在 cs-user 层（identity 行挂载/转移）承接「一人多登录方式」的统一；
2. **配额合并已休眠**：`config.yaml` 的 `quotaManager.baseURL` 为空串，`MergeUserQuota` 直接跳过（quota_service.go:63-66）——绑定流程唯一的外部副作用实际不生效；
3. **与形态 A 方向一致**：Casdoor 账号合并是身份职责（改写「一个人是谁」的 IdP 层事实），与「身份职责全权委托 cs-user」相悖；cs-user 多身份模型（决策 1）本就支持一个 subject 挂多条 identity 行，人的统一应发生在 cs-user 层而非 IdP 合并。

前一轮「保留但瘦身」结论中「不可委托给 cs-user bind/transfer」仍然成立（IdP 层账号合并与 identity 行挂载不同层，后者不删除 Casdoor 账号），但那是「保留时的正确实现方式」；能力整体退役后该问题不复存在，且不再产生合并死行 `casdoor:<deleted_uid>`（原方案每执行一次绑定即留一条，将来需 cs-user 增设 purge RPC 清理）。

下线后的残留影响——均为**既有现状的延续，非新退化**（今天未做绑定的分裂账号用户本就如此）：

- 一人多登录方式 = 多个 Casdoor 账号并存。oidc-auth 本地按 universal_id 分行（§2.1 锚点不变）：设备列表按行独立；邀请码判新改用 cs-user `is_new_user`（决策 9 修订，见 §6.1），本地不再查库判新；
- cs-user 侧经 server 绑定后，两条 identity 行归同一 subject，两种登录的 get-or-create 均返回该 subject——oidc-auth 本地行的合并要等 Phase 3 锚点迁至 subject_id 后自然消解；
- 已完成合并的历史用户不受影响（Casdoor 已是一账号）。

实施：`bindAccount` / `bindAccountCallback`（manger_handler.go:61/107）、`determineMainAccount` 与 12 字段 coalesce、`MergeByCasdoor` / `MergeUserQuota` 调用点整体删除；`BindAccountCallbackURI` 常量（const.go:37）与相关 RedirectURL 配置清理。**前置动作**：本地可见仓库虽无调用方，发布前仍需向插件端 / 外部管理页确认无硬编码引用；未确认前可先以 410 过渡一个版本。

## 3. 改动清单

### 3.1 `internal/config/config.go`：新增 userCenter 配置段

```go
type UserCenterConfig struct {
    BaseURL       string        `json:"baseURL"       mapstructure:"baseURL"`       // 如 http://cs-user:8082
    InternalToken string        `json:"internalToken" mapstructure:"internalToken"` // X-Internal-Token
    Timeout       time.Duration `json:"timeout"       mapstructure:"timeout"`
}
```

挂到 `AppConfig`（config.go:15），`config/config.yaml` 增加对应段。超时与登录回调的 15s 上下文匹配（三次 RPC 串行，单次建议 3–5s）。配置缺失时**启动即报错**（fail-closed 原则，不静默走旧路径）。

同时删除 star 相关配置：`AppConfig.GithubConfig` 字段（config.go:19）与 `GithubStarConfig` 结构（config.go:69-76）、`config/config.yaml` 的 `syncStar` 段（决策 10，见 §3.7）。

### 3.2 新增 `internal/usercenter` 客户端包

- `Client.ParseIdentity(ctx, rawJWT) (*IdentityProfile, error)`
- `Client.GetOrCreate(ctx, claims) (*User, bool, error)`（取 `subject_id` + `is_new_user`，供邀请码判新，决策 9 修订）
- `Client.ListIdentities(ctx, subjectID) ([]Identity, error)`
- `Client.GetProfile(ctx, subjectID) (*Profile, error)`（userinfo 软 TTL 刷新用，§3.6）
- 共享 `http.Client`（复用 `Server.HTTP` 的传输配置），错误分类：网络/超时 vs 401 验签失败 vs 5xx，统一向上返回「cs-user 不可用」或「身份无效」，handler 映射 503 / 401。登录链路（三连 RPC）按 §2.1 映射错误；userinfo 刷新失败按软 TTL 降级为旧值（不向用户报错）。

### 3.3 `internal/handler/login_handler.go` + `web_handler.go`：替换身份获取步骤

- `GetUserByOauth`（login_handler.go:216）与 `GetWebUserByOauth` 中 `providerInstance.GetUserInfo(ctx, token.AccessToken)` 替换为 `usercenter` 三连调用；
- 由 `profile` + github 身份组装 `repository.AuthUser`，字段映射见 §4；
- `user.ID == uuid.Nil → uuid.New()` 的兜底（login_handler.go:243-245）删除：parse-identity 成功即保证 universal_id 存在（缺失时按登录失败处理）；
- 邀请码处理保留但改造（决策 9 修订，见 §3.7）：`parseInviterCodeFromState`（login_handler.go:265-278）与两个调用点（:105、:175-180）保留；`handleInviterCodeValidation`（:280-317）的本地判新（github_id / phone / email 三次查库）与「you have registered」拒绝登录逻辑删除，改为消费 `get-or-create` 返回的 `is_new_user`：仅 `is_new_user=true` 时 `ValidateInviteCode` 并写 `InviterID`，否则忽略邀请码、登录正常完成。

### 3.4 `internal/providers/casdoor.go`：GetUserInfo 彻底退役（决策 5）

- `GetUserInfo` 与 `DecodeJWTPayloadUnverified` 依赖一并删除，不留 dev-only 旁路（与 fail-closed 决策一致；dev 环境用 docker-compose 起 cs-user）；
- `OAuthProvider` 接口（`internal/providers/interface.go:30`）中 `GetUserInfo` 方法随之删除；
- **EmployeeNumber 停止同步**：本地 `employee_number` 列与已有数据保留不动，新用户该字段为空（原 manger_handler.go:199 绑定流程中对它的 coalesce 兼容逻辑随流程一并删除，见 §2.4）；
- `ExchangeToken`、`RefreshToken`、`GetAuthURL`、`GetEndpoint` 不变。

### 3.5 `Update()` 收敛 + 零值覆盖 bug 修复（决策 3）

现状 bug（casdoor.go:126-133）：每次登录无条件覆盖 `Vip / Company / Location`，但上游从不提供 → 每次登录刷成零值。收敛规则：

- **有源字段**（name / phone / email / github_id / github_name，来源为 cs-user profile 与 auth-identities）：每次登录覆盖写回本地；
- **本地自有字段**（Vip / Company / Location 等）：不再覆盖，保留本地值。

（决策 3 此前的论述「有源字段保持新鲜以供邀请码判新本地查询」已失效——判新改用 cs-user `is_new_user`，本地字段不再参与判新，见 §6.1。）

### 3.6 `internal/handler/manger_handler.go` userinfo：软 TTL 缓存（决策 6）

`userInfoHandler`（manger_handler.go:258）读的字段分两类，处置不同：

- **身份类**（username / email / phone / githubID / githubName）：本地行为缓存，登录时写入（决策 3 为主刷新路径），叠加**软 TTL** 兜长会话漂移（refresh token 用户可数周不重登）；
- **本地产品字段**（Vip / isPrivate）：不走 TTL，仍读本地列；isStar / GithubStar 随决策 10 移除（§3.7），`userInfoHandler` 不再计算与输出 `isStar`。

设计要点：

1. `AuthUser` 新增 `identity_synced_at` 列（不可复用 `UpdatedAt`——设备写入也会动它），另增存 `subject_id`（get-or-create 响应已含）；
2. TTL 到期 → 读路径上惰性刷新（或后台刷新），刷新源为现有内部端点，cs-user 仍零改动：`GET /api/internal/users/{subject_id}/profile`（username / display_name / email / phone）+ `/:subject_id/auth-identities`（github 行）；
3. **软过期（stale-while-revalidate），禁止硬过期**：刷新失败时返回旧值 + 告警，不报错。决策 2 的 fail-closed 只覆盖登录（建立信任边界的步骤）；userinfo 是纯读、数据在登录时已验证过，把 cs-user 可用性扩大到稳态读路径是无安全收益的可用性回归；
4. TTL 取值 ≥ cs-user 自身 `GetOrCreateUser` 的 `syncInterval` 门控（service.go:477）——同一模式的先例，两层节流不应倒挂；
5. 语义增益：cs-user 用户字段是 user-owned（重新登录不覆盖，service.go:470-475 防 auto-clobber），`profile` 端点返回用户自编辑后的值——比登录链路的 Casdoor claims 值更权威，恰合「以 cs-user 为主」的定位。

### 3.7 邀请码减负（决策 9 修订）与 GitHub star 机制移除（决策 10）

本次对接退役 GitHub star 本地账号体系能力；邀请码按决策 9 修订**保留但减负**（判新改用 cs-user `is_new_user`），oidc-auth 收窄为「device token 轮询 + Login flow」。

**邀请码（决策 9 修订：保留但减负）**

- **保留**（外部协议与文件不变，仅判新语义变化）：
  - `internal/handler/login_handler.go`：`parseInviterCodeFromState`（:265-278）保留，回调内调用点（:105、:175-180）保留；
  - `internal/handler/web_handler.go`：`webLoginHandler` 的 `inviter_code` query 参数（:28、:39、:49）保留，`getUserInviteCodeHandler`（:226-273）保留；
  - `internal/handler/router.go`：`GET invite-code` 路由（:56）保留；
  - `pkg/utils/invite_code.go`：整文件保留（`GenerateInviteCode` / `ValidateInviteCode` / `GenerateUniqueInviteCode`）；
  - `internal/constants/const.go`：`InviteCodeLength / InviteCodeChars / InviteCodeSeparator`（:48-53）保留；
  - `internal/repository/models.go`：`AuthUser.InviteCode`（:43）、`AuthUser.InviterID`（:44）保留；`operator.go:27` 的 `invite_code` 白名单条目保留。
- **判新更换（本次改造的核心改动）**：
  - `handleInviterCodeValidation`（login_handler.go:280-317）的本地判新（按 github_id / phone / email 三次 `GetUserByField`）与「you have registered」拒绝登录删除，改消费 cs-user `get-or-create` 响应的 `is_new_user`：`inviterCode != ""` 且 `is_new_user=true` 时 `ValidateInviteCode` 并写 `InviterID`；否则忽略邀请码，正常完成登录（已注册用户带邀请码不再被拒）；
  - web 侧 `webLoginCallbackHandler` 的判新 / 写 `InviterID`（:93-124）同款改造；
  - 登录路径查库从最多 4 次（判新 3 + `ValidateInviteCode` 1）降为 1 次（仅 `ValidateInviteCode`）。

**GitHub star（决策 10）**

- `internal/sync/sync_star.go`：整个包删除（`SyncStar` / `StarCount` / `Stargazers` / `withSyncLock` / `StarSyncTimer` / `Owner` / `Repo` 全局）；
- `cmd/main.go`：star 启动逻辑（:138-141）与 `internal/sync` import（:20）删除；
- `internal/config/config.go`：`GithubConfig` 字段（:19）与 `GithubStarConfig`（:69-76）删除；`config/config.yaml` 的 `syncStar` 段（:74-89）删除；
- `internal/repository/models.go`：`StarUser`（:10-16）、`SyncLock`（:19-22）、`AuthUser.GithubStar`（:40）删除；`operator.go` 的 StarUser / SyncLock 白名单（:29-38）与 `AddSyncLock` / `RemoveSyncLock`（:303-319）删除；`datebase.go:97` AutoMigrate 的 `&SyncLock{}` 移除（`StarUser` 本就不在 AutoMigrate 列表）；
- `internal/constants/const.go`：`GitHubStarBaseURL / DefaultPageSize / MaxPageLimit`（:16-22）删除；
- `internal/handler/manger_handler.go`：`userInfoHandler` 的 isStar 计算与输出（:277-283、:294）删除，`internal/sync` import（:19）移除。

**外部协议影响**

- 插件端：`state` 参数携带的 `__inviter_code=XXXX` 后缀（`InviteCodeSeparator`）**保留不变**，无协议变化；
- web 端：`GET /manager/login?inviter_code=XXXX`（邀请分享链接落地参数）、`GET /manager/invite-code` API（「我的邀请码」页面）**保留不变**；
- **credit-manager 前端（`D:\DEV\credit-manager`）是上述端点与字段的直接消费者**，本次仅随 star 移除（决策 10）与绑定下线（决策 7）同步改动，邀请码相关（决策 9 修订后保留）不动：
  - 邀请码（保留，不动）：`src/api/mods/quota.mod.ts` 的 `getInviteCode()` 与 `getLoginUrl()` 的 `inviter_code` 参数；`useProfile.fetchInviteCode`（home 个人信息）；`activity-card.toInvite`（活动页「邀请有礼」卡）；`annual-summary-page` 的邀请码拉取与 `:invite-code` 传参链；`login-page` 对 `inviteCode` URL query 的透传（`isShare` 保留）；`credit-reward-plan` 邀请人/被邀请人 UI（页面去留需运营确认）；
  - 绑定（决策 7）：`getBindAccount()` 与 profile 的 bind-github / bind-phone 按钮（GitHub 名/手机号保留为只读展示）；
  - star 门控（决策 10）：`useProfile` 的 `isStar`/`transferData`、`credit-transfer-modal` 的 `starred` 提示、`usage-section` / `home-page` 的 `:is-star` 透传（`is_star` 数据源见下一条）；
  - 测试与 mock：`mock/auth.ts` 的 `/invite-code` mock 保留，`/bind/account` mock 删除；`useProfile.spec.ts` / `annual-summary-page.spec.ts` 的 invite_code 夹具保留；
- **quota-manager 直读库依赖（新增发现，重点）**：`quota-manager/internal/services/quota.go:364-393` 的 `GetUserQuota` 经 `auth_database`（`config.yaml` `dbname: "auth"`，即 oidc-auth 库）**直读 `users.GithubStar` 列**计算 `is_star`（受 `config.yaml` `github_star_check.enabled` 门控，当前 false），并经 `GET /quota` 下发给 credit-manager（前端 `is_star` 转移门控的消费方）。oidc-auth 停写该列后其值冻结为旧值 → **须与 quota-manager / credit-manager 同步删除 is_star 计算与前端门控，不能只依赖 `enabled: false` 配置**；quota-manager 另有的 star-check-permission 服务（`internal/services/star_check_permission.go`，经 Higress StarCheckClient 下发企业权限）是独立能力，与 is_star 转移门控无关，不受影响；
- userinfo 响应 `isStar` 字段不再输出，需确认管理页无硬编码依赖；另注意 credit-manager `services/user.ts` 消费的 `employee_number` 后端现状本就未下发（前端 `|| ''` 兜底），cs-user 集成后新响应 shape 需与该映射对齐；
- star 移除未确认前可先以 410 过渡一个版本；确认后随本次改造一并删除（邀请码不涉及 410）。

**数据处置**：仅 `GithubStar` 列停止读写与输出（**不做 DROP**，回滚友好）；`InviteCode / InviterID` 列继续读写（判新改用 cs-user `is_new_user` 后，`InviterID` 仅在 `is_new_user=true` 时写入）；`star_users` / `sync_locks` 表数据留存无副作用。

## 4. 字段映射

| AuthUser 字段 | 新来源 | 说明 |
|---|---|---|
| `ID` | `profile.universal_id`（uuid.Parse） | 与现状一致；同时等于 cs-user `external_key` 的主体段 |
| `Name` | `profile.name` | |
| `Phone` | `profile.phone` | `+86` 前缀剥离逻辑保留在 oidc-auth 侧 |
| `Email` | `profile.email` | |
| `GithubID` | auth-identities 中 `provider=="github"` 行的 `provider_user_id` | 见 §2.3 ③ |
| `GithubName` | 同上行的 `display_name` | |
| `EmployeeNumber` | **停止同步**（决策 5） | 本地列保留，新值为空 |
| `Vip / Company / Location` | 本地保留，不再被覆盖 | 修复 3.5 的 bug |
| `UserCode / DeviceCode / Devices` | oidc-auth 本地生成 | 不变 |
| `InviteCode / InviterID` | 保留（决策 9 修订） | 判新改用 cs-user `is_new_user`，`InviterID` 仅 `is_new_user=true` 时写入 |
| `GithubStar` | **移除**（决策 10） | 列停止读写，`isStar` 不再输出 |
| `SubjectID` | `get-or-create` 响应的 `subject_id` | 新增列，供 userinfo TTL 刷新与 auth-identities 查询（决策 6） |
| `IdentitySyncedAt` | 登录链路 / TTL 刷新时刻 | 新增列，软 TTL 的时间戳（决策 6） |

## 5. 存量用户与 ID 对齐

- 本地 `users.id` 本就是 Casdoor `universal_id`，cs-user `external_key` 坍缩后也是 `casdoor:<universal_id>`——**两侧天然按同一锚点对齐，无需批量 ETL**。
- 懒加载迁移：存量用户下次登录时经 `get-or-create` 在 cs-user 侧补建账号；`get-or-create` 幂等，重复调用无副作用。
- 顺序要求：先发 cs-user（含 2026-08-18 external_key 坍缩迁移已生效的环境），后发 oidc-auth 新版。

## 6. 风险与对策

### 6.1 邀请码「新用户」判定语义（2026-08-21 修订：改用 cs-user `is_new_user`）

原风险（2026-08-18 落定）：`handleInviterCodeValidation`（login_handler.go:280-317）与 web 侧同款逻辑（web_handler.go:93-124）按 **github_id / phone / email 查本地库**判定新老用户，且不能改用 cs-user 的 `is_new_user`——迁移窗口内所有存量 oidc-auth 用户在 cs-user 侧都是「新用户」（首次 get-or-create），用它会导致老用户误用邀请码成功。

**2026-08-21 修订（决策 9 改判）**：用户拍板判新改用 cs-user `get-or-create` 响应的 `is_new_user`，邀请码机制保留但减负（§3.7）。原「不能改用」决策作废，风险转化为**接受项**：

- **迁移窗口语义变化（接受）**：存量 oidc-auth 用户迁 cs-user 后首次登录 `is_new_user=true`，老用户携带邀请码会绑定成功（原「你已注册」拒绝登录 → 邀请码被绑定、正常完成登录）。收益归邀请人，无安全 / 计费风险；已注册用户不再被拒登录。如需严格语义，可后续补「本地用户行创建时间」双条件（本次不做）；
- **判新粒度分歧消除**：「本地判新 vs cs-user 按 universal_id 去重的粒度分歧」（同 email 不同 universal_id 的账号）不再存在——是否「新用户」由 cs-user 按身份锚点（universal_id / external_key）判定，与本地账号分裂现状解耦；
- **判新与 userinfo 的列/流耦合**：判新职责移交 cs-user 后，决策 3 有源字段保持新鲜、决策 6 TTL 互不影响的论述不再依赖本地判新列，边界更清晰；
- 登录路径查库从最多 4 次降为 1 次（仅 `ValidateInviteCode`）。

### 6.2 可用性耦合（已决策：全 fail-closed）

登录链路串行依赖 cs-user 三次同步调用，cs-user 成为登录路径的硬依赖：

- 任一 RPC 失败（网络/超时/5xx/验签失败）即登录失败，不做本地降级——身份信任边界已在 cs-user，降级等于放弃本次改造的安全收益；
- 配套要求：cs-user 部署需多实例保障；oidc-auth 侧监控登录成功率与三次 RPC 的错误分类（401 / 503 / 超时）；
- 单次 RPC 超时预算受登录回调 15s 上下文约束，三次串行建议各 3–5s。

### 6.3 GitHub 身份的时序特性

auth-identities 反映的是**登录时刻累积**的 identity 行。已核实 Casdoor 无论以何种方式登录，JWT properties 都带全部已关联身份，且 cs-user normalize 据此推导 github provider 并落行——因此「oidc-auth 旧路径能拿到 GithubID 的场合，新路径同样能拿到」。残留边界：若某账号在 Casdoor 侧从未通过 GitHub 授权（properties 无 oauth_GitHub_*），两条路径拿到的都是空，语义一致。

### 6.4 cs-user 侧仍在快速演进

external_key 坍缩迁移 2026-08-18 刚落地。对接期间冻结对 `parse-identity` / `get-or-create` / `auth-identities` 契约的破坏性变更，或约定契约版本号。

### 6.5 软删除与解绑语义

cs-user 用户带 `gorm.DeletedAt` 软删除；identity 带 `explicitly_unbound` 标记（主动解绑后静默重登会被 `ErrExplicitlyUnbound` 拒绝）。oidc-auth 侧统一按登录失败兜底，错误信息透传给用户引导重新绑定。软删用户经 get-or-create 的具体行为（拒绝 / 重建）需联调实测确认。

### 6.6 测试基建缺失（最大执行风险）

oidc-auth 全仓库当前**没有任何单元测试**，而本次改动位于登录关键路径（token 哈希握手、踢下线、设备合并）。执行顺序要求：先为 `usercenter` 客户端与 `GetUserByOauth` 组装逻辑补单测（cs-user 接口可用 httptest 桩模拟），再动主链路；跨三方联调（Casdoor + cs-user + oidc-auth）作为最终验收。

### 6.7 待测试验证场景用例表（2026-08-21 增补）

按 §10 执行顺序拆为可验证场景。**层级**：单测 = oidc-auth 侧 httptest 桩模拟 cs-user 接口（P0-2 先行补齐）；联调 = Casdoor + cs-user + oidc-auth 三方真环境（P0-1 就绪后执行，作为 P1/P2/P3 验收门）。

**A. usercenter 客户端（P1-1）**

| 编号 | 场景 | 前置 / 输入 | 预期结果 | 层级 |
|---|---|---|---|---|
| U-01 | ParseIdentity 成功 | 合法 JWT，桩返回 200 | 返回 IdentityProfile（含 universal_id），无错误 | 单测 |
| U-02 | ParseIdentity 验签失败 | 桩返回 401 | 错误分类「身份无效」，handler 映射 401，登录失败 | 单测 |
| U-03 | ParseIdentity 不可达 | 桩不监听 / 连接超时 | 错误分类「cs-user 不可用」，映射 503 | 单测 |
| U-04 | ParseIdentity 5xx | 桩返回 500 | 同 U-03，映射 503 | 单测 |
| U-05 | GetOrCreate 新用户 | 桩返回 `is_new_user=true` | 客户端透传 `(user, true)` | 单测 |
| U-06 | GetOrCreate 老用户 | 桩返回 `is_new_user=false` | 客户端透传 `(user, false)` | 单测 |
| U-07 | GetOrCreate 软删用户 | 桩模拟 cs-user 软删行为 | 行为以联调实测为准（§6.5），错误按登录失败兜底 | 联调 |
| U-08 | GetOrCreate explicitly_unbound | 桩返回 `ErrExplicitlyUnbound` | 按登录失败兜底，错误透传给用户引导重新绑定 | 单测 |
| U-09 | GetOrCreate 并发竞态 | 同 universal_id 两并发请求 | cs-user 恢复 existing 行并返回 `false`（cs-user 侧行为） | 联调 |
| U-10 | GetProfile 成功 / 失败 | 200 / 401 / 500 / 超时 | 成功返回 profile；失败按错误分类（userinfo 路径降级旧值） | 单测 |
| U-11 | ListIdentities 无 github 行 | 桩返回空 / 无 `provider=="github"` | 返回空列表，语义与旧路径一致（§6.3） | 单测 |
| U-12 | 超时预算 | 三次 RPC 串行计时 | 单次 3–5s、合计在登录回调 15s 上下文内；任一超时整体失败 | 单测 |
| U-13 | 配置缺失 | `userCenter` 段未配置 | 启动即报错（fail-closed），不静默走旧路径 | 单测 |

**B. 登录链底座（P1-2）**

| 编号 | 场景 | 前置 / 输入 | 预期结果 | 层级 |
|---|---|---|---|---|
| L-01 | 新用户全链路组装 | 三连 RPC 桩成功，`is_new_user=true` | 本地 `AuthUser` 正确组装（含 github 身份），`SubjectID` 写入 | 单测 |
| L-02 | 老用户全链路 | get-or-create 命中，`is_new_user=false` | 本地行更新，不重建 | 单测 |
| L-03 | parse-identity 失败 | 桩返回 401 | 登录失败 401，token 不签发 | 单测 |
| L-04 | get-or-create 失败 | 桩返回 503 | 登录失败 503，无本地降级（fail-closed） | 单测 |
| L-05 | list-identities 失败 | 桩返回 503 | 登录失败 503（fail-closed） | 单测 |
| L-06 | uuid.Nil 兜底删除 | parse-identity 返回无 universal_id | 按登录失败处理，不再生成新 uuid（:243-245 兜底已删） | 单测 |
| L-07 | github 时序对照 | Casdoor properties 带 `oauth_GitHub_*` | 新路径同样拿到 GithubID（§6.3） | 联调 |
| L-08 | 无 github 授权账号 | properties 无 `oauth_GitHub_*` | 两条路径均空，语义一致 | 联调 |

**C. 登录邀请机制（P2，目标 ①）**

| 编号 | 场景 | 前置 / 输入 | 预期结果 | 层级 |
|---|---|---|---|---|
| I-01 | 新用户 + 合法邀请码 | `is_new_user=true` 且码有效 | `ValidateInviteCode` 通过，写 `InviterID` | 单测 |
| I-02 | 新用户 + 非法 / 过期邀请码 | `is_new_user=true` 且码无效 | 登录失败（沿用 `ValidateInviteCode` 原语义） | 单测 |
| I-03 | 老用户 + 邀请码（迁移窗口语义变化） | `is_new_user=false` + 码 | 忽略邀请码、正常登录、不写 `InviterID`（§6.1 接受项） | 单测 |
| I-04 | 新用户无邀请码 | `is_new_user=true`，`inviterCode=""` | 正常登录，不调 `ValidateInviteCode` | 单测 |
| I-05 | 「you have registered」删除 | 原本地判新（github_id / phone / email 三次查库） | 不再触发拒绝登录 | 单测 |
| I-06 | 查库次数收敛 | 登录全流程 DB 访问计数 | 最多 1 次（仅 `ValidateInviteCode` 查邀请码表） | 单测 |
| I-07 | web 侧同款 | `webLoginCallbackHandler`（web_handler.go:93-124） | 判新 / 写 `InviterID` 与 P2-1 一致 | 单测 |
| I-08 | 邀请码存量回归 | `parseInviterCodeFromState`、`inviter_code` query、`GET /invite-code`、`Generate/ValidateInviteCode`、`InviteCode/InviterID` 列、`__inviter_code=` 插件协议 | 全部保留可用，协议不变 | 单测 + 联调 |

**D. userinfo 软 TTL（P3，目标 ②）**

| 编号 | 场景 | 前置 / 输入 | 预期结果 | 层级 |
|---|---|---|---|---|
| F-01 | 未过期 | `identity_synced_at` 新鲜 | 直接读本地列，零 RPC | 单测 |
| F-02 | 过期 + 刷新成功 | TTL 过期，`GetProfile` + `ListIdentities` 成功 | 回写本地列并更新 `identity_synced_at` | 单测 |
| F-03 | 过期 + 刷新失败 | cs-user 不可达 | 返回旧值 + 告警，不向用户报错 | 单测 |
| F-04 | 字段分流 | userinfo 请求 | 身份字段（username/email/phone/githubID/githubName）走 TTL；`Vip` / `isPrivate` 永远读本地列 | 单测 |
| F-05 | 并发单飞 | 同一用户并发过期请求 | 仅触发一次刷新（stale-while-revalidate，避免击穿） | 单测 |
| F-06 | `isStar` 输出保持 | P4 前 star 同步运行中 | `isStar` 列与输出一致，不因 P3 改动漂移 | 联调 |
| F-07 | 刷新持久化 | TTL 刷新后重启服务 | `identity_synced_at` 已落库，重启后按新时间判断 | 单测 + 联调 |

**E. 发布 / 回滚（P5 / P6）**

| 编号 | 场景 | 前置 / 输入 | 预期结果 | 层级 |
|---|---|---|---|---|
| R-01 | 回滚兼容 | 旧代码 + 已新增 `subject_id` / `identity_synced_at` 列 | 旧代码不读不写两列，登录正常 | 联调 |
| R-02 | 软删用户经 get-or-create 实测 | 软删用户再次登录 | 拒绝 / 重建行为与 §6.5 联调结论一致 | 联调 |
| R-03 | star 三端同发（P4/P5 执行时） | quota-manager `is_star` 计算移除、credit-manager 门控移除、插件端无 `isStar` 硬编码 | 三端一致，无「前端删门控、后端仍读冻结列」中间态 | 联调 |

## 7. 发布与回滚

1. **前置**：cs-user 上线且 `CS_USER_CASDOOR_JWKS_URL` 已配置、多实例部署；联调环境验证 §6.5 语义；向管理页确认对 `isStar` 无硬编码引用（§3.7，未确认可先 410 过渡一个版本；`__inviter_code=` / `inviter_code` query / `/invite-code` 因邀请码保留，协议不变，无需确认）；**与 quota-manager / credit-manager 同步落地 is_star 移除**——quota-manager 删 `GetUserQuota` 的 is_star 计算（§3.7 直读库依赖），credit-manager 删 `isStar`/`transferData` 门控，三端同发避免「前端删门控、后端仍读冻结列」的中间态。
2. **发布**：oidc-auth 新版带 `userCenter` 配置上线；配置缺失时启动即报错；`syncStar` 配置段同步移除。
3. **观察**：登录成功率、三次 RPC 的 401/503/超时比例、本地 Upsert 冲突、移除端点的 404/410 与 `isStar` 缺失反馈。
4. **回滚**：回滚 oidc-auth 版本 + 移除 `userCenter` 配置即可——本地库仅**新增**了 `subject_id` / `identity_synced_at` 两列（旧代码不读不写，向后兼容），旧代码按旧路径读写不受影响；cs-user 侧已补建的账号与 identity 行留存无副作用。star 移除**不涉及数据迁移**（`GithubStar` 列停读写、未 DROP），随版本回滚整体恢复；邀请码保留，判新改造随登录链路一起回滚。

## 8. 工作量与涉及文件

| 文件 | 改动 |
|---|---|
| `internal/config/config.go` | 新增 `UserCenterConfig`；删除 `GithubConfig` / `GithubStarConfig`（决策 10） |
| `internal/usercenter/`（新） | 四方法客户端（§3.2）+ httptest 桩单测 |
| `internal/handler/login_handler.go` | `GetUserByOauth` 身份获取替换为三连调用；`handleInviterCodeValidation` 判新改造（本地判新 + 「you have registered」删除，改消费 `is_new_user`，决策 9 修订）；`parseInviterCodeFromState` 保留 |
| `internal/handler/web_handler.go` | `GetWebUserByOauth` 同款替换；`webLoginCallbackHandler` 判新 / 写 `InviterID` 同款改造（决策 9 修订）；`inviter_code` 参数与 `getUserInviteCodeHandler` 保留 |
| `internal/handler/router.go` | 删除绑定路由 `bind/account`（§2.4，决策 7）；`invite-code` 路由（:56）保留（决策 9 修订） |
| `internal/handler/manger_handler.go` | 绑定流程整体删除：`bindAccount` / `bindAccountCallback` 及合并逻辑下线（§2.4，含 router.go:50-51 路由与 const.go:37 常量清理）；`userInfoHandler` 增加软 TTL 读逻辑（§3.6）、删除 isStar 输出（决策 10） |
| `internal/repository/models.go` | `AuthUser` 增 `SubjectID` / `IdentitySyncedAt` 列；`InviteCode` / `InviterID` 保留（决策 9 修订）；删除 `GithubStar` 与 `StarUser` / `SyncLock`（决策 10） |
| `internal/repository/operator.go` | allowedFields 清理（`github_star` / StarUser / SyncLock 条目；`invite_code` 条目保留）；删除 `AddSyncLock` / `RemoveSyncLock` |
| `internal/repository/datebase.go` | AutoMigrate 移除 `&SyncLock{}` |
| `internal/sync/sync_star.go` | 整包删除（决策 10） |
| `cmd/main.go` | 删除 star 启动逻辑（:138-141）与 `internal/sync` import |
| `internal/providers/casdoor.go` | `GetUserInfo`/非验证解码删除、`Update` 按决策 3 收敛、EmployeeNumber 停同步 |
| `internal/providers/interface.go` | 接口删 `GetUserInfo` |
| `internal/constants/const.go` | `InviteCodeLength/Chars/Separator` 保留（决策 9 修订）；删除 `GitHubStarBaseURL/DefaultPageSize/MaxPageLimit`（决策 10） |
| `pkg/utils/invite_code.go` | 整文件保留（决策 9 修订，仅判新改造不涉及该文件） |
| `config/config.yaml` | 新增 `userCenter` 段；删除 `syncStar` 段（决策 10） |

cs-user 侧：零改动。

主要成本不在编码：测试基建从零补起（§6.6）、三方联调环境、cs-user 成为登录硬依赖后的运维配套，是本次改造的真正大头。

## 9. 决策记录（2026-08-18 落定，2026-08-21 增补决策 9/10 并修订决策 8/9）

| # | 事项 | 决策 | 理由摘要 |
|---|---|---|---|
| 1 | `properties` claim 透传缺口 | **经 cs-user 多身份接口取 GitHub 身份**：parse-identity 验签 → get-or-create → auth-identities 查 `provider=="github"` 行 | cs-user 零改动；normalize 已从 properties 推导并落 identity 行；EmployeeNumber 无对应字段，归入决策 5 处置 |
| 2 | cs-user 不可达策略 | **三次 RPC 全部 fail-closed**（503），无降级旁路 | 安全边界一致性优先；cs-user 多实例 + 监控配套 |
| 3 | `Update()` 收敛方式 | **有源字段每次登录覆盖，本地自有字段不动** | 修复零值覆盖 bug；本地查库判定不因陈旧数据误判 |
| 4 | 租户映射 | **维持 default**，不发 `X-Tenant-Id` | 与单租户现状一致；将来需要时加请求头即可 |
| 5 | `GetUserInfo` 退役边界 | **彻底删除**（含非验证解码，不留 dev 旁路）；EmployeeNumber 停止同步、本地列保留 | 与 fail-closed 原则一致；dev 用 docker-compose 起 cs-user |
| 6（2026-08-18 续） | userinfo 读策略 | **本地列 = 软 TTL 缓存**：登录时写入为主刷新，`identity_synced_at` 过期后经 `profile` + `auth-identities` 惰性刷新，刷新失败返回旧值 + 告警（stale-while-revalidate） | fail-closed 只覆盖登录；稳态纯读不应扩大 cs-user 故障半径。cs-user 用户字段 user-owned，profile 返回自编辑后的值，比 claims 更权威；cs-user 自身 syncInterval 为同模式先例 |
| 7（2026-08-18 修订） | 绑定流程归属 | **整体下线，不做迁移**：`/manager/bind/account` 两端点、`MergeByCasdoor` / `MergeUserQuota` 随绑定流程退役（§2.4） | 原「保留但瘦身」修正：「不可委托给 cs-user」仍成立（IdP 层合并 vs identity 行操作，层不同），但能力本身可退役——costrict-web 无调用方、server 侧 bind 已在 cs-user 层承接统一、quota 合并休眠（baseURL 空）。残留本地行分裂与未绑定用户现状一致，Phase 3 迁锚点后自然消解 |
| 8（2026-08-18 续） | 邀请码判新的数据依赖 | **维持本地查库不变**，但其依赖仅为登录时写路径（决策 3 保证新鲜），与 userinfo 读路径解耦，TTL 设计互不影响 | 按 github_id/phone/email 判新是「按人判新」产品语义；cs-user `is_new_user` 在迁移窗口内不可用（存量用户均为首见）。**（已失效：2026-08-21 决策 9 改判为「保留但减负」，判新改用 cs-user `is_new_user`，本地查库判新删除，迁移窗口语义变化接受为产品语义，见 §6.1）** |
| 9（2026-08-21 修订） | 邀请码机制 | **保留但减负**（§3.7）：`parseInviterCodeFromState`、web `inviter_code` query、`GET /invite-code`、`GenerateInviteCode` / `ValidateInviteCode` 全部保留；判新改用 cs-user `get-or-create` 响应的 `is_new_user`——`is_new_user=true` 且 `inviterCode != ""` 时 `ValidateInviteCode` 并写 `InviterID`，否则忽略邀请码正常登录（「you have registered」拒绝删除）；`InviteCode / InviterID` 列继续读写（不 DROP） | 原「整体移除」改判：邀请码是获客能力，保留但把判新职责外移给 cs-user，符合「身份职责外移」方向；登录路径查库从最多 4 次降为 1 次；§6.1 迁移窗口语义变化（存量老用户首登 `is_new_user=true` 可绑定邀请码）接受为产品语义；外部协议不变，插件端 / 管理页无需改动 |
| 10（2026-08-21） | GitHub star 同步 | **整体移除**（§3.7）：`sync_star` 包、`isStar` / `GithubStar`、`syncStar` 配置段、`StarUser` / `SyncLock` 模型全部下线；`GithubStar` 列停读写（不 DROP）；**下游直读依赖同步处置**——quota-manager `GetUserQuota` 对 `users.GithubStar` 列的 is_star 计算、credit-manager 前端门控一并移除（§3.7） | oidc-auth 收窄为「device token 轮询 + Login flow」；star 态是 oidc-auth 本地产品字段，非身份事实，不属于「以 cs-user 为主」范围；userinfo 的 `isStar` 字段停发；quota-manager 直读库依赖需三端同发，不能依赖 `enabled:false` 配置 |

## 10. 任务执行优先级（2026-08-21 落定；同日修订：登录邀请机制与 userinfo 获取提前为首要目标）

本阶段目标收敛为两项可先交付的能力，均落在「三连 RPC 登录链」底座之上，且 cs-user 侧已确认零改动可用：

- **目标 ①：登录邀请机制**（决策 9 修订）——邀请码保留、判新改用 cs-user `is_new_user`；
- **目标 ②：userinfo 获取**（决策 6）——本地列软 TTL 缓存 + 惰性刷新。

已核实（2026-08-21）：cs-user `POST /api/internal/users/get-or-create` 响应带 `is_new_user`（仅本次实际新建行时 true，契约与 §2.3 ② 一致）；`GetUserInfo` 仅被两个登录 handler 调用（login_handler.go:227、web_handler.go:174），star 同步不依赖它——故目标 ①② 可脱离 star 移除 / 绑定下线独立交付，A 类纯删减降级为延后批次（P4）。

### P0 前置准备（串行起点，无代码风险）

- P0-1 **cs-user 就绪**：`CS_USER_CASDOOR_JWKS_URL` 已配置、多实例部署；联调环境验证 §6.5 语义（软删用户经 get-or-create 的行为、`explicitly_unbound` 拒绝路径）。**已核实 get-or-create 返回 `is_new_user`，零改动可用。**
- P0-2 **测试基建**（§6.6，最大执行风险）：先为 `usercenter` 客户端（§3.2）与 `GetUserByOauth` 组装逻辑补 httptest 桩单测（cs-user 接口用桩模拟），再动主链路；用例清单见 §6.7（U / L / I / F 组单测先行，联调作验收门）。
- P0-3 **userCenter 配置段**（§3.1）先行落地，`config/config.go` + `config.yaml` 加 `userCenter` 段（此时无消费方，纯新增）。

### P1 登录链底座（目标 ①② 的共同依赖，串行主线）

- P1-1 **usercenter 客户端**（§3.2）：四方法 `ParseIdentity` / `GetOrCreate` / `GetProfile` / `ListIdentities` + 配置注入 + 单测。
- P1-2 **三连 RPC 替换身份获取**（§3.3）：`GetUserByOauth`（login_handler.go:216）与 web 侧 `GetWebUserByOauth` 的 `GetUserInfo` 非验证解码替换为 parse-identity → get-or-create → auth-identities 串行调用（§2.1 / §2.2）；`user.ID == uuid.Nil → uuid.New()` 兜底（login_handler.go:243-245）删除。
- P1-3 **GetUserInfo 退役**（决策 5）：`internal/providers/casdoor.go` 删除 `GetUserInfo` / 非验证解码、`interface.go` 删接口方法；EmployeeNumber 停止同步、本地列保留。
- P1-4 **Update 收敛**（决策 3）：`casdoor.go` 的 `Update` 改为有源字段覆盖、`Vip / Company / Location` 本地保留。
- P1-5 **新列**：`AuthUser` 增 `SubjectID` / `IdentitySyncedAt`（§4），AutoMigrate 落地。

### P2 登录邀请机制（决策 9 修订，本阶段目标 ①，依赖 P1）

- P2-1 **判新改造**（login_handler.go:280-317）：删本地判新（github_id / phone / email 三次 `GetUserByField`）与「you have registered」拒绝登录，改消费 get-or-create 响应 `is_new_user`——仅 `is_new_user=true` 且 `inviterCode != ""` 时 `ValidateInviteCode` 并写 `InviterID`，否则忽略邀请码、登录正常完成。
- P2-2 **web 侧同款改造**（§3.7）：`webLoginCallbackHandler`（web_handler.go:93-124）判新 / 写 `InviterID` 同款。
- P2-3 **邀请码存量零改动**：`parseInviterCodeFromState`（login_handler.go:265-278）与调用点（:105、:175-180）、`inviter_code` query、`GET /invite-code`、`GenerateInviteCode` / `ValidateInviteCode`、`InviteCode / InviterID` 列全部保留；`__inviter_code=` 插件端协议不变。

### P3 userinfo 软 TTL（决策 6，本阶段目标 ②，依赖 P1-5）

- P3-1 `userInfoHandler`（manger_handler.go:258）：身份类字段（username / email / phone / githubID / githubName）读本地列 + `identity_synced_at` 软过期判断（stale-while-revalidate），过期后经 `GetProfile` + `ListIdentities` 惰性刷新并回写本地列，刷新失败返回旧值 + 告警；本地产品字段（Vip / isPrivate）不走 TTL，仍读本地列。**`isStar` 输出暂保持现状**（star 同步仍在运行，列与输出维持一致，移除随 P4）。

### P4 A 类纯删减（延后批次，不阻塞目标 ①②；原「先行」降级）

- P4-1 **GitHub star 同步移除**（决策 10）：删 `internal/sync/sync_star.go` 整包、`isStar` / `GithubStar`、`syncStar` 配置段、`StarUser` / `SyncLock` 模型、`github_star` allowedFields 条目、`cmd/main.go` 启动逻辑（§8 对应行）；`GithubStar` 列停读写不 DROP；`userInfoHandler` 的 `isStar` 计算与输出一并删除（manger_handler.go:277-283、:294）。
- P4-2 **绑定流程下线**（决策 7）：删 `/manager/bind/account` 两端点、`bindAccount` / `bindAccountCallback`、`MergeByCasdoor` / `MergeUserQuota`、对应路由与常量（§8）。

> 依赖注意：P4 的**发布**被 P5 三端同发阻塞；代码可先行，上线排 P5 之后。

### P5 下游三端协调（依赖 P4 代码 + P1）

- P5-1 **quota-manager**：删 `GetUserQuota` 对 `users.GithubStar` 列的 is_star 计算（§3.7 直读库依赖）。
- P5-2 **credit-manager**：删 `isStar` / `transferData` 门控（与 P4-2 绑定下线协同）；邀请码相关**保留不动**。
- P5-3 **插件端 / 管理页确认**：确认无 `isStar` 硬编码引用（§7-1，未确认可先 410 过渡一个版本）；`employee_number` 字段对齐；`__inviter_code=` / `inviter_code` query / `/invite-code` 因邀请码保留，协议不变，无需确认。

### P6 发布与回滚（§7）

- P6-1 **目标 ①② 可先行单独发布**（带 `userCenter` 配置，缺失启动即报错；star 同步不依赖新登录链，独立运行，行为不变）；is_star 移除（P4/P5）三端同发。
- P6-2 观察：登录成功率、三次 RPC 的 401/503/超时比例、本地 Upsert 冲突、TTL 刷新失败告警、`isStar` 缺失反馈。
- P6-3 回滚：回滚 oidc-auth 版本 + 移除 `userCenter` 配置即可——本地库仅**新增**了 `subject_id` / `identity_synced_at` 两列（旧代码不读不写，向后兼容）；cs-user 侧已补建的账号与 identity 行留存无副作用；star 移除不涉及数据迁移，随版本整体恢复；邀请码判新改造随登录链路一起回滚。

### 依赖链速览

```
P0（cs-user 就绪 + 测试基建 + 配置段）
 └─→ P1（登录链底座：三连 RPC + 新列）─→ P2（登录邀请机制 ①）
     └─→ P3（userinfo 软 TTL ②）          目标①② 可先行单独发布
P4（A 类纯删减，延后）──→ P5（下游三端同发）──→ P6（发布与回滚）
```

### 待拍板项

1. **目标 ①② 是否开始实施**——需 P0-2 测试基建先行，P0-1 cs-user 就绪确认。
2. **P4 A 类纯删减**（star 移除 + 绑定下线）已降级为延后批次，是否仍保留在本文档范围内。
3. 其他协调：插件端 / 管理页确认 `isStar`、credit-manager 前端同步、`employee_number` 字段对齐、credit-reward-plan 页面去留（需运营确认）。
