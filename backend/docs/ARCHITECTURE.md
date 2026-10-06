# Arsitektur Backend Transhome

Backend memakai **DDD modular monolith** di atas framework **go-zero**.
Satu aplikasi, satu database PostgreSQL (`transhome`, schema `crm_schema`),
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
    ├── middleware/            Adapter HTTP untuk autentikasi JWT, otorisasi, logging
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
Request publik:
HTTP -> handler -> logic -> application use case -> domain/repository port
                                                   ^                 |
                                                   |                 v
                                    infrastructure adapter <- PostgreSQL

Request terlindungi:
HTTP -> authentication middleware -> authorization middleware -> handler -> logic
             |                               |
             v                               v
       verifikasi JWT              application query/policy
                                             |
                                             v
                                    permission repository
```

## Authentication dan authorization

Authentication dan authorization adalah tanggung jawab berbeda:

- **Authentication** membuktikan identitas pemanggil. Login admin mencari akun,
  memverifikasi password terhadap hash bcrypt, lalu menerbitkan access token JWT.
- **Authorization** memutuskan apakah identitas yang sudah terautentikasi boleh
  menjalankan aksi tertentu. Keputusan dibuat per route/aksi berdasarkan permission,
  bukan hanya karena pengguna berhasil login atau memiliki token.

### Alur yang direncanakan untuk request terlindungi

1. Route menyatakan permission code yang diwajibkan, misalnya permission untuk
   membaca katalog. Permission code adalah identifier stabil dari konfigurasi/data,
   bukan teks pesan yang ditampilkan kepada pengguna.
2. Authentication middleware mengambil bearer token dan memverifikasi signature
   dengan algoritma yang diizinkan, expiration, serta klaim identitas. Middleware
   menaruh principal terverifikasi ke request context. Klaim dari token yang belum
   diverifikasi tidak boleh dipakai.
3. Authorization middleware mengambil principal dan permission yang diwajibkan
   route, lalu meminta keputusan kepada application query/policy Identity & Access.
4. Application layer memakai repository interface milik domain. Adapter PostgreSQL
   Identity & Access membaca relasi `user_roles`, `roles`, `role_permissions`, dan
   `permissions`, termasuk status aktif yang relevan.
5. Middleware melanjutkan request hanya jika permission diberikan; jika tidak,
   request ditolak. Ketiadaan/kegagalan autentikasi dibedakan dari penolakan akses.

Middleware adalah adapter HTTP: ia boleh membaca request dan menulis response,
tetapi tidak memuat query SQL atau aturan role/permission. Handler dan logic tetap
tipis. Aturan keputusan berada di application/domain; akses tabel hanya berada pada
adapter PostgreSQL dalam modul Identity & Access. Permission tidak dipercaya dari
body request atau klaim role yang dikirim client. Jika permission di-cache kelak,
aturan invalidasi cache harus memastikan perubahan role/permission berlaku sesuai
kebijakan keamanan.

### Kondisi implementasi saat ini

- Login memeriksa email/password, status aktif, dan role staf/admin aktif
  (`SUPER_ADMIN`, `ADMIN_KATALOG`, `ADMIN_MEMBERSHIP`, atau `ADMIN_KONTEN`)
  sebelum menerbitkan token admin. Gerbang ini hanya menentukan kelayakan masuk
  ke aplikasi administratif; tidak memilih API atau aksi yang boleh diakses.
- JWT saat ini adalah access token HS256 dengan expiration; belum ada refresh token
  atau mekanisme logout/revocation.
- Authentication middleware dan authorization per permission **belum
  diimplementasikan**. Karena itu, selain endpoint login, route belum mendapat
  perlindungan dari middleware autentikasi/otorisasi.
- Authorization tetap terpisah dan nantinya membatasi API per route/aksi
  berdasarkan role/permission, misalnya katalog, membership, atau akses
  administratif yang lebih luas.
- Detail kebijakan logout/token revocation dan penyimpanan status sesi masih perlu
  diputuskan sebelum implementasi.

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
