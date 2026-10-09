<script setup>
import { computed } from 'vue'
import { Minus, Plus } from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import { bandingkanData } from '@/utils/log'

const props = defineProps({
  sebelum: { type: Object, default: null },
  sesudah: { type: Object, default: null },
})

const baris = computed(() => bandingkanData(props.sebelum, props.sesudah))
// Data baru (create) atau hanya sebelum: satu kolom nilai saja.
const duaKolom = computed(() => Boolean(props.sebelum && props.sesudah))
const judulNilai = computed(() => (props.sesudah ? 'Data' : 'Data sebelum'))
const jumlahBerubah = computed(() => baris.value.filter((b) => b.berubah).length)
</script>

<template>
  <div class="flex flex-col gap-2">
    <p v-if="duaKolom" class="text-sm text-muted">{{ jumlahBerubah }} kolom berubah. Kolom tanpa tanda hanya penanda data.</p>
    <div class="relative overflow-x-auto rounded-lg border border-border">
      <table class="w-full border-collapse text-sm">
        <caption class="sr-only">Perbandingan data sebelum dan sesudah</caption>
        <thead class="bg-background text-left text-xs font-semibold tracking-wide text-muted uppercase">
          <tr>
            <th scope="col" class="px-3 py-2.5">Kolom</th>
            <template v-if="duaKolom">
              <th scope="col" class="px-3 py-2.5">Sebelum</th>
              <th scope="col" class="px-3 py-2.5">Sesudah</th>
            </template>
            <th v-else scope="col" class="px-3 py-2.5">{{ judulNilai }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in baris" :key="b.kunci" class="border-t border-border align-top" :class="b.berubah && 'bg-warning-soft/60'" :data-kolom="b.kunci">
            <th scope="row" class="px-3 py-2.5 text-left font-semibold">
              <span class="flex flex-col items-start gap-1">
                {{ b.label }}
                <Badge v-if="b.berubah" variant="warning">Berubah</Badge>
              </span>
            </th>
            <template v-if="duaKolom">
              <td class="px-3 py-2.5 break-words" data-sebelum><span class="whitespace-pre-wrap">{{ b.sebelum }}</span></td>
              <td class="px-3 py-2.5 break-words" data-sesudah>
                <span class="whitespace-pre-wrap">{{ b.sesudah }}</span>
                <span v-if="b.ditambah?.length || b.dicabut?.length" class="mt-2 flex flex-col gap-1">
                  <span v-for="x in b.ditambah" :key="`+${x}`" class="inline-flex items-center gap-1 text-success-ink">
                    <Plus class="size-3.5" aria-hidden="true" />Ditambah: <code class="font-mono text-xs">{{ x }}</code>
                  </span>
                  <span v-for="x in b.dicabut" :key="`-${x}`" class="inline-flex items-center gap-1 text-danger-ink">
                    <Minus class="size-3.5" aria-hidden="true" />Dicabut: <code class="font-mono text-xs">{{ x }}</code>
                  </span>
                </span>
              </td>
            </template>
            <td v-else class="px-3 py-2.5 break-words"><span class="whitespace-pre-wrap">{{ sesudah ? b.sesudah : b.sebelum }}</span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
