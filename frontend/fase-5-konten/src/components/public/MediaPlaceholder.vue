<script setup>
import { computed } from 'vue'

const props = defineProps({
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
  tone: { type: String, default: 'sand' },
})

const TONE = {
  sand: 'bg-th-sand',
  blue: 'bg-th-blue',
  yellow: 'bg-th-yellow',
  peach: 'bg-th-peach',
  sage: 'bg-th-sage',
  navy: 'bg-th-navy',
}
const kelas = computed(() => TONE[props.tone] ?? TONE.sand)
</script>

<template>
  <div class="relative isolate flex min-h-40 items-center justify-center overflow-hidden" :class="kelas">
    <img v-if="src" :src="src" :alt="alt" loading="lazy" class="absolute inset-0 size-full object-cover" />
    <slot v-else>
      <div aria-hidden="true" class="pointer-events-none absolute inset-0 overflow-hidden opacity-60">
        <span class="absolute -top-1/3 -right-[12%] size-[75%] rounded-full border border-white/60"></span>
        <span class="absolute -right-[4%] -bottom-1/2 size-[64%] rounded-full border border-white/50"></span>
        <span class="absolute right-[18%] bottom-[-55%] size-[78%] rounded-[40%] border border-white/40"></span>
      </div>
    </slot>
  </div>
</template>
