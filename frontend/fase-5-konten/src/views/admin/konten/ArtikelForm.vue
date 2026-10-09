<script setup>
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Eye, Send, Undo2 } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import Textarea from '@/components/ui/Textarea.vue'
import FormSection from '@/components/admin/FormSection.vue'
import GambarField from '@/components/admin/GambarField.vue'
import StatusArtikel from '@/components/admin/konten/StatusArtikel.vue'
import { useFormArtikel } from '@/composables/useFormArtikel'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { usePermission } from '@/composables/usePermission'
import { tujuanKembali } from '@/composables/useTujuanKembali'
import { useUiStore } from '@/stores/ui'
import { TIPE_ARTIKEL } from '@/utils/artikel'
import { formatTanggalJam } from '@/utils/format'

const router = useRouter()
const ui = useUiStore()
const { can } = usePermission()
const f = useFormArtikel()
const { form, errors } = f
const keluar = useKonfirmasiKeluar(f.kotor, { saatSesiHabis: f.simpanDraf })
const formEl = ref(null)
const konfirmasi = ref('') // '' | terbit | batal-terbit

const OPSI_TIPE = Object.entries(TIPE_ARTIKEL).map(([value, t]) => ({ value, label: t.label }))
const judulPendek = computed(() => form.title.trim() || 'Artikel ini')
const pesanKonfirmasi = computed(() =>
  konfirmasi.value === 'terbit'
    ? `“${judulPendek.value}” akan tampil di website mulai sekarang.${f.kotor.value ? ' Perubahan yang belum disimpan ikut disimpan.' : ''}`
    : `“${judulPendek.value}” tidak tampil lagi di website. Isinya tetap tersimpan sebagai draf.`,
)

async function fokusGalatPertama() {
  await nextTick()
  const el = formEl.value?.querySelector('[aria-invalid="true"], [data-galat]')
  el?.scrollIntoView({ block: 'center' })
  if (el?.matches('input, select, textarea')) el.focus({ preventScroll: true })
}

// Artikel baru yang sudah tersimpan pindah ke alamat editornya (/admin/artikel/:id), termasuk bila
// terbit gagal setelah draf tersimpan; pesan gagal dibawa lewat toast karena halaman dimuat ulang.
async function selesai(berhasil, pesanBerhasil) {
  if (berhasil) ui.tampilkanToast({ pesan: pesanBerhasil, jenis: 'success' })
  if (!f.id && f.tersimpan.value) {
    if (!berhasil && f.galatUmum.value) ui.tampilkanToast({ pesan: f.galatUmum.value, jenis: 'danger' })
    return router.replace(`/admin/artikel/${f.tersimpan.value.id}`)
  }
  if (berhasil) return
  if (Object.keys(f.errors.value).length) return fokusGalatPertama()
  formEl.value?.scrollIntoView({ block: 'start' })
}

async function simpan() {
  if (!f.validasi()) return fokusGalatPertama()
  selesai(await f.simpan(), f.terbit.value ? 'Perubahan disimpan dan langsung tampil di website' : 'Draf disimpan')
}

function mintaTerbit() {
  if (!f.validasi()) return fokusGalatPertama()
  konfirmasi.value = 'terbit'
}

async function konfirmasiAksi() {
  const terbit = konfirmasi.value === 'terbit'
  const berhasil = terbit ? await f.terbitkan() : await f.batalkanTerbit()
  konfirmasi.value = ''
  await selesai(berhasil, terbit ? 'Artikel terbit dan tampil di website' : 'Artikel kembali jadi draf')
  // Tombol Terbitkan berganti Batalkan terbit (atau sebaliknya): fokus pindah ke tombol penggantinya.
  if (berhasil && f.id) {
    await nextTick()
    formEl.value?.querySelector('[data-tombol-terbit]')?.focus()
  }
}

const kembali = () => router.push(tujuanKembali('artikel', '/admin/artikel'))

onMounted(f.muat)
</script>

