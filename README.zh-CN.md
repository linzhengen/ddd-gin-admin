# ddd-gin-admin
[English](./README.md) | 简体中文 | [日本語](./README.ja.md)

`ddd-gin-admin` 采用 DDD 架构并提供构建 CMS 所需的组件。通过 RBAC 进行用户权限管理。

[![golangci-lint](https://github.com/linzhengen/ddd-gin-admin/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/linzhengen/ddd-gin-admin/actions/workflows/golangci-lint.yml)

## DDD Architecture
+ Domain: 这是定义应用程序的域和业务逻辑的地方。
+ Infrastructure: 这一层依赖于Domain层，例如数据库访问等, 对Domain层定义的接口的实现。
+ Application：这一层作为Domain层和Interfaces层之间的通道。将请求从Interfaces层发送到Domain层，Domain层对其进行处理并返回响应。
+ Interfaces：该层包含与其他系统交互的所有内容，例如Web应用程序以及批处理等。
<div>
    <img height="400" src="docs/img/ddd_architecture.png">
</div>

## Swagger UI
- GitHub Page: https://linzhengen.github.io/ddd-gin-admin/docs/swagger-ui/
- Localhost: http://localhost:8080/swagger/index.html
<div align="center">
    <img src="docs/img/swagger.png">
</div>

## Compiles and hot-reloads for development
```
make skaffold-dev
```
## Starting tunnel via minikube for service ddd-gin-admin-web
```
minikube service ddd-gin-admin-web --url -n ddd-gin-admin
```
## Lint
```
make lint
```
## Build binary
```
make build
```

## references
+ https://dev.to/stevensunflash/using-domain-driven-design-ddd-in-golang-3ee5
+ https://github.com/LyricTian/gin-admin
