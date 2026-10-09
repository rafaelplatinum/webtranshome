<script setup>
import { computed, nextTick, ref } from 'vue'
import { UserPlus, Users } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'
import SkeletonTabel from '@/components/admin/SkeletonTabel.vue'
import PenggunaFilterBar from '@/components/admin/akses/PenggunaFilterBar.vue'
import PenggunaTabel from '@/components/admin/akses/PenggunaTabel.vue'
import TambahAdminModal from '@/components/admin/akses/TambahAdminModal.vue'
import UbahRoleModal from '@/components/admin/akses/UbahRoleModal.vue'
import { useDaftar } from '@/composables/useDaftar'
import { useFilter } from '@/composables/useFilter'
import { usePermission } from '@/composables/usePermission'
import { listAdmins, resetAdminPassword, setAdminActive } from '@/services/adminAksesService'
import { pesanError } from '@/services/errors'
import { useAuthStore } from '@/stores/auth'
import { useUiStore } from '@/stores/ui'

const KUNCI_FILTER = ['q', 'role', 'active']

const { nilai, terapkan, reset } = useFilter()
const { can } = usePermission()
const auth = useAuthStore()
const ui = useUiStore()

const params = computed(() => {
  const hasil = { per_page: 20, page: nilai.value.page || undefined }
  for (const k of KUNCI_FILTER) if (nilai.value[k]) hasil[k] = nilai.value[k]
  return hasil
})
const adaFilter = computed(() => KUNCI_FILTER.some((k) => nilai.value[k]))
const { status, baris: admins, meta, pesan, muat, rentang } = useDaftar(listAdmins, params, { satuan: 'admin' })
const roles = computed(() => meta.value?.roles ?? [])

const tambahBuka = ref(false)
const roleBuka = ref(false)
const terpilih = ref(null)
// Konfirmasi: { jenis: 'nonaktif' | 'reset', admin }
const konfirmasi = ref(null)
const memproses = ref(false)

const nama = (a) => a?.full_name ?? a?.email ?? ''
const judulKonfirmasi = computed(() =>
  konfirmasi.value?.jenis === 'nonaktif' ? `Nonaktifkan ${nama(konfirmasi.value.admin)}?` : `Kirim link reset password ke ${konfirmasi.value?.admin.email}?`,
)
const pesanKonfirmasi = computed(() =>
  konfirmasi.value?.jenis === 'nonaktif'
    ? 'Admin ini tidak bisa masuk lagi dan sesi yang sedang berjalan berakhir. Log aktivitasnya tetap tersimpan, dan akunnya bisa diaktifkan lagi.'
    : 'Link sekali pakai, berlaku 1 jam. Password lama tetap berlaku sampai admin itu membuat password baru.',
)
const konfirmasiBuka = computed({ get: () => Boolean(konfirmasi.value), set: (v) => !v && (konfirmasi.value = null) })

async function tersimpanBaru(admin) {
  ui.tampilkanToast({ pesan: `${nama(admin)} ditambahkan sebagai ${admin.roles[0]?.name}.`, jenis: 'success' })
  await muat({ diam: true })
}

function bukaRole(admin) {
  terpilih.value = admin
  roleBuka.value = true
}

function roleTersimpan(baru) {
  const baris = admins.value.find((a) => a.id === baru.id)
  if (baris) Object.assign(baris, baru)
  ui.tampilkanToast({ pesan: `Role ${nama(baru)} sekarang ${baru.roles.map((r) => r.name).join(', ')}.`, jenis: 'success' })
}

// Mengaktifkan tanpa konfirmasi; menonaktifkan dengan konfirmasi (aksi yang memutus akses).
function ubahAktif(admin, nilaiBaru) {
  if (nilaiBaru) jalankanAktif(admin, true)
  else konfirmasi.value = { jenis: 'nonaktif', admin }
}

