# MoeFile

MoeFile is a simple self-hosted file listing service. As an replacement to `index` options on common HTTP server, it does not targeting to be a full-featured file server, but a better-looking file listing page.

> **Note**: See chinese version of this document at [README.zh-Hans.md](README.zh-Hans.md).

MoeFile highlights on
 - I18n, multiple timezone, dark mode, responsive design, better UI/UX.
 - File sort and search, with file type icons.
 - Player: play video files with subtitles and danmaku support.
 - S3-compatible `ListBucketResult` output, but both machine-readable and human-readable.
 - Go & React wrttien, portable, without environment setup or dependencies.
 - No CDN or 3rd party resources.
 - Docker ready, with pre-built images.
 - Stateless & zero configuration, no database, no dashboard, fully read-only.

MoeFile will never be:
 - A multi-user file sharing or listing service.
 - A file manager or explorer.
 - A net disk redirector or proxy.
 - A full-functional HTTP server.

## Usage
To get started, download the latest release from [GitHub Releases](https://github.com/baobao1270/moefile/releases) page, and extract the archive. Open the terminal to the directory where the binary is located, and run the following command:

```bash
./moefile
```

On Windows, please replace `./moefile` with `moefile.exe`.

You can type `./moefile -h` to see all available options. By default, MoeFile will listen on port `3328`.

> **Warning**: On public server, it is recommended to set the `-origins` and `-proxies` to configure CORS and reverse proxy for better security.

### Run with Docker
You can also run MoeFile with Docker. The following command will start MoeFile with the default configuration:

```bash
docker run -d -p 3328:3328 -v $(pwd):/data --name moefile baobao1270/moefile
```

> **Note**: This command will expose the current directory to the container, which is, the root directory of the file listing. You can replace `$(pwd)` with the path to the directory you want to expose.

### Tags
MoeFile docker image follows the tag rules following:
 - `latest`: The latest stable release.
 - `dev`: The latest commit on the `main` Git branch.
 - `v<version>`: The specific version of MoeFile related to a specific Git tag.
 - `<commit-sha>`: The CI-built image of a specific commit.

### Configuration
**Environment Variables**
| Environment Variable | Default Value     | Equivalent Flag | Description                                                    |
| -------------------- | ----------------- | --------------- | -------------------------------------------------------------- |
| `LEVEL`              | `inf`             | `-level`        | The log level.                                                 |
| `LISTEN`             | `:3328`           | `-listen`       | The address and port to listen on.                             |
| `ORIGINS`            | `*`               | `-origins`      | The allowed origins for CORS, separated by comma.              |
| `PROXIES`            | `127.0.0.1`       | `-proxies`      | The trusted proxies CIDR, separated by comma.                  |
| `ROOT`               | `/data`           | `-root`         | The root directory (in container) to serve listing service on. |
| `SERVER`             | `MoeFile`         | `-server`       | The server name or page title.                                 |
| `XMLTAB`             | `true`            | `-xmltab`       | Whether to add tab space in XML output.                        |
| `TZ`                 | (server)          | N/A             | The timezone to use and shown as _Server Time_ on web page.    |

**Volumes**
 - `/data`: The root directory to serve listing service on.

**Ports**
 - `3328`: The port to listen on. You can change it by setting `LISTEN` environment variable.

### Run with Docker Compose
You can also use Docker Compose to run MoeFile. Here is an example `compose.yaml` file:

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

## Build & Development
To build or start developing MoeFile, you need dependencies following:
 - [Bun](https://bun.sh) v1.x
 - [Go](https://golang.org) 1.23 or later
 - GNU Make
 - [Docker](https://www.docker.com) _(optional, for build docker image only)_

> **Note**: Build is **on POSIX-like system ONLY**, such as Linux, macOS, or WSL. Building on Windows is not supported.

You can obtain the source code by cloning the repository:

```bash
git clone https://github.com/baobao1270/moefile.git
cd moefile
```

Then, install the dependencies and build the project:

```bash
bun install
bun run build
```

### Build Commands

Following is a table explaining all build and development related commands:

| Command                  | Description                                                                            |
| ------------------------ | -------------------------------------------------------------------------------------- |
| `bun install`            | Install dependencies for both frontend and backend.                                    |
| `bun run dev:frontend`   | Start the development server for frontend.                                             |
| `bun run dev:backend`    | Run the go program at debug mode with existing frontend build in code tree.            |
| `bun run dev`            | Build the frontend, then start the backend in debug mode.                              |
| `bun run build:frontend` | Build the frontend in production mode.                                                 |
| `bun run build:backend`  | Build the frontend and backend (native) in production mode.                            |
| `bun run build`          | Build the frontend and cross-compile the backend in production mode for all platforms. |
| `bun run tmplgen`        | Generate fake testing data for development, or real HTML template for production.      |
| `bun run clean`          | Clean up the build directory.                                                          |

You can found these commands in `package.json` file.

Frontend build output will be in `dist/` directory, and backend build output will be in `bin/` directory.

### Build Flags & Customization

MoeFile gives you the ability to customize for your branding. You can (and allowed by license to) change the display or behavior of MoeFile by setting environment variables.

For example, if you want to change the banner copyright holder name, you can set the `COPYRIGHT_HOLDER` environment variable to your name.

```bash
export COPYRIGHT_HOLDER="Some Company"
bun run build
```

Following is a table of environment variables that can be used to customize MoeFile:

| Flag              | Default       | Frontend | Backend | Runtime | Notes                               |
|-------------------|---------------|----------|---------|---------|-------------------------------------|
| `APP_NAME`        | `MoeFile`     | Y        | Y       | ?       | Change by `-server` flag in runtime |
| `APP_VERSION`     | `DEV`         | Y        | Y       | N       | from git                            |
| `APP_AUTHOR`      | `MoeFile`     | N        | N       | N       | const in code                       |
| `APP_COPYRIGHT`   | (year)        | N        | N       | N       | const in code                       |
| `APP_LICENSE`     | `MIT`         | N        | N       | N       | const in code                       |
| `COPYRIGHT_HOLDER`| `MoeFile`     | Y        | N       | N       | change by rebuild                   |
| `BUILD_TIMESTAMP` | (now)         | Y        | Y       | N       | set by build script                 |
| `BUILD_MODE`      | `production`  | N        | Y       | N       |                                     |
| `NODE_ENV`        | `=BUILD_MODE` | Y        | N       | N       |                                     |
| `TZ`              | (server)      | N        | N*      | Y       | runtime only, set by env            |

**Notes**:
 1. Frontend: stands for the environment variable is used in frontend (React & Vite) build process.
 2. Backend: stands for the environment variable is used in backend build (Go) process.
 3. Runtime: stands for the environment variable is used in runtime (by command line).
 4. _(*)_ In development mode, the backend will use the TZ environment variable to set the timezone. This does not affect the production build.

**If you are using docker**, you also need to pass `--build-arg` to set the same environment variable when building the image.

### Code Style: A Note of Spelling
The word _**Danmaku**_ or _**Danmuku**__ is a transliteration of the Japanese word **「弾幕」**, which has two variants of spellings. However, different upstream libraries and projects may use different spellings.

Here we set the rule of spelling for code clarity and readability:

| Condition                | Spelling                                                 |
|--------------------------|----------------------------------------------------------|
| Frontend code            | `danmuku` (follow upstream library)                      |
| I18n key                 | `danmuku` (follow code)                                  |
| I18n English Translation | `danmaku` (follow common usage)                          |
| Backend code             | `danmaku`                                                |
| API                      | backend provides `danmaku`, frontend should support both |

## License
MoeFile is licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.
