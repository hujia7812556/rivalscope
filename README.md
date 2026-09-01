# RivalScope 竞品观察

查看竞品社媒账号数据的内部工具。React 前端构建产物通过 `go:embed` 嵌入 Go 二进制,**单可执行文件 + 单配置文件**即可部署,无需 nginx。

由旧项目 `web-tool-front`(Vue3)+ `web-tool-go`(Nunu 脚手架)重写合并而来。

## 功能

- **统计结果**:按微博账号聚合的转发/评论/点赞 SUM/AVG/MAX/MIN 统计,支持账号与日期范围筛选、前端排序、导出 Excel
- **详细数据**:微博文明细列表,服务端分页 + 服务端排序
- **账号管理**:微博竞品账号的增删改(数据由外部爬虫写入 `weibo_account` / `weibo_feed` 表)
- **登录**:用户名 + 密码(JWT);账号直接配置在 `config.yaml` 的 `auth.users`,无用户表、无注册

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite + Ant Design 5 + TanStack Query v5 + React Router 6 |
| 后端 | Go + Gin + GORM(PostgreSQL,可切 MySQL/SQLite)+ Viper + Zap + golang-jwt v5 |
| 部署 | `go:embed` 嵌入前端,单二进制部署 |

## 快速开始

### 1. 准备配置

```bash
cp config/config.example.yaml config/config.yaml
# 编辑 config.yaml:
#   - auth.users:登录账号(用户名/密码/昵称),需要几个人就写几条
#   - 数据库 DSN、JWT 密钥等
```

密码支持明文或 bcrypt hash(`$2a$` 等前缀自动识别,推荐 hash):

```bash
htpasswd -bnBC 10 "" '你的密码' | tr -d ':\n'
```

### 2. 本地开发

```bash
make dev-api    # 终端 1:后端 127.0.0.1:17317
make dev-web    # 终端 2:前端 127.0.0.1:5173,/api 代理到后端
```

浏览器访问 http://127.0.0.1:5173

### 3. 构建与部署

```bash
make build      # 构建前端 → 嵌入 → bin/rivalscope
./bin/rivalscope -config config/config.yaml   # 单文件启动,浏览器访问 http://127.0.0.1:17317
```

对外 HTTPS 部署时用 nginx 反向代理(配置见 `deploy/nginx/rivalscope.conf`):前端与 API 已同源,nginx 全量反代 `127.0.0.1:17317` 即可,无需静态目录与路径重写;Go 服务监听 `127.0.0.1` 不直接暴露公网。

### 4. 服务器进程管理(systemd)

服务模板见 `deploy/systemd/rivalscope.service`(含 `{{DEPLOY_PATH}}`/`{{DEPLOY_USER}}` 占位符,CI 部署时自动渲染)。手动部署时先渲染再安装:

```bash
sed -e 's|{{DEPLOY_PATH}}|/实际/部署/路径|g' -e 's|{{DEPLOY_USER}}|实际用户|g' \
    deploy/systemd/rivalscope.service | sudo tee /etc/systemd/system/rivalscope.service >/dev/null
sudo systemctl daemon-reload && sudo systemctl enable --now rivalscope

# 日常操作
systemctl status rivalscope        # 状态
systemctl restart rivalscope       # 重启(更新二进制/配置后)
journalctl -u rivalscope -f        # 跟踪日志
```

nginx 模板同理(占位符 `{{NGINX_SERVER_NAME}}`,见 `deploy/nginx/rivalscope.conf`)。推荐直接使用 GitHub Actions 自动部署(见下节),模板渲染无需手工参与。

### 5. 版本发布与自动部署(GitHub Actions)

见 `.github/workflows/release.yml`,三种触发方式:

| 触发 | 版本号 | 类型 | 部署 |
|---|---|---|---|
| push 到 `main` | 自动 `v1.0.<run_number>` | 预发布 | 自动部署 |
| 推送 `v*` tag(如 `v1.0.0`) | 即 tag 名 | 正式版 | 自动部署 |
| Actions 页面手动触发 | 可输入(留空自动生成) | 正式版 | 自动部署 |

```bash
git tag v1.0.0 && git push origin v1.0.0   # 正式发布
git push origin main                        # 日常合并自动发布+部署
```

产物:`linux/amd64` 单平台二进制(前端已嵌入、静态编译,对应部署服务器架构),附 sha256 校验;本机 macOS 调试用 `make build` 自行构建。

**部署参数在 GitHub 仓库配置**(改配置不用改代码)。**注意:仓库是 public 的,全部参数放 Secrets**(Secrets 加密存储且日志自动脱敏;Variables 对 public 仓库可被读取、且会在 workflow 日志中明文展示):

