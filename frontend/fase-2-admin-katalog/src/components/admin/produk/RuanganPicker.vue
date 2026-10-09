<script setup>
// Produk bisa cocok untuk beberapa ruangan (tabel product_rooms).
const pilihan = defineModel({ type: Array, default: () => [] })

defineProps({
  // GET /admin/rooms: [{ id, name, is_active }]
  ruangan: { type: Array, required: true },
})

function ubah(id, dicentang) {
  pilihan.value = dicentang ? [...new Set([...pilihan.value, id])] : pilihan.value.filter((x) => x !== id)
}
</script>

<template>
  <fieldset class="flex flex-col">
    <legend class="sr-only">Ruangan yang cocok</legend>
    <p v-if="!ruangan.length" class="text-sm text-muted">Belum ada data ruangan.</p>
    <label
      v-for="r in ruangan"
      :key="r.id"
      class="-mx-2 flex min-h-11 cursor-pointer items-center gap-3 rounded-md px-2 transition-colors duration-150 hover:bg-subtle"
    >
      <input type="checkbox" class="size-5 shrink-0 cursor-pointer accent-primary" :checked="pilihan.includes(r.id)" @change="ubah(r.id, $event.target.checked)" />
      <span class="text-sm">{{ r.name }}<span v-if="!r.is_active" class="text-muted"> (nonaktif)</span></span>
    </label>
  </fieldset>
</template>
