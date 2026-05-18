---
title: Reverse Proxy
---

# Reverse Proxy

Snipraw is plain HTTP on localhost. A reverse proxy gives you TLS, a custom domain, and optionally basic auth.

## Caddy

[Caddy](https://caddyserver.com) handles TLS automatically via Let's Encrypt.

```txt
snippets.example.com {
    reverse_proxy localhost:8245
}
```

Reload:

```bash
caddy reload
```

## nginx

Obtain a certificate first:

```bash
certbot certonly --nginx -d snippets.example.com
```

Then configure the server block:

```nginx
server {
    listen 80;
    server_name snippets.example.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name snippets.example.com;

    ssl_certificate     /etc/letsencrypt/live/snippets.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/snippets.example.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8245;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Basic auth

Add access control at the proxy layer. Snipraw has no built-in authentication.

**Caddy:**

```txt
snippets.example.com {
    basicauth {
        pat $2a$14$...hashed-password...
    }
    reverse_proxy localhost:8245
}
```

Generate a password hash:

```bash
caddy hash-password
```

**nginx:**

```nginx
location / {
    auth_basic "Snippets";
    auth_basic_user_file /etc/nginx/.htpasswd;
    proxy_pass http://127.0.0.1:8245;
}
```

Generate a password file:

```bash
htpasswd -c /etc/nginx/.htpasswd your-username
```

::: warning
Basic auth transmits credentials in base64. Always use it with TLS.
:::
