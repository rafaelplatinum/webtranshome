<script setup>
import { computed } from 'vue'
import { KOLOM_IZIN, kolomDari, labelKhusus } from '@/composables/useMatriksIzin'

const props = defineProps({
  // GET /admin/permissions: [{ menu_id, menu, group, menu_permission, permissions: [{ code, action, description }] }]
  katalog: { type: Array, required: true },
  // Hasil useMatriksIzin(): pilihan, toggle, toggleBaris, toggleKolom, izinKolom, keadaan
  matriks: { type: Object, required: true },
  namaRole: { type: String, required: true },
  // true: hanya dibaca (role terkunci atau tidak punya izin role.update)
  terkunci: { type: Boolean, default: false },
})

const grup = computed(() => {
  const hasil = []
  for (const m of props.katalog) {
    const g = hasil.find((x) => x.nama === m.group) ?? hasil[hasil.push({ nama: m.group, menu: [] }) - 1]
    g.menu.push(m)
  }
  return hasil
})

const izinDi = (menu, kunci) => menu.permissions.filter((p) => kolomDari(p) === kunci)
const punya = (kode) => props.matriks.pilihan.value.has(kode)
const kelasKotak = 'size-5 cursor-pointer accent-primary disabled:cursor-not-allowed'
const kelasLabel = 'inline-flex min-h-11 min-w-11 cursor-pointer items-center rounded-md focus-within:ring-2 focus-within:ring-ring'
// Nama aksesibel memuat teks yang terlihat (judul kolom / label) lalu keterangan izinnya.
const namaIzin = (menu, kunci, izin) => `${menu.menu}: ${kunci === 'lainnya' ? labelKhusus(izin) : KOLOM_IZIN.find((k) => k.kunci === kunci).label} — ${izin.description}`
</script>

<template>
  <!-- Matriks lebar bergulir di wadahnya sendiri di HP, bukan halaman. -->
  <div class="relative overflow-x-auto rounded-lg border border-border">
    <table class="w-full min-w-176 border-collapse text-sm">
      <caption class="sr-only">Izin role {{ namaRole }} per menu dan aksi</caption>
      <thead class="bg-background text-xs font-semibold tracking-wide text-muted uppercase">
        <tr>
          <th scope="col" class="px-3 py-2 text-left">Menu</th>
          <th v-for="k in KOLOM_IZIN" :key="k.kunci" scope="col" class="px-2 py-2" :class="k.kunci === 'lainnya' ? 'text-left' : 'text-center'">
            <label :class="[kelasLabel, 'flex-col justify-center gap-1 px-1']">
              <span>{{ k.label }}</span>
              <input
                type="checkbox"
                :class="kelasKotak"
                :checked="matriks.keadaan(matriks.izinKolom(k.kunci)) === 'semua'"
                :indeterminate="matriks.keadaan(matriks.izinKolom(k.kunci)) === 'sebagian'"
                :disabled="terkunci"
                :aria-label="`Pilih semua kolom ${k.label}`"
                @change="matriks.toggleKolom(k.kunci, $event.target.checked)"
              />
            </label>
          </th>
          <th scope="col" class="px-2 py-2 text-center">Semua</th>
        </tr>
      </thead>
      <tbody v-for="g in grup" :key="g.nama">
        <tr>
          <th colspan="6" scope="colgroup" class="border-t border-border bg-subtle px-3 py-1.5 text-left text-xs font-bold tracking-wide text-muted uppercase">{{ g.nama }}</th>
        </tr>
        <tr v-for="menu in g.menu" :key="menu.menu_id" class="border-t border-border" :data-menu="menu.menu">
          <th scope="row" class="px-3 py-1 text-left font-semibold">{{ menu.menu }}</th>
          <td v-for="k in KOLOM_IZIN" :key="k.kunci" class="px-2 py-1" :class="k.kunci === 'lainnya' ? 'text-left' : 'text-center'">
            <template v-if="izinDi(menu, k.kunci).length">
              <label
                v-for="izin in izinDi(menu, k.kunci)"
                :key="izin.code"
                :class="[kelasLabel, k.kunci === 'lainnya' ? 'gap-2 pr-2' : 'justify-center']"
                :title="izin.description"
              >
                <input
                  type="checkbox"
                  :class="kelasKotak"
                  :checked="punya(izin.code)"
                  :disabled="terkunci"
                  :aria-label="namaIzin(menu, k.kunci, izin)"
                  :data-izin="izin.code"
                  @change="matriks.toggle(izin.code, $event.target.checked)"
                />
                <span v-if="k.kunci === 'lainnya'" aria-hidden="true">{{ labelKhusus(izin) }}</span>
              </label>
            </template>
            <span v-else class="text-faint"><span aria-hidden="true">—</span><span class="sr-only">Tidak ada</span></span>
          </td>
          <td class="px-2 py-1 text-center">
            <label :class="[kelasLabel, 'justify-center']">
              <input
                type="checkbox"
                :class="kelasKotak"
                :checked="matriks.keadaan(menu.permissions) === 'semua'"
                :indeterminate="matriks.keadaan(menu.permissions) === 'sebagian'"
                :disabled="terkunci"
                :aria-label="`Semua izin ${menu.menu}`"
                @change="matriks.toggleBaris(menu, $event.target.checked)"
              />
            </label>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
