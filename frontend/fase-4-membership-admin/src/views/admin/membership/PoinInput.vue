<script setup>
import { computed, nextTick, onMounted, ref, useId } from 'vue'
import { useRoute } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import CariMember from '@/components/admin/poin/CariMember.vue'
import IsianBelanja from '@/components/admin/poin/IsianBelanja.vue'
import IsianTukarPoin from '@/components/admin/poin/IsianTukarPoin.vue'
import MemberTerpilih from '@/components/admin/poin/MemberTerpilih.vue'
import PilihanJenis from '@/components/admin/poin/PilihanJenis.vue'
import RiwayatRingkas from '@/components/admin/poin/RiwayatRingkas.vue'
import { useFormPoin } from '@/composables/useFormPoin'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { useMemberTerpilih } from '@/composables/useMemberTerpilih'
import { usePermission } from '@/composables/usePermission'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'
import { formatAngka } from '@/utils/format'
import { tandaPoin } from '@/utils/poin'

// Alur sketsa "AdminPoin": 1. cari member → 2. catat transaksi, dengan riwayat member di samping.
const route = useRoute()
const { can } = usePermission()
const ui = useUiStore()
const settings = useSettingsStore()
const idLangkah1 = useId()
const idLangkah2 = useId()

const pilihan = useMemberTerpilih()
const { member } = pilihan
const konfirmasi = ref(false)
const cari = ref(null)
const riwayat = ref(null)
const formEl = ref(null)

const rasio = computed(() => Number(settings.data?.point_ratio_rupiah) || null)
const f = useFormPoin(member, rasio)
const { form, errors } = f
const keluar = useKonfirmasiKeluar(f.kotor, { saatSesiHabis: f.simpanDraf })

// Penyesuaian (ADJUST) hanya untuk izin point.adjust (usulan: Super Admin).
const opsiJenis = computed(() => [
  { value: 'EARN', label: 'Belanja (EARN)' },
  { value: 'REDEEM', label: 'Tukar poin (REDEEM)' },
  ...(can('point.adjust') ? [{ value: 'ADJUST', label: 'Penyesuaian (ADJUST)' }] : []),
])

async function fokusIsianPertama() {
  await nextTick()
  formEl.value?.querySelector('input[type="text"], input:not([type])')?.focus()
}

async function fokusCari() {
  await nextTick()
  cari.value?.fokus()
}

async function fokusGalat() {
  await nextTick()
  formEl.value?.querySelector('[aria-invalid="true"]')?.focus()
}

// Setelah member terpilih: langsung ke isian pertama (akun nonaktif: ke tombol Ganti member).
async function fokusSetelahPilih(m) {
  if (m.is_active) return fokusIsianPertama()
  await nextTick()
  document.querySelector('[data-ganti-member]')?.focus()
}

async function pilih(m) {
  await pilihan.pilih(m)
  fokusSetelahPilih(m)
}

async function gantiMember() {
  await pilihan.lepas()
  fokusCari()
}

onMounted(async () => {
  settings.load()
  if (!route.query.member) {
    f.pulihkanDraf(null)
    fokusCari()
    return
  }
  const m = await pilihan.muatDariUrl()
  if (!m) return fokusCari()
  f.pulihkanDraf(m.id)
  fokusSetelahPilih(m)
})

async function simpan() {
  if (f.menyimpan.value || konfirmasi.value) return
  if (await f.siapkan()) konfirmasi.value = true
  else fokusGalat()
}

async function konfirmasiSimpan() {
  const hasil = await f.kirim()
  konfirmasi.value = false
  if (!hasil) {
    if (Object.keys(errors.value).length) pilihan.perbaruiSaldo()
    fokusGalat()
    return
  }
  const t = hasil.transaction
  pilihan.aturSaldo(hasil.balance)
  riwayat.value?.tambah(t)
  const apa = { EARN: `${tandaPoin(t.points)} poin`, REDEEM: `Penukaran ${Math.abs(t.points)} poin`, ADJUST: `Penyesuaian ${tandaPoin(t.points)} poin` }[t.type]
  ui.tampilkanToast({ pesan: `${apa} dicatat untuk ${member.value.member_code}. Saldo sekarang ${formatAngka(hasil.balance)} poin.`, jenis: 'success' })
  // Form kosong untuk nota berikutnya, member tetap terpilih.
  f.kosongkan()
  fokusIsianPertama()
}

