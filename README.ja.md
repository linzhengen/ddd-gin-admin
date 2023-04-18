# ddd-gin-admin
[English](./README.md) | [简体中文](./README.zh-CN.md) | 日本語

このリポジトリは、DDD（Domain-Driven Design）アーキテクチャとGinフレームワークを使用したWebアプリケーションのサンプルです。

[![golangci-lint](https://github.com/linzhengen/ddd-gin-admin/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/linzhengen/ddd-gin-admin/actions/workflows/golangci-lint.yml)

## 機能
以下の機能を提供しています。

- ユーザーの登録、ログイン、ログアウト
- ユーザーの一覧表示、詳細表示、編集、削除
- ロール（管理者、一般ユーザー）に基づくアクセス制御
- Swaggerを使用したAPIドキュメント

## 技術スタック
以下の技術スタックを使用しています。

- Golang
- Gin - Webフレームワーク
- GORM - ORMライブラリ
- MySQL - データベース
- Swagger - APIドキュメント生成ツール
- K8s / Skaffold / Docker - コンテナ化

## DDD Architecture
+ Domain: Domain層は、アプリケーションのドメインとビジネスロジックが定義されます。
+ Infrastructure: Infrastructure層は、DBアクセスなどの技術的関心を記述します。この層はDomain層に依存しています。 そのためInfrastructure層はDomain層のrepositoryで定義したインタフェースを実装します。
+ Application: Application層は、ドメインとインターフェース層の間の通路として機能します。は、インターフェイス層からドメイン層に要求を送信し、ドメイン層はそれを処理して応答を返します。
+ Interfaces: Interfaces層は、Webアプリケーションやバッチ処理など、他のシステムと対話するすべてのものを保持します。
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
