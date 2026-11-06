<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Inbox, RefreshCw } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import ErrorState from '@/components/ui/ErrorState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { getMockEmails } from '@/services/devService'
import { pesanError } from '@/services/errors'
import { formatWaktuRelatif } from '@/utils/format'

// Email yang "dikirim" mode mock: link verifikasi dan reset password bisa dibuka dari sini.
const status = ref('loading')
const email = ref([])
const pesan = ref('')

async function muat() {
  status.value = 'loading'
  try {
    email.value = await getMockEmails()
    status.value = 'ready'
  } catch (error) {
    pesan.value = pesanError(error)
    status.value = 'error'
  }
}

onMounted(muat)
</script>

<template>
  <section class="mx-auto flex max-w-3xl flex-col gap-5 px-4 py-10 sm:px-6">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div class="flex flex-col gap-1">
        <p class="text-sm font-semibold text-muted">Mode mock</p>
        <h1 class="text-2xl font-extrabold">Kotak email tiruan</h1>
        <p class="text-sm text-muted">Email verifikasi dan reset password tidak benar-benar dikirim. Buka link-nya dari sini.</p>
      </div>
      <Button variant="outline" size="sm" @click="muat">
        <template #icon><RefreshCw class="size-4" aria-hidden="true" /></template>
        Muat ulang
      </Button>
    </div>

    <div v-if="status === 'loading'" class="flex flex-col gap-2" aria-busy="true">
      <Skeleton v-for="n in 3" :key="n" class="h-20 rounded-xl" />
    </div>
    <ErrorState v-else-if="status === 'error'" :message="pesan" @retry="muat" />
    <EmptyState v-else-if="!email.length" title="Belum ada email" description="Coba daftar member atau minta link reset password dulu.">
      <template #icon><Inbox class="size-10" aria-hidden="true" /></template>
    </EmptyState>
    <ul v-else class="flex flex-col gap-2">
      <li v-for="m in email" :key="m.id" class="flex flex-wrap items-center gap-3 rounded-xl border border-border bg-surface p-4" data-email-mock>
        <div class="flex min-w-0 flex-1 flex-col gap-0.5">
          <span class="font-semibold">{{ m.subject }}</span>
          <span class="truncate text-sm text-muted">Kepada {{ m.to }} · {{ formatWaktuRelatif(m.created_at) }}</span>
        </div>
        <Button size="sm" :to="m.link">{{ m.jenis === 'verifikasi' ? 'Verifikasi email' : 'Atur ulang password' }}</Button>
      </li>
    </ul>
    <RouterLink to="/dev/components" class="self-start rounded-sm text-sm font-semibold text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
      Kembali ke galeri komponen
    </RouterLink>
  </section>
</template>
