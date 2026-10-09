import { computed, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { createArticle, getAdminArticle, setArticleStatus, updateArticle } from '@/services/adminKontenService'
import { errorField, pesanError } from '@/services/errors'
import { tautanArtikel } from '@/utils/artikel'
import { buatSlug, POLA_SLUG } from '@/utils/slug'

const KUNCI_DRAF = 'transhome.drafArtikel'
const formKosong = () => ({ type: 'ARTIKEL', title: '', slug: '', thumbnail_url: '', content: '' })

/**
 * State dan aksi editor artikel/promo/event: muat, slug otomatis, simpan isi, terbitkan (simpan dulu bila
 * ada perubahan), batalkan terbit, tautan pratinjau, dan draf saat sesi habis. Tampilan: ArtikelForm.vue.
 */
export function useFormArtikel() {
  const route = useRoute()
  const id = route.params.id ? Number(route.params.id) : null

  const form = reactive(formKosong())
  const tersimpan = ref(null) // versi di server: status, published_at, preview_token, slug
  const errors = ref({})
  const status = ref(id ? 'loading' : 'ready') // loading | ready | tidak-ada | error
  const pesan = ref('')
  const galatUmum = ref('')
  const menyimpan = ref('') // '' | simpan | terbit | batal-terbit
  const dipulihkan = ref(false)
  const slugManual = ref(Boolean(id))
  const awal = ref('')

  const potret = () => JSON.stringify(form)
  const kotor = computed(() => status.value === 'ready' && potret() !== awal.value)
  const terbit = computed(() => tersimpan.value?.status === 'PUBLISHED')

  // Pratinjau memakai versi tersimpan. Draf butuh token; yang sudah terbit cukup alamat publiknya.
  const tautanPratinjau = computed(() => {
    const a = tersimpan.value
    if (!a) return null
    return terbit.value ? tautanArtikel(a.slug) : tautanArtikel(a.slug, a.preview_token)
  })

  function isi(a) {
    tersimpan.value = a
    Object.assign(form, { type: a.type, title: a.title, slug: a.slug, thumbnail_url: a.thumbnail_url ?? '', content: a.content ?? '' })
  }

  function simpanDraf() {
    try {
      sessionStorage.setItem(KUNCI_DRAF, JSON.stringify({ path: route.path, form, slugManual: slugManual.value }))
    } catch {
      // Penyimpanan penuh atau diblokir: draf tidak tersimpan.
    }
  }

  function pulihkanDraf() {
    try {
      const draf = JSON.parse(sessionStorage.getItem(KUNCI_DRAF) ?? 'null')
      if (!draf || draf.path !== route.path) return
      sessionStorage.removeItem(KUNCI_DRAF)
      Object.assign(form, draf.form)
      slugManual.value = draf.slugManual
      dipulihkan.value = true
    } catch {
      // Draf rusak: abaikan.
    }
  }

  async function muat() {
    if (!id) return
    status.value = 'loading'
    try {
      isi(await getAdminArticle(id))
      awal.value = potret()
      pulihkanDraf()
      status.value = 'ready'
    } catch (error) {
      if (error?.response?.status === 404) status.value = 'tidak-ada'
      else {
        pesan.value = pesanError(error)
        status.value = 'error'
      }
    }
  }

  function ubahJudul(nilai) {
    form.title = nilai
    if (!slugManual.value) form.slug = buatSlug(nilai)
  }

  function ubahSlug(nilai) {
    slugManual.value = true
    form.slug = nilai
  }

  function validasi() {
    const e = {}
    if (!form.title.trim()) e.title = 'Judul wajib diisi'
    if (!form.slug) e.slug = 'Slug wajib diisi'
    else if (!POLA_SLUG.test(form.slug)) e.slug = 'Slug hanya boleh huruf kecil, angka, dan tanda hubung'
    if (!form.content.trim()) e.content = 'Isi wajib diisi'
    errors.value = e
    return Object.keys(e).length === 0
  }

  const payload = () => ({
    type: form.type, title: form.title.trim(), slug: form.slug, thumbnail_url: form.thumbnail_url || null, content: form.content.trim(),
  })

  // Isi disimpan tanpa mengubah status terbit. Artikel baru selalu mulai sebagai draf.
  async function kirimIsi() {
    const hasil = tersimpan.value ? await updateArticle(tersimpan.value.id, payload()) : await createArticle(payload())
    tersimpan.value = hasil
    awal.value = potret()
    return hasil
  }

  /** Jalankan satu aksi; galat 422/409 dipetakan ke field, sisanya jadi pesan di atas form. Hasil: true bila berhasil. */
  async function jalankan(mode, aksi) {
    if (menyimpan.value) return false
    galatUmum.value = ''
    menyimpan.value = mode
    try {
      await aksi()
      return true
    } catch (error) {
      const statusHttp = error?.response?.status
      if (statusHttp === 422 || statusHttp === 409) {
        errors.value = errorField(error)
        galatUmum.value = error.response.data?.message ?? 'Periksa kembali isian yang ditandai.'
      } else if (statusHttp !== 401) {
        galatUmum.value = pesanError(error)
      }
      return false
    } finally {
      menyimpan.value = ''
    }
  }

  const simpan = () => jalankan('simpan', kirimIsi)

  /** Terbitkan: artikel baru atau yang punya perubahan disimpan dulu, lalu status → PUBLISHED. */
  const terbitkan = () =>
    jalankan('terbit', async () => {
      if (!tersimpan.value || kotor.value) await kirimIsi()
      tersimpan.value = await setArticleStatus(tersimpan.value.id, 'PUBLISHED')
    })

  /** Kembali jadi draf: tidak tampil di website, isi tetap tersimpan. */
  const batalkanTerbit = () =>
    jalankan('batal-terbit', async () => {
      tersimpan.value = await setArticleStatus(tersimpan.value.id, 'DRAFT')
    })

  if (!id) {
    awal.value = potret()
    pulihkanDraf()
  }

  return {
    id, form, tersimpan, errors, status, pesan, galatUmum, menyimpan, dipulihkan, kotor, terbit, tautanPratinjau,
    muat, ubahJudul, ubahSlug, validasi, simpan, terbitkan, batalkanTerbit, simpanDraf,
  }
}
