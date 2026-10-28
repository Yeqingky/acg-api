# Changelog

## 2026-10-28

### Added

- 新增 20 张横图、47 张竖图。

### Changed

- GitHub Actions 工作流忽略仅改动 `data` 目录或 Markdown 文件的提交, tag 推送仍正常触发发布.

## 2026-10-22

### Changed

- 使用 Go 重写, 移除 PHP 与 PHP-FPM 依赖, 部署改为单个静态二进制与多阶段构建镜像.
- 接口路径改为 `api`, `pc`, `pe`, `test`, 保留 `return=json` 与 `return=all` 参数.
- 新增 `-php-compat` 参数, 开启后额外注册 `api.php`, `pc.php`, `pe.php`, `test.php`, 默认关闭.
- 图片列表移入 `data` 目录, 由 `-pc` 与 `-pe` 指定路径, 容器内通过 `./data:/app/data:ro` 挂载.
- CORS 与 `OPTIONS` 由服务自身处理, Nginx 改为纯反向代理.
- 图片列表文件按 `mtime` 与文件大小自动重载, 更新后无需重启容器.
- 图片列表中的空行被忽略.
- 工具链升级至 Go 1.27.3, `golang.org/x/image` 升级至 v0.47.0.
- 默认监听端口改为 8080.
- `docker-compose.yml` 增加注释形式的 `-php-compat` 示例.
- 新增 GitHub Actions 工作流: 提交与 PR 构建 `linux/amd64` 与 `linux/arm64` 镜像并使用短 hash 标签, tag 构建 `vX.X.X` 与 `latest` 镜像, 并发布 Windows 与 Linux 的 amd64/arm64 二进制到 Release, 资产为未压缩可执行文件.

### Fixed

- `return=json` 的图片尺寸查询增加 10 秒超时, 慢速 CDN 不再无限阻塞请求.
- 修复 `GO-2026-6809` 证书校验, `GO-2026-6222` VP8L 解码内存分配, `GO-2026-5061` 与 `GO-2026-4961` webp 解码 panic.

## 2026-10-12

### Changed

- 图片链接由 `list.yppp.net` 更换为 `cdn.yppp.net/acg/`。
- 图片统一为 webp 格式,横图 420 张使用 `pc/`,竖图 929 张使用 `pe/`。
- 调整横竖图归类:`pc/60.webp` 与 `pe/595.webp` 互换内容,原 `pe/798.webp` 删除,原 `pe/930.webp` 改名补入 `pe/798.webp`。

## 2026-06-20

### Added

- 新增 19 张横图、39 张竖图。

## 2026-04-14

### Added

- 新增 8 张横图、51 张竖图。
- 新增 7 张横图、61 张竖图。

## 2026-04-03

### Added

- 新增 24 张横图、29 张竖图。

## 2026-03-22

### Added

- 新增 10 张横图、49 张竖图。

## 2026-03-15

### Changed

- 更新项目文档，添加可用性监测徽章。

## 2026-02-15

### Added

- 新增 15 张横图、66 张竖图。

## 2026-01-21

### Removed

- 移除部分轻 NSFW 图片。

## 2026-01-13

### Added

- 新增 8 张横图、43 张竖图。

## 2025-12-27

### Added

- 新增 8 张横图、24 张竖图。

## 2025-12-20

### Added

- 新增 20 张竖图。

## 2025-12-19

### Added

- 添加更新日志。