## 文档

[https://blog.yeqing.net/acg-api/](https://blog.yeqing.net/acg-api/)

## 图片

QQ群: [959731247](https://qm.qq.com/q/oOJXedC0BG)

进群后查看群文件 有全部图片的压缩包 请不要频繁调用api接口来抓取图片

## 部署教程

1. 安装Docker和Nginx

```bash
apt install nginx
curl -fsSL https://get.docker.com | bash -s docker
```

2. Git Clone

```bash
git clone https://gitea.890100.xyz/yeqing/acg-api.git
```

4. 启用容器

```
cd acg-api
docker compose up -d
```

5. 配置Nginx

这是与一份Nginx监听80端口 套CDN的配置文件 不包含HTTPS/SSL

```conf
server {
    listen 80 ; 
    server_name example.com; 

    # 获取真实IP
    real_ip_recursive on; 
    set_real_ip_from 0.0.0.0/0; 
    real_ip_header X-Forwarded-For; 

    # CORS
    add_header Access-Control-Allow-Origin * always; 
    add_header Access-Control-Allow-Methods GET,POST always; 
    if ( $request_method = 'OPTIONS' ) {
        return 204; 
    }

    # 处理API
    location ~ ^/(api|test|pc|pe)\.php(/|$) {
        include fastcgi_params;
        fastcgi_param SCRIPT_FILENAME /var/www/html$fastcgi_script_name;
        fastcgi_pass 127.0.0.1:9000; # 对应Docker-compose中的端口 
        set $real_script_name $fastcgi_script_name; 
    }

    # 其他请求
    location / {
        return 404; 
    }
}
```
