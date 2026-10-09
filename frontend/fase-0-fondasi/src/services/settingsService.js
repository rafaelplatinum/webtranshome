import http from './http'

/** Pengaturan toko yang boleh dilihat publik: nomor WhatsApp, alamat, jam buka, link peta. */
export async function getPublicSettings() {
  const { data } = await http.get('/settings/public')
  return data.data
}
