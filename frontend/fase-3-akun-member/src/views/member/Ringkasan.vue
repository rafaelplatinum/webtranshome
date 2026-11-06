<script setup>
import { computed, onMounted, ref, useId } from 'vue'
import { RouterLink } from 'vue-router'
import Skeleton from '@/components/ui/Skeleton.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import DaftarPoin from '@/components/akun/DaftarPoin.vue'
import KartuMember from '@/components/akun/KartuMember.vue'
import { useWhatsApp } from '@/composables/useWhatsApp'
import { pesanError } from '@/services/errors'
import { getPoints } from '@/services/memberService'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { formatRupiah } from '@/utils/format'

const auth = useAuthStore()
const settings = useSettingsStore()
// Pesan WhatsApp diawali nama + member_code (FR-S).
const { generalLink } = useWhatsApp()

const idSaldo = useId()
const idTransaksi = useId()
const status = ref('loading') // loading | ready | error
const transaksi = ref([])
const saldo = ref(0)
const pesan = ref('')
const angka = new Intl.NumberFormat('id-ID')

async function muat() {
  status.value = 'loading'
  try {
    const hasil = await getPoints({ per_page: 3 })
    transaksi.value = hasil.data
    saldo.value = hasil.meta.balance
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

onMounted(muat)

// Rasio dari pengaturan toko (site_settings.point_ratio_rupiah), tidak di-hardcode.
const teksRasio = computed(() => {
  const rasio = settings.data?.point_ratio_rupiah
  return rasio ? `1 poin tiap belanja ${formatRupiah(rasio)} di toko.` : 'Poin dihitung dari total belanja di toko.'
})

const kelasTautan = 'rounded-sm text-sm font-bold text-primary transition-colors duration-150 hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'
</script>

<template>
  <div class="flex flex-col gap-6">
    <h1 class="text-2xl font-extrabold sm:text-3xl">Halo, {{ auth.user?.full_name }}</h1>

    <div class="grid gap-4 md:grid-cols-2">
      <KartuMember v-if="auth.user" :user="auth.user" />
      <section :aria-labelledby="idSaldo" class="flex flex-col gap-2.5 rounded-2xl border border-border bg-surface p-6 sm:p-7">
        <h2 :id="idSaldo" class="text-sm font-semibold text-muted">Saldo poin</h2>
        <Skeleton v-if="status === 'loading'" class="h-11 w-36" />
        <p v-else-if="status === 'error'" class="text-sm text-muted">
          Saldo belum bisa dimuat.
          <button type="button" :class="kelasTautan" @click="muat">Coba lagi</button>
        </p>
        <p v-else class="text-5xl leading-none font-extrabold tabular-nums" data-saldo>
          {{ angka.format(saldo) }}<span class="text-lg font-semibold text-muted"> poin</span>
        </p>
        <p class="text-sm text-muted">{{ teksRasio }} Tukar poin langsung di kasir.</p>
        <RouterLink to="/akun/poin" :class="[kelasTautan, 'mt-auto inline-flex min-h-11 items-center self-start md:min-h-0']">Lihat riwayat poin</RouterLink>
      </section>
    </div>

    <section :aria-labelledby="idTransaksi" class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-5 sm:p-6">
      <div class="flex items-center justify-between gap-3">
        <h2 :id="idTransaksi" class="text-lg font-extrabold">Transaksi poin terakhir</h2>
        <RouterLink v-if="transaksi.length" to="/akun/poin" :class="[kelasTautan, 'inline-flex min-h-11 items-center md:min-h-0']">Lihat semua</RouterLink>
      </div>
      <div v-if="status === 'loading'" class="flex flex-col gap-2" aria-busy="true">
        <span class="sr-only">Memuat transaksi…</span>
        <Skeleton v-for="n in 3" :key="n" class="h-10" />
      </div>
      <p v-else-if="status === 'error'" role="alert" class="text-sm text-muted">{{ pesan }}</p>
      <p v-else-if="!transaksi.length" class="text-sm text-muted">Belum ada transaksi poin. Sebutkan kode member ke kasir setiap belanja untuk mengumpulkan poin.</p>
      <DaftarPoin v-else :items="transaksi" />
    </section>

    <section aria-label="Bantuan poin" class="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-primary-soft px-5 py-4 sm:px-6">
      <div class="flex flex-col">
        <strong>Poin belum masuk?</strong>
        <span class="text-sm text-primary-hover">Kirim foto nota ke CS, kami cek manual.</span>
      </div>
      <WhatsAppButton :href="generalLink()" label="Chat WhatsApp CS" />
    </section>
  </div>
</template>
