# Panduan kerja Transhome

- Baca `docs/redesign-spec.md` di awal setiap sesi sebelum mengubah kode.
- Area terkunci: top bar, main navigation, sub navigation, footer, dan tampilan katalog; ubah hanya pengecualian yang disebut spesifikasi.
- Gunakan design tokens dari bagian B dan font yang sudah digunakan project.
- Jangan tampilkan teks pengembangan seperti TODO, LOGO, GAMBAR CONTOH, atau Banner Transhome di UI. Placeholder administrasi hanya boleh ada di kode/dokumentasi sesuai spesifikasi.
- Perbarui `docs/redesign-progress.md` setelah setiap tahap, dan salin checkpoint-nya ke `C:\Users\LENOVO\Desktop\Transhome-progress`.
- Pekerjaan sesi ini dilanjutkan pada `main` sesuai instruksi eksplisit terbaru pemilik project; catatan commit/push mengikuti instruksi terbaru, tidak mengunggah perubahan tanpa otorisasi.
- Jalankan lint/build/test yang diminta spesifikasi dan catat hasil nyata saja. Jangan jalankan migrasi ke database produksi.
- Keputusan pemilik 6 Nov 2026: `frontend/` memakai tampilan lama + fitur pilihan (`frontend/docs/KEPUTUSAN.md` #24); redesign frontend di `docs/redesign-spec.md` tidak dipakai. Jangan mengganti tampilan `frontend/` tanpa persetujuan pemilik. Kontrak yang perlu disesuaikan backend: `frontend/README.md`.
