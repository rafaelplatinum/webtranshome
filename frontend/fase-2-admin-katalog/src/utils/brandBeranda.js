// Lingkaran "Brand pilihan" di Beranda (KEPUTUSAN #21). Dipakai admin (Brand, Brand di Beranda), API tiruan (validasi 422),
// dan Beranda. Backend memakai batas yang sama.
// Frame desktop: tiap sisi pola sarang lebah 3–2–3 = 8 lingkaran utuh (tidak terpotong) → 16. HP/tablet: 8 teratas.
export const MAKS_LINGKARAN_BRAND = 16
export const LINGKARAN_HP = 8

const urutOtomatis = (a, b) =>
  Number(Boolean(b.logo_url)) - Number(Boolean(a.logo_url)) || (b.product_count ?? 0) - (a.product_count ?? 0) || a.name.localeCompare(b.name, 'id')

/**
 * Brand untuk lingkaran, sudah berurutan (maks 16): yang dicentang admin menurut `home_order`. Belum ada yang dicentang
 * sama sekali → otomatis dari semua brand (berlogo dulu, produk terbanyak, nama). Hanya kirim brand AKTIF ke fungsi ini.
 */
export function brandLingkaran(brands) {
  const dipilih = brands
    .filter((b) => b.show_on_home)
    .sort((a, b) => (a.home_order ?? Infinity) - (b.home_order ?? Infinity) || a.name.localeCompare(b.name, 'id'))
  if (dipilih.length) return { brands: dipilih.slice(0, MAKS_LINGKARAN_BRAND), otomatis: false }
  return { brands: [...brands].sort(urutOtomatis).slice(0, MAKS_LINGKARAN_BRAND), otomatis: true }
}
