<script setup>
import { computed, useId } from 'vue'
import { ChevronDown } from 'lucide-vue-next'

defineOptions({ inheritAttrs: false })

const model = defineModel({ type: [String, Number], default: '' })

const props = defineProps({
  label: { type: String, required: true },
  // [{ value, label }]
  options: { type: Array, default: () => [] },
  // Opsi berkelompok (<optgroup>): [{ label, options: [{ value, label }] }], tampil setelah `options`.
  groups: { type: Array, default: () => [] },
  placeholder: { type: String, default: '' },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
  required: { type: Boolean, default: false },
  hideLabel: { type: Boolean, default: false },
})

const id = useId()
const keterangan = computed(() => {
  if (props.error) return `${id}-error`
  if (props.hint) return `${id}-hint`
  return undefined
})
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label :for="id" :class="hideLabel ? 'sr-only' : 'text-sm font-semibold'">
      {{ label }}<span v-if="required" class="text-danger" aria-hidden="true"> *</span>
    </label>
    <div class="relative">
      <select
        :id="id"
        v-model="model"
        :required="required"
        :aria-invalid="error ? 'true' : undefined"
        :aria-describedby="keterangan"
        v-bind="$attrs"
        class="h-11 w-full min-w-0 cursor-pointer appearance-none rounded-md border bg-surface pr-10 pl-3 text-base text-foreground transition-colors duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:bg-subtle"
        :class="error ? 'border-danger' : 'border-border-strong focus-visible:border-primary'"
      >
        <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
        <option v-for="opsi in options" :key="opsi.value" :value="opsi.value">{{ opsi.label }}</option>
        <optgroup v-for="grup in groups" :key="grup.label" :label="grup.label">
          <option v-for="opsi in grup.options" :key="opsi.value" :value="opsi.value">{{ opsi.label }}</option>
        </optgroup>
      </select>
      <ChevronDown class="pointer-events-none absolute top-1/2 right-3 size-4 -translate-y-1/2 text-muted" aria-hidden="true" />
    </div>
    <p v-if="error" :id="`${id}-error`" class="text-sm font-medium text-danger-ink">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-hint`" class="text-sm text-muted">{{ hint }}</p>
  </div>
</template>
