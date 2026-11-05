<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { Check, MessageCircle, X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { useAuthStore } from '@/stores/auth'
import { useWhatsAppContactStore } from '@/stores/whatsappContact'
import { CONTACT } from '@/config/contact'

const store = useWhatsAppContactStore()
const auth = useAuthStore()
const dialog = ref(null)
const nama = ref('')
const nomor = ref('')
const topik = ref('Lainnya')
const produk = ref('')
const pesan = ref('')
const galat = ref('')
const selesai = ref(false)
const linkCadangan = ref('')
const memproses = ref(false)

const OPSI_TOPIK = ['Tanya produk', 'Minta penawaran/harga', 'Proyek/renovasi', 'Jasa', 'Status pesanan', 'Lainnya']
const phonePattern = /^(?:08\d{8,11}|\+62\d{8,12})$/
const pesanWhatsapp = computed(() => [
  'Halo Trans Home, saya ingin bertanya.',
  '',
  `Nama: ${nama.value.trim()}`,
  `No. WhatsApp: ${nomor.value.trim()}`,
  `Keperluan: ${topik.value}`,
  `Produk/SKU/Jasa: ${produk.value.trim() || '-'}`,
  `Pesan: ${pesan.value.trim()}`,
  '',
  'Dikirim dari website Trans Home',
].join('\n'))

function bersihkan() {
  nama.value = auth.user?.full_name || ''
  const nomorTersimpan = auth.user?.phone_number || ''
  nomor.value = nomorTersimpan.startsWith('62') ? `+${nomorTersimpan}` : nomorTersimpan
  topik.value = store.konteks.topic
  produk.value = store.konteks.product
  pesan.value = store.konteks.message
  galat.value = ''
  selesai.value = false
  linkCadangan.value = ''
}

watch(() => store.terbuka, async (aktif) => {
  if (aktif) {
    bersihkan()
    await nextTick()
    if (!dialog.value?.open) dialog.value?.showModal()
    await nextTick()
    dialog.value?.querySelector('input')?.focus()
  } else if (dialog.value?.open) dialog.value.close()
})

function tutup() {
  store.tutup()
}

function kirim() {
  galat.value = ''
  if (!nama.value.trim()) { galat.value = 'Nama lengkap wajib diisi.'; return }
  const angka = nomor.value.replace(/[\s()-]/g, '')
  if (!phonePattern.test(angka)) { galat.value = 'Masukkan nomor 08xx atau +62xx yang valid.'; return }
  if (!pesan.value.trim()) { galat.value = 'Pesan atau pertanyaan wajib diisi.'; return }

  memproses.value = true
  nomor.value = angka.startsWith('+62') ? `0${angka.slice(3)}` : angka
  const url = `https://wa.me/${CONTACT.whatsappNumber}?text=${encodeURIComponent(pesanWhatsapp.value)}`
  const popup = window.open(url, '_blank')
  selesai.value = true
  linkCadangan.value = url
  memproses.value = false
  if (popup) popup.opener = null
}

function saatTutupNative() { store.tutup() }
onBeforeUnmount(() => { if (dialog.value?.open) dialog.value.close() })
</script>

<template>
  <dialog
    ref="dialog"
    aria-labelledby="judul-form-whatsapp"
    aria-describedby="deskripsi-form-whatsapp"
    class="m-auto max-h-[min(92dvh,52rem)] w-[calc(100%-2rem)] max-w-xl overflow-y-auto rounded-2xl border border-border bg-surface p-0 text-foreground backdrop:bg-foreground/50"
    @close="saatTutupNative"
    @cancel.prevent="tutup"
  >
    <section class="p-5 sm:p-7">
      <div class="mb-5 flex items-start justify-between gap-4">
        <div>
          <span class="mb-3 inline-grid size-10 place-items-center rounded-full bg-success-soft text-success-ink"><MessageCircle class="size-5" aria-hidden="true" /></span>
          <h2 id="judul-form-whatsapp" class="text-xl font-extrabold">Hubungi Transhome</h2>
          <p id="deskripsi-form-whatsapp" class="mt-1 text-sm text-muted">Isi data berikut untuk melanjutkan percakapan di WhatsApp CS.</p>
        </div>
        <button type="button" aria-label="Tutup form" class="grid size-10 shrink-0 place-items-center rounded-full text-muted transition-colors hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" @click="tutup"><X class="size-5" aria-hidden="true" /></button>
      </div>

      <div v-if="selesai" role="status" class="rounded-xl border border-success/20 bg-success-soft p-5">
        <span class="mb-3 inline-grid size-10 place-items-center rounded-full bg-success text-success-foreground"><Check class="size-5" aria-hidden="true" /></span>
        <h3 class="font-bold text-success-ink">Terima kasih, kamu akan diarahkan ke WhatsApp.</h3>
        <p class="mt-1 text-sm text-success-ink">Jika WhatsApp tidak terbuka otomatis, gunakan tombol berikut.</p>
        <Button class="mt-4" :href="linkCadangan">Buka WhatsApp</Button>
      </div>

      <form v-else class="flex flex-col gap-4" novalidate @submit.prevent="kirim">
        <p v-if="galat" role="alert" class="rounded-md bg-danger-soft px-3 py-2 text-sm font-semibold text-danger-ink">{{ galat }}</p>
        <label class="flex flex-col gap-1.5 text-sm font-semibold">
          Nama lengkap <span class="sr-only">wajib</span>
          <input v-model="nama" required autocomplete="name" class="h-11 rounded-md border border-border-strong bg-surface px-3 text-base font-normal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" />
        </label>
        <label class="flex flex-col gap-1.5 text-sm font-semibold">
          Nomor WhatsApp pengirim <span class="sr-only">wajib</span>
          <input v-model="nomor" required autocomplete="tel" inputmode="tel" placeholder="08xx atau +62xx" class="h-11 rounded-md border border-border-strong bg-surface px-3 text-base font-normal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" />
        </label>
        <label class="flex flex-col gap-1.5 text-sm font-semibold">
          Topik / keperluan
          <select v-model="topik" class="h-11 rounded-md border border-border-strong bg-surface px-3 text-base font-normal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
            <option v-for="opsi in OPSI_TOPIK" :key="opsi">{{ opsi }}</option>
          </select>
        </label>
        <label class="flex flex-col gap-1.5 text-sm font-semibold">
          Nama produk / SKU / jasa <span class="font-normal text-muted">(opsional)</span>
          <input v-model="produk" class="h-11 rounded-md border border-border-strong bg-surface px-3 text-base font-normal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring" />
        </label>
        <label class="flex flex-col gap-1.5 text-sm font-semibold">
          Pesan / pertanyaan <span class="sr-only">wajib</span>
          <textarea v-model="pesan" required rows="4" class="rounded-md border border-border-strong bg-surface px-3 py-2 text-base font-normal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"></textarea>
        </label>
        <div class="flex flex-wrap justify-end gap-2 pt-1">
          <Button type="button" variant="outline" :disabled="memproses" @click="tutup">Batal</Button>
          <Button type="submit" :loading="memproses" :disabled="memproses">Lanjut ke WhatsApp</Button>
        </div>
      </form>
    </section>
  </dialog>
</template>
