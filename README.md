# go でWebフレームワークを再発明した

## .env を作成
```
cp .env.example .env
```

## main関数から実行する場合
```
go run cmd/server/main.go

# or

# air をインストール
go install github.com/cosmtrek/air@latest

air -c .air.toml 
```

## docker compose で起動する場合
```
make up
```
