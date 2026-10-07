import { onBeforeUnmount } from 'vue'
import { pesanError } from '@/services/errors'
import { getContact, retryContact } from '@/services/qontakService'
import { useUiStore } from '@/stores/ui'

const tunggu = (ms) => new Promise((selesai) => setTimeout(selesai, ms))

/**
 * "Coba lagi" sync kontak ke Qontak. Sync berjalan di background: status langsung PENDING,
 * lalu dicek tiap detik (maks 15 kali) sampai SYNCED / FAILED. Pengecekan berhenti saat halaman ditutup.
 * `kontak` objek reaktif { id (qontak_contacts.id), sync_status, … } yang ikut diperbarui; `sedang` = proses berjalan.
 */
export function useSyncQontak() {
  const ui = useUiStore()
  let aktif = true
  onBeforeUnmount(() => (aktif = false))

  async function cobaLagi(kontak, nama) {
    kontak.sedang = true
    let terkirim = false
    try {
      Object.assign(kontak, await retryContact(kontak.id))
      terkirim = true
      for (let i = 0; i < 15 && kontak.sync_status === 'PENDING' && aktif; i++) {
        await tunggu(1000)
        if (aktif) Object.assign(kontak, await getContact(kontak.id))
      }
      if (kontak.sync_status === 'SYNCED') ui.tampilkanToast({ pesan: `Kontak ${nama} berhasil tersinkron.`, jenis: 'success' })
      else if (kontak.sync_status === 'FAILED') ui.tampilkanToast({ pesan: `Sync ${nama} gagal lagi. Coba beberapa saat lagi.`, jenis: 'warning' })
    } catch (error) {
      const status = error?.response?.status
      // 401: interceptor sudah mengarahkan ke halaman masuk.
      if (status === 401) return
      // Retry sudah terkirim tapi status gagal dicek (GET tidak ditoast interceptor).
      if (terkirim) ui.tampilkanToast({ pesan: 'Sync sedang diproses, tapi statusnya belum bisa dicek. Muat ulang halaman sebentar lagi.', jenis: 'warning' })
      // 5xx/jaringan/403 pada POST sudah ditoast interceptor; selain itu tampilkan pesannya.
      else if (status && status < 500 && status !== 403) ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    } finally {
      kontak.sedang = false
    }
  }

  return { cobaLagi }
}
