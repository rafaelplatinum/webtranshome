<script setup>
import { nextTick, ref } from 'vue'
import { Search } from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { searchMembers } from '@/services/adminMemberService'
import { errorField, pesanError } from '@/services/errors'
import { formatNomorHp } from '@/utils/format'

// Cari (Enter): 1 hasil → langsung dipilih; lebih dari 1 → daftar pilihan; 0 → "Member tidak ditemukan".
const emit = defineEmits(['pilih'])

const q = ref('')
const galat = ref('')
const mencari = ref(false)
const hasil = ref([])
const wadah = ref(null)

async function cari() {
  if (mencari.value) return
  const kata = q.value.trim()
  hasil.value = []
  if (kata.length < 3) {
    galat.value = kata ? 'Isi minimal 3 karakter.' : 'Isi kode member, email, atau nomor HP dulu.'
    return
  }
  galat.value = ''
  mencari.value = true
  try {
    const data = await searchMembers(kata)
    if (!data.length) galat.value = 'Member tidak ditemukan. Periksa lagi kode member (contoh TF-000123), email, atau nomor HP.'
    else if (data.length === 1) pilih(data[0])
    else hasil.value = data
  } catch (error) {
    galat.value = errorField(error).q ?? pesanError(error)
  } finally {
    mencari.value = false
  }
}

function pilih(member) {
  hasil.value = []
  q.value = ''
  emit('pilih', member)
}

async function fokus() {
  await nextTick()
  wadah.value?.querySelector('input')?.focus()
}

defineExpose({ fokus })
</script>

<template>
  <div ref="wadah" class="flex flex-col gap-3">
    <form class="flex flex-wrap items-start gap-2.5" role="search" novalidate @submit.prevent="cari">
      <div class="min-w-0 flex-[1_1_20rem]">
        <Input
          v-model="q"
          label="Kode member, email, atau nomor HP"
          type="search"
          placeholder="TF-000123"
          autocomplete="off"
          enterkeyhint="search"
          class="font-mono"
          :error="galat"
          data-cari-member
        />
      </div>
      <!-- Mulai 640px sebaris dengan isian: geser setinggi label supaya tetap sejajar walau ada pesan galat. -->
      <Button type="submit" variant="secondary" class="w-full sm:mt-6.5 sm:w-auto" :loading="mencari">
        <template #icon><Search class="size-4.5" aria-hidden="true" /></template>
        Cari
      </Button>
    </form>

    <div v-if="hasil.length" class="flex flex-col gap-2">
      <p class="text-sm font-semibold" role="status">Ditemukan {{ hasil.length }} member. Pilih salah satu:</p>
      <ul class="flex flex-col gap-2" data-hasil-cari>
        <li v-for="m in hasil" :key="m.id">
          <button
            type="button"
            class="flex w-full flex-wrap items-center gap-x-4 gap-y-1 rounded-lg border border-border bg-surface px-4 py-3 text-left transition-colors duration-150 hover:border-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            @click="pilih(m)"
          >
            <span class="flex min-w-0 flex-[1_1_12rem] flex-col">
              <strong class="font-semibold">{{ m.full_name }}</strong>
              <span class="font-mono text-xs text-muted">{{ m.member_code }}</span>
            </span>
            <span class="min-w-0 flex-[1_1_12rem] text-sm text-muted [overflow-wrap:anywhere]">{{ m.email }}</span>
            <span class="text-sm text-muted tabular-nums">{{ m.phone_number ? formatNomorHp(m.phone_number) : 'Tanpa nomor HP' }}</span>
            <Badge v-if="!m.is_active" variant="danger">Nonaktif</Badge>
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>
