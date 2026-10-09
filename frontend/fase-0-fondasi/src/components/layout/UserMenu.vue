<script setup>
import { computed, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { ChevronDown, History, LayoutDashboard, LogOut, User } from 'lucide-vue-next'
import { useAuth } from '@/composables/useAuth'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const { logout } = useAuth()
const route = useRoute()
const router = useRouter()

const terbuka = ref(false)
const tombol = ref(null)
const idMenu = useId()

const nama = computed(() => auth.user?.full_name ?? auth.user?.email ?? 'Akun')
const inisial = computed(() =>
  nama.value
    .split(/[\s@.]+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((kata) => kata[0].toUpperCase())
    .join(''),
)

const tautan = computed(() =>
  auth.isAdmin
    ? [{ label: 'Admin panel', to: '/admin', ikon: LayoutDashboard }]
    : [
        { label: 'Akun saya', to: '/akun', ikon: User },
        { label: 'Riwayat poin', to: '/akun/poin', ikon: History },
      ],
)

function tutup({ kembalikanFokus = false } = {}) {
  terbuka.value = false
  if (kembalikanFokus) tombol.value?.focus()
}

async function keluar() {
  tutup()
  await logout()
  router.push('/')
}

function klikDiLuar(event) {
  if (terbuka.value && !event.target.closest?.('[data-user-menu]')) tutup()
}

watch(() => route.fullPath, () => tutup())
onMounted(() => document.addEventListener('mousedown', klikDiLuar))
onBeforeUnmount(() => document.removeEventListener('mousedown', klikDiLuar))
</script>

<template>
  <div data-user-menu class="relative" @keydown.esc="tutup({ kembalikanFokus: true })">
    <button
      ref="tombol"
      type="button"
      :aria-expanded="terbuka ? 'true' : 'false'"
      :aria-controls="idMenu"
      class="inline-flex h-11 items-center gap-2 rounded-full pr-2 pl-1 transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      @click="terbuka = !terbuka"
    >
      <span class="grid size-9 place-items-center rounded-full bg-primary-soft text-sm font-bold text-primary-hover" aria-hidden="true">{{ inisial }}</span>
      <span class="hidden max-w-40 truncate text-sm font-semibold sm:inline">{{ nama }}</span>
      <span class="sr-only sm:hidden">Menu akun</span>
      <ChevronDown class="size-4 text-muted" aria-hidden="true" />
    </button>

    <div v-show="terbuka" :id="idMenu" class="absolute top-full right-0 z-40 mt-2 w-60 rounded-xl border border-border bg-surface p-2 shadow-lg">
      <p v-if="auth.user?.member_code" class="px-3 pt-1 pb-2 text-xs text-muted">
        Member <span class="font-mono font-semibold text-foreground">{{ auth.user.member_code }}</span>
      </p>
      <RouterLink
        v-for="item in tautan"
        :key="item.to"
        :to="item.to"
        class="flex h-11 items-center gap-3 rounded-lg px-3 text-sm font-semibold transition-colors duration-150 hover:bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <component :is="item.ikon" class="size-4.5 text-muted" aria-hidden="true" />
        {{ item.label }}
      </RouterLink>
      <hr class="my-1 border-border" />
      <button
        type="button"
        class="flex h-11 w-full items-center gap-3 rounded-lg px-3 text-sm font-semibold text-danger transition-colors duration-150 hover:bg-danger-soft focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        @click="keluar"
      >
        <LogOut class="size-4.5" aria-hidden="true" />
        Keluar
      </button>
    </div>
  </div>
</template>