<template>
  <div v-if="f.status.value === 'loading'" class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]" aria-busy="true">
    <span class="sr-only">Memuat artikel…</span>
    <Skeleton class="h-120 rounded-xl" />
    <Skeleton class="h-64 rounded-xl" />
  </div>

  <section v-else-if="f.status.value === 'tidak-ada'" class="mx-auto flex max-w-md flex-col items-center gap-4 py-16 text-center">
    <h2 class="text-2xl font-bold">Artikel tidak ditemukan</h2>
    <p class="text-muted">Artikel ini tidak ada atau alamatnya salah.</p>
    <Button to="/admin/artikel">Kembali ke daftar artikel</Button>
  </section>

  <ErrorState v-else-if="f.status.value === 'error'" :message="f.pesan.value" @retry="f.muat" />

  <form v-else ref="formEl" novalidate class="flex flex-col gap-6" @submit.prevent="simpan">
    <p v-if="f.dipulihkan.value" role="status" class="rounded-md bg-info-soft px-4 py-3 text-sm font-medium text-info-ink">
      Isian sebelum sesi habis sudah dipulihkan. Periksa lagi, lalu simpan.
    </p>
    <p v-if="f.galatUmum.value" role="alert" class="rounded-md bg-danger-soft px-4 py-3 text-sm font-semibold text-danger-ink">{{ f.galatUmum.value }}</p>

    <div class="grid items-start gap-6 xl:grid-cols-[minmax(0,1fr)_20rem]">
      <div class="flex min-w-0 flex-col gap-6">
        <FormSection judul="Isi">
          <Input :model-value="form.title" label="Judul" required maxlength="150" autocomplete="off" :error="errors.title ?? ''" @update:model-value="f.ubahJudul" />
          <Input
            :model-value="form.slug"
            label="Slug (alamat halaman)"
            required
            autocomplete="off"
            prefix="/artikel/"
            :error="errors.slug ?? ''"
            hint="Dibuat otomatis dari judul. Boleh diubah: huruf kecil, angka, tanda hubung."
            @update:model-value="f.ubahSlug"
          />
          <Textarea
            v-model="form.content"
            label="Isi"
            required
            :rows="16"
            :error="errors.content ?? ''"
            hint="Pisahkan paragraf dengan satu baris kosong. Baris diawali ## menjadi subjudul, baris diawali - menjadi daftar."
          />
        </FormSection>
        <FormSection judul="Gambar sampul" deskripsi="Opsional. Tampil di kartu artikel dan di atas isi.">
          <GambarField v-model="form.thumbnail_url" label="Gambar" rasio="aspect-video" lebar="w-48" hint="JPG, PNG, atau WebP, maksimal 2 MB. Ukuran disarankan 1200 × 675 px." />
        </FormSection>
      </div>

      <div class="flex flex-col gap-6 xl:sticky xl:top-24">
        <FormSection judul="Status">
          <div class="flex flex-wrap items-center gap-2" data-status-editor>
            <StatusArtikel :status="f.tersimpan.value?.status ?? 'DRAFT'" />
            <span class="text-sm text-muted">
              {{ f.tersimpan.value?.published_at ? `Terbit ${formatTanggalJam(f.tersimpan.value.published_at)}` : 'Belum terbit' }}
            </span>
          </div>
          <template v-if="f.tautanPratinjau.value">
            <Button variant="outline" :href="f.tautanPratinjau.value" data-tombol-pratinjau>
              <template #icon><Eye class="size-4" aria-hidden="true" /></template>
              Pratinjau <span class="sr-only">(tab baru)</span>
            </Button>
            <p v-if="f.kotor.value" class="text-xs text-muted">Pratinjau menampilkan versi terakhir yang disimpan.</p>
          </template>
          <p v-else class="text-sm text-muted">Pratinjau tersedia setelah draf disimpan.</p>
          <template v-if="can('article.publish')">
            <Button v-if="!f.terbit.value" :disabled="Boolean(f.menyimpan.value)" data-tombol-terbit @click="mintaTerbit">
              <template #icon><Send class="size-4" aria-hidden="true" /></template>
              Terbitkan
            </Button>
            <Button v-else variant="outline" :disabled="Boolean(f.menyimpan.value)" data-tombol-terbit @click="konfirmasi = 'batal-terbit'">
              <template #icon><Undo2 class="size-4" aria-hidden="true" /></template>
              Batalkan terbit
            </Button>
          </template>
          <p v-else class="text-sm text-muted">Penerbitan dilakukan admin yang punya izin terbit.</p>
        </FormSection>
        <FormSection judul="Tipe" deskripsi="Promo dan event tampil di halaman Promo, artikel di halaman Artikel.">
          <Select v-model="form.type" label="Tipe" hide-label :options="OPSI_TIPE" :error="errors.type ?? ''" />
        </FormSection>
      </div>
    </div>

    <div class="sticky bottom-0 z-20 -mx-4 flex flex-wrap items-center justify-end gap-2 border-t border-border bg-surface px-4 py-3 sm:-mx-7 sm:px-7">
      <span v-if="f.kotor.value" class="mr-auto hidden text-sm text-muted sm:inline">Ada perubahan yang belum disimpan</span>
      <Button variant="outline" :disabled="Boolean(f.menyimpan.value)" @click="kembali">Kembali</Button>
      <Button type="submit" :loading="f.menyimpan.value === 'simpan'" :disabled="Boolean(f.menyimpan.value)">
        {{ f.terbit.value ? 'Simpan perubahan' : 'Simpan draf' }}
      </Button>
    </div>

    <ConfirmModal
      :open="Boolean(konfirmasi)"
      :title="konfirmasi === 'terbit' ? 'Terbitkan artikel?' : 'Batalkan terbit?'"
      :message="pesanKonfirmasi"
      :confirm-label="konfirmasi === 'terbit' ? 'Terbitkan' : 'Batalkan terbit'"
      cancel-label="Batal"
      :loading="f.menyimpan.value === 'terbit' || f.menyimpan.value === 'batal-terbit'"
      @update:open="(buka) => !buka && (konfirmasi = '')"
      @confirm="konfirmasiAksi"
    />
    <ConfirmModal
      :open="keluar.terbuka.value"
      title="Buang perubahan?"
      message="Perubahan yang belum disimpan akan hilang."
      confirm-label="Buang perubahan"
      cancel-label="Tetap di sini"
      danger
      @confirm="keluar.buang"
      @cancel="keluar.tetap"
    />
  </form>
</template>
