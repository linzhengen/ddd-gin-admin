# ddd-gin-admin
English | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md)

`ddd-gin-admin` adopts the DDD architecture and provides the components needed to build a CMS. User permission management by RBAC.

[![golangci-lint](https://github.com/linzhengen/ddd-gin-admin/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/linzhengen/ddd-gin-admin/actions/workflows/golangci-lint.yml)

## DDD Architecture
+ Domain: This is where the domain and business logic of the application is defined.
+ Infrastructure: The infrastructure layer describes technical concerns such as DB access. This layer depends on the domain layer. Therefore, the infrastructure layer implements the interface defined in the repository of the domain layer.
+ Application: This layer serves as a passage between the domain and the interface layer. The sends the requests from the interface layer to the domain layer, which processes it and returns a response.
+ Interfaces: This layer holds everything that interacts with other systems, such as web services, RMI interfaces or web applications, and batch processing frontend.
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
