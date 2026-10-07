# File per Fase

Kode frontend dipisah per fase. Setiap folder fase berisi `src/` dengan susunan yang sama
(`components/`, `views/`, `services/`, ...), tetapi aplikasinya tetap satu: `npm run dev` dan `npm run build`
dijalankan dari folder `frontend/`.

```
frontend/
  index.html, vite.config.js, jsconfig.json, package.json, .env.example   pengaturan proyek (dipakai semua fase)
  CLAUDE.md, README.md, docs/, design-system/                             aturan dan dokumen
  fase-0-fondasi/src/            74 file
  fase-1-katalog/src/            30 file
  fase-2-admin-katalog/src/      41 file
  fase-3-akun-member/src/        24 file
  fase-4-membership-admin/src/   24 file
  fase-5-konten/src/             42 file
  fase-6-sistem/src/             27 file
```

## Aturan

- Import antarfile memakai `@/...`, misalnya `@/components/ui/Button.vue`. Alias `@` mencari di `src/` semua
  folder fase secara berurutan (diatur di `vite.config.js` dan `jsconfig.json`), jadi file di Fase 4 bisa memakai
  komponen dari Fase 0 tanpa menulis nama foldernya.
- Satu path di bawah `src/` hanya boleh ada di satu folder fase (misalnya `components/ui/Button.vue` hanya di Fase 0).
- Import relatif (`./`, `../`) hanya untuk file di folder fase yang sama.
- File baru ditaruh di folder fase yang sedang dikerjakan. File lama diubah di tempatnya, di folder fase pembuatnya,
  walaupun perubahannya untuk fase lain.

## Isi tiap folder

### Fase 0 — Fondasi

`fase-0-fondasi/src/` · 74 file. Pengaturan proyek di `frontend/` (lihat di atas) juga dibuat di fase ini.

| Folder | File |
|---|---|
| `src/` | App.vue, main.js |
| `assets/` | main.css |
| `components/katalog/` | ProductCard.vue, ProductListRow.vue |
| `components/layout/` | AppFooter.vue, AppHeader.vue, MegaMenu.vue, SearchBox.vue, UserMenu.vue, WhatsAppButton.vue, WhatsAppFloat.vue |
| `components/ui/` | Badge.vue, Button.vue, Checkbox.vue, ConfirmModal.vue, EmptyState.vue, ErrorState.vue, Input.vue, Modal.vue, PaginationBar.vue, PriceTag.vue, Select.vue, Skeleton.vue, SkeletonCard.vue, StockBadge.vue, Switch.vue, Table.vue, Toast.vue |
| `composables/` | useAnimasi.js (8 Okt 2026, setelah Fase 6), useAuth.js, useFilter.js, usePermission.js, useSeo.js, useWhatsApp.js |
| `layouts/` | AdminLayout.vue, AuthLayout.vue, MemberLayout.vue, PublicLayout.vue |
| `router/` | guard.js, index.js |
| `router/routes/` | admin.js, auth.js, member.js, public.js |
| `services/` | authService.js, categoryService.js, errors.js, http.js, productService.js, settingsService.js |
| `services/mock/` | adapter.js, routes.js |
| `services/mock/data/` | akun.js, catalog.js, konten.js |
| `stores/` | auth.js, menu.js, settings.js, ui.js |
| `utils/` | format.js, whatsapp.js |
| `views/dev/` | Components.vue |
| `views/dev/bagian/` | DemoAkun.vue, DemoData.vue, DemoForm.vue, DemoKatalog.vue, DemoKeadaan.vue, DemoStatus.vue, DemoTombol.vue, DemoUmpanBalik.vue |
| `views/public/` | SegeraHadir.vue, TidakDitemukan.vue, TidakPunyaAkses.vue |

### Fase 1 — Katalog, pencarian, detail produk, brand, ruangan

`fase-1-katalog/src/` · 30 file

