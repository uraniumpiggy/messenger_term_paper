run: 
	go run cmd/main/main.go

db:
	docker exec -it postgres-db psql -U user service-db

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .