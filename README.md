# 開発環境
- Go 1.21.4
- MySQL 8.0.33
- Redis 7.2.1

# 使用しているGoモジュール
- github.com/go-sql-driver/mysql v1.7.1
- github.com/golang-jwt/jwt v3.2.2+incompatible
- github.com/google/uuid v1.3.0
- github.com/redis/go-redis/v9 v9.2.1

# データベース
## ToDoTable
- ID(PK)
- UserID(UserTable)
- ToDoTitle
- ToTo
- Status
- CreateAt
- UpdateAt
## UserTable
- UserID(PK)
- signin_password
## Redis
- KEY:UserID
- VALUE:JWT

# メモ
- SQL
    - リレーションの仕方
    - SQLインジェクション
        - ソルト、ストレッチング
- MVCモデル(Model,View,Controller)
- JWT(JSON Web Token)
    - none攻撃

# 修正予定
- 他人のToDOのIDを指定しても削除できてしまうので、ToDoのIDのユーザーが正しいかも確認する機能を追加する
- エラー内容をレスポンスボディに入れて返す
- ログをファイルに出力する

# やりたいこと
- DockerBuildする
- Docker-Composeで起動できるにする(imageをBuidして起動するとMySQL-Redisと接続できない)
- GHCR(GitHub Container Registry)にpushする
- SQLインジェクション対策
