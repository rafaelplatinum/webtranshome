# Inventaris gambar Transhome — Fase 1

**Tanggal pemeriksaan:** 5 November 2026
**Batas pekerjaan:** inventaris read-only. Tidak ada fitur gambar yang diimplementasikan pada fase ini. Beberapa perubahan lokal yang belum di-commit dicatat sebagai kondisi working tree saat audit, bukan sebagai bagian dari commit laporan.

## Ringkasan status

| Status | Jumlah slot | Makna |
| --- | ---: | --- |
| SUDAH ADA | 0 | Belum ditemukan siklus gambar lengkap pada backend dan admin produksi. |
| SEBAGIAN | 13 | UI/field atau URL gambar tersedia, tetapi alur produksi belum lengkap atau hanya berjalan di mock. |
| BELUM ADA | 4 | Belum ada slot gambar yang dikelola admin. |
| **Total** | **17** | Slot yang dipakai atau diminta dalam lingkup situs. |

## Tabel slot gambar

Ukuran rekomendasi mengikuti rasio komponen yang ada. Untuk logo/favicons, gunakan aset transparan dan hindari `object-cover`.

| ID slot | Halaman/komponen dan path pemakai | Sumber saat ini | Jumlah | Rasio / rekomendasi unggah | Wajib, opsional, fallback | Status admin dan bukti |
| --- | --- | --- | --- | --- | --- | --- |
| `home.hero` | `fase-5-konten/src/components/beranda/HeroBanner.vue`; admin `fase-5-konten/src/views/admin/konten/BannerForm.vue` | Banner `HOME_HERO`; mock SVG `/contoh-gambar/banner-hero.svg`; video terpisah opsional | Satu aktif utama | Tampilan media sekitar 5:4; form menyarankan 12:5, 1920×800 px; perlu diselaraskan | Gambar wajib; placeholder/latar dasar sebagai fallback, gambar juga cadangan video | **SEBAGIAN** — form/posisi ada di `fase-5-konten/src/utils/banner.js`; upload mengarah ke `/admin/uploads`, tetapi route itu tidak ada di `backend/internal/handler/routes.go`. Banner di mock (`fase-5-konten/src/services/mock/routesKonten.js`). Alt text tidak ada. |
| `home.promo_banner` | `fase-5-konten/src/components/beranda/BannerSlider.vue`; admin `fase-5-konten/src/views/admin/konten/BannerForm.vue`, `BannerDaftar.vue` | Record `HOME_SLIDER`; SVG seed/mock di `fase-0-fondasi/src/services/mock/data/konten.js` | Banyak, berurutan | 5:2; 1600×640 px (`fase-5-konten/src/utils/banner.js`) | Wajib per banner; state kosong/skeleton bila tidak ada | **SEBAGIAN** — form/daftar/pengurutan frontend mock tersedia; backend Go tidak memiliki endpoint konten/banner/upload. Alt eksplisit tidak ada. |
| `home.side_banner` | `fase-5-konten/src/components/beranda/BannerPromo.vue`; form/daftar banner admin | Record `HOME_SIDE`; seed/mock | Banyak, berurutan | 20:13; 800×520 px | Wajib untuk banner aktif; tidak ditampilkan jika kosong | **SEBAGIAN** — posisi, urutan dan UI mock ada; backend belum menyediakan API konten/upload. Alt tidak ada. |
| `category.tile` | `fase-5-konten/src/components/beranda/KategoriTahap.vue`; admin `fase-2-admin-katalog/src/components/admin/MasterFormModal.vue` | `categories.image_url`; mock membuat SVG contoh berdasar nama kategori | Per kategori | Sekitar 5:4; 800×640 px | Opsional; tile tetap bisa tampil tanpa foto | **SEBAGIAN** — field kategori ada; `categories.image_url` dan CRUD katalog ada di `backend/migrations/000001_init_schema.up.sql` serta `backend/internal/handler/routes.go`; upload/storage dan alt tidak ada. |
| `product.main` | `fase-0-fondasi/src/components/katalog/ProductCard.vue`, `ProductListRow.vue`; detail `fase-1-katalog/src/components/produk/ProductGallery.vue` | Gambar utama dari `product_images.image_url`; mock SVG contoh per produk | Satu utama per produk | 1:1; 1200×1200 px | Opsional di mock; placeholder/area kosong bila tidak ada | **SEBAGIAN** — admin foto di `fase-2-admin-katalog/src/components/admin/produk/FotoProduk.vue`; tabel `product_images` tersedia, tetapi route backend tidak memiliki API foto/upload. Alt tidak seragam. |
| `product.gallery` | `fase-1-katalog/src/components/produk/ProductGallery.vue`, `Lightbox.vue`; admin `fase-2-admin-katalog/src/components/admin/produk/FotoProduk.vue` | Banyak `product_images.image_url`; mock SVG contoh | Banyak per produk, berurutan dan satu ditandai utama | 1:1; 1200×1200 px | Opsional; placeholder bila galeri kosong | **SEBAGIAN** — UI dapat unggah banyak, hapus, urutkan, pilih utama; pengunggahan memakai route `/admin/uploads` yang tidak tersedia di backend. Tabel menyimpan URL/urutan, tanpa alt. |
| `brand.logo` | `fase-5-konten/src/components/beranda/LingkaranBrand.vue`; `fase-1-katalog/src/views/public/BrandDaftar.vue`, `Brand.vue`; admin `fase-2-admin-katalog/src/views/admin/katalog/Brand.vue`, `MasterFormModal.vue` | `brands.logo_url`; mock SVG nama/inisial brand | Satu per brand | 1:1; 512×512 px transparan (SVG/PNG/WebP) | Opsional; nama brand tetap tampil | **SEBAGIAN** — field/CRUD ada (`brands.logo_url`, migration awal, repository brand); upload backend dan alt tersimpan tidak ada. Daftar publik menggunakan nama brand untuk alt. |
| `room.cover` | `fase-1-katalog/src/components/katalog/RoomCard.vue`, `views/public/RuanganDetail.vue`; admin `fase-2-admin-katalog/src/components/admin/MasterFormModal.vue` | `rooms.image_cover`; mock SVG berdasarkan nama ruangan | Satu per ruangan | 4:3; 1200×900 px | Wajib pada form; detail menyembunyikan gambar jika kosong | **SEBAGIAN** — field wajib, tabel dan CRUD room tersedia; upload/storage dan alt tidak ada. Detail memakai alt kosong. |
| `service.cover` | `fase-5-konten/src/components/jasa/KartuJasa.vue`, `views/public/JasaDetail.vue`; admin `views/admin/konten/JasaForm.vue` | `services.image_url`; data mock SVG | Satu per jasa | 16:9; 1200×675 px (hint form) | Opsional; blok warna bila kosong | **SEBAGIAN** — `GambarField` dan mock CRUD tersedia; belum ditemukan tabel/repository/API jasa atau upload di backend Go. Alt kosong. |
| `article.cover` | `fase-5-konten/src/components/konten/KartuArtikel.vue`, `views/public/ArtikelDetail.vue`; admin `views/admin/konten/ArtikelForm.vue` | `articles.thumbnail_url`; mock SVG berdasarkan judul | Satu per artikel | 16:9; 1200×675 px (hint form) | Opsional; `MediaPlaceholder`/skeleton | **SEBAGIAN** — form frontend mock dan kolom `articles.thumbnail_url` pada migration awal; backend tidak mendaftarkan API artikel/upload. Alt kosong. |
| `promo.cover` | `fase-5-konten/src/views/public/Promo.vue`, `components/konten/KartuArtikel.vue`; form artikel | Artikel bertipe `PROMO`; mock SVG | Satu per promo | 16:9; 1200×675 px | Opsional; placeholder kartu | **SEBAGIAN** — berbagi thumbnail dan mock admin artikel; tidak ada API promo/upload backend. Alt kosong. |
| `event.cover` | Halaman konten publik dan `fase-5-konten/src/components/konten/KartuArtikel.vue`; form artikel | Artikel bertipe `EVENT`; mock SVG | Satu per event | 16:9; 1200×675 px | Opsional; placeholder kartu | **SEBAGIAN** — berbagi field thumbnail/mock artikel; tidak ada API event/upload backend. Alt kosong. |
| `store.front` | Info Toko `fase-5-konten/src/views/public/InfoToko.vue`; banner admin `fase-5-konten/src/views/admin/konten/BannerForm.vue`, `BannerDaftar.vue` | Posisi `ABOUT_HERO`; halaman membaca banner pertama. Dukungan ini ada pada working tree lokal saat audit, belum di-commit | Satu banner aktif utama | 12:5; 1920×800 px | Opsional; latar dasar jika belum diisi | **SEBAGIAN** — field posisi banner ada pada perubahan lokal; belum ada endpoint/storage backend. Gambar diperlakukan dekoratif (`alt=""`, `aria-hidden`). |
| `store.map` | `fase-5-konten/src/views/public/InfoToko.vue` | `maps_url` mengarah ke Google Maps; visual map saat ini bukan aset gambar upload (permukaan/ilustrasi CSS/HTML) | Satu preview jika kelak memakai gambar | Foto/preview 4:3; 1200×900 px. Peta interaktif dapat tetap berupa embed/link | Opsional; tautan lokasi tersedia | **BELUM ADA** — tidak ada slot gambar peta atau upload; `maps_url` adalah tautan tujuan, bukan gambar yang di-host. |
| `site.logo` | `fase-0-fondasi/src/components/layout/LogoTranshome.vue` (digunakan header/footer) | Aset statis `fase-0-fondasi/src/assets/logo-transhome.png` | Satu aset dipakai berulang | Lambang 1:1 (tampil 36×36 px); sumber minimal 512×512 px transparan | Wajib; aset statis menjadi fallback | **BELUM ADA** — tidak ada UI admin/endpoint untuk ganti logo. Audit tidak mengubah tampilan header/nav/footer. |
| `site.favicon` | `frontend/index.html`; `frontend/public/favicon.ico`, `favicon-16x16.png`, `favicon-32x32.png`, `apple-touch-icon.png` | Aset statis | Beberapa ukuran dari satu identitas | 16×16, 32×32, ICO multi-size, Apple touch 180×180 px | Wajib untuk ikon tab; fallback ikon browser | **BELUM ADA** — hanya dapat diganti lewat kode/deploy, tanpa UI admin. |
| `site.og_image` | `fase-0-fondasi/src/composables/useSeo.js`; default tags `frontend/index.html` | `useSeo({ image })` dapat mengatur OG image per halaman; tidak ada default `og:image` di index | Satu per halaman atau default situs | 1.91:1; 1200×630 px | Opsional; tanpa nilai, crawler memakai preview standar | **BELUM ADA** — composable bukan slot admin; tidak ada form/tabel/upload/default image. Catatan `useSeo.js`: crawler sosial masih membaca meta awal SPA sampai tersedia prerender/SSR. |

