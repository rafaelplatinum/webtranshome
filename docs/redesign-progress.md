# Progress redesign Transhome

Terakhir diperbarui: 2026-11-06
Branch kerja sesuai instruksi pengguna: `main` lokal. `origin/main` saat mulai berada di `bf6964c`; local `main` berisi fitur OTP `f858a67` (belum berada di `origin/main`). Perubahan UI terbaru `5dde9b5` ada di `feat/home-hero-image-carousel` dan belum digabung ke `main`.

## Checklist tahap

- [x] Step 0: prompt induk, AGENTS, dan checklist tersimpan di repo dan checkpoint Desktop. di `docs/redesign-spec.md`.
- [x] Ringkasan aturan kerja disimpan di `AGENTS.md`.
- [ ] A+B+D+E: implementasi ada di commit `5dde9b5` pada branch UI, belum diintegrasikan ke `main`; perlu merge dan cek final.
- [ ] C: redesign halaman publik ada di `5dde9b5`, belum diintegrasikan ke `main`; perlu merge dan cek visual/responsif.
- [x] F1: inventaris aset/admin ada di `docs/image-inventory.md` pada branch `docs/image-inventory-phase-1`; inventaris sudah tersedia pada branch dokumen; perlu ikut diintegrasikan.
- [ ] F2: manajemen gambar admin untuk slot yang BELUM/SEBAGIAN — belum dikerjakan. Pengguna meminta melanjutkan progress; fase ini akan menjadi fokus berikutnya.
- [ ] F3: aset dummy/fallback dan alt — belum dikerjakan.
- [ ] G: Google Identity Services/backend — belum ditemukan implementasi pada commit terakhir.
- [x] H: lupa password OTP diimplementasikan pada commit lokal `f858a67`; butuh verifikasi setelah integrasi.

## Checkpoint awal

- URL aplikasi lokal yang tercatat: `http://127.0.0.1:5173/`; jalur preview perubahan terakhir OTP: `/lupa-password` dan `/dev/email`.
- Build frontend lulus setelah penggabungan sementara; `go test ./...` lulus setelah memakai cache Go lokal. GET `http://127.0.0.1:5173/lupa-password` mengembalikan HTTP 200. Pemeriksaan listener port 5173 tidak menemukan proses dari sesi shell ini.
- Jangan migrasikan database produksi. File `backend/.cache/` muncul sebagai file lokal tak terlacak pada worktree utama fitur OTP; jangan masukkan ke commit.
- Seluruh pekerjaan berikutnya dicatat di file ini dan dicopy sebagai checkpoint ke Desktop.