- Secrets:`DEPLOY_SSH_KEY`(服务器 SSH 私钥)、`DEPLOY_HOST`(服务器地址)、`DEPLOY_USER`(SSH 用户)、`DEPLOY_PATH`(部署目录)、`NGINX_SERVER_NAME`(对外域名)、`DEPLOY_CONFIG`(完整的 config.yaml 内容,多行)

未配置 `DEPLOY_HOST` 时只构建发布、跳过部署。部署内容:替换 `bin/rivalscope`、安装 systemd 服务与 nginx 配置(均按上述 Secrets 渲染,**每次部署覆盖**)、重启并做健康检查。`config.yaml` 由 `DEPLOY_CONFIG` 控制:配置了则每次覆盖写入(旧配置自动备份为 `config/config.yaml.bak`);未配置则仅首次部署从模板生成,之后不覆盖。服务器一次性前提:部署用户具备 `systemctl`/`nginx` 相关命令的 sudo 免密,DNS 已解析,证书已签发。

## 目录结构

```
├── main.go               # 入口:服务启动 + useradd 子命令
├── webembed.go           # go:embed 前端 + SPA 路由回退
├── internal/
│   ├── auth/             # JWT 签发与解析
│   ├── config/           # 配置加载(Viper)
│   ├── dto/              # 请求/响应结构 + 统一响应封装
│   ├── handler/          # HTTP 处理层
│   ├── log/              # Zap 日志(文件轮转)
│   ├── middleware/       # JWT 鉴权中间件
│   ├── model/            # GORM 模型
│   ├── repo/             # 数据访问层
│   ├── server/           # 路由组装 + HTTP 启动
│   └── service/          # 业务逻辑层
├── web/                  # React 前端源码(pnpm)
└── config/               # 配置文件(example 提交,真实配置 gitignore)
```

## 接口清单

统一前缀 `/api/v1`,响应格式 `{code, message, data}`(成功 code=0)。除登录外均需 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | /api/v1/login | 登录,返回 accessToken |
| GET | /api/v1/user | 当前用户信息 |
| GET | /api/v1/weibo/accounts | 微博账号列表(status=1) |
| POST | /api/v1/weibo/accounts | 新增账号(name/uid 必填) |
| PATCH | /api/v1/weibo/accounts/:id | 更新账号 |
| DELETE | /api/v1/weibo/accounts/:id | 删除账号 |
| GET | /api/v1/weibo/messages | 微博文列表(服务端分页/排序) |
| GET | /api/v1/weibo/statistics | 统计聚合(按账号,不分页) |

**查询参数**(messages/statistics 通用):`account_id`、`start_date`、`end_date`(YYYY-MM-DD,含当天);
messages 另支持 `sort_field`(白名单:id/mid/forward/comment/like/pubtime)、`sort_order`(asc/desc)、`page`、`page_size`。

## 数据库迁移(goose)

使用 [goose](https://github.com/pressly/goose) 管理迁移,迁移 SQL 通过 `go:embed` 打进二进制,**部署时无需安装额外工具**。

- **自动迁移**:服务启动时自动执行 `up`(幂等),见 `internal/migration/`
- **手动管理**:`./bin/rivalscope migrate up|down|status|version -config config/config.yaml`
- **迁移文件**:`internal/migration/migrations/{postgres,mysql,sqlite}/`,按部署方言各维护一套(本地测试用 PostgreSQL,线上旧库为 MySQL)

### 与线上已有数据的兼容约定

线上库中 `weibo_account` / `weibo_feed` 由爬虫程序创建并持续写入,**本项目只读消费、严禁为其编写迁移**。RivalScope 当前没有自有表(登录账号在配置文件里,不占数据库),goose 机制保留待用:未来新增自有表(如微信模块)时,在三个方言目录追加迁移 SQL 即可。goose 版本记录独立存放于 `goose_db_version` 表,首次部署线上不会影响任何已有数据。

### 新增迁移的规范

1. 在三个方言目录新增 `NNNNN_描述.sql`(版本号递增,保持三套同步);
2. 必须同时写 `-- +goose Up` 与 `-- +goose Down`;
3. **禁止修改已发布的历史迁移文件**,只允许追加;
4. 涉及已有数据的表结构变更时,优先用 `IF NOT EXISTS` / 幂等写法,保证可重入。

## 数据依赖

`weibo_account`、`weibo_feed` 两张表由外部爬虫程序(旧 jennycrawl 体系)写入,本项目只读消费(账号管理保留增删改);本项目不建任何自有表。统计口径:按账号聚合,INNER JOIN,时间范围内无微博的账号不出现。

## 后续扩展(微信等平台)

后端按模块组织(model/repo/service/handler 各一个文件),新增平台时复制微博模块模式,路由加 `wechat` 分组;前端在侧边菜单与 `api/` 中按同样方式扩展。