function batal() {
  f.kosongkan()
  fokusIsianPertama()
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <p class="sr-only" aria-live="polite">{{ pilihan.pengumuman.value }}</p>

    <section :aria-labelledby="idLangkah1" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-5">
      <h2 :id="idLangkah1" class="text-base font-extrabold">1. Cari member</h2>
      <div v-if="pilihan.memuat.value" class="flex flex-col gap-2" aria-busy="true">
        <span class="sr-only">Memuat member…</span>
        <Skeleton class="h-20 rounded-lg" />
      </div>
      <MemberTerpilih v-else-if="member" :member="member" @ganti="gantiMember" />
      <template v-else>
        <p v-if="pilihan.galat.value" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ pilihan.galat.value }}</p>
        <CariMember ref="cari" @pilih="pilih" />
      </template>
    </section>

    <section v-if="!member" :aria-labelledby="idLangkah2" class="flex flex-col gap-2 rounded-xl border border-dashed border-border-strong bg-surface p-5">
      <h2 :id="idLangkah2" class="text-base font-extrabold text-muted">2. Catat transaksi</h2>
      <p class="text-sm text-muted">Cari dan pilih member dulu. Kode member disebutkan pelanggan ke kasir saat belanja.</p>
    </section>

    <div v-else class="grid items-start gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
      <section :aria-labelledby="idLangkah2" class="flex flex-col gap-4 rounded-xl border border-border bg-surface p-5">
        <h2 :id="idLangkah2" class="text-base font-extrabold">2. Catat transaksi</h2>

        <p v-if="!member.is_active" class="text-sm text-muted">Poin tidak bisa dicatat untuk akun nonaktif.</p>

        <form v-else ref="formEl" class="flex flex-col gap-4" novalidate data-form-poin @submit.prevent="simpan">
          <p v-if="f.dipulihkan.value" role="status" class="rounded-md bg-info-soft px-3 py-2.5 text-sm font-medium text-info-ink">
            Isian sebelum sesi habis sudah dipulihkan. Periksa lagi, lalu simpan.
          </p>
          <p v-if="f.galatUmum.value" role="alert" class="rounded-md bg-danger-soft px-3 py-2.5 text-sm font-semibold text-danger-ink">{{ f.galatUmum.value }}</p>

          <PilihanJenis v-model="f.jenis.value" :opsi="opsiJenis" />

          <IsianBelanja
            v-if="f.jenis.value === 'EARN'"
            v-model:nota="form.reference_no"
            v-model:total="form.purchase_amount"
            v-model:catatan="form.note"
            :errors="errors"
            :status-nota="f.statusNota.value"
            :poin="f.poinBelanja.value"
            :total-angka="f.total.value"
            :rasio="rasio"
            @periksa-nota="f.periksaNota"
            @rapikan-total="f.rapikanTotal"
          />
          <IsianTukarPoin
            v-else
            v-model:poin="form.points"
            v-model:catatan="form.note"
            :jenis="f.jenis.value"
            :errors="errors"
            :saldo="member.balance"
            :perubahan="f.perubahan.value"
            :saldo-setelah="f.saldoSetelah.value"
            :poin-valid="!Number.isNaN(f.angkaPoin.value) && f.angkaPoin.value !== 0"
          />

          <div class="flex flex-wrap gap-2">
            <Button type="submit" :loading="f.menyimpan.value">Simpan transaksi</Button>
            <Button variant="outline" :disabled="f.menyimpan.value" @click="batal">Batal</Button>
          </div>
          <p class="text-xs text-muted">Notifikasi WhatsApp hanya terkirim jika member punya nomor HP dan menyetujui komunikasi.</p>
        </form>
      </section>

      <RiwayatRingkas ref="riwayat" :user-id="member.id" />
    </div>

    <ConfirmModal
      v-model:open="konfirmasi"
      :title="f.ringkasan.value.judul"
      :message="f.ringkasan.value.pesan"
      confirm-label="Ya, simpan"
      cancel-label="Periksa lagi"
      :loading="f.menyimpan.value"
      @confirm="konfirmasiSimpan"
    />
    <ConfirmModal
      :open="keluar.terbuka.value"
      title="Buang isian transaksi?"
      message="Transaksi yang belum disimpan akan hilang."
      confirm-label="Buang"
      cancel-label="Tetap di sini"
      danger
      @confirm="keluar.buang"
      @cancel="keluar.tetap"
    />
  </div>
</template>
