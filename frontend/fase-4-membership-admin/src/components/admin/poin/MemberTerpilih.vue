<script setup>
import { RouterLink } from 'vue-router'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import { usePermission } from '@/composables/usePermission'
import { formatAngka, formatNomorHp, inisialNama } from '@/utils/format'

defineProps({
  // Member terpilih: { id, member_code, full_name, email, phone_number, tier, is_active, balance }
  member: { type: Object, required: true },
})

const emit = defineEmits(['ganti'])

const LABEL_TIER = { TAHAP_1: 'Tahap 1' }
const { can } = usePermission()
</script>

<template>
  <div class="flex flex-col gap-3 rounded-lg border border-border bg-background p-4" data-member-terpilih>
    <div class="flex flex-wrap items-center gap-x-8 gap-y-4">
      <div class="flex min-w-0 flex-[1_1_14rem] items-center gap-3">
        <span class="flex size-11 shrink-0 items-center justify-center rounded-full bg-primary-soft font-extrabold text-primary-hover" aria-hidden="true">
          {{ inisialNama(member.full_name) }}
        </span>
        <span class="flex min-w-0 flex-col">
          <component
            :is="can('member.view') ? RouterLink : 'strong'"
            :to="can('member.view') ? `/admin/member/${member.id}` : undefined"
            class="text-base font-bold [overflow-wrap:anywhere]"
            :class="can('member.view') && 'rounded-sm transition-colors duration-150 hover:text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring'"
          >
            {{ member.full_name }}
          </component>
          <span class="font-mono text-sm text-muted">{{ member.member_code }} · {{ LABEL_TIER[member.tier] ?? member.tier }}</span>
        </span>
      </div>
      <div class="flex min-w-0 flex-col">
        <span class="text-xs text-muted">Email</span>
        <span class="text-sm [overflow-wrap:anywhere]">{{ member.email }}</span>
      </div>
      <div class="flex flex-col">
        <span class="text-xs text-muted">Nomor HP</span>
        <span v-if="member.phone_number" class="text-sm tabular-nums">{{ formatNomorHp(member.phone_number) }}</span>
        <span v-else class="text-sm text-faint">Belum diisi</span>
      </div>
      <div class="flex flex-col">
        <span class="text-xs text-muted">Saldo poin</span>
        <strong class="text-2xl leading-tight tabular-nums" data-saldo-terpilih>{{ formatAngka(member.balance) }}</strong>
      </div>
      <Button variant="outline" size="sm" class="sm:ml-auto" data-ganti-member @click="emit('ganti')">Ganti member</Button>
    </div>
    <p v-if="!member.is_active" role="status" class="flex flex-wrap items-center gap-2 text-sm font-medium text-danger-ink">
      <Badge variant="danger">Nonaktif</Badge>
      Akun ini nonaktif, jadi poin tidak bisa dicatat. Aktifkan dulu dari detail member.
    </p>
  </div>
</template>
