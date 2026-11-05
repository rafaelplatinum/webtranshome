// Katalog publik dimuat dari snapshot yang diekspor dari database backup Transhome.
// Gambar eksternal yang belum tersedia diganti SVG dummy lokal agar tidak ada broken image.
import backup from './backup.json'

const warna = ['#E7E5E0', '#DCE3D2', '#D4E2EA', '#E7D7C1', '#FDF0C2']
const xml = (teks) => String(teks ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

function gambarContoh(judul, keterangan, latar = warna[0], lebar = 800, tinggi = 800) {
  const svg =
    `<svg xmlns="http://www.w3.org/2000/svg" width="${lebar}" height="${tinggi}" viewBox="0 0 ${lebar} ${tinggi}">` +
    `<rect width="${lebar}" height="${tinggi}" fill="${latar}"/>` +
    `<rect x="4%" y="4%" width="92%" height="92%" fill="none" stroke="#FFFFFF" stroke-opacity=".65" stroke-width="4"/>` +
    `<text x="50%" y="47%" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 13)}" font-weight="700" fill="#1C2B3A">${xml(judul)}</text>` +
    `<text x="50%" y="56%" text-anchor="middle" font-family="sans-serif" font-size="${Math.round(tinggi / 24)}" fill="#57534E">${xml(keterangan)}</text>` +
    `<text x="95%" y="95%" text-anchor="end" font-family="sans-serif" font-size="${Math.round(tinggi / 32)}" fill="#57534E">GAMBAR CONTOH</text></svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

function logoContoh(nama, indeks) {
  const latar = warna[indeks % warna.length]
  const inisial = nama.split(/\s+/).map((bagian) => bagian[0]).join('').slice(0, 2).toUpperCase()
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" viewBox="0 0 200 200"><rect width="200" height="200" rx="24" fill="${latar}"/><circle cx="100" cy="78" r="44" fill="#1C2B3A"/><text x="100" y="91" text-anchor="middle" font-family="sans-serif" font-size="38" font-weight="700" fill="#fff">${xml(inisial)}</text><text x="100" y="158" text-anchor="middle" font-family="sans-serif" font-size="15" font-weight="700" fill="#1C2B3A">${xml(nama.toUpperCase())}</text></svg>`
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

export const categories = backup.categories.map((item, i) => ({
  ...item,
  image_url: gambarContoh(item.name, 'Kategori produk', warna[i % warna.length]),
}))

export const brands = backup.brands.map((item, i) => ({
  ...item,
  logo_url: item.logo_url ? logoContoh(item.name, i) : null,
}))

export const rooms = backup.rooms.map((item, i) => ({
  ...item,
  image_cover: gambarContoh(item.name, 'Inspirasi ruangan', warna[i % warna.length], 800, 600),
}))

export const products = backup.products
export const productRooms = backup.product_rooms

const fotoBawaan = backup.product_images.map((foto, i) => {
  const product = products.find((item) => item.id === foto.product_id)
  return {
    ...foto,
    image_url: gambarContoh(product?.name ?? 'Produk', i % 4 === 0 ? 'Foto utama' : 'Detail produk', warna[i % warna.length]),
  }
})

for (const product of products) {
  if (!fotoBawaan.some((foto) => foto.product_id === product.id)) {
    fotoBawaan.push({
      id: product.id * 100 + 1,
      product_id: product.id,
      image_url: gambarContoh(product.name, product.sku, warna[product.id % warna.length]),
      is_primary: true,
      sort_order: 0,
      created_at: product.created_at,
    })
  }
}

export const productImages = fotoBawaan
