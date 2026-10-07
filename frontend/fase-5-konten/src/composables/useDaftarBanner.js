import { computed, nextTick, onMounted, ref } from 'vue'
import { listAdminBanners, patchBanner } from '@/services/adminKontenService'
import { pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'
import { POSISI_BANNER } from '@/utils/banner'

/**
 * Halaman /admin/banner: muat semua banner, aktif/nonaktif (konfirmasi saat menonaktifkan), dan geser urutan
 * di dalam posisi yang sama. Status (Tayang, Terjadwal, ...) selalu diambil dari jawaban API.
 */
export function useDaftarBanner() {
  const ui = useUiStore()
  const status = ref('loading') // loading | ready | error
  const daftar = ref([])
  const pesan = ref('')
  const konfirmasi = ref(null) // banner yang akan dinonaktifkan
  const sedangNonaktif = ref(false)

  async function muat() {
    status.value = 'loading'
    try {
      daftar.value = await listAdminBanners()
      status.value = 'ready'
    } catch (error) {
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }
  onMounted(muat)

  const perPosisi = computed(() =>
    Object.fromEntries(
      Object.keys(POSISI_BANNER).map((posisi) => [
        posisi,
        daftar.value.filter((b) => b.position === posisi).sort((a, b) => a.sort_order - b.sort_order || a.id - b.id),
      ]),
    ),
  )

  // Data baris diganti jawaban API (status ikut dihitung ulang server).
  const ganti = (baris, hasil) => Object.assign(baris, hasil)

  async function aktifkan(baris) {
    baris.is_active = true
    baris.sibuk = true
    try {
      ganti(baris, await patchBanner(baris.id, { is_active: true }))
      ui.tampilkanToast({ pesan: `Banner “${baris.title}” aktif lagi.`, jenis: 'success' })
    } catch (error) {
      baris.is_active = false
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    } finally {
      baris.sibuk = false
    }
  }

  function ubahAktif(baris, nilai) {
    if (nilai) aktifkan(baris)
    else konfirmasi.value = baris
  }

  async function nonaktifkan() {
    const baris = konfirmasi.value
    sedangNonaktif.value = true
    try {
      ganti(baris, await patchBanner(baris.id, { is_active: false }))
      konfirmasi.value = null
      ui.tampilkanToast({ pesan: `Banner “${baris.title}” tidak tampil lagi di Beranda.`, jenis: 'success' })
    } catch (error) {
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    } finally {
      sedangNonaktif.value = false
    }
  }

  /**
   * Tukar posisi dengan tetangga, lalu beri nomor urut 1..n untuk posisi itu (aman bila ada nilai kembar).
   * Hanya baris yang nomornya berubah yang dikirim. Fokus ikut baris yang dipindah.
   */
  async function geser(baris, arah) {
    const urutan = [...perPosisi.value[baris.position]]
    const i = urutan.indexOf(baris)
    if (!urutan[i + arah]) return
    ;[urutan[i], urutan[i + arah]] = [urutan[i + arah], urutan[i]]
    const berubah = urutan.map((b, n) => ({ b, lama: b.sort_order, baru: n + 1 })).filter((x) => x.lama !== x.baru)
    for (const x of berubah) x.b.sort_order = x.baru
    await nextTick()
    const tombol = (a) => document.querySelector(`[data-urut="${baris.id}:${a}"]:not(:disabled)`)
    ;(tombol(arah) ?? tombol(-arah))?.focus()
    try {
      await Promise.all(berubah.map((x) => patchBanner(x.b.id, { sort_order: x.baru })))
    } catch (error) {
      for (const x of berubah) x.b.sort_order = x.lama
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    }
  }

  return { status, daftar, pesan, muat, perPosisi, konfirmasi, sedangNonaktif, ubahAktif, nonaktifkan, geser }
}
