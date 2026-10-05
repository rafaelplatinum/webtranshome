# Backend Transhome

Backend website Transhome Muka 1 (Katalog & Membership).
Go + go-zero, PostgreSQL, arsitektur DDD modular monolith.

Penjelasan struktur dan aturan arsitektur: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Kebutuhan

- Go 1.24+
- PostgreSQL 14+
- goctl: `go install github.com/zeromicro/go-zero/tools/goctl@latest`
- golang-migrate: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

## Menjalankan pertama kali

```bash
cd backend
cp .env.example .env        # lalu isi DATABASE_URL dan lainnya
make migrate-up             # buat tabel di database transhome
make gen                    # generate handler, logic, types, dan transhome.go
go mod tidy                 # unduh dependency go-zero
make run                    # server jalan di http://localhost:8888
```

Cek: `GET http://localhost:8888/api/v1/health`

## Perintah lain

| Perintah | Fungsi |
|---|---|
| `make worker` | Jalankan proses background |
| `make build` | Compile semua package |
| `make test` | Jalankan test |
| `make migrate-down` | Batalkan 1 migration terakhir |
| `make migrate-new name=xxx` | Buat file migration baru |

## Menambah endpoint baru

1. Tambahkan endpoint di `api/transhome.api` dengan `group` sesuai nama modul.
2. Jalankan `make gen`. File yang sudah ada tidak akan ditimpa.
3. Tulis use case di `internal/modules/<modul>/application/`.
4. Isi `internal/logic/<modul>/...` hanya untuk memanggil use case tersebut.
