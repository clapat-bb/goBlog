# GoBlog 📝 – 使用 Go + Gin 实现的博客系统

GoBlog 是一个使用 Golang 编写的全栈式博客系统，支持文章发布、评论、点赞、标签管理，采用 PostgreSQL 作为数据库，Redis 作为缓存系统，并已通过 Docker 容器化部署，支持一键启动开发环境。

---

## 🚀 技术栈

- **Golang** 1.21 + Gin 框架
- **PostgreSQL** 15
- **Redis** 7
- **GORM** ORM
- **JWT** 用户认证
- **Swagger** API 文档
- **Docker + docker-compose** 一键部署
- **pgcli**（推荐）数据库交互工具

---

## 📦 功能模块

- 🧑 用户注册、登录（JWT）
- 📄 文章发布、更新、删除、置顶、推荐
- 💬 评论（支持子评论结构）
- ❤️ 点赞系统（支持取消）
- 🏷️ 标签系统（多对多关联）
- 🔍 文章分页、标签筛选
- ⚡ Redis 缓存加速：文章列表、点赞计数等
- 📃 Swagger UI 接口文档

---

## 🐳 快速开始（Docker 启动）

### 1️⃣ 克隆项目

```bash
git clone https://github.com/clapat-bb/goblog.git
cd goblog
````

### 2️⃣ 构建 & 启动容器

```bash
docker-compose up --build
```

启动后会运行：

* 应用服务（端口 `8080`）
* PostgreSQL 数据库（端口 `5432`）
* Redis 缓存（端口 `6379`）

---

## 📚 API 文档

浏览器访问：

```
http://localhost:8080/swagger/index.html
```

你将看到自动生成的 Swagger API 文档，支持接口测试。

---

## ⚙️ 数据库配置（自动迁移）

应用会自动从 `.env` 读取配置，连接 PostgreSQL 并自动创建表结构。

你可以使用 `pgcli` 测试连接：

```bash
PGPASSWORD=postgres pgcli -h localhost -p 5432 -U postgres -d goblog
```

默认数据库配置如下：

```
host=localhost
port=5432
user=postgres
password=postgres
dbname=goblog
```

---

## 🌍 目录结构

```plaintext
goblog/
├── controllers/        # 控制器层（处理路由请求）
├── models/             # GORM 数据模型
├── database/           # 数据库连接逻辑
├── middlewares/        # JWT 等中间件
├── routes/             # 路由注册
├── config/             # 配置加载（支持 .env）
├── pkg/cache/          # Redis 缓存封装
├── docs/               # Swagger 文档
├── main.go             # 应用入口
├── Dockerfile          # 应用构建镜像配置
├── docker-compose.yml  # 一键部署数据库 + Redis + 应用
├── .env                # 环境变量配置（别提交到生产）
```

---

## ✅ 已完成特性

* [x] 用户模块
* [x] 文章模块（含分页）
* [x] 评论系统（含子评论）
* [x] 点赞模块
* [x] 标签管理
* [x] Redis 缓存加速（文章列表、详情、点赞）
* [x] Swagger API 文档
* [x] Docker 一键部署

---

## 🛠 TODO（可扩展）

* [ ] 后台管理界面（前端）
* [ ] 图片上传（集成 OSS 或本地存储）
* [ ] 消息通知（如新评论、点赞）
* [ ] WebSocket 推送
* [ ] 单元测试 + 接口测试（Ginkgo / Postman）

---

## 🧑‍💻 作者

* GitHub: [@clapat-bb](https://github.com/clapat-bb)
```
独立开发基于 Golang 的博客系统，支持文章、评论、标签、点赞等功能，后端使用 Gin + PostgreSQL + Redis + JWT，支持 Swagger 文档与 Docker 一键部署，前后端接口解耦，具备良好的代码结构与扩展能力。
```

---


```
go build. go blog. go big.
```
