<script setup>
import { nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import { X } from 'lucide-vue-next'

const open = defineModel('open', { type: Boolean, default: false })

const props = defineProps({
  title: { type: String, required: true },
  description: { type: String, default: '' },
  size: { type: String, default: 'md', validator: (v) => ['sm', 'md', 'lg', 'xl'].includes(v) },
  // false: Esc dan klik latar tidak menutup (mis. saat proses simpan berjalan).
  dismissible: { type: Boolean, default: true },
  hideClose: { type: Boolean, default: false },
  // true: di HP tampil sebagai bottom sheet selebar layar (mis. filter katalog); mulai 640px tetap modal biasa.
  sheet: { type: Boolean, default: false },
  // true: panel di sisi kanan setinggi layar (drawer), mis. detail log aktivitas. Di HP selebar layar.
  side: { type: Boolean, default: false },
})

const emit = defineEmits(['close'])

const LEBAR = { sm: 'max-w-sm', md: 'max-w-lg', lg: 'max-w-2xl', xl: 'max-w-5xl' }
const SELEKTOR_FOKUS =
  'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'

const panel = ref(null)
const idJudul = useId()
const idDeskripsi = useId()
let pemicu = null

function tutup() {
  open.value = false
  emit('close')
}

function elemenFokus() {
  return [...(panel.value?.querySelectorAll(SELEKTOR_FOKUS) ?? [])]
}

// Fokus terkunci di dalam modal; Esc menutup bila boleh.
function onKeydown(event) {
  if (event.key === 'Escape') {
    event.stopPropagation()
    if (props.dismissible) tutup()
    return
  }
  if (event.key !== 'Tab') return
  const daftar = elemenFokus()
  if (daftar.length === 0) {
    event.preventDefault()
    return
  }
  const pertama = daftar[0]
  const terakhir = daftar[daftar.length - 1]
  if (event.shiftKey && document.activeElement === pertama) {
    event.preventDefault()
    terakhir.focus()
  } else if (!event.shiftKey && document.activeElement === terakhir) {
    event.preventDefault()
    pertama.focus()
  }
}

watch(
  open,
  async (buka) => {
    if (buka) {
      pemicu = document.activeElement
      document.body.style.overflow = 'hidden'
      await nextTick()
      ;(elemenFokus()[0] ?? panel.value)?.focus()
    } else {
      document.body.style.overflow = ''
      pemicu?.focus?.()
      pemicu = null
    }
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  document.body.style.overflow = ''
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-150 ease-out"
      enter-from-class="opacity-0"
      leave-active-class="transition-opacity duration-150 ease-in"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex bg-foreground/50"
        :class="side ? 'items-stretch justify-end' : ['items-end justify-center sm:items-center sm:p-4', sheet ? 'p-0' : 'p-4']"
        @click.self="dismissible && tutup()"
        @keydown="onKeydown"
      >
        <!-- Judul dan footer tetap terlihat; hanya isi yang bergulir bila panjang. -->
        <div
          ref="panel"
          data-panel-modal
          :data-sheet="sheet ? '' : undefined"
          role="dialog"
          aria-modal="true"
          :aria-labelledby="idJudul"
          :aria-describedby="description ? idDeskripsi : undefined"
          tabindex="-1"
          class="flex w-full flex-col overflow-hidden bg-surface shadow-xl focus:outline-none"
          :class="
            side
              ? [LEBAR[size], 'h-dvh sm:border-l sm:border-border']
              : [LEBAR[size], sheet ? 'max-h-[85dvh] rounded-t-2xl sm:max-h-[90dvh] sm:rounded-2xl' : 'max-h-[90dvh] rounded-2xl']
          "
        >
          <div class="flex items-start justify-between gap-4 px-6 pt-6" :class="!$slots.default && !$slots.footer && 'pb-6'">
            <div class="flex flex-col gap-1">
              <h2 :id="idJudul" class="text-lg font-bold">{{ title }}</h2>
              <p v-if="description" :id="idDeskripsi" class="text-sm text-muted">{{ description }}</p>
            </div>
            <button
              v-if="!hideClose && dismissible"
              type="button"
              class="-m-2 inline-flex size-11 shrink-0 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              aria-label="Tutup"
              @click="tutup"
            >
              <X class="size-5" aria-hidden="true" />
            </button>
          </div>
          <div v-if="$slots.default" class="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-4" :class="$slots.footer ? 'pb-1' : 'pb-6'">
            <slot />
          </div>
          <div
            v-if="$slots.footer"
            class="gap-2 px-6"
            :class="
              sheet
                ? 'grid grid-cols-[auto_minmax(0,1fr)] border-t border-border pt-4 pb-[calc(1rem+env(safe-area-inset-bottom))] sm:flex sm:justify-end sm:pb-6'
                : 'flex flex-col-reverse pt-4 pb-6 sm:flex-row sm:justify-end'
            "
          >
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
