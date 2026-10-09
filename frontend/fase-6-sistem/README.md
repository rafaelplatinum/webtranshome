# Fase 6 — Sistem Admin dan Integrasi Qontak

Selesai 7 Oktober 2026 (27 file di `src/`). Hasil QA dan catatan ada di `docs/FRONTEND_BUILD_PLAN.md`, FASE 6.

Isi:
- Pengguna admin dan role & izin (`views/admin/akses/`): tambah admin, ubah role, reset password, aktif/nonaktif; matriks izin menu × aksi.
- Log aktivitas (`views/admin/audit/`): filter admin/modul/tanggal, drawer perbandingan data sebelum dan sesudah.
- Sync Qontak (`views/admin/integrasi/`): status PENDING/SYNCED/FAILED, coba lagi per kontak atau semua yang gagal, log sync.
- Mock API: `services/mock/routesSistem.js`, `routesQontak.js`, `antreanQontak.js` (tombol "Matikan Qontak" di `/dev/components`).
