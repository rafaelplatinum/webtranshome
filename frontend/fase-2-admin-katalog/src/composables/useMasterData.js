import { nextTick, onMounted, ref } from 'vue'
import { listMaster, patchMaster } from '@/services/adminMasterService'
import { pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'

/**
 * Halaman master katalog (kategori, brand, ruangan): muat daftar, aktif/nonaktif (konfirmasi saat
 * menonaktifkan), geser urutan, dan form modal tambah/ubah. `jenis` = 'categories' | 'brands' | 'rooms'.
 */
export function useMasterData(jenis) {
  const ui = useUiStore()
  const status = ref('loading')
  const daftar = ref([])
  const pesan = ref('')
  const konfirmasi = ref(null)
  const sedangNonaktif = ref(false)
  const formBuka = ref(false)
  const dipilih = ref(null)

  async function muat() {
    status.value = 'loading'
    try {
      daftar.value = await listMaster(jenis)
      status.value = 'ready'
    } catch (error) {
      pesan.value = pesanError(error)
      status.value = 'error'
    }
  }
  onMounted(muat)

  async function aktifkan(baris) {
    baris.is_active = true
    baris.sibuk = true
    try {
      await patchMaster(jenis, baris.id, { is_active: true })
      ui.tampilkanToast({ pesan: `${baris.name} tampil di website.`, jenis: 'success' })
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
      await patchMaster(jenis, baris.id, { is_active: false })
      baris.is_active = false
      konfirmasi.value = null
      ui.tampilkanToast({ pesan: `${baris.name} tidak tampil lagi di website.`, jenis: 'success' })
    } catch (error) {
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    } finally {
      sedangNonaktif.value = false
    }
  }

  /**
   * Tukar sort_order dengan tetangga di `saudara` (daftar yang sudah urut). Fokus ikut baris yang dipindah:
   * tombol arah yang sama, atau arah sebaliknya bila baris sudah di ujung. Tombol memakai data-urut="<id>:<arah>".
   */
  async function geser(baris, arah, saudara) {
    const tetangga = saudara[saudara.indexOf(baris) + arah]
    if (!tetangga) return
    const [lamaA, lamaB] = [baris.sort_order, tetangga.sort_order]
    baris.sort_order = lamaB
    tetangga.sort_order = lamaA
    await nextTick()
    const tombol = (a) => document.querySelector(`[data-urut="${baris.id}:${a}"]:not(:disabled)`)
    ;(tombol(arah) ?? tombol(-arah))?.focus()
    try {
      await Promise.all([patchMaster(jenis, baris.id, { sort_order: lamaB }), patchMaster(jenis, tetangga.id, { sort_order: lamaA })])
    } catch (error) {
      baris.sort_order = lamaA
      tetangga.sort_order = lamaB
      ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
    }
  }

  function bukaForm(baris = null) {
    dipilih.value = baris
    formBuka.value = true
  }

  function tersimpan(hasil, { baru }) {
    if (baru) daftar.value = [...daftar.value, hasil]
    else daftar.value = daftar.value.map((x) => (x.id === hasil.id ? { ...x, ...hasil } : x))
    ui.tampilkanToast({ pesan: `${hasil.name} disimpan.`, jenis: 'success' })
  }

  return { status, daftar, pesan, muat, konfirmasi, sedangNonaktif, ubahAktif, nonaktifkan, geser, formBuka, dipilih, bukaForm, tersimpan }
}
