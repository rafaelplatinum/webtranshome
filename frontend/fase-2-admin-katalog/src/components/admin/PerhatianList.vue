<script setup>
import { computed, ref, useId, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { CircleCheck, CircleX, LoaderCircle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { usePermission } from '@/composables/usePermission'
import { useSyncQontak } from '@/composables/useSyncQontak'
import { formatTanggal, selisihHari } from '@/utils/format'

const props = defineProps({
  // attention dari GET /admin/dashboard: [{ type, id, title, ... }]
  items: { type: Array, required: true },
  // attention_total: jumlah lengkap per jenis (daftar hanya memuat 3 per jenis)
  total: { type: Object, default: () => ({}) },
})

const JENIS = {
  sync_gagal: { tag: 'Sync gagal', kelas: 'bg-info-soft text-info-ink', semua: '/admin/sync-qontak?status=FAILED', labelSemua: 'sync gagal' },
  stok_habis: { tag: 'Stok habis', kelas: 'bg-danger-soft text-danger-ink', semua: '/admin/produk?stock=HABIS&active=1', labelSemua: 'produk stok habis' },
  banner_berakhir: { tag: 'Banner', kelas: 'bg-warning-soft text-warning-ink', semua: '/admin/banner', labelSemua: 'banner' },
  tanpa_foto: { tag: 'Produk', kelas: 'bg-subtle text-foreground', semua: '/admin/produk?foto=0&active=1', labelSemua: 'produk tanpa foto' },
}

const { can } = usePermission()
const { cobaLagi } = useSyncQontak()
const id = useId()
const daftar = ref([])

watch(() => props.items, (items) => (daftar.value = items.map((x) => ({ ...x }))), { immediate: true })

const lainnya = computed(() =>
  Object.entries(props.total)
    .map(([jenis, n]) => ({ jenis, sisa: n - daftar.value.filter((x) => x.type === jenis).length, n }))
    .filter((x) => x.sisa > 0 && JENIS[x.jenis]),
)

function keterangan(item) {
  if (item.type === 'stok_habis') return `${item.sku} · stok habis`
  if (item.type === 'tanpa_foto') return `${item.sku} · belum ada foto utama`
  if (item.type === 'banner_berakhir') {
    const hari = selisihHari(item.end_at)
    return `Berakhir ${hari <= 0 ? 'hari ini' : `dalam ${hari} hari`} (${formatTanggal(item.end_at)})`
  }
  if (item.sync_status === 'PENDING') return 'Sedang dikirim ulang ke Qontak…'
  if (item.sync_status === 'SYNCED') return 'Berhasil tersinkron'
  return `Percobaan ke-${item.retry_count} gagal · status ${item.status_code ?? '-'}`
}

const kelasTautan =
  'inline-flex h-10 shrink-0 items-center rounded-md border border-border-strong bg-surface px-3 text-sm font-semibold transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:h-9'
</script>

<template>
  <section :aria-labelledby="id" class="flex flex-col gap-3.5 rounded-xl border border-border bg-surface p-5">
    <h2 :id="id" class="text-lg font-extrabold">Perlu perhatian</h2>

    <p v-if="!daftar.length" class="flex items-center gap-2 text-sm text-muted">
      <CircleCheck class="size-5 text-success" aria-hidden="true" />
      Semua beres. Tidak ada yang perlu ditindaklanjuti.
    </p>

    <ul v-else class="flex flex-col gap-2">
      <li v-for="item in daftar" :key="`${item.type}-${item.id}`" class="flex flex-wrap items-center gap-3 rounded-lg border border-subtle p-3">
        <span class="rounded px-2 py-0.5 text-xs font-bold whitespace-nowrap" :class="JENIS[item.type]?.kelas">{{ JENIS[item.type]?.tag }}</span>
        <span class="flex min-w-0 flex-[1_1_10rem] flex-col">
          <strong class="font-semibold">{{ item.title }}</strong>
          <span class="flex items-center gap-1 text-xs text-muted" aria-live="polite">
            <LoaderCircle v-if="item.sync_status === 'PENDING'" class="size-3.5 animate-spin" aria-hidden="true" />
            <CircleCheck v-else-if="item.sync_status === 'SYNCED'" class="size-3.5 text-success" aria-hidden="true" />
            <CircleX v-else-if="item.sync_status === 'FAILED'" class="size-3.5 text-danger" aria-hidden="true" />
            {{ keterangan(item) }}
          </span>
        </span>
        <template v-if="item.type === 'sync_gagal'">
          <Button
            v-if="can('qontak.retry') && item.sync_status !== 'SYNCED'"
            variant="outline"
            size="sm"
            :loading="item.sedang"
            @click="cobaLagi(item, item.title)"
          >
            Coba lagi
          </Button>
        </template>
        <RouterLink v-else-if="item.type === 'stok_habis' && can('product.update')" :to="`/admin/produk/${item.id}#stok`" :class="kelasTautan">Ubah stok</RouterLink>
        <RouterLink v-else-if="item.type === 'tanpa_foto' && can('product.update')" :to="`/admin/produk/${item.id}#foto`" :class="kelasTautan">Lengkapi</RouterLink>
        <RouterLink v-else-if="item.type === 'banner_berakhir' && can('banner.update')" :to="`/admin/banner/${item.id}#periode`" :class="kelasTautan">Perpanjang</RouterLink>
      </li>
    </ul>

    <ul v-if="lainnya.length" class="flex flex-col gap-1 text-sm">
      <li v-for="x in lainnya" :key="x.jenis">
        <RouterLink :to="JENIS[x.jenis].semua" class="rounded-sm font-semibold text-primary hover:text-primary-hover hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
          Lihat semua {{ JENIS[x.jenis].labelSemua }} ({{ x.n }})
        </RouterLink>
      </li>
    </ul>
  </section>
</template>
