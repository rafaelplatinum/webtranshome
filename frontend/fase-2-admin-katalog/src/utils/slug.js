/** "Keramik Lantai 40×40 (Abu)" → "keramik-lantai-40x40-abu". Slug tetap bisa diedit admin; backend menolak duplikat (409). */
export function buatSlug(teks) {
  return String(teks ?? '')
    .normalize('NFD')
    .replace(/[̀-ͯ]/g, '')
    .toLowerCase()
    .replace(/×/g, 'x')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

export const POLA_SLUG = /^[a-z0-9]+(?:-[a-z0-9]+)*$/
