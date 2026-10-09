<script setup>
import { ref } from 'vue'
import { MessageCircle, PackageSearch } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import SkeletonCard from '@/components/ui/SkeletonCard.vue'

const mencobaLagi = ref(false)

function cobaLagi() {
  mencobaLagi.value = true
  setTimeout(() => (mencobaLagi.value = false), 1200)
}
</script>

<template>
  <div class="grid gap-6 lg:grid-cols-2">
    <EmptyState title="Belum ada produk yang cocok" description="Coba kurangi filter, atau tanyakan langsung ke tim kami.">
      <template #icon><PackageSearch class="size-10" aria-hidden="true" /></template>
      <template #actions>
        <Button variant="outline">Reset filter</Button>
        <Button variant="whatsapp">
          <template #icon><MessageCircle class="size-4.5" aria-hidden="true" /></template>
          Tanya via WhatsApp
        </Button>
      </template>
    </EmptyState>
    <ErrorState message="Tidak bisa terhubung ke server. Periksa koneksi, lalu coba lagi." :retrying="mencobaLagi" @retry="cobaLagi" />
    <div class="flex flex-col gap-3 rounded-xl border border-border bg-surface p-5">
      <p class="text-sm font-semibold text-muted">Skeleton baris</p>
      <Skeleton class="h-5 w-1/2" />
      <Skeleton class="h-4 w-full" />
      <Skeleton class="h-4 w-5/6" />
    </div>
    <div class="max-w-64">
      <SkeletonCard />
    </div>
  </div>
</template>
