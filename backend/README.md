# backend
```
cp .envrc-example .envrc
```

起動
```
docker compose up -d
go run cmd/server/main.go
```

確認
```
curl -X GET http://localhost:8080

curl -X GET http://localhost:8080/files/:id

curl -X DELETE http://localhost:8080/files/:id
```

TODO
- [x] ユーザ認証
  - リファクタリングした方が良いが、とりあえずやった
  - セキュリティの穴は埋めた方が良い
- [ ] HTTPハンドラーのエラー処理を改善
