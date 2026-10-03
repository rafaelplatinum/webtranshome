# Arsitektur Backend Transhome

Backend memakai **DDD modular monolith** di atas framework **go-zero**.
Satu aplikasi, satu database PostgreSQL (`CRM_TH`, schema `"CRM_Schema"`),
dibagi menjadi beberapa modul (bounded context) yang tidak saling mengakses
tabel milik modul lain.

## Struktur folder

```
backend/
├── api/transhome.api          Definisi endpoint HTTP (sumber untuk goctl)
├── etc/transhome-api.yaml     Konfigurasi server
├── migrations/                File SQL golang-migrate (000001_init_schema.up/down.sql)
├── cmd/worker/                Proses background (sync Qontak, retry, notifikasi)
├── transhome.go               Entry point API (dibuat oleh `make gen`)
├── docs/ARCHITECTURE.md       Dokumen ini
└── internal/
    ├── config/                Struct konfigurasi (dibuat goctl)
    ├── svc/                   ServiceContext: tempat merakit repository & service
    ├── middleware/            Auth JWT, cek permission, logging
    ├── types/                 Request/response HTTP (dibuat goctl dari .api)
    ├── handler/<modul>/       Handler HTTP (dibuat goctl)
    ├── logic/<modul>/         Penghubung HTTP -> use case (TIPIS, tanpa aturan bisnis)
    ├── shared/
    │   ├── domain/            Shared kernel: base error, domain event, value object umum
    │   ├── eventbus/          Bus event antar modul (in-process)
    │   ├── tx/                Transaksi / unit of work
    │   ├── mailer/            Kirim email
    │   └── storage/           Simpan foto produk & datasheet PDF
    └── modules/
        ├── identityaccess/    users, password_resets, roles, menus, permissions, user_roles, role_permissions
        ├── membership/        member_profiles, point_transactions
        ├── catalog/           categories, brands, products (+ product_images, product_rooms), rooms
        ├── content/           articles, banners, site_settings
        ├── audit/             activity_logs
        └── integration/
            └── qontak/        qontak_contacts, sync_logs
```

Nama folder di `handler/` dan `logic/` **sama persis** dengan nama modul di
`modules/`, supaya mudah dicari.

## Lapisan di dalam setiap modul

```
modules/<modul>/
├── domain/
│   ├── <aggregate>/   Entity, value object, aturan bisnis, INTERFACE repository
│   └── service/       Domain service (aturan yang melibatkan lebih dari 1 aggregate)
├── application/
│   ├── command/       Use case yang mengubah data (create, update, assign)
│   ├── query/         Use case yang hanya membaca data
│   ├── dto/           Input/output use case
│   ├── listener/      Menangani domain event dari modul lain
│   └── port/          Interface ke modul lain / layanan luar yang dibutuhkan modul ini
└── infrastructure/
    ├── postgres/      Implementasi repository (GORM/SQL) + model tabel
    └── <adapter>/     Implementasi port: google, qontakapi, notifier, identityaccess
```

Folder `listener/` dan `port/` hanya dibuat di modul yang memang membutuhkannya.

## Alur request

```
HTTP -> handler/<modul> -> logic/<modul> -> modules/<modul>/application/command|query
                                                   |
                                                   v
                                        domain (entity + interface repository)
                                                   ^
                                                   | implementasi
                                        infrastructure/postgres
```

## 5 aturan wajib

1. **`domain/` tidak boleh meng-import** GORM, go-zero, `database/sql`, atau package
   `infrastructure` apa pun. Domain hanya berisi Go murni.
2. **Interface repository** ada di `domain/<aggregate>/`, **implementasinya** ada di
   `infrastructure/postgres/`.
3. **Aturan bisnis** (nota tidak boleh dobel, saldo poin harus cukup, role nonaktif
   tidak boleh di-assign) ditulis di entity atau domain service. Bukan di handler,
   logic, atau query SQL.
4. **`logic/` harus tipis**: ubah `types` -> `dto`, panggil use case, ubah hasilnya
   kembali ke `types`. Tidak ada `if` aturan bisnis di sini.
5. **Modul tidak boleh mengakses tabel milik modul lain.** Gunakan `port/` + adapter
   (contoh: `membership/infrastructure/identityaccess`) atau domain event lewat
   `shared/eventbus`.

## `types/` vs `application/dto/`

| | `internal/types/` | `modules/*/application/dto/` |
|---|---|---|
| Dibuat oleh | goctl dari `transhome.api` | Ditulis manual |
| Isi | Bentuk JSON request/response | Input/output use case |
| Boleh tahu HTTP? | Ya (tag `json`, `form`, `path`) | Tidak |

Dengan begitu use case tetap bisa dipakai dari worker atau test tanpa HTTP.

## Contoh alur antar modul (event)

```
identityaccess: RegisterUser  -> publish UserRegistered
membership:     listener      -> buat member_profiles (member_code, tier TAHAP_1)
qontak:         listener      -> qontak_contacts PENDING -> worker kirim ke API Qontak
audit:          listener      -> tulis activity_logs (untuk aksi admin)
```

## Kenapa `integration/` punya sub-folder?

Muka 2 (e-commerce) akan menambah integrasi lain, misalnya payment gateway dan
kurir. Semuanya diletakkan sejajar dengan `qontak/` di dalam `integration/`.
