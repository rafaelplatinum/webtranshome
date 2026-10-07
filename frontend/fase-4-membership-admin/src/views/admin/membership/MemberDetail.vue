<script setup>
import { computed, nextTick, onMounted, ref, useId } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeft, Coins } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import DataMember from '@/components/admin/member/DataMember.vue'
import KartuProfilMember from '@/components/admin/member/KartuProfilMember.vue'
import RiwayatMember from '@/components/admin/member/RiwayatMember.vue'
import SyncMember from '@/components/admin/member/SyncMember.vue'
import { usePermission } from '@/composables/usePermission'
import { useSyncQontak } from '@/composables/useSyncQontak'
import { tujuanKembali } from '@/composables/useTujuanKembali'
import { getMember, patchMember } from '@/services/adminMemberService'
import { pesanError } from '@/services/errors'
import { useSettingsStore } from '@/stores/settings'
import { useUiStore } from '@/stores/ui'
import { formatAngka, formatRupiah } from '@/utils/format'

const route = useRoute()
const { can } = usePermission()
const { cobaLagi } = useSyncQontak()
const ui = useUiStore()
const settings = useSettingsStore()
const idSaldo = useId()

const status = ref('loading') // loading | ready | tidak-ada | error
const member = ref(null)
const pesan = ref('')
const konfirmasiNonaktif = ref(false)
const mengubahStatus = ref(false)

const kembali = computed(() => tujuanKembali('member', '/admin/member'))
const rasio = computed(() => settings.data?.point_ratio_rupiah)

async function muat() {
  status.value = 'loading'
  try {
    member.value = await getMember(route.params.id)
    status.value = 'ready'
  } catch (error) {
    if (error?.response?.status === 404) {
      status.value = 'tidak-ada'
      return
    }
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

onMounted(() => {
  settings.load()
  muat()
})

// Nonaktif: member tidak bisa masuk. Mengaktifkan lagi tidak perlu konfirmasi.
async function ubahStatus(aktif) {
  mengubahStatus.value = true
  try {
    const baru = await patchMember(member.value.id, { is_active: aktif })
    member.value = { ...baru, qontak: member.value.qontak }
    konfirmasiNonaktif.value = false
    mengubahStatus.value = false
    ui.tampilkanToast({ pesan: aktif ? `Akun ${baru.full_name} aktif lagi.` : `Akun ${baru.full_name} dinonaktifkan.`, jenis: 'success' })
    // Tombol lawannya menggantikan tombol yang diklik: fokus dipindah ke sana (tombol loading tidak bisa difokus).
    await nextTick()
    document.querySelector(aktif ? '[data-tombol-nonaktif]' : '[data-tombol-aktif]')?.focus()
  } catch (error) {
    // 401/403/5xx/jaringan pada PATCH sudah ditangani interceptor.
    const kode = error?.response?.status
    if (kode && kode < 500 && ![401, 403].includes(kode)) ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    mengubahStatus.value = false
  }
}

async function syncUlang() {
  await cobaLagi(member.value.qontak, member.value.full_name)
  // Tombol hilang bila berhasil: fokus pindah ke judul kartu supaya tidak lepas ke awal halaman.
  if (member.value?.qontak?.sync_status !== 'FAILED') document.querySelector('[data-kartu-sync]')?.focus()
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <RouterLink
      :to="kembali"
      class="inline-flex min-h-11 items-center gap-1.5 self-start rounded-sm text-sm font-semibold text-muted transition-colors duration-150 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-0"
    >
      <ArrowLeft class="size-4" aria-hidden="true" />
      Daftar member
    </RouterLink>

    <div v-if="status === 'loading'" class="flex flex-col gap-4" aria-busy="true">
      <span class="sr-only">Memuat member…</span>
      <Skeleton class="h-24 rounded-xl" />
      <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_20rem]">
        <Skeleton class="h-56 rounded-xl" />
        <Skeleton class="h-56 rounded-xl" />
      </div>
    </div>

    <section v-else-if="status === 'tidak-ada'" class="mx-auto flex max-w-md flex-col items-center gap-4 py-16 text-center">
      <h2 class="text-2xl font-bold">Member tidak ditemukan</h2>
      <p class="text-muted">Member ini tidak ada atau alamatnya salah.</p>
      <Button to="/admin/member">Kembali ke daftar member</Button>
    </section>

    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

    <template v-else>
      <KartuProfilMember :member="member">
        <template #aksi>
          <Button v-if="member.is_active && can('point.create')" :to="{ path: '/admin/poin/input', query: { member: member.id } }">
            <template #icon><Coins class="size-4.5" aria-hidden="true" /></template>
            Input poin untuk member ini
          </Button>
          <template v-if="can('member.update')">
            <Button v-if="member.is_active" variant="outline" data-tombol-nonaktif @click="konfirmasiNonaktif = true">Nonaktifkan akun</Button>
            <Button v-else variant="outline" data-tombol-aktif :loading="mengubahStatus" @click="ubahStatus(true)">Aktifkan akun</Button>
          </template>
        </template>
      </KartuProfilMember>

      <p v-if="!member.is_active" role="status" class="rounded-md bg-danger-soft px-4 py-3 text-sm font-medium text-danger-ink">
        Akun ini nonaktif: member tidak bisa masuk dan poin tidak bisa dicatat. Saldo dan riwayat poin tetap tersimpan.
      </p>

      <div class="grid items-start gap-4 lg:grid-cols-[minmax(0,1fr)_20rem]">
        <DataMember :member="member" />
        <div class="flex flex-col gap-4">
          <section :aria-labelledby="idSaldo" class="flex flex-col gap-1.5 rounded-xl border border-border bg-surface p-5">
            <h2 :id="idSaldo" class="text-sm font-semibold text-muted">Saldo poin</h2>
            <p class="text-4xl leading-none font-extrabold tabular-nums" data-saldo>
              {{ formatAngka(member.balance) }}<span class="text-base font-semibold text-muted"> poin</span>
            </p>
            <p v-if="rasio" class="text-sm text-muted">1 poin tiap belanja {{ formatRupiah(rasio) }}.</p>
          </section>
          <SyncMember :kontak="member.qontak" :punya-hp="Boolean(member.phone_number)" @coba-lagi="syncUlang" />
        </div>
      </div>

      <RiwayatMember :user-id="member.id" />
    </template>

    <ConfirmModal
      v-if="member"
      v-model:open="konfirmasiNonaktif"
      title="Nonaktifkan akun member?"
      :message="`${member.full_name} (${member.member_code}) tidak bisa masuk ke akunnya lagi. Saldo dan riwayat poin tetap tersimpan, dan akun bisa diaktifkan kembali kapan saja.`"
      confirm-label="Nonaktifkan"
      danger
      :loading="mengubahStatus"
      @confirm="ubahStatus(false)"
    />
  </div>
</template>