| Folder | File |
|---|---|
| `components/katalog/` | ActiveFilterChips.vue, ChipNav.vue, FilterOpsi.vue, FilterPanel.vue, FilterSheet.vue, KatalogDaftar.vue, KatalogToolbar.vue, RentangHarga.vue, RoomCard.vue, RoomCarousel.vue |
| `components/produk/` | Lightbox.vue, ProductGallery.vue, ProductSummary.vue, ProductTabs.vue, RelatedProducts.vue |
| `components/ui/` | Breadcrumb.vue, CopyButton.vue, NotFoundState.vue |
| `composables/` | useKatalog.js |
| `services/` | catalogService.js |
| `utils/` | clipboard.js, gerak.js, produk.js |
| `views/public/` | Brand.vue, BrandDaftar.vue (8 Okt 2026), Cari.vue, Katalog.vue, ProdukDetail.vue, Ruangan.vue, RuanganDetail.vue |

### Fase 2 — Admin: masuk, dashboard, master katalog

`fase-2-admin-katalog/src/` · 41 file

| Folder | File |
|---|---|
| `components/admin/` | AdminSidebar.vue, AktivitasTerbaru.vue, BrandCepatModal.vue, FormSection.vue, GambarField.vue, MasterFormModal.vue, PerhatianList.vue, TombolUrut.vue |
| `components/admin/produk/` | DatasheetField.vue, FotoProduk.vue, ProdukFilterBar.vue, ProdukTabel.vue, RuanganPicker.vue, SpesifikasiEditor.vue |
| `components/ui/` | Combobox.vue (8 Okt 2026), PasswordInput.vue, Textarea.vue |
| `composables/` | useFormProduk.js, useKonfirmasiKeluar.js, useMasterData.js, useTujuanKembali.js |
| `services/` | adminMasterService.js, adminProductService.js, adminService.js, cache.js, qontakService.js, uploadService.js |
| `services/mock/` | db.js, routesAdmin.js, sesi.js |
| `services/mock/data/` | admin.js |
| `utils/` | slug.js |
| `views/admin/` | Dashboard.vue, Masuk.vue, TidakDitemukan.vue |
| `views/admin/katalog/` | Brand.vue, Kategori.vue, ProdukDaftar.vue, ProdukForm.vue, ProdukPratinjau.vue, Ruangan.vue |

### Fase 3 — Auth dan akun member

`fase-3-akun-member/src/` · 24 file

| Folder | File |
|---|---|
| `assets/` | logo-google.svg |
| `components/akun/` | DaftarPoin.vue, KartuAuth.vue, KartuMember.vue, PetunjukMock.vue, TombolGoogle.vue |
| `composables/` | useJeda.js |
| `services/` | devService.js, memberService.js |
| `services/mock/` | routesAkun.js |
| `utils/` | validasi.js |
| `views/auth/` | Daftar.vue, LupaPassword.vue, Masuk.vue, MasukGoogle.vue, ResetPassword.vue, Verifikasi.vue |
| `views/dev/` | KotakEmail.vue, SimulasiGoogle.vue |
| `views/member/` | Keamanan.vue, Profil.vue, Ringkasan.vue, RiwayatPoin.vue, TidakDitemukan.vue |

### Fase 4 — Membership admin

`fase-4-membership-admin/src/` · 24 file

| Folder | File |
|---|---|
| `components/admin/` | StatusSync.vue |
| `components/admin/member/` | DataMember.vue, KartuProfilMember.vue, MemberFilterBar.vue, MemberTabel.vue, RiwayatMember.vue, RiwayatPoinAdmin.vue, SyncMember.vue |
| `components/admin/poin/` | CariMember.vue, IsianBelanja.vue, IsianTukarPoin.vue, MemberTerpilih.vue, PilihanJenis.vue, RingkasanPoin.vue, RiwayatRingkas.vue |
| `composables/` | useFormPoin.js, useMemberTerpilih.js, useSyncQontak.js |
| `services/` | adminMemberService.js |
| `services/mock/` | routesMember.js |
| `utils/` | poin.js |
| `views/admin/membership/` | MemberDaftar.vue, MemberDetail.vue, PoinInput.vue |

### Fase 5 — Konten, Beranda, dan Landing 3D

`fase-5-konten/src/` · 42 file