## Audit alur admin dan backend

| Area | Temuan dan bukti |
| --- | --- |
| Lihat, unggah, ganti, hapus | `fase-2-admin-katalog/src/components/admin/GambarField.vue` menangani pratinjau, unggah/ganti, dan hapus gambar opsional. `components/admin/produk/FotoProduk.vue` menangani banyak foto, hapus, urut, dan pilih utama. Upload menuju `POST /admin/uploads` dari `fase-2-admin-katalog/src/services/uploadService.js`, tetapi route tersebut tidak ada pada `backend/internal/handler/routes.go`. Mock routes menyimpan data URL di localStorage. |
| API produksi | `backend/internal/handler/routes.go` mendaftarkan API katalog dan RBAC; tidak ditemukan upload, banner, artikel, jasa, atau media-slot. CRUD katalog mengubah URL pada record, bukan menerima dan menyimpan file. |
| Skema data | `backend/migrations/000001_init_schema.up.sql` memiliki `categories.image_url`, `brands.logo_url`, `rooms.image_cover`, `product_images.image_url`, `banners.image_url`, `articles.thumbnail_url`. Kolom tersebut menyimpan URL; tidak ditemukan entitas aset generik atau kolom `alt_text`, MIME, bytes, dimensi. Tidak ditemukan tabel jasa pada migration backend yang diperiksa. |
| Urutkan | Galeri produk dapat diurutkan di UI; `product_images` memiliki `sort_order`. Kategori/ruangan memiliki `sort_order`; banner/jasa mock memakai urutan. Belum ada endpoint media untuk backend. Logo brand bukan galeri. |
| Validasi | `fase-2-admin-katalog/src/services/uploadService.js` hanya mengecek `file.type` dan ukuran di browser: JPG/PNG/WebP maksimal 2 MB untuk gambar. Tidak ada pemeriksaan magic bytes/konten, dimensi minimum, atau nama aman di sisi backend karena route upload belum ada. Nilai MIME browser dapat dipalsukan. |
| Penyimpanan | Mock menggunakan data URL di localStorage. Upload API sungguhan diharapkan mengembalikan URL `/admin/uploads`, tetapi route dan adapter storage tidak ditemukan. Logo/favicon disimpan sebagai file repo. |
| Alt text | Sebagian brand dan foto produk punya alt. Banner, kategori, detail ruangan, jasa, artikel/promo/event, serta preview admin memiliki alt kosong. Beberapa gambar memang dekoratif, tetapi tidak ada field/kolom alt yang dikelola admin untuk gambar konten. |
| RBAC | Mock routes memakai izin seperti `category.view`, `brand.update`, `banner.create`, `article.update`, dan `service.update` melalui `cekAdmin(...)`. Backend memiliki RBAC katalog, tetapi belum memiliki upload/media endpoint atau permission `image.manage`. |
| Audit log | Mock routes mencatat perubahan katalog/konten dengan `catat(...)`; tidak ada operasi file backend untuk merekam unggah/ganti/hapus. Reorder media juga belum memiliki operasi backend nyata. |
| Fallback/layout | Fallback berbeda: katalog mock menyediakan SVG, artikel memakai `MediaPlaceholder`, komponen lain menyembunyikan gambar atau menampilkan area warna/skeleton. Banyak `<img>` tidak memakai handler `@error`; gambar rusak dapat meninggalkan area kosong. Rasio tidak sepenuhnya konsisten. |

