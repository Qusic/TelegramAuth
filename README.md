# TelegramAuth

[![Go](https://github.com/Qusic/TelegramAuth/workflows/Go/badge.svg?branch=master)](https://github.com/Qusic/TelegramAuth/actions?query=workflow%3AGo)

Authorization server using [Telegram](https://core.telegram.org/widgets/login) as the authentication provider.

Works with:

* [NGINX](http://nginx.org/en/docs/http/ngx_http_auth_request_module.html)
* [NGINX Ingress controller](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#external-authentication)

Usage:

1. Create [`config.yaml`](config.example.yaml).
2. Run the executable in the directory containing the config file.

Endpoints:

* `/prefix/`  \
  Reverse proxy should send subrequest with cookies here to check the auth status.  \
  If the response is 200, proceed with the original request.  \
  If the response is 401, redirect to the login page.
* `/prefix/login`  \
  Unauthorized users should be redirected here to login with Telegram.
* `/prefix/callback`  \
  Telegram redirects authenticated users here to further redirect them to the app if authorized.

Example with NGINX Ingress controller:

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: myauth
spec:
  rules:
    - host: example.com
      http:
        paths:
          - path: /auth/
            backend: # service of TelegramAuth deployment
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: myapp
  annotations:
    nginx.ingress.kubernetes.io/auth-url: "https://$host/auth/?role=abc"
    nginx.ingress.kubernetes.io/auth-signin: "https://$host/auth/login?role=abc"
    nginx.ingress.kubernetes.io/auth-response-headers: "X-Telegram-Auth"
spec:
  rules:
    - host: example.com
      http:
        paths:
          - path: /
            backend: # service of upstream app
```

You can define multiple roles in the config file and refer to one of them using a configurable query parameter `role` in `auth-url` and `auth-signin`. This way, one instance of TelegramAuth is capable to handle multiple upstream apps with different sets of authorized users.

NGINX Ingress controller use `rd` query parameter to pass the original request url for redirecting after auth completion. The parameter for redirect url is configurable.

After `myapp` is set up with Telegram Login, you can read the `X-Telegram-Auth` request header in the upstream server to know who is using your app.