| Folder | File |
|---|---|
| `components/admin/konten/` | ArtikelFilterBar.vue, StatusArtikel.vue, StatusBanner.vue |
| `components/beranda/` | AjakanMember.vue, BannerPromo.vue, BannerSlider.vue, BrandPilihan.vue, HeroRumah.vue, InspirasiRuangan.vue, JudulSeksi.vue, KategoriTahap.vue, Keunggulan.vue, PesanGagal.vue, ProdukUnggulan.vue, PromoEvent.vue, RumahTigaDimensi.vue, bagianRumah.js, geometriRumah.js |
| `components/konten/` | ArtikelLainnya.vue, DaftarArtikel.vue, IsiArtikel.vue, KartuArtikel.vue |
| `composables/` | useDaftarBanner.js, useFormArtikel.js, useFormBanner.js, useMuat.js |
| `services/` | adminKontenService.js, kontenService.js |
| `services/mock/` | routesKonten.js |
| `services/mock/data/` | beranda.js |
| `utils/` | artikel.js, banner.js |
| `views/admin/konten/` | ArtikelDaftar.vue, ArtikelForm.vue, BannerDaftar.vue, BannerForm.vue, Pengaturan.vue |
| `views/public/` | ArtikelDaftar.vue, ArtikelDetail.vue, Beranda.vue, InfoToko.vue, Promo.vue |

### Fase 6 — Sistem Admin dan Integrasi Qontak

`fase-6-sistem/src/` · 27 file

| Folder | File |
|---|---|
| `components/admin/` | SkeletonTabel.vue |
| `components/admin/akses/` | DaftarRole.vue, MatriksIzin.vue, PenggunaFilterBar.vue, PenggunaTabel.vue, TambahAdminModal.vue, UbahRoleModal.vue |
| `components/admin/audit/` | DetailLog.vue, LogFilterBar.vue, LogTabel.vue, PerbandinganData.vue |
| `components/admin/integrasi/` | KontakFilterBar.vue, KontakTabel.vue, LogSyncDrawer.vue |
| `composables/` | useDaftar.js, useMatriksIzin.js |
| `services/` | adminAksesService.js, adminLogService.js |
| `services/mock/` | antreanQontak.js, routesQontak.js, routesSistem.js |
| `services/mock/data/` | sistem.js |
| `utils/` | log.js |
| `views/admin/akses/` | PenggunaDaftar.vue, Role.vue |
| `views/admin/audit/` | LogAktivitas.vue |
| `views/admin/integrasi/` | SyncQontak.vue |

## File yang diubah lintas fase

File tetap berada di folder fase pembuatnya; fase berikutnya mengubahnya di tempat.

