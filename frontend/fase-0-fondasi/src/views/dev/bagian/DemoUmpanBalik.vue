<script setup>
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import Input from '@/components/ui/Input.vue'
import Modal from '@/components/ui/Modal.vue'
import { useUiStore } from '@/stores/ui'

const ui = useUiStore()
const modalBuka = ref(false)
const konfirmasiBuka = ref(false)
const menyimpan = ref(false)
const namaBrand = ref('')

const TOAST = [
  { jenis: 'success', pesan: 'Produk disimpan.' },
  { jenis: 'info', pesan: 'Membuka WhatsApp…' },
  { jenis: 'warning', pesan: 'Terlalu banyak percobaan. Coba lagi beberapa saat.' },
  { jenis: 'danger', pesan: 'Kamu tidak punya akses untuk aksi ini.' },
]

function nonaktifkan() {
  menyimpan.value = true
  setTimeout(() => {
    menyimpan.value = false
    konfirmasiBuka.value = false
    ui.tampilkanToast({ pesan: 'Produk dinonaktifkan.', jenis: 'success' })
  }, 1200)
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap gap-3">
      <Button variant="outline" @click="modalBuka = true">Buka modal form</Button>
      <Button variant="danger" @click="konfirmasiBuka = true">Nonaktifkan produk</Button>
    </div>
    <div class="flex flex-wrap gap-3">
      <Button v-for="t in TOAST" :key="t.jenis" variant="ghost" @click="ui.tampilkanToast(t)">Toast {{ t.jenis }}</Button>
    </div>

    <Modal v-model:open="modalBuka" title="Tambah brand" description="Brand baru langsung bisa dipilih di form produk.">
      <Input v-model="namaBrand" label="Nama brand" required />
      <template #footer>
        <Button variant="outline" @click="modalBuka = false">Batal</Button>
        <Button @click="modalBuka = false">Simpan</Button>
      </template>
    </Modal>

    <ConfirmModal
      v-model:open="konfirmasiBuka"
      title="Nonaktifkan produk ini?"
      message="Produk tidak akan tampil di website, tapi tetap tersimpan dan bisa diaktifkan lagi."
      confirm-label="Ya, nonaktifkan"
      danger
      :loading="menyimpan"
      @confirm="nonaktifkan"
    />
  </div>
</template>
