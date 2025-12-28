# Project API

プロジェクト管理を担当するサンプルサービス。

## 機能

- プロジェクトのCRUD操作
- メンバー管理

## 起動方法

### ローカルで起動

```bash
# プロジェクトルートディレクトリで
make dev

# または
cd project-api && go run cmd/server/main.go
```

### Docker で起動

```bash
# プロジェクトルートディレクトリで
make up
```

## API エンドポイント

サーバーは `:8082` で起動します。

### ヘルスチェック

```bash
curl http://localhost:8082/health
```

### プロジェクト一覧取得

```bash
curl -X POST "http://localhost:8082/project.v1.ProjectService/ListProjects" \
  -H "Content-Type: application/json" \
  -d '{"workspace_id": "ws_1"}'
```

### プロジェクト取得

```bash
curl -X POST "http://localhost:8082/project.v1.ProjectService/GetProject" \
  -H "Content-Type: application/json" \
  -d '{"id": "proj_1"}'
```

### プロジェクト作成

```bash
curl -X POST "http://localhost:8082/project.v1.ProjectService/CreateProject" \
  -H "Content-Type: application/json" \
  -d '{
    "workspace_id": "ws_1",
    "name": "New Project",
    "description": "Project description"
  }'
```

### プロジェクトメンバー一覧取得

```bash
curl -X POST "http://localhost:8082/project.v1.ProjectService/ListMembers" \
  -H "Content-Type: application/json" \
  -d '{"project_id": "proj_1"}'
```

## 開発

### Protocol Buffers のコード生成

```bash
# プロジェクトルートディレクトリで
make proto
```

### 依存関係のインストール

```bash
make install-deps
```

## モックデータ

サーバー起動時に以下のサンプルデータが登録されます：

### プロジェクト

- `proj_1`: Sample Project 1 (ws_1)
- `proj_2`: Sample Project 2 (ws_1)

### メンバー

- `proj_1`:
  - `user_1` (owner)
  - `user_2` (member)
- `proj_2`:
  - `user_1` (admin)
