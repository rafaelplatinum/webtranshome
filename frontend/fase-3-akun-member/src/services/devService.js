import http from '@/services/http'

// Khusus mode mock + development. Tidak ada di API asli.

/** Email terakhir yang "dikirim" mock (verifikasi & reset password): [{ id, to, jenis, subject, link, created_at }]. */
export async function getMockEmails() {
  const { data } = await http.get('/dev/emails')
  return data.data
}
