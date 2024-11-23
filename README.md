# wallet

I've spent 2 days to finish this test

The most important part is to handle transfer correctly in highly concurrent scenario

To make code testable, most dependencies are injected rather than hard coded

In this way, I can inject suitable dependencies to cover more test cases

To review code, start with wallet/controller.go file

Postman collection file is listed at bottom of this page

## Run
Create `config.yaml` file

`cp config.example.yaml config.yaml`

Create `.env` file and edit like following
```shell
PG_USER=postgres
PG_PASS=pass
PG_DB=mydb
```
Run by docker-compose
```shell
docker-compose up 
```

When the program starts, it will insert two demo data entries for demonstration purposes.

| username | balance |
|----------|---------|
| user1    | 100.0   |
| user2    | 100.0   |

## Table design

[db/migrations/migration.sql](db/migrations/migration.sql)

## Folder structure

| folder        | usage                                                  |
|---------------|--------------------------------------------------------|
| config        | parse config file                                      |
| db            | connect postgres and redis                             |
| db/migrations | create table in postgres and insert some testing data  |
| factory       | a singleton to get components                          |
| route         | http router                                            |
| wallet        | core business logic                                    |

## Test
Prepare redis and postgres
```shell
docker run --name my-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=mydb -p 5432:5432 -d postgres

docker run --name my-redis -p 6379:6379 -d redis
```

```shell
go test ./...  -race -cover -coverprofile=coverage.out
```

### Coverage

```shell
 go tool cover -func=coverage.out
github.com/bitmyth/walletserivce/config/config.go:10:           SetConfigPath           100.0%
github.com/bitmyth/walletserivce/config/config.go:34:           NewConfig               81.8%
github.com/bitmyth/walletserivce/db/db.go:16:                   IsAlive                 100.0%
github.com/bitmyth/walletserivce/db/db.go:21:                   Open                    87.5%
github.com/bitmyth/walletserivce/db/db.go:41:                   IsAlive                 100.0%
github.com/bitmyth/walletserivce/db/db.go:46:                   OpenRedis               100.0%
github.com/bitmyth/walletserivce/db/migration.go:14:            Migrate                 81.2%
github.com/bitmyth/walletserivce/factory/factory.go:29:         RegisterRoutes          100.0%
github.com/bitmyth/walletserivce/factory/factory.go:33:         WalletController        100.0%
github.com/bitmyth/walletserivce/factory/factory.go:40:         New                     71.4%
github.com/bitmyth/walletserivce/factory/factory.go:54:         Config                  100.0%
github.com/bitmyth/walletserivce/factory/factory.go:58:         DB                      62.5%
github.com/bitmyth/walletserivce/factory/factory.go:72:         Redis                   62.5%
github.com/bitmyth/walletserivce/factory/factory.go:86:         Logger                  100.0%
github.com/bitmyth/walletserivce/factory/factory.go:90:         logger                  100.0%
github.com/bitmyth/walletserivce/factory/factory.go:102:        RegisterRoutes          100.0%
github.com/bitmyth/walletserivce/factory/factory.go:106:        WalletController        100.0%
github.com/bitmyth/walletserivce/factory/factory.go:113:        NewTesting              83.3%
github.com/bitmyth/walletserivce/factory/factory.go:126:        Config                  100.0%
github.com/bitmyth/walletserivce/factory/factory.go:130:        DB                      100.0%
github.com/bitmyth/walletserivce/factory/factory.go:134:        Redis                   100.0%
github.com/bitmyth/walletserivce/factory/factory.go:138:        Logger                  0.0%
github.com/bitmyth/walletserivce/main.go:11:                    main                    0.0%
github.com/bitmyth/walletserivce/route/router.go:11:            checkDB                 100.0%
github.com/bitmyth/walletserivce/route/router.go:22:            checkRedis              100.0%
github.com/bitmyth/walletserivce/route/router.go:33:            Router                  100.0%
github.com/bitmyth/walletserivce/wallet/controller.go:14:       NewController           100.0%
github.com/bitmyth/walletserivce/wallet/controller.go:26:       Deposit                 88.9%
github.com/bitmyth/walletserivce/wallet/controller.go:80:       Withdraw                88.9%
github.com/bitmyth/walletserivce/wallet/controller.go:140:      Transfer                90.0%
github.com/bitmyth/walletserivce/wallet/controller.go:204:      GetBalance              80.0%
github.com/bitmyth/walletserivce/wallet/controller.go:215:      GetTransactionHistory   78.6%
github.com/bitmyth/walletserivce/wallet/controller.go:239:      logTransaction          63.6%
github.com/bitmyth/walletserivce/wallet/controller.go:257:      RegisterRoutes          100.0%
github.com/bitmyth/walletserivce/wallet/controller.go:265:      handleError             50.0%
github.com/bitmyth/walletserivce/wallet/fixtures/balance.go:13: PreloadTestingData      90.9%
github.com/bitmyth/walletserivce/wallet/service.go:26:          NewService              100.0%
github.com/bitmyth/walletserivce/wallet/service.go:30:          GetBalance              86.4%
total:                                                          (statements)            80.5%
```

### Benchmark

`go test ./... -bench . -run none`

```shell
goos: darwin
goarch: arm64
pkg: github.com/bitmyth/walletserivce/wallet
BenchmarkWithdraw-8                 1402            960807 ns/op
BenchmarkDeposit-8                  1664            786380 ns/op
BenchmarkTransfer-8                 1231           1132481 ns/op
BenchmarkGetBalance-8               8811            125804 ns/op
```

## Postman
[wallet.postman_collection.json](wallet.postman_collection.json)