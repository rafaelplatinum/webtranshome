<script setup>
import { useId } from 'vue'

const model = defineModel({ type: String, required: true })

defineProps({
  // [{ value: 'EARN', label: 'Belanja (EARN)' }]
  opsi: { type: Array, required: true },
})

const nama = useId()
</script>

<template>
  <!-- Radio asli: Tab masuk ke grup, panah kiri/kanan memilih. Tepi tebal lewat ring supaya ukuran tidak bergeser. -->
  <fieldset class="flex flex-col">
    <legend class="mb-2 text-sm font-semibold">Jenis transaksi</legend>
    <div class="flex flex-wrap gap-2">
      <label
        v-for="o in opsi"
        :key="o.value"
        class="inline-flex min-h-11 cursor-pointer items-center gap-2 rounded-md border px-3.5 text-sm font-semibold transition-colors duration-150"
        :class="model === o.value ? 'border-primary bg-primary-soft ring-1 ring-primary ring-inset' : 'border-border-strong bg-surface hover:border-foreground'"
      >
        <input v-model="model" type="radio" :name="nama" :value="o.value" class="size-4 cursor-pointer accent-primary" />
        {{ o.label }}
      </label>
    </div>
  </fieldset>
</template>
