<script setup>
import { computed, useId } from 'vue'

defineOptions({ inheritAttrs: false })

const model = defineModel({ type: String, default: '' })

const props = defineProps({
  label: { type: String, required: true },
  rows: { type: Number, default: 4 },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
  required: { type: Boolean, default: false },
  // Tampilkan hitungan karakter bila ada batas, mis. meta description 160.
  maxlength: { type: Number, default: null },
})

const id = useId()
const keterangan = computed(() => [props.error ? `${id}-error` : props.hint ? `${id}-hint` : null, props.maxlength ? `${id}-hitung` : null].filter(Boolean).join(' ') || undefined)
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label :for="id" class="text-sm font-semibold">
      {{ label }}<span v-if="required" class="text-danger" aria-hidden="true"> *</span>
    </label>
    <textarea
      :id="id"
      v-model="model"
      :rows="rows"
      :required="required"
      :maxlength="maxlength ?? undefined"
      :aria-invalid="error ? 'true' : undefined"
      :aria-describedby="keterangan"
      v-bind="$attrs"
      class="w-full min-w-0 rounded-md border bg-surface px-3 py-2.5 text-base leading-relaxed text-foreground transition-colors duration-150 placeholder:text-faint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      :class="error ? 'border-danger' : 'border-border-strong focus-visible:border-primary'"
    />
    <div class="flex justify-between gap-3">
      <p v-if="error" :id="`${id}-error`" class="text-sm font-medium text-danger-ink">{{ error }}</p>
      <p v-else-if="hint" :id="`${id}-hint`" class="text-sm text-muted">{{ hint }}</p>
      <p v-if="maxlength" :id="`${id}-hitung`" class="ml-auto shrink-0 text-xs text-muted tabular-nums">{{ (model ?? '').length }}/{{ maxlength }}</p>
    </div>
  </div>
</template>
