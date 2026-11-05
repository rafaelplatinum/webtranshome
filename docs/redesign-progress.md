# Progress redesign Transhome

Terakhir diperbarui: 2026-11-06
Branch kerja sesuai instruksi pengguna: `main` lokal. Perubahan hero `5dde9b5`, reset password OTP `f858a67`, dan integrasi frontend sudah berada di `main` lokal. Pengguna meminta perubahan dikirim langsung ke `main`; push sebelumnya terhenti karena Git tidak menemukan kredensial (`SEC_E_NO_CREDENTIALS`).

## Checklist tahap

- [x] Step 0: prompt induk, AGENTS, dan checklist tersimpan di repo dan checkpoint Desktop. di `docs/redesign-spec.md`.
- [x] Ringkasan aturan kerja disimpan di `AGENTS.md`.
- [x] A+B+D+E: implementasi UI/WhatsApp/config kontak dari `5dde9b5` diintegrasikan ke `main`; header mempertahankan top bar, main bar, sub-navigation, dan footer terbaru sesuai arahan pengguna.
- [x] C: redesign halaman publik dari `5dde9b5` diintegrasikan ke `main`; build frontend lulus. Pemeriksaan visual mobile/desktop masih perlu dilakukan.
- [x] F1: inventaris aset/admin sudah terintegrasi di `docs/image-inventory.md`.
- [ ] F2: panel admin untuk 17 slot, upload, alt text, status aktif, hapus, dan urut sudah dibuat. API Go, penyimpanan lokal, validasi JPG/PNG/WebP (maks. 2 MiB, minimum 64×64), migrasi 000008, izin SUPER_ADMIN, serta audit log sudah ditambahkan; build dan `go test ./...` lulus. Belum dijalankan migrasi database, dan pemakaian aset baru pada seluruh komponen publik serta verifikasi visual masih TODO.
- [ ] F3: aset dummy/fallback dan alt — belum dikerjakan.
- [ ] G: Google Identity Services/backend — belum ditemukan implementasi pada commit terakhir.
- [x] H: lupa password OTP diimplementasikan pada commit lokal `f858a67`; backend penuh dan build frontend sudah tervalidasi setelah integrasi.

## Checkpoint awal

- URL aplikasi lokal yang tercatat: `http://127.0.0.1:5173/`; jalur preview perubahan terakhir OTP: `/lupa-password` dan `/dev/email`.
- Build frontend lulus, termasuk halaman Pengelolaan gambar; `go test ./...` lulus setelah memakai cache Go lokal. GET `http://127.0.0.1:5173/lupa-password` mengembalikan HTTP 200. Pemeriksaan listener port 5173 tidak menemukan proses dari sesi shell ini.
- Jangan migrasikan database produksi. File `backend/.cache/` muncul sebagai file lokal tak terlacak pada worktree utama fitur OTP; jangan masukkan ke commit.
- Seluruh pekerjaan berikutnya dicatat di file ini dan dicopy sebagai checkpoint ke Desktop.

## Checkpoint 2026-11-06

- Koreksi pengguna: pertahankan Top Bar, Main Bar, dan Sub-Navigation yang sudah ada; footer menggunakan desain terbaru. Sudah diverifikasi di `frontend/fase-0-fondasi/src/components/layout/AppHeader.vue` dan `AppFooter.vue`; membership juga tetap ada di sub-navigation.
- F2: API admin gambar dan storage lokal sudah diimplementasikan, termasuk endpoint file publik `/media/:filename`. Rute file dipisahkan dari prefix `/api/v1` agar sesuai URL yang disimpan di database.
- Validasi: `npm run build` frontend lulus; `go test ./...` backend lulus, termasuk tes dimensi WebP lossy, lossless, dan extended.
- Migrasi `000008_image_assets` belum dijalankan terhadap database. Jangan jalankan ke database produksi.
- File `backend/.cache/` adalah cache lokal dan tidak boleh ikut commit.