async function jalankanAktif(admin, nilaiBaru) {
  admin.sedang = true
  try {
    Object.assign(admin, await setAdminActive(admin.id, nilaiBaru))
    ui.tampilkanToast({ pesan: `Akun ${nama(admin)} ${nilaiBaru ? 'diaktifkan' : 'dinonaktifkan'}.`, jenis: 'success' })
  } catch (error) {
    const kode = error?.response?.status
    if (kode && kode < 500 && kode !== 401 && kode !== 403) ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    admin.sedang = false
  }
}

async function konfirmasiYa() {
  const { jenis, admin } = konfirmasi.value
  memproses.value = true
  try {
    if (jenis === 'nonaktif') {
      await jalankanAktif(admin, false)
    } else {
      ui.tampilkanToast({ pesan: await resetAdminPassword(admin.id), jenis: 'success' })
    }
  } catch (error) {
    const kode = error?.response?.status
    if (kode && kode < 500 && kode !== 401 && kode !== 403) ui.tampilkanToast({ pesan: pesanError(error), jenis: 'danger' })
  } finally {
    memproses.value = false
    konfirmasi.value = null
    await nextTick()
    document.querySelector(`[data-${jenis === 'nonaktif' ? 'aktif' : 'reset'}-admin="${admin.id}"]`)?.focus()
  }
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="max-w-2xl text-sm text-muted">Akun yang bisa masuk ke panel admin. Izin tiap role diatur di halaman Role dan izin.</p>
      <Button v-if="can('user.create')" :disabled="!roles.length" @click="tambahBuka = true">
        <template #icon><UserPlus class="size-4.5" aria-hidden="true" /></template>
        Tambah admin
      </Button>
    </div>

    <PenggunaFilterBar :nilai="nilai" :roles="roles" @ubah="terapkan" @reset="reset()" />

    <section aria-label="Daftar pengguna admin" :aria-busy="status === 'loading' ? 'true' : 'false'">
      <SkeletonTabel v-if="status === 'loading'" label="Memuat pengguna admin…" :baris="5" />
      <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
      <EmptyState
        v-else-if="!admins.length"
        :title="adaFilter ? 'Tidak ada admin yang cocok' : 'Belum ada admin'"
        :description="adaFilter ? 'Periksa lagi nama atau email, atau kurangi filter.' : 'Tambahkan admin pertama untuk mulai membagi tugas.'"
      >
        <template #icon><Users class="size-10" aria-hidden="true" /></template>
        <template v-if="adaFilter" #actions>
          <Button variant="outline" @click="reset()">Reset filter</Button>
        </template>
      </EmptyState>
      <div v-else class="overflow-hidden rounded-xl border border-border bg-surface">
        <PenggunaTabel :admins="admins" :id-saya="auth.user?.id ?? null" @aktif="ubahAktif" @role="bukaRole" @reset="(a) => (konfirmasi = { jenis: 'reset', admin: a })" />
        <div class="flex flex-wrap items-center justify-between gap-3 border-t border-border px-4 py-3">
          <span class="text-sm text-muted" aria-live="polite">{{ rentang }}</span>
          <PaginationBar v-if="meta?.total_pages > 1" :page="meta.page" :total-pages="meta.total_pages" @update:page="(n) => terapkan({ page: n > 1 ? n : null })" />
        </div>
      </div>
    </section>

    <TambahAdminModal v-model:open="tambahBuka" :roles="roles" @tersimpan="tersimpanBaru" />
    <UbahRoleModal v-model:open="roleBuka" :admin="terpilih" :roles="roles" @tersimpan="roleTersimpan" />
    <ConfirmModal
      v-model:open="konfirmasiBuka"
      :title="judulKonfirmasi"
      :message="pesanKonfirmasi"
      :confirm-label="konfirmasi?.jenis === 'nonaktif' ? 'Ya, nonaktifkan' : 'Kirim link'"
      :danger="konfirmasi?.jenis === 'nonaktif'"
      :loading="memproses"
      @confirm="konfirmasiYa"
    />
  </div>
</template>
