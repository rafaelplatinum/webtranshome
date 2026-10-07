import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMember } from '@/services/adminMemberService'
import { pesanError } from '@/services/errors'

/**
 * Member terpilih di halaman input poin. Disimpan di URL (?member=<id>) supaya tetap terpilih setelah refresh
 * dan bisa dibuka langsung dari detail member ("Input poin untuk member ini").
 */
export function useMemberTerpilih() {
  const route = useRoute()
  const router = useRouter()
  const member = ref(null)
  const memuat = ref(false)
  const galat = ref('')
  // Diumumkan pembaca layar saat member dipilih dari hasil cari.
  const pengumuman = ref('')

  async function pilih(m, { umumkan = true } = {}) {
    member.value = m
    galat.value = ''
    if (String(route.query.member) !== String(m.id)) await router.replace({ query: { ...route.query, member: m.id } })
    if (umumkan) pengumuman.value = `Member dipilih: ${m.full_name}, ${m.member_code}.`
  }

  /** "Ganti member": kosongkan pilihan dan hapus ?member dari URL. */
  async function lepas() {
    member.value = null
    pengumuman.value = ''
    const query = { ...route.query }
    delete query.member
    await router.replace({ query })
  }

  /** Muat member dari ?member=<id>. Hasil member, atau null dengan pesan di `galat`. */
  async function muatDariUrl() {
    memuat.value = true
    try {
      const m = await getMember(route.query.member)
      await pilih(m, { umumkan: false })
      return m
    } catch (error) {
      galat.value = error?.response?.status === 404 ? 'Member tidak ditemukan. Cari dengan kode member, email, atau nomor HP.' : pesanError(error)
      return null
    } finally {
      memuat.value = false
    }
  }

  /** Setelah transaksi tersimpan: saldo dari respons server. */
  function aturSaldo(balance) {
    member.value = { ...member.value, balance }
  }

  /** Ditolak server (nota ganda / saldo kurang): saldo mungkin baru diubah admin lain, jadi dimuat ulang. */
  async function perbaruiSaldo() {
    const id = member.value?.id
    try {
      const terbaru = await getMember(id)
      if (member.value?.id === id) member.value = { ...member.value, balance: terbaru.balance, is_active: terbaru.is_active }
    } catch {
      // Saldo lama tetap tampil; backend tetap memeriksa saat simpan berikutnya.
    }
  }

  return { member, memuat, galat, pengumuman, pilih, lepas, muatDariUrl, aturSaldo, perbaruiSaldo }
}
