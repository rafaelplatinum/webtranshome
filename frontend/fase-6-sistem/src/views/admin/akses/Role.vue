<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Info, Lock } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import DaftarRole from '@/components/admin/akses/DaftarRole.vue'
import MatriksIzin from '@/components/admin/akses/MatriksIzin.vue'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { useMatriksIzin } from '@/composables/useMatriksIzin'
import { useMuat } from '@/composables/useMuat'
import { usePermission } from '@/composables/usePermission'
import { listPermissions, listRoles, saveRolePermissions } from '@/services/adminAksesService'
import { errorField, pesanError } from '@/services/errors'
import { useUiStore } from '@/stores/ui'
import { perilakuGulir } from '@/utils/gerak'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()
const { can } = usePermission()

const roles = ref([])
const katalog = ref([])
const { status, pesan, muat } = useMuat(async () => {
  const [daftarRole, daftarIzin] = await Promise.all([listRoles(), listPermissions()])
  roles.value = daftarRole
  katalog.value = daftarIzin
})

const m = useMatriksIzin(katalog)
// Role terpilih disimpan di URL (?role=ADMIN_KATALOG): tetap sama setelah refresh.
const role = computed(() => roles.value.find((r) => r.code === route.query.role) ?? roles.value[0] ?? null)
const terkunci = computed(() => Boolean(role.value?.is_locked) || !can('role.update'))
watch(role, (r) => r && m.setel(r.permissions), { immediate: true })

const menyimpan = ref(false)
const galat = ref('')
const pindahKe = ref(null)
const keluar = useKonfirmasiKeluar(m.kotor)

const ringkasan = computed(() => {
  const bagian = []
  if (m.ditambah.value.length) bagian.push(`${m.ditambah.value.length} izin ditambah`)
  if (m.dicabut.value.length) bagian.push(`${m.dicabut.value.length} izin dicabut`)
  return bagian.join(', ')
})

function pilih(kode) {
  if (kode === role.value?.code) return
  if (m.kotor.value) pindahKe.value = kode
  else gantiRole(kode)
}

const panelIzin = ref(null)

async function gantiRole(kode) {
  pindahKe.value = null
  galat.value = ''
  await router.replace({ query: { ...route.query, role: kode } })
  // Di bawah 1024 px matriks ada di bawah daftar role: gulir ke sana.
  if (!window.matchMedia('(min-width: 1024px)').matches) {
    await nextTick()
    panelIzin.value?.scrollIntoView({ behavior: perilakuGulir(), block: 'start' })
  }
}

async function simpan() {
  if (menyimpan.value || !role.value) return
  menyimpan.value = true
  galat.value = ''
  try {
    const baru = await saveRolePermissions(role.value.id, m.hasilAkhir())
    Object.assign(role.value, baru)
    m.setel(baru.permissions)
    ui.tampilkanToast({ pesan: `Izin ${baru.name} disimpan. Berlaku setelah admin dengan role ini memuat ulang halaman.`, jenis: 'success' })
  } catch (error) {
    const kode = error?.response?.status
    if (kode === 422) galat.value = errorField(error).permissions ?? pesanError(error)
    else if (kode && kode < 500 && kode !== 401 && kode !== 403) galat.value = pesanError(error)
  } finally {
    menyimpan.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <p class="max-w-3xl text-sm text-muted">
      Atur menu dan aksi yang boleh dipakai tiap role. Tombol tanpa izin disembunyikan, dan server tetap menolak aksinya. Perubahan berlaku setelah admin terkait memuat ulang halaman.
    </p>

    <div v-if="status === 'loading'" class="grid gap-4 lg:grid-cols-[18rem_minmax(0,1fr)]" aria-busy="true">
      <span class="sr-only">Memuat role dan izin…</span>
      <div class="flex flex-col gap-2"><Skeleton v-for="n in 4" :key="n" class="h-24 w-full rounded-xl" /></div>
      <Skeleton class="h-96 w-full rounded-xl" />
    </div>

    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

    <div v-else-if="role" class="grid items-start gap-4 lg:grid-cols-[18rem_minmax(0,1fr)]">
      <DaftarRole :roles="roles" :aktif="role.code" @pilih="pilih" />

      <section ref="panelIzin" :aria-label="`Izin ${role.name}`" class="flex min-w-0 scroll-mt-24 flex-col gap-4 rounded-xl border border-border bg-surface p-4 md:p-5">
        <div class="flex flex-col gap-1">
          <h2 class="text-lg font-extrabold">Izin {{ role.name }}</h2>
          <p class="text-sm text-muted">{{ role.description }}</p>
        </div>
        <p v-if="role.is_locked" class="flex items-start gap-2 rounded-md bg-subtle px-3 py-2 text-sm">
          <Lock class="mt-0.5 size-4 shrink-0" aria-hidden="true" />
          Super Admin selalu punya semua izin supaya selalu ada yang bisa mengatur akses. Izinnya tidak bisa diubah.
        </p>
        <p v-else-if="!can('role.update')" class="flex items-start gap-2 rounded-md bg-subtle px-3 py-2 text-sm">
          <Info class="mt-0.5 size-4 shrink-0" aria-hidden="true" />
          Kamu hanya bisa melihat izin. Mengubahnya butuh izin "Mengubah izin role".
        </p>
        <p v-else class="text-sm text-muted">Mencentang Tambah, Ubah, atau Lainnya ikut mencentang Lihat. Mencabut Lihat mencabut seluruh baris.</p>

        <MatriksIzin :katalog="katalog" :matriks="m" :nama-role="role.name" :terkunci="terkunci" />

        <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2 text-sm font-medium text-danger-ink">{{ galat }}</p>

        <div
          v-if="!terkunci"
          class="sticky bottom-0 -mx-4 -mb-4 flex flex-wrap items-center justify-between gap-3 border-t border-border bg-surface px-4 py-3 md:-mx-5 md:-mb-5 md:px-5"
        >
          <span class="text-sm" aria-live="polite">{{ m.kotor.value ? ringkasan : 'Belum ada perubahan.' }}</span>
          <span class="flex gap-2">
            <Button variant="outline" :disabled="!m.kotor.value || menyimpan" @click="m.setel(role.permissions)">Batal</Button>
            <Button :disabled="!m.kotor.value" :loading="menyimpan" @click="simpan">Simpan izin</Button>
          </span>
        </div>
      </section>
    </div>

    <ConfirmModal
      :open="Boolean(pindahKe)"
      :title="`Buang perubahan izin ${role?.name ?? ''}?`"
      message="Perubahan yang belum disimpan akan hilang."
      confirm-label="Buang perubahan"
      cancel-label="Tetap di sini"
      danger
      @confirm="gantiRole(pindahKe)"
      @cancel="pindahKe = null"
    />
    <ConfirmModal
      :open="keluar.terbuka.value"
      title="Buang perubahan?"
      message="Perubahan izin yang belum disimpan akan hilang."
      confirm-label="Buang perubahan"
      cancel-label="Tetap di sini"
      danger
      @confirm="keluar.buang"
      @cancel="keluar.tetap"
    />
  </div>
</template>
