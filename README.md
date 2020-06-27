# TelegramAuth

[![Go](https://github.com/Qusic/TelegramAuth/workflows/Go/badge.svg?branch=master)](https://github.com/Qusic/TelegramAuth/actions?query=workflow%3AGo)

Authorization server using [Telegram](https://core.telegram.org/widgets/login) as the authentication provider.

Works with:

* [NGINX](http://nginx.org/en/docs/http/ngx_http_auth_request_module.html)
* [NGINX Ingress controller](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#external-authentication).

Usage:

1. Create [`config.yaml`](config.example.yaml).
2. Run the executable in the directory containing the config file.

Endpoints:

* `/prefix/app`  \
  Reverse proxy should send subrequest with cookies here to check the auth status.  \
  If the response is 200, proceed with the original request.  \
  If the response is 401, redirect to the login page.
* `/prefix/app/login`  \
  Unauthorized users should be redirected here to login with Telegram.
* `/prefix/app/callback`  \
  Telegram redirects authenticated users here to further redirect them to the app if authorized.