## Temuan masalah

1. Form upload frontend belum tersambung ke route `/admin/uploads` di backend Go. Mode mock menyimpan data URL di localStorage, bukan penyimpanan file produksi.
2. Gambar contoh katalog/konten banyak dibuat sebagai SVG inline dalam mock, belum menjadi pustaka aset dummy yang konsisten dan relevan dengan barangnya.
3. Fallback gambar tidak konsisten antar halaman dan penanganan gambar gagal (`@error`) belum menyeluruh.
4. Banyak gambar konten memakai `alt=""`; tidak ada input/kolom `alt_text` yang tersimpan. Alt kosong hanya tepat untuk gambar dekoratif.
5. Validasi file berjalan di browser saja, tanpa pemeriksaan konten file, dimensi minimum, nama aman, atau batas backend.
6. Logo, favicon, peta, dan gambar Open Graph belum mempunyai alur pengelolaan admin. `og:image` per halaman tidak sama dengan slot media; belum ada OG image default.
7. Hero Beranda pada working tree membagi copy dan media; rasio display yang terlihat tidak cocok sepenuhnya dengan rekomendasi upload banner 12:5.
8. **Permintaan lanjutan yang menunggu persetujuan Fase 2:** ubah hero Beranda menjadi foto penuh seperti hero Tentang Kami, dukung beberapa slide/banner, dan samakan rasio/alignment saat berpindah halaman navigasi. Screenshot yang dikirim menunjukkan teks banner bertabrakan dengan bidang media. Perubahan ini tidak dikerjakan pada fase inventarisasi.

## Asumsi audit

- **SUDAH ADA** mensyaratkan siklus produksi yang lengkap: admin dapat melihat, mengunggah, mengganti, menghapus dan mengurutkan sesuai kebutuhan, dengan validasi server, storage nyata, RBAC, audit log, dan alt text. Tidak ada slot yang memenuhi seluruh kriteria tersebut.
- Ukuran unggah adalah rekomendasi awal berdasarkan rasio komponen/hint form, belum divalidasi oleh backend.
- Promo/event dicatat sebagai slot berbeda karena jenis konten dan kebutuhan bisnisnya berbeda, walau memakai field thumbnail artikel yang sama.
- Perubahan pada working tree dicatat sebagai kondisi checkout saat audit; laporan ini tidak menyatakan bahwa perubahan tersebut sudah ada di commit dasar atau produksi.

**Fase 1 selesai. Berhenti di sini dan menunggu persetujuan sebelum Fase 2.**

**Catatan kartu member:** ilustrasi kartu Trans Family pada rontend/fase-5-konten/src/views/public/Membership.vue dibuat dengan elemen/CSS, bukan aset gambar upload. Karena belum ada slot gambar aktual, kartu tersebut tidak dihitung dalam 17 slot di atas.
