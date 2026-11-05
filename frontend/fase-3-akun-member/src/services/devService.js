import http from '@/services/http'

// Khusus mode mock + development. Tidak ada di API asli.

/** Email terakhir yang "dikirim" mock (verifikasi & reset password): [{ id, to, jenis, subject, link, created_at }]. */
export async function getMockEmails() {
  const { data } = await http.get('/dev/emails')
  return data.data
}

/** OTP mock hanya ditampilkan pada panel dev lokal (bukan respons publik API). */
export async function getMockOtps() {
  const { data } = await http.get('/dev/otps')
  return data.data
}
