<script setup>
import { ref } from 'vue'
import { Eye, EyeOff } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'

// Isian password dengan tombol mata (tampil / sembunyikan). Atribut lain diteruskan ke <input>.
defineOptions({ inheritAttrs: false })

const model = defineModel({ type: String, default: '' })

defineProps({
  label: { type: String, default: 'Password' },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
  required: { type: Boolean, default: false },
})

const tampil = ref(false)
</script>

<template>
  <Input v-model="model" :label="label" :type="tampil ? 'text' : 'password'" :error="error" :hint="hint" :required="required" v-bind="$attrs">
    <template #akhir>
      <button
        type="button"
        class="inline-flex size-10 items-center justify-center rounded-md text-muted transition-colors duration-150 hover:bg-subtle hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        :aria-label="tampil ? 'Sembunyikan password' : 'Tampilkan password'"
        :aria-pressed="tampil ? 'true' : 'false'"
        @click="tampil = !tampil"
      >
        <component :is="tampil ? EyeOff : Eye" class="size-4.5" aria-hidden="true" />
      </button>
    </template>
  </Input>
</template>
