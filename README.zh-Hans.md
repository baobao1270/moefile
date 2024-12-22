# MoeFile

MoeFile 是一个简单的自托管文件列表服务。其目的不是成为一个功能齐全的文件服务器，而是各种 Web 服务器内置 `index` 功能的更加美观的替代。

> **注意**: English version of this document is available [here](README.md).

> **注意**: 此文档使用 ChatGPT 翻译，部分内容可能不符合您所在地区的语言习惯，敬请谅解。

MoeFile 支持：
- I18n、多时区支持、Dark Mode、响应式设计、文件类型图标
- 文件排序、文件搜索
- 视频播放器：支持多语言字幕和弹幕
- 与 S3 兼容的 `ListBucketResult` 输出，同时提供机读和人类可读的格式。
- 使用 Go 和 React 编写，无需复杂的环境配置和依赖，可快速部署。
- 不使用任何 CDN 或第三方资源。
- 支持 Docker，提供预构建镜像。
- 无状态、零配置、完全只读，无需数据库等配套设施。

MoeFile 不是：
- 类似 CloudReve 的多人网盘或者文件共享服务
- 类似 ZFile 的文件管理器
- 类似 AList 的网盘理
- 类似 Nginx 的功能完整的 HTTP 服务器

## 使用方法
您可以从 [GitHub Releases](https://github.com/baobao1270/moefile/releases) 下载 MoeFile 的最新版本。解压后，打开文件所在目录，然后运行以下命令：

```bash
./moefile
```

在 Windows 上，请将 `./moefile` 替换为 `moefile.exe`。

您可以键入 `./moefile -h` 查看帮助。默认情况下，MoeFile 将监听端口 `3328`。

> **提醒**: `3328 = 0x0d00` - Ciallo～(∠・ω< )⌒★

> **警告**: 如果您的服务器对公网暴露，建议设置 `-origins` 和 `-proxies` 来配置 CORS 和反向代理，以提高安全性。

### 使用 Docker 运行
您也可以使用 Docker 运行 MoeFile。以下命令将以默认配置启动 MoeFile：

```bash
docker run -d -p 3328:3328 -v $(pwd):/data --name moefile baobao1270/moefile
```

> **注意**: 此命令会将当前目录暴露在 3328 端口上。您需要将 `$(pwd)` 替换为您想要公开的目录的路径。

### Tags
MoeFile Docker 镜像遵循以下 Tag 规则：
- `latest`: 最新的稳定版本
- `dev`: Git 仓库 `main` 分支上的最新提交
- `v<版本>`: 与 Git tag 关联的特定版本的
- `<commit-sha>`: 从该 commit 构建的 CI 镜像

### 配置
**环境变量**
| 环境变量         | 默认值       | 等效命令行    | 描述 |
| --------------- | ----------- | ----------- | --- |
| `LEVEL`         | `inf`       | `-level`    | 日志级别 |
| `LISTEN`        | `:3328`     | `-listen`   | HTTP 监听地址和端口 |
| `ORIGINS`       | `*`         | `-origins`  | CORS 允许的 origin (逗号分隔) |
| `PROXIES`       | `127.0.0.1` | `-proxies`  | 可信的反向代理 CIDR (逗号分隔) |
| `ROOT`          | `/data`     | `-root`     | 服务器根目录 |
| `SERVER`        | `MoeFile`   | `-server`   | 服务器名 (用于显示页面标题) |
| `XMLTAB`        | `true`      | `-xmltab`   | XML 输出时是否加上缩进 |
| `TZ`            | (server)    | N/A         | 服务器时区，用于在客户端进行按时间排序 |

**卷**
- `/data`: 服务器根目录

**端口**
- `3328`: 要监听的端口，通过设置 `LISTEN` 环境变量更改

### 使用 Docker Compose 运行
您也可以使用 Docker Compose 运行 MoeFile。以下是一个示例 `compose.yaml` 文件：

```yaml
services:
  app:
    image: baoabao1270/moefile
    restart: unless-stopped
    container_name: moefile
    environment:
      - TZ=Asia/Hong_Kong
      - SERVER=MyFile
    volumes:
      - ./data:/data
    ports:
      - 127.0.0.1:3328:3328/tcp
```

## 构建和开发
要构建或开始开发 MoeFile，您需要以下依赖项：
- [Bun](https://bun.sh) v1.x
- [Go](https://golang.org) 1.23 或更高版本
- GNU Make
- [Docker](https://www.docker.com) _（可选，仅用于构建 Docker 镜像）_

> **注意**: 构建**仅在类 POSIX 系统上进行**，例如 Linux、macOS 或 WSL。不支持在 Windows 上构建。

首先 Clone 本仓库：

```bash
git clone https://github.com/baobao1270/moefile.git
cd moefile
```

然后，安装依赖并构建项目：

```bash
bun install
bun run build
```

### 相关命令

下表列出了所有与构建和开发相关的命令：

| 命令                      | 描述 |
| ------------------------ | ---- |
| `bun install`            | 安装所有依赖，包括前端和后端的依赖 |
| `bun run dev:frontend`   | 启动前端的本地开发服务器 |
| `bun run dev:backend`    | 运行后端服务器，要求目录树下已经有编译成功的前端 `dist` 目录 |
| `bun run dev`            | 先编译前端，再运行后端服务器 |
| `bun run build:frontend` | 生产模式编译前端 |
| `bun run build:backend`  | 生产模式编译前端和后端 (仅编译当前平台二进制) |
| `bun run build`          | 生产模式编译前端和后端 (编译所有平台二进制) |
| `bun run tmplgen`        | 生成开发用的测试数据，生产模式下为生成 `go embed` 的模板文件 |
| `bun run clean`          | 清理构建结果和 `node_modules` 目录 |

您可以在 `package.json` 文件中找到这些命令。

前端编译结果在 `dist/` 目录中，后端编译结果在 `bin/` 目录中。

### 自定义 MoeFile

MoeFile 支持您自定义显示和程序行为，这样的操作也受许可证允许。

例如，如果您想更改页脚版权持有者名称，您可以将 `COPYRIGHT_HOLDER` 环境变量设置为您的名字。

```bash
export COPYRIGHT_HOLDER="米忽悠二游有限公司"
bun run build
```

下表列出了可用于自定义 MoeFile 的环境变量：

| Flag               | 默认值         | 用于前端 | 用于后端 | 运行时可更改 | 说明                                   |
| ------------------ | ------------- | ------  | ------ | ---------- | -------------------------------------- |
| `APP_NAME`         | `MoeFile`     | Y       | Y      | ?          | 运行时使用 `-server` 命令行参数更改        |
| `APP_VERSION`      | `DEV`         | Y       | Y      | N          | 构建时根据 Git 自动设置                   |
| `APP_AUTHOR`       | `MoeFile`     | N       | N      | N          | 代码中 const, 不可更改                   |
| `APP_COPYRIGHT`    | (year)        | N       | N      | N          | 代码中 const, 不可更改                   |
| `APP_LICENSE`      | `MIT`         | N       | N      | N          | 代码中 const, 不可更改                   |
| `COPYRIGHT_HOLDER` | `MoeFile`     | Y       | N      | N          | 只支持通过自行构建更改                    |
| `BUILD_TIMESTAMP`  | (now)         | Y       | Y      | N          | 构建脚本自动设置                         |
| `BUILD_MODE`       | `production`  | N       | Y      | N          |                                       |
| `NODE_ENV`         | `=BUILD_MODE` | Y       | N      | N          |                                       |
| `TZ`               | (server)      | N       | N*     | Y          | 仅运行时可修改，只能通过环境变量设置        |

> **注释**：_(*)_ 在开发模式下，构建脚本会将 TZ 环境变量设置为 `Asia/Hong_Kong` 来测试多时区功能。这不会影响生产构建。

**如果您使用 Docker**，在构建 Docker 镜像时要通过 `--build-arg` 传递一模一样的参数。

### 代码风格：关于拼写的说明
_**Danmaku**_ 或 _**Danmuku**_ 是日语词 **「弾幕」** 的音译，它有两种拼写变体。但是，不同的类库可能会使用不同的拼写。

为了代码的清晰性和可读性，现规定如下：

| 出现场合 | 拼写 |
| ------- | --- |
| 前端代码 | `danmuku`（遵循上游） |
| 国际化 Key | `danmuku`（遵循前端代码） |
| 国际化英文翻译 | `danmaku`（遵循常用用法） |
| 后端代码 | `danmaku` |
| API | 后端拼写为 `danmaku`，前端应同时支持两种拼写 |

## 许可
MoeFile 采用 MIT 许可证。请参阅 [LICENSE](LICENSE) 获取完整的许可证文本。
