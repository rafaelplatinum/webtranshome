# Perintah sehari-hari backend Transhome.
# Butuh: Go 1.24+, goctl (go-zero), golang-migrate.
#   go install github.com/zeromicro/go-zero/tools/goctl@latest
#   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

-include .env
export

.PHONY: gen run worker build test migrate-up migrate-down migrate-new

## gen: generate handler, logic, types, config, dan main (transhome.go) dari api/transhome.api
gen:
	goctl api go -api api/transhome.api -dir . --style go_zero

## run: jalankan server API
run:
	go run transhome.go -f etc/transhome-api.yaml

## worker: jalankan proses background
worker:
	go run ./cmd/worker

## build: compile semua package
build:
	go build ./...

## test: jalankan semua test
test:
	go test ./...

## migrate-up: terapkan semua migration ke database
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

## migrate-down: batalkan 1 migration terakhir
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

## migrate-new: buat file migration baru, contoh: make migrate-new name=add_member_address
migrate-new:
	migrate create -ext sql -dir migrations -seq $(name)