| File (folder fase) | Diubah lagi di fase |
|---|---|
| `router/routes/*.js` (0) | 1, 2, 3, 4, 5: halaman asli menggantikan placeholder |
| `services/mock/routes.js`, `services/mock/data/*.js` (0) | 1, 2, 3, 4, 5 |
| `layouts/*Layout.vue` (0) | 1, 2, 3 |
| `utils/format.js` (0) | 1, 3, 4 |
| `stores/auth.js`, `services/authService.js` (0) | 3 |
| `services/categoryService.js` (0), `services/catalogService.js` (1) | 2: cache dibersihkan setelah admin mengubah data |
| `components/ui/Input.vue`, `Select.vue`, `Table.vue` (0) | 2 |
| `components/ui/Modal.vue` (0) | 1 |
| `components/ui/CopyButton.vue` (1) | 3 |
| `components/admin/PerhatianList.vue` (2) | 4: memakai `useSyncQontak` |
| `ProdukTabel.vue`, `Kategori.vue`, `Brand.vue`, `Ruangan.vue` (2), `RelatedProducts.vue` (1) | 4: perbaikan label pembaca layar |
| `services/mock/routesAdmin.js` (2), `services/mock/routesAkun.js` (3) | 4 |
| `router/routes/public.js`, `admin.js` (0) | 5: Beranda, artikel, promo, info toko, dan halaman admin konten; `views/public/Home.vue` (0) dihapus, diganti `Beranda.vue` (5) |
| `assets/main.css` (0) | 5: token warna rumah 3D `--color-rumah-*` (`@theme static`) |
| `components/ui/Table.vue` (0) | 5: wadah `relative` (teks sr-only tidak melebarkan halaman) |
| `utils/whatsapp.js`, `composables/useWhatsApp.js` (0) | 5: pesan & tautan WhatsApp untuk promo/artikel |
| `services/mock/data/konten.js`, `data/akun.js` (0), `services/mock/db.js` (2) | 5: banner & artikel contoh, izin `setting.point_ratio`, banner/artikel/pengaturan ikut tersimpan |
| `services/catalogService.js` (1) | 5: `listBrands()` untuk Brand pilihan |
| `components/admin/GambarField.vue`, `PerhatianList.vue` (2) | 5: pratinjau gambar lebar (banner); "Perpanjang" ke `#periode` |
| `router/routes/admin.js`, `services/mock/routes.js` (0) | 6: halaman Sistem asli menggantikan "sedang dibangun"; rute mock sistem & Qontak |
| `components/ui/Modal.vue` (0) | 6: varian `side` (drawer kanan) |
| `views/dev/bagian/DemoAkun.vue` (0) | 6: tombol "Matikan Qontak" (mode mock) |
| `services/mock/routesAdmin.js`, `db.js`, `data/admin.js`, `services/qontakService.js` (2) | 6: `perubahan()` untuk before/after log produk & master, antrean Qontak, izin role tersimpan, log contoh berisi data, layanan daftar/log kontak |
| `components/admin/AktivitasTerbaru.vue` (2) | 6: label aksi yang mudah dibaca |
| `AppHeader.vue`, `MegaMenu.vue`, `WhatsAppFloat.vue`, `Button.vue`, `Modal.vue`, `main.css`, `main.js`, `router/index.js` (0) | Setelah Fase 6 (8 Okt 2026): header publik menempel (baris alamat ikut tergulir di bawah 1024 px) dan animasi situs publik (`useAnimasi.js`) |
| `KatalogDaftar.vue` (1), `Beranda.vue`, `DaftarArtikel.vue` (5) | Setelah Fase 6: `v-reveal` (scroll reveal) |
| `router/routes/public.js`, `AppHeader.vue`, `AppFooter.vue`, `services/mock/routes.js`, `data/catalog.js`, `data/akun.js`, `main.css` (0) | Brand (8 Okt 2026): halaman `/brand`, tautan Brand, `product_count` di `GET /brands`, 25 brand contoh tambahan (24 aktif tanpa produk + "Merek Lama" nonaktif), izin `brand.delete`, animasi logo berjalan |
| `FilterOpsi.vue`, `services/catalogService.js` (1) | Brand: kotak cari bila opsi > 10, daftar lengkap bergulir |
| `views/admin/katalog/Brand.vue`, `ProdukForm.vue`, `ProdukFilterBar.vue`, `adminMasterService.js`, `routesAdmin.js` (2) | Brand: cari + hapus brand nonaktif, pilihan brand jadi `Combobox`, `DELETE /admin/brands/:id` |
| `BrandPilihan.vue` (5) | Brand: semua brand aktif, logo berjalan (jeda, berhenti saat disorot/difokus) |
| `data/sistem.js`, `useMatriksIzin.js`, `utils/log.js` (6) | Brand: izin `brand.delete` (kolom Lainnya: "Hapus"), label log `brand.delete` |
| `services/mock/routesAkun.js` (3), `services/mock/routesMember.js` (4) | 6: member baru masuk antrean sync, `kirimLinkReset()` & `hashPassword()` diekspor; status sync diproses sebelum dibaca |

## Push dan pull per folder

Belum ada yang di-commit. Kalau sudah siap, setiap fase bisa di-commit dari foldernya sendiri (dari root repo):

```bash
git switch -c feat/frontend-fase-0-4
git add frontend/fase-0-fondasi frontend/index.html frontend/vite.config.js frontend/jsconfig.json \
  frontend/package.json frontend/package-lock.json frontend/.gitignore frontend/.env.example \
  frontend/CLAUDE.md frontend/README.md frontend/docs frontend/design-system README.md
git commit -m "feat(frontend): fase 0 fondasi"
git add frontend/fase-1-katalog && git commit -m "feat(frontend): fase 1 katalog"
git add frontend/fase-2-admin-katalog && git commit -m "feat(frontend): fase 2 admin katalog"
git add frontend/fase-3-akun-member && git commit -m "feat(frontend): fase 3 auth dan akun member"
git add frontend/fase-4-membership-admin && git commit -m "feat(frontend): fase 4 membership admin"
git add frontend/fase-5-konten frontend/fase-6-sistem && git commit -m "chore(frontend): folder fase 5 dan 6"
git push -u origin feat/frontend-fase-0-4
```

Hanya commit terakhir yang dijamin bisa di-build, karena file lintas fase (tabel di atas) sudah berisi versi akhir.
