# RivalScope(竞品观察)项目说明

竞品社媒数据内部工具:旧 `web-tool-front`(Vue3)+ `web-tool-go`(Nunu 脚手架,均在 `/Users/hujia/myproject/jenny/` 下)重写合并而来,旧项目仅作参考,不要往里加新功能。

**单仓库单二进制**架构:React 前端(`web/`,pnpm)构建产物经 `//go:embed all:web/dist` 嵌入 Go 二进制(`webembed.go`),一个二进制 + 一份配置即可部署。

## 目录与分层

- 后端:`main.go`(入口 + migrate 子命令)、`internal/{config,log,auth,dto,handler,middleware,model,repo,server,service,migration}`、`webembed.go`
- 分层规则:handler(参数解析/响应)→ service(业务与 DTO 映射)→ repo(GORM 查询);统一响应 `{code,message,data}` 见 `internal/dto/response.go`
- 前端:`web/src/{api,layouts,pages,components,utils}`;请求走 `src/api/client.ts`(axios,baseURL `/api/v1`,拦截器统一弹错与 401 跳登录)

## 常用命令(项目根目录)

```bash
make dev-api     # 后端开发(127.0.0.1:17317,go run)
make dev-web     # 前端开发(5173,/api 代理到 17317)
make build       # pnpm build → go build → bin/rivalscope(完整二进制)
make migrate CMD=status   # goose 迁移(up/down/status/version)
go vet ./...     # 静态检查
```

## 关键约定与坑

- **登录无用户表**:账号写死在 `config.yaml` 的 `auth.users`(用户名+密码,密码支持明文或 `$2a$` bcrypt hash 自动识别);身份在 JWT claims 里自包含,`GET /user` 不查库。
- **数据边界**:`weibo_account`/`weibo_feed` 表归外部爬虫写入,本项目只读,严禁为其写迁移;goose 只管自有表(当前无迁移文件,空目录为预期,`internal/migration` 对空目录做 no-op)。
- **`like` 是保留字**:统计聚合 SQL(`internal/repo/weibo_statistics.go`)中 like 列必须按方言引用(MySQL 反引号,PG/SQLite 双引号),裸列名在两个方言下都会炸。
- **接口前缀统一 `/api/v1`**,排序字段有白名单(`internal/repo/weibo_msg.go` 的 `NormalizeWeiboMsgSort`)。
- 本地开发库:PostgreSQL `localhost:54321/jenny_crawl_test`(配置在 `config/config.yaml`,gitignore;模板 `config.example.yaml`)。
- 前端一律 pnpm(不用 npm);表格密集页面用 AntD 5,搜索区复用 `components/SearchFilter.tsx`(支持 `extra` 插槽)。
- 代码注释、提交信息用中文;函数级注释。

## 部署相关

- `deploy/nginx/rivalscope.conf`、`deploy/systemd/rivalscope.service` 是**模板**,含 `{{DEPLOY_PATH}}`/`{{DEPLOY_USER}}`/`{{NGINX_SERVER_NAME}}` 占位符,由 GitHub Actions 渲染,不能直接 cp 使用。
- CI(`.github/workflows/release.yml`):push main / v* tag / 手动触发均发布并自动部署;部署参数全在 GitHub **Secrets**(仓库为 public,禁用 Variables):`DEPLOY_SSH_KEY`、`DEPLOY_HOST`、`DEPLOY_USER`、`DEPLOY_PATH`、`NGINX_SERVER_NAME`、`DEPLOY_CONFIG`(完整 config.yaml,每次部署覆盖,旧配置备份 `.bak`)。
- 真实部署路径:`/home/ecs-user/myapp/rivalscope`(服务器上)。

## 其他目录(~/myproject/jenny/)

- `web-tool-front`、`web-tool-go`:旧版前后端,已废弃,仅参考表结构与统计口径。
- `competitor_crawl`、`web_tool`、`web_tool_old`:旧 PHP 爬虫体系(jennycrawl,线上 jennycrawl.jerehu.com),与 RivalScope 共享数据库;改动前先确认是否影响线上。
