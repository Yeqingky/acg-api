## 文档

[https://blog.yeqing.net/acg-api/](https://blog.yeqing.net/acg-api/)

## 图片

QQ群: [959731247](https://qm.qq.com/q/oOJXedC0BG)

进群后查看群文件 有全部图片的压缩包 请不要频繁调用api接口来抓取图片

## 接口

| 路径 | 说明 |
| --- | --- |
| `GET /api` | 按 `User-Agent` 跳转到 `/pc` 或 `/pe` |
| `GET /pc` | 随机横图, 302 重定向到图片地址 |
| `GET /pe` | 随机竖图, 302 重定向到图片地址 |
| `GET /pc?return=json` | 随机横图并返回 JSON, 含 `code`, `acgurl`, `width`, `height`, `size` |
| `GET /pe?return=json` | 随机竖图并返回 JSON, 字段同上 |
| `GET /pc?return=all` | 返回全部横图地址, 每行一个 |
| `GET /pe?return=all` | 返回全部竖图地址, 每行一个 |
| `GET /test` | 返回 `Hello World`, 用于健康检查 |

`width` 与 `height` 由服务实时拉取图片头部解析, 失败时返回 `"0"`.

传入 `-php-compat` 后额外注册 `/api.php`, `/pc.php`, `/pe.php`, `/test.php`, 行为与无后缀路径一致. 该开关默认关闭.

## 部署教程

1. 安装Docker和Nginx

```bash
apt install nginx
curl -fsSL https://get.docker.com | bash -s docker
```

2. Git Clone

```bash
git clone https://github.com/Yeqingky/acg-api.git
```

3. 启用容器

```bash
cd acg-api
docker compose up -d --build
```

图片列表位于 `data/acg_横图.txt` 与 `data/acg_竖图.txt`, 由 `./data:/app/data:ro` 挂载进容器, 不打进镜像.

4. 配置Nginx

这是一份Nginx监听80端口 套CDN的配置文件 不包含HTTPS/SSL

```conf
server {
    listen 80 ;

    # 获取真实IP
    real_ip_recursive on ;
    set_real_ip_from 0.0.0.0/0 ;
    real_ip_header X-Forwarded-For ;

    # 全部请求交给 acg-api
    location / {
        proxy_pass http://127.0.0.1:8080 ;
        proxy_set_header Host $host ;
        proxy_set_header X-Real-IP $remote_addr ;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for ;
        proxy_set_header X-Forwarded-Proto $scheme ;
    }
}
```

## 配置

| 参数 | 默认值 | 说明 |
| --- | --- | --- |
| `-pc` | `data/acg_横图.txt` | 横图地址列表文件路径 |
| `-pe` | `data/acg_竖图.txt` | 竖图地址列表文件路径 |
| `-addr` | `:8080` | 监听地址 |
| `-php-compat` | `false` | 额外注册 `.php` 后缀的旧路径 |

图片列表文件按 `mtime` 与文件大小自动重载, 更新列表后无需重启容器.

## 本地运行

```bash
go run .
```

或构建后运行:

```bash
go build -o acg-api .
./acg-api -pc data/acg_横图.txt -pe data/acg_竖图.txt
```
