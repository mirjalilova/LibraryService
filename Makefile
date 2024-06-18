CURRENT_DIR=$(shell pwd)

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

exp:
	export DBURL='postgres://postgres:feruza1727@localhost:5432/library_service?sslmode=disable'

mig-up:
	migrate -path migrations -database 'postgres://postgres:feruza1727@localhost:5432/library_service?sslmode=disable' -verbose up

mig-down:
	migrate -path migrations -database ${DBURL} -verbose down

mig-create:
	migrate create -ext sql -dir migrations -seq create_table

run-server:
	go run server/main.go

run-client:
	go run client/main.go

swag_init:
	swag init -g api/api.go -o api/docs