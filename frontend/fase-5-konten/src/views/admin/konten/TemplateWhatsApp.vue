<script setup>
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import FormSection from '@/components/admin/FormSection.vue'
import EditorTemplateWa from '@/components/admin/konten/EditorTemplateWa.vue'
import PratinjauWa from '@/components/admin/konten/PratinjauWa.vue'
import { useFormTemplateWa } from '@/composables/useFormTemplateWa'
import { useKonfirmasiKeluar } from '@/composables/useKonfirmasiKeluar'
import { useSettingsStore } from '@/stores/settings'
import { pesanKeranjang, pesanProduk } from '@/utils/templateWhatsApp'
import { buatLinkWhatsApp } from '@/utils/whatsapp'

// Admin › Konten › Template WhatsApp (KEPUTUSAN #23): pesan per produk dan pesan keranjang dibedakan.
// Pratinjau memakai produk contoh, jadi admin bisa melihat hasilnya untuk setiap keadaan sebelum menyimpan.
const formEl = ref(null)
const settings = useSettingsStore()
const { status, pesan, muat, form, kotor, galat, galatUmum, menyimpan, simpan, batal } = useFormTemplateWa(formEl)
const keluar = useKonfirmasiKeluar(kotor)

const contoh = (data) => ({ unit_sale: 'sak', min_order: 1, show_price: true, stock_status: 'TERSEDIA', brand: { name: 'Merek A' }, ...data })
const PRODUK_CONTOH = {
  biasa: contoh({ name: 'Semen PCC 50 kg', sku: 'SMN-PCC-50', slug: 'semen-pcc-50-kg', price_general: 68000 }),
  harga: contoh({ name: 'Keramik lantai abu doff 40×40', sku: 'KRM-4040-01', slug: 'keramik-lantai-abu-doff-40x40', price_general: null, show_price: false, unit_sale: 'dus' }),
  habis: contoh({ name: 'Cat tembok putih 5 kg', sku: 'CAT-PTH-05', slug: 'cat-tembok-putih-5-kg', price_general: 185000, unit_sale: 'galon', stock_status: 'HABIS' }),
}
const KEADAAN = [
  { nilai: 'biasa', label: 'Biasa' },
  { nilai: 'harga', label: 'Harga disembunyikan' },
  { nilai: 'habis', label: 'Stok habis' },
]
const KERANJANG_CONTOH = [
  { produk: PRODUK_CONTOH.biasa, jumlah: 20 },
  { produk: PRODUK_CONTOH.harga, jumlah: 12 },
  { produk: { ...PRODUK_CONTOH.habis, stock_status: 'PRE_ORDER' }, jumlah: 3 },
]
const PELANGGAN = [
  { nilai: '', label: 'Tamu' },
  { nilai: 'Budi Santoso (member TF-000123)', label: 'Member masuk' },
]

const keadaan = ref('biasa')
const pelanggan = ref('')
const linkContoh = (produk) => new URL(`/produk/${produk.slug}`, window.location.origin).href
const kosong = (teks) => !String(teks ?? '').trim()

const pratinjauProduk = computed(() =>
  kosong(form.wa_template_product)
    ? ''
    : pesanProduk(PRODUK_CONTOH[keadaan.value], { template: form.wa_template_product, link: linkContoh(PRODUK_CONTOH[keadaan.value]), pelanggan: pelanggan.value }),
)
const pratinjauKeranjang = computed(() =>
  kosong(form.wa_template_cart) || kosong(form.wa_template_cart_item)
    ? ''
    : pesanKeranjang(KERANJANG_CONTOH, { template: form.wa_template_cart, templateBaris: form.wa_template_cart_item, linkProduk: linkContoh, pelanggan: pelanggan.value }),
)
// "Coba di WhatsApp": buka pesan pratinjau ke nomor toko (belum terkirim sampai admin menekan kirim).
const coba = (teks) => (settings.waNumber && teks ? buatLinkWhatsApp(settings.waNumber, teks) : null)

const kelasPilihan = (aktif) => [
  'inline-flex min-h-11 cursor-pointer items-center rounded-md px-3 text-sm font-semibold transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring md:min-h-9',
  aktif ? 'bg-secondary text-secondary-foreground' : 'text-foreground hover:bg-subtle',
]
</script>

<template>
  <div v-if="status === 'loading'" class="flex max-w-6xl flex-col gap-6" aria-busy="true">
    <span class="sr-only">Memuat template…</span>
    <Skeleton v-for="n in 2" :key="n" class="h-72 rounded-xl" />
  </div>

  <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />

  <form v-else ref="formEl" novalidate class="flex flex-col gap-6" @submit.prevent="simpan">
    <div class="flex max-w-6xl flex-col gap-3">
      <p class="max-w-3xl text-sm text-muted">
        Isi pesan otomatis saat pelanggan menekan tombol WhatsApp. Pesan per produk dan pesan keranjang diatur terpisah; pelanggan masih bisa
        mengubah pesannya sebelum mengirim. Nomor tujuan diatur di
        <RouterLink to="/admin/pengaturan" class="rounded-sm font-semibold text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">Pengaturan toko</RouterLink>.
      </p>
      <div class="flex flex-wrap items-center gap-2">
        <span id="label-contoh-pelanggan" class="text-sm font-semibold">Pratinjau sebagai</span>
        <div class="inline-flex gap-1 rounded-lg border border-border bg-surface p-1" role="group" aria-labelledby="label-contoh-pelanggan">
          <button v-for="p in PELANGGAN" :key="p.label" type="button" :class="kelasPilihan(pelanggan === p.nilai)" :aria-pressed="pelanggan === p.nilai" @click="pelanggan = p.nilai">
            {{ p.label }}
          </button>
        </div>
      </div>
      <p v-if="galatUmum" role="alert" class="max-w-3xl rounded-md bg-danger-soft px-4 py-3 text-sm font-semibold text-danger-ink">{{ galatUmum }}</p>
    </div>

    <FormSection
      class="max-w-6xl"
      judul="Pesan per produk"
      deskripsi='Tombol "Tanya via WhatsApp", "Tanya harga", dan "Tanya ketersediaan" di kartu produk, daftar, dan halaman produk.'
      data-bagian-template="produk"
    >
      <div class="grid items-start gap-6 xl:grid-cols-2">
        <EditorTemplateWa v-model="form.wa_template_product" jenis="produk" label="Template pesan per produk" :rows="6" :error="galat.wa_template_product ?? ''" />
        <div class="flex min-w-0 flex-col gap-3">
          <div class="flex flex-wrap items-center gap-2">
            <span id="label-keadaan-produk" class="text-sm font-semibold">Contoh produk</span>
            <div class="inline-flex flex-wrap gap-1 rounded-lg border border-border bg-surface p-1" role="group" aria-labelledby="label-keadaan-produk">
              <button v-for="k in KEADAAN" :key="k.nilai" type="button" :class="kelasPilihan(keadaan === k.nilai)" :aria-pressed="keadaan === k.nilai" @click="keadaan = k.nilai">
                {{ k.label }}
              </button>
            </div>
          </div>
          <PratinjauWa :pesan="pratinjauProduk" />
          <Button variant="outline" size="sm" class="self-start" :href="coba(pratinjauProduk)">Coba di WhatsApp</Button>
        </div>
      </div>
    </FormSection>

    <FormSection
      class="max-w-6xl"
      judul="Pesan keranjang belanja"
      deskripsi='Tombol "Pesan via WhatsApp" di halaman Keranjang. Semua produk di keranjang masuk ke {daftar_produk}, masing-masing memakai format baris per produk.'
      data-bagian-template="keranjang"
    >
      <div class="grid items-start gap-6 xl:grid-cols-2">
        <div class="flex min-w-0 flex-col gap-6">
          <EditorTemplateWa v-model="form.wa_template_cart" jenis="keranjang" label="Template pesan keranjang" :rows="10" :error="galat.wa_template_cart ?? ''" />
          <EditorTemplateWa
            v-model="form.wa_template_cart_item"
            jenis="baris"
            label="Format baris per produk"
            :rows="4"
            hint="Diulang untuk setiap produk di keranjang."
            :error="galat.wa_template_cart_item ?? ''"
          />
        </div>
        <div class="flex min-w-0 flex-col gap-3 xl:sticky xl:top-6">
          <PratinjauWa :pesan="pratinjauKeranjang" label="Pratinjau (contoh 3 produk: harga tampil, harga disembunyikan, pre-order)" />
          <Button variant="outline" size="sm" class="self-start" :href="coba(pratinjauKeranjang)">Coba di WhatsApp</Button>
        </div>
      </div>
    </FormSection>

    <div class="sticky bottom-0 z-20 -mx-4 flex flex-wrap items-center justify-end gap-2 border-t border-border bg-surface px-4 py-3 sm:-mx-7 sm:px-7">
      <span v-if="kotor" class="mr-auto hidden text-sm text-muted sm:inline">Ada perubahan yang belum disimpan</span>
      <Button v-if="kotor" variant="outline" :disabled="menyimpan" @click="batal">Batalkan perubahan</Button>
      <Button type="submit" :loading="menyimpan" data-simpan-template>Simpan</Button>
    </div>

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
