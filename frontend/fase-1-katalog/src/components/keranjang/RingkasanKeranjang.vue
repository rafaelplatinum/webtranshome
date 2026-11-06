<script setup>
import { computed, ref } from 'vue'
import { CircleCheck } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import WhatsAppButton from '@/components/layout/WhatsAppButton.vue'
import { formatRupiah } from '@/utils/format'
import { hitungTotal } from '@/utils/templateWhatsApp'

// Ringkasan /keranjang + tombol "Pesan via WhatsApp" (template admin "Pesan keranjang", KEPUTUSAN #23).
// Perkiraan total hanya dari harga yang tampil di website; harga final tetap dari tim Transhome.
const props = defineProps({
  // [{ produk, jumlah }] yang masih dijual (ikut dikirim).
  items: { type: Array, required: true },
  // Isi pesan WhatsApp (untuk "Lihat isi pesan").
  pesan: { type: String, default: '' },
  // Link wa.me; null selama nomor toko belum dimuat / tidak ada produk yang bisa dikirim.
  link: { type: String, default: null },
  jumlahTidakDijual: { type: Number, default: 0 },
})
const emit = defineEmits(['kosongkan'])

const total = computed(() => hitungTotal(props.items))
const terkirim = ref(false)
</script>

<template>
  <aside class="flex flex-col gap-4 rounded-xl border border-border bg-surface p-4 sm:p-5" aria-labelledby="judul-ringkasan-keranjang" data-ringkasan-keranjang>
    <h2 id="judul-ringkasan-keranjang" class="text-lg font-bold">Ringkasan pesanan</h2>

    <dl class="flex flex-col gap-2 text-sm">
      <div class="flex items-baseline justify-between gap-3">
        <dt class="text-muted">Produk dikirim</dt>
        <dd class="font-semibold tabular-nums">{{ items.length }} produk</dd>
      </div>
      <div class="flex items-baseline justify-between gap-3 border-t border-border pt-2">
        <dt class="font-semibold">Perkiraan total</dt>
        <dd class="text-right text-xl font-extrabold tabular-nums" data-perkiraan-total>
          <template v-if="total.adaHarga">{{ formatRupiah(total.rupiah) }}</template>
          <span v-else class="text-base text-muted">Via WhatsApp</span>
        </dd>
      </div>
    </dl>

    <ul class="flex flex-col gap-1.5 text-sm text-muted">
      <li v-if="total.tanpaHarga && total.adaHarga">Belum termasuk {{ total.tanpaHarga }} produk dengan harga via WhatsApp.</li>
      <li v-if="jumlahTidakDijual" class="font-semibold text-danger-ink">{{ jumlahTidakDijual }} produk sudah tidak dijual dan tidak ikut dikirim.</li>
      <li>Harga dan stok final dikonfirmasi tim Transhome lewat WhatsApp. Belum termasuk ongkos kirim.</li>
    </ul>

    <div class="flex flex-col gap-2">
      <WhatsAppButton size="lg" block wrap :href="link" label="Pesan via WhatsApp" data-pesan-keranjang @click="terkirim = true" />
      <p class="text-sm text-muted">Pesan berisi semua produk di atas. Kamu masih bisa mengubahnya di WhatsApp sebelum dikirim.</p>
    </div>

    <div v-if="terkirim" class="flex flex-col gap-2 rounded-lg bg-success-soft p-3 text-sm text-success-ink" role="status">
      <p class="flex items-start gap-2 font-semibold">
        <CircleCheck class="mt-0.5 size-4 shrink-0" aria-hidden="true" />
        WhatsApp sudah dibuka. Setelah pesanan terkirim, keranjang bisa dikosongkan.
      </p>
      <Button variant="outline" size="sm" class="self-start" @click="emit('kosongkan')">Kosongkan keranjang</Button>
    </div>

    <details v-if="pesan" class="group rounded-lg border border-border">
      <summary class="flex min-h-11 cursor-pointer items-center rounded-lg px-3 text-sm font-semibold transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
        Lihat isi pesan
      </summary>
      <pre class="max-h-80 overflow-auto border-t border-border px-3 py-2.5 font-sans text-sm leading-relaxed break-words whitespace-pre-wrap" data-isi-pesan-keranjang>{{ pesan }}</pre>
    </details>
  </aside>
</template>
