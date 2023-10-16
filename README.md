# 開発環境
- Go 1.20.6
- MySQL 8.0.33
- Redis 7.2.1

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

#
- SQL
    - リレーションの仕方
    - SQLインジェクション
        - ソルト、ストレッチング
- MVCモデル(Model,View,Controller)
- JWT(JSON Web Token)
    - none攻撃
