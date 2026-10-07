# Audit Backend Transhome

## Keputusan (final, jangan diubah Copilot)
- Arsitektur tetap: go-zero shell + modules/<context>/{domain,application,infrastructure}.
- Update role memakai semantik: PATCH dengan field opsional (*bool/*string).
- Format code role: trim + upper-case, pola ^[A-Z][A-Z0-9_]*$, maks 50 karakter.
- Persistence fitur baru: modules/.../infrastructure/postgres. internal/model (goctl) tidak disentuh.
- Role bawaan (isi kode-kodenya, mis. ADMIN): tidak boleh dinonaktifkan atau dihapus.
- Perubahan kontrak API (request/response): boleh tanpa bertanya.
- Perubahan schema DB: [boleh lewat migration baru 0000NN | tidak boleh].

## Temuan
| ID | Approve | Severity | File:line | Problem | Fix | Status |
|----|---------|----------|-----------|---------|-----|--------|
| A1 | Y | High | .../manage_roles.go:42 | IsActive bool menonaktifkan role tanpa sengaja | Ubah ke *bool | Blocked: generated request type cannot preserve omitted field |
| A2 | Y | Med  | ... | ... | ... | Blocked: incomplete finding lacks file, problem, and fix details |
| A3 | N | Low  | ... | ... | ... | |

## Di luar scope (jangan dikerjakan)
- Semua baris Approve = N
- Folder .gitkeep, file generated, modul selain identityaccess

## Standar Arsitektur dan Prinsip (binding, struktur TIDAK diubah)

### A. Clean Architecture (sesuai struktur repo saat ini)
Arah dependensi: handler/logic -> application -> domain <- infrastructure.

Aturan yang bisa dicek:
- `domain/**` tidak mengimpor: `application`, `infrastructure`, `internal/types`, `internal/svc`, `go-zero`, `database/sql`, `net/http`.
- `application/**` tidak mengimpor: `infrastructure`, `internal/types`, `internal/svc`, `go-zero`, `database/sql`, `net/http`.
- `infrastructure/**` boleh mengimpor domain dan `application/port`, tidak boleh mengimpor `handler`, `logic`, `svc`.
- `internal/logic` dan `internal/handler` tidak mengimpor `infrastructure` dan tidak berisi SQL.
- Satu modul tidak mengimpor `domain` atau `infrastructure` modul lain. Jalurnya lewat port milik pemakai (contoh: `membership/infrastructure/identityaccess`).
- Wiring konkret hanya di `internal/svc/service_context.go`.

Cara memeriksa (jalankan dan laporkan hasilnya sebagai temuan):
    grep -rn "internal/types\|go-zero\|database/sql\|net/http" backend/internal/modules/*/domain backend/internal/modules/*/application
    grep -rn "/infrastructure" backend/internal/modules/*/domain backend/internal/modules/*/application backend/internal/logic backend/internal/handler

Pembagian tanggung jawab:
| Lapisan | Boleh | Tidak boleh |
|---------|-------|-------------|
| handler | parse request, panggil logic, tulis response | aturan bisnis, SQL |
| logic | mapping types <-> dto, panggil SATU use case | aturan bisnis, memanggil banyak repository |
| application | orkestrasi use case, transaksi, otorisasi use case | HTTP, SQL, tipe framework |
| domain | entity, aturan bisnis, error domain, interface repository | I/O apa pun, tag JSON/DB |
| infrastructure | SQL, bcrypt, JWT, API eksternal, mapping error teknis -> error domain | aturan bisnis |

Kebocoran yang harus dilaporkan:
- Aturan bisnis (validasi, status, perhitungan) berada di handler, logic, middleware, atau repository.
- Entity domain membawa tag `json`/`db` atau tipe dari driver.
- Use case mengembalikan model DB atau `types.*` ke pemanggil.
- Logic yang sama ada di `internal/logic/*` dan `application/command/*` (duplikasi jalur).

### B. SOLID (versi praktis, bisa diukur)
- **S (Single Responsibility):** satu use case = satu file = satu aksi bisnis. Fungsi yang melakukan validasi + persist + notifikasi sekaligus dipecah. Tanda pelanggaran: nama berisi "And", struct `Manage*` dengan > 5 method yang tidak berhubungan.
- **O (Open/Closed):** perilaku baru ditambah lewat implementasi atau fungsi baru, bukan menambah cabang `switch/if` pada tipe string di banyak tempat. Tanda pelanggaran: `switch` atas jenis/role/status yang tersebar di > 2 file.
- **L (Liskov):** implementasi interface tidak menambah syarat tersembunyi (mis. repository fake dan postgres harus mengembalikan error domain yang sama untuk kasus yang sama).
- **I (Interface Segregation):** interface milik pemakai (application/domain), 1-3 method. Interface > 4 method atau didefinisikan di sebelah implementasinya dilaporkan.
- **D (Dependency Inversion):** use case menerima interface lewat constructor `NewX(...)`. Tidak ada `new`/konstruksi dependensi konkret di dalam use case, tidak ada global, tidak ada `init()` dengan side effect.

### C. Prinsip lain
- **KISS:** pilih solusi paling sederhana yang lolos aturan di atas. Jika dua solusi setara, pilih yang diffnya lebih kecil.
- **YAGNI:** tidak ada field, parameter, method, config, atau interface yang belum dipakai kode saat ini.
- **DRY (dengan batas):** ekstrak setelah duplikasi KETIGA dan hanya jika maknanya sama. Dua potongan yang mirip tapi alasan berubahnya berbeda JANGAN digabung.
- **Law of Demeter:** tidak ada rantai panggilan dalam seperti `a.B().C().D()` lintas lapisan.
- **Tell, don't ask:** aturan yang menyangkut state entity berada sebagai method entity (`role.Deactivate()`), bukan di use case yang memeriksa field lalu menulis field.
- **Composition over inheritance:** gunakan struct embedding hanya untuk komposisi yang jelas, tidak membuat "base struct".
- **Fail fast:** validasi di awal fungsi, kembalikan error secepatnya.
- **Make invalid states unrepresentable:** gunakan constructor/validator (`role.New`) agar entity tidak bisa dibuat dalam keadaan tidak valid.
- **Principle of least privilege:** komponen hanya menerima dependensi yang benar-benar dipakai.

### D. Anti Over-Engineering
Sebuah temuan "Over-engineered" valid bila memenuhi salah satu tanda berikut:
1. Interface hanya punya SATU implementasi dan tidak dipakai untuk mengisolasi I/O di test.
2. Abstraksi generik (generics, base repository, factory, builder, registry, plugin) yang dipakai < 3 tempat.
3. Layer/lapisan perantara yang hanya meneruskan panggilan tanpa menambah logika (pass-through).
4. Struct/DTO/mapper yang bentuknya identik dengan yang lain dan hanya disalin field per field tanpa alasan batas lapisan.
5. Parameter, field config, atau opsi yang tidak pernah diisi nilai berbeda.
6. Event, queue, cache, atau goroutine untuk operasi yang sinkron dan murah.
7. Pola desain dengan nama pola di nama tipe (`Factory`, `Strategy`, `Manager`, `Handler` ganda) tanpa kebutuhan variasi nyata.
8. Kode untuk "kebutuhan masa depan" (folder, file, method) yang belum ada use case-nya.

Aturan perbaikan untuk over-engineering:
- Hapus atau inline, jangan diganti dengan abstraksi lain.
- Pengecualian yang SAH (jangan dilaporkan): port yang memang sengaja untuk memisahkan modul atau infrastruktur (`application/port`), interface yang dipakai fake di test, dan folder `.gitkeep` yang sudah direncanakan.
- Sebelum menambah abstraksi baru, jawab tiga pertanyaan; jika ada jawaban "tidak", jangan tambah:
  1. Apakah ada minimal 2 pemakai/implementasi nyata sekarang?
  2. Apakah tanpa abstraksi ini aturan dependensi (bagian A) dilanggar?
  3. Apakah abstraksi ini membuat diff lebih kecil dan test lebih mudah?

### E. Kompleksitas yang dibatasi
- Fungsi: <= 30 baris, <= 4 parameter, nesting <= 3, cyclomatic complexity <= 10.
- File: <= 300 baris. Lebih dari itu, pecah per use case atau per tanggung jawab (tanpa membuat paket baru yang tidak diminta).
- Struct dependensi use case: <= 5 dependensi. Lebih dari itu biasanya tanda use case terlalu besar.

### F. Keseimbangan: kapan JANGAN memperbaiki
- Jika memperbaiki temuan butuh perubahan struktur folder atau pola arsitektur: tandai "Blocked: needs decision", jangan dikerjakan.
- Jika perbaikan membuat kode lebih panjang dan lebih abstrak tetapi tidak menghilangkan pelanggaran nyata: jangan dikerjakan.
- Jika kode sudah memenuhi aturan: jangan "dipercantik". Tidak ada temuan = tidak ada perubahan.