# 社区系统（C 端 + B 端）

基于计划书实现的社区系统：C 端（用户端）与 B 端（管理后台）用户体系隔离，Vue 3 + Golang + MySQL + Kafka + Docker。

## 技术栈

- **后端**: Golang (Gin)、MySQL 8、Kafka、JWT（C 端 / B 端分离）、Swagger
- **C 端**: Vue 3 + Vite，极光蓝主色（#81D8CF），响应式（PC + H5）
- **B 端**: Vue 3 + Vite + Element Plus
- **部署**: Docker Compose（MySQL、Kafka、Zookeeper、Backend、Nginx）

## 快速启动

### 1. 环境要求

- Docker、Docker Compose
- Node 18+（本地开发前端）
- Go 1.21+（本地开发后端）

### 2. 后端与数据库（Docker 一键）

```bash
# 复制环境变量
cp .env.example .env

# 启动 MySQL、Kafka、后端、Nginx（需先构建前端产物，见下）
docker-compose up -d mysql zookeeper kafka
# 等待 MySQL 就绪后，执行 migrations（MySQL 会自动执行 /docker-entrypoint-initdb.d 下的 01_schema.sql）
# 然后启动后端
docker-compose up -d backend
```

### 3. 前端构建（供 Nginx 挂载）

```bash
# C 端
cd frontend-c && npm install && npm run build

# B 端
cd frontend-admin && npm install && npm run build
```

### 4. 启动 Nginx（挂载前端与上传目录）

```bash
docker-compose up -d nginx
```

### 5. 本地开发（不依赖 Docker 全量）

- **后端**: `cd backend && go run ./cmd/api`（需本地 MySQL：`community` 库，且执行过 `migrations/01_schema.sql`；首次运行会自动插入管理员 admin/admin123）
- **C 端**: `cd frontend-c && npm run dev` → http://localhost:5173
- **B 端**: `cd frontend-admin && npm run dev` → http://localhost:5174

## 访问地址（Docker 全量启动后）

**方式一：Nginx 提供构建后的静态页面**

- **C 端**: http://localhost/c/
- **B 端**: http://localhost/admin/

**方式二：前端开发服务（热更新，`frontend-c` / `frontend-admin` 容器）**

- **C 端**: http://localhost:5173/c/
- **B 端**: http://localhost:5174/admin/

**其他**

- **API**: http://localhost/api/v1 或 http://localhost:8080/api/v1
- **Swagger 文档**: http://localhost:8080/swagger/index.html（直连后端）；若用 Nginx：http://localhost/swagger/index.html
  - 若页面空白或「Failed to load API definition」：确认后端已启动且重启过（`swag init` 后需重新编译/启动）

### Swagger 文档生成

接口注解已写在各 handler 中，如需重新生成文档：

```bash
cd backend
go install github.com/swaggo/swag/cmd/swag@latest   # 首次需安装
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

会更新 `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml`。

## 默认账号

- **C 端**: 手机号任意，验证码固定 `123456`（mock）
- **B 端**: 用户名 `admin`，密码 `admin123`（首次启动由后端 Seed 写入 admins 表）

## 目录结构

```
├── backend/           # Golang API
│   ├── cmd/api/       # 入口
│   ├── internal/      # config, handler/c, handler/admin, model, repository, service, middleware
│   ├── migrations/    # 01_schema.sql（无外键，每列 NOT NULL+默认值+COMMENT）
│   └── docs/          # Swagger（可运行 swag init -g cmd/api/main.go -o docs 重新生成）
├── frontend-c/        # C 端 Vue 3 + Vite，极光蓝风格
├── frontend-admin/    # B 端 Vue 3 + Element Plus
├── nginx/conf.d/      # 反向代理与静态资源
├── docker-compose.yml
└── .env.example
```

## 接口说明

- C 端：`/api/v1/auth/login`、`/api/v1/posts`、`/api/v1/upload`、`/api/v1/users/me/posts`、点赞/评论、通知等；鉴权 Header：`Authorization: Bearer <C端JWT>`
- B 端：`/api/v1/admin/auth/login`、`/api/v1/admin/posts`、`/api/v1/admin/users`、`/api/v1/admin/stats`；鉴权 Header：`Authorization: Bearer <B端JWT>`

详见 Swagger 或计划书中的「五、后端 API 设计要点」。
