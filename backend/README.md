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
- ユーザ認証
- HTTPハンドラーのエラー処理を改善
