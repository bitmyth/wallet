TAG=bitmyth/walletservice

image:
	docker build -t $(TAG) .

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOGC=off  go build

pg:
	docker run --name my-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=mydb -p 5432:5432 -d postgres
