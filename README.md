# webtranshome

Website Transhome, retailer bahan bangunan dan elektronik di Yogyakarta.

- **Muka 1**: Katalog produk & Membership (poin per pembelian), terhubung ke Mekari Qontak.
- **Muka 2** (rencana): E-commerce.

## Struktur repo

| Folder | Isi |
|---|---|
| [`backend/`](backend/) | API Go (go-zero), arsitektur DDD modular monolith. Lihat [backend/README.md](backend/README.md) |
| [`frontend/`](frontend/) | Frontend Vue 3 + Vite + Tailwind CSS, hanya memakai REST API backend. Lihat [frontend/README.md](frontend/README.md) dan [frontend/CLAUDE.md](frontend/CLAUDE.md) |


## Menjalankan frontend

```bash
cd frontend
npm install
npm run dev
```
