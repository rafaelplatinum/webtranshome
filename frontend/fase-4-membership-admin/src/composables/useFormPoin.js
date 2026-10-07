import { computed, reactive, ref, watch } from 'vue'
import { checkReference, createPoint } from '@/services/adminMemberService'
import { errorField, pesanError } from '@/services/errors'
import { formatAngka, formatRupiah, formatTanggal, parseRupiah } from '@/utils/format'
import { hitungPoinBelanja, parsePoin, rapikanNota, tandaPoin } from '@/utils/poin'

const KUNCI_DRAF = 'transhome.drafPoin'
const KOSONG = { reference_no: '', purchase_amount: '', points: '', note: '' }

/**
 * State dan aturan form input poin (/admin/poin/input). `member` (terpilih, berisi `balance`) dan `rasio`
 * (site_settings.point_ratio_rupiah) berupa ref/computed dari halaman. Pengecekan di browser untuk umpan balik
 * cepat; backend tetap menolak nota ganda (409) dan saldo kurang (422), dan pesannya dipetakan ke field.
 */
export function useFormPoin(member, rasio) {
  const jenis = ref('EARN') // EARN | REDEEM | ADJUST
  const form = reactive({ ...KOSONG })
  const errors = ref({})
  const galatUmum = ref('')
  const statusNota = ref('') // '' | memeriksa | tersedia | dipakai
  const notaDiperiksa = ref('')
  const pesanNotaDipakai = ref('')
  const menyimpan = ref(false)
  const dipulihkan = ref(false)

  const total = computed(() => parseRupiah(form.purchase_amount))
  // null bila rasio belum dimuat: poin tetap dihitung backend saat disimpan.
  const poinBelanja = computed(() => hitungPoinBelanja(total.value, rasio.value))
  const angkaPoin = computed(() => parsePoin(form.points))
  // Perubahan saldo yang akan dicatat (negatif untuk penukaran).
  const perubahan = computed(() => {
    if (jenis.value === 'EARN') return poinBelanja.value ?? 0
    if (Number.isNaN(angkaPoin.value)) return 0
    return jenis.value === 'REDEEM' ? -Math.abs(angkaPoin.value) : angkaPoin.value
  })
  const saldo = computed(() => member.value?.balance ?? 0)
  const saldoSetelah = computed(() => saldo.value + perubahan.value)
  const kotor = computed(() => Object.keys(KOSONG).some((k) => String(form[k]).trim() !== ''))

  // Isi ConfirmModal sebelum simpan, mis. "Tambah +4 poin ke TF-000123?"
  const ringkasan = computed(() => {
    const m = member.value
    if (!m) return { judul: '', pesan: '' }
    const jadi = `Saldo ${m.full_name} menjadi ${formatAngka(saldoSetelah.value)} poin.`
    const catatan = form.note.trim()
    if (jenis.value === 'EARN') {
      const nota = `Nota ${form.reference_no} · belanja ${formatRupiah(total.value)}.`
      if (poinBelanja.value == null) return { judul: `Catat nota untuk ${m.member_code}?`, pesan: `${nota} Poin dihitung otomatis saat disimpan.` }
      return { judul: `Tambah ${tandaPoin(poinBelanja.value)} poin ke ${m.member_code}?`, pesan: `${nota} ${jadi}` }
    }
    if (jenis.value === 'REDEEM') return { judul: `Kurangi ${Math.abs(perubahan.value)} poin dari ${m.member_code}?`, pesan: `Penukaran: ${catatan}. ${jadi}` }
    return { judul: `Sesuaikan ${tandaPoin(perubahan.value)} poin untuk ${m.member_code}?`, pesan: `Alasan: ${catatan}. ${jadi}` }
  })

  /** Keluar dari field total belanja: "4250000" → "4.250.000". */
  function rapikanTotal() {
    if (total.value) form.purchase_amount = formatAngka(total.value)
  }

  function hapusGalat(kunci) {
    if (!(kunci in errors.value)) return
    const sisa = { ...errors.value }
    delete sisa[kunci]
    errors.value = sisa
  }

  // Nilai yang sama setelah dirapikan (huruf besar, "4250000" → "4.250.000") tidak dihitung sebagai perubahan.
  const bentukBaku = { reference_no: rapikanNota, purchase_amount: parseRupiah }
  const baku = (kunci, nilai) => (bentukBaku[kunci] ?? String)(nilai)

  // Isian diubah → pesan galat field itu hilang; nota yang diubah perlu dicek ulang.
  watch(
    () => ({ ...form }),
    (baru, lama) => {
      for (const kunci of Object.keys(KOSONG)) if (baku(kunci, baru[kunci]) !== baku(kunci, lama[kunci])) hapusGalat(kunci)
      if (rapikanNota(baru.reference_no) !== notaDiperiksa.value && statusNota.value !== 'memeriksa') statusNota.value = ''
    },
  )
  watch(jenis, () => {
    errors.value = {}
    galatUmum.value = ''
  })

  async function cekKeServer(nota) {
    statusNota.value = 'memeriksa'
    try {
      const hasil = await checkReference(nota)
      if (rapikanNota(form.reference_no) !== nota) return true // diubah lagi selama dicek
      notaDiperiksa.value = nota
      if (hasil.available) {
        statusNota.value = 'tersedia'
        return true
      }
      statusNota.value = 'dipakai'
      pesanNotaDipakai.value = `Nomor nota ini sudah pernah diinput pada ${formatTanggal(hasil.used_at)}${hasil.member_code ? ` (${hasil.member_code})` : ''}.`
      errors.value = { ...errors.value, reference_no: pesanNotaDipakai.value }
      return false
    } catch {
      // Belum bisa dicek (mis. jaringan): backend tetap menolak nota ganda saat disimpan.
      statusNota.value = ''
      return true
    }
  }

  let cekBerjalan = null // { nota, janji }: keluar field lalu langsung klik Simpan memakai pengecekan yang sama

  /** Keluar dari field nota: rapikan lalu cek ke server. Hasil false bila nota sudah pernah diinput. */
  async function periksaNota() {
    const nota = rapikanNota(form.reference_no)
    if (form.reference_no !== nota) form.reference_no = nota
    if (!nota) {
      statusNota.value = ''
      return true
    }
    if (cekBerjalan?.nota === nota) return cekBerjalan.janji
    if (nota === notaDiperiksa.value && ['tersedia', 'dipakai'].includes(statusNota.value)) return statusNota.value === 'tersedia'
    const janji = cekKeServer(nota)
    cekBerjalan = { nota, janji }
    try {
      return await janji
    } finally {
      if (cekBerjalan?.janji === janji) cekBerjalan = null
    }
  }

  function validasi() {
    const e = {}
    if (jenis.value === 'EARN') {
      const nota = rapikanNota(form.reference_no)
      if (!nota) e.reference_no = 'Nomor nota wajib diisi'
      else if (statusNota.value === 'dipakai' && nota === notaDiperiksa.value) e.reference_no = pesanNotaDipakai.value
      if (!total.value) e.purchase_amount = 'Total belanja wajib diisi'
      else if (poinBelanja.value === 0) e.purchase_amount = `Total belanja di bawah ${formatRupiah(rasio.value)}, jadi belum dapat poin`
    } else {
      const n = angkaPoin.value
      if (jenis.value === 'REDEEM') {
        if (Number.isNaN(n) || n <= 0) e.points = 'Isi jumlah poin yang ditukar, angka bulat lebih dari 0'
        else if (n > saldo.value) e.points = `Saldo poin tidak cukup. Saldo saat ini ${formatAngka(saldo.value)} poin.`
      } else if (Number.isNaN(n) || n === 0) e.points = 'Isi jumlah poin, mis. 2 untuk menambah atau -2 untuk mengurangi'
      else if (saldo.value + n < 0) e.points = `Saldo tidak boleh kurang dari 0. Saldo saat ini ${formatAngka(saldo.value)} poin.`
      if (!form.note.trim()) e.note = jenis.value === 'REDEEM' ? 'Catatan wajib diisi' : 'Alasan wajib diisi'
    }
    errors.value = e
    return Object.keys(e).length === 0
  }

  /** Sebelum ringkasan konfirmasi: validasi + nota dicek ke server. */
  async function siapkan() {
    galatUmum.value = ''
    if (!validasi()) return false
    if (jenis.value === 'EARN' && !(await periksaNota())) return false
    return true
  }

  /** POST /admin/points → { transaction, balance }, atau null bila gagal (pesan sudah dipasang). */
  async function kirim() {
    if (menyimpan.value) return null
    menyimpan.value = true
    galatUmum.value = ''
    const payload =
      jenis.value === 'EARN'
        ? { user_id: member.value.id, type: 'EARN', reference_no: rapikanNota(form.reference_no), purchase_amount: total.value, note: form.note.trim() || null }
        : { user_id: member.value.id, type: jenis.value, points: jenis.value === 'REDEEM' ? Math.abs(angkaPoin.value) : angkaPoin.value, note: form.note.trim() }
    try {
      return await createPoint(payload)
    } catch (error) {
      const status = error?.response?.status
      if (status === 409 || status === 422) {
        const e = errorField(error)
        if (status === 409) {
          // Nota ganda lolos dari cek browser (mis. diinput admin lain barusan): ditolak server.
          notaDiperiksa.value = rapikanNota(form.reference_no)
          statusNota.value = 'dipakai'
          pesanNotaDipakai.value = e.reference_no ?? pesanError(error)
          e.reference_no = pesanNotaDipakai.value
        }
        errors.value = e
        if (e.user_id || !Object.keys(e).length) galatUmum.value = e.user_id ?? pesanError(error)
      } else if (status !== 401 && status !== 403) {
        // 403 & sesi habis sudah ditangani interceptor; selain itu jelaskan bahwa data belum tersimpan.
        galatUmum.value = `Transaksi belum tersimpan. ${pesanError(error)}`
      }
      return null
    } finally {
      menyimpan.value = false
    }
  }

  function kosongkan() {
    Object.assign(form, KOSONG)
    errors.value = {}
    galatUmum.value = ''
    statusNota.value = ''
    notaDiperiksa.value = ''
    dipulihkan.value = false
  }

  /** Sesi habis di tengah mengisi: isian disimpan di tab ini, dipulihkan setelah masuk lagi. */
  function simpanDraf() {
    try {
      sessionStorage.setItem(KUNCI_DRAF, JSON.stringify({ user_id: member.value?.id ?? null, jenis: jenis.value, form: { ...form } }))
    } catch {
      // Storage diblokir: isian tidak bisa diselamatkan.
    }
  }

  function pulihkanDraf(userId) {
    try {
      const draf = JSON.parse(sessionStorage.getItem(KUNCI_DRAF) ?? 'null')
      sessionStorage.removeItem(KUNCI_DRAF)
      if (!draf?.form || (draf.user_id && String(draf.user_id) !== String(userId))) return
      jenis.value = draf.jenis
      Object.assign(form, KOSONG, draf.form)
      dipulihkan.value = true
    } catch {
      // Draf rusak: mulai dari form kosong.
    }
  }

  return {
    jenis, form, errors, galatUmum, statusNota, menyimpan, dipulihkan,
    total, poinBelanja, angkaPoin, perubahan, saldoSetelah, kotor, ringkasan,
    periksaNota, rapikanTotal, siapkan, kirim, kosongkan, simpanDraf, pulihkanDraf,
  }
}
