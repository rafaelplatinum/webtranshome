<script setup>
import { useId } from 'vue'
import { BadgeCheck } from 'lucide-vue-next'
import Badge from '@/components/ui/Badge.vue'
import { formatNomorHp, formatTanggal, formatTanggalJam } from '@/utils/format'

defineProps({
  // Detail GET /admin/members/:id
  member: { type: Object, required: true },
})

const idJudul = useId()
</script>

<template>
  <section :aria-labelledby="idJudul" class="flex flex-col gap-4 rounded-xl border border-border bg-surface p-5">
    <h2 :id="idJudul" class="text-base font-bold">Data akun</h2>
    <dl class="grid gap-4 sm:grid-cols-2">
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs font-semibold text-muted">Email</dt>
        <dd class="flex flex-wrap items-center gap-2">
          <span class="[overflow-wrap:anywhere]">{{ member.email }}</span>
          <Badge v-if="member.email_verified" variant="success"><BadgeCheck class="size-3.5" aria-hidden="true" />Terverifikasi</Badge>
          <Badge v-else variant="warning">Belum diverifikasi</Badge>
        </dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs font-semibold text-muted">Nomor HP</dt>
        <dd v-if="member.phone_number" class="tabular-nums" data-nomor-hp>{{ formatNomorHp(member.phone_number) }}</dd>
        <dd v-else class="text-faint">Belum diisi</dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs font-semibold text-muted">Info promo dan poin</dt>
        <dd>{{ member.communication_consent ? `Setuju${member.consent_at ? ` sejak ${formatTanggal(member.consent_at)}` : ''}` : 'Tidak setuju' }}</dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs font-semibold text-muted">Akun Google</dt>
        <dd>{{ member.google_connected ? 'Terhubung' : 'Tidak terhubung' }}</dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs font-semibold text-muted">Terdaftar</dt>
        <dd>{{ formatTanggal(member.created_at) }}</dd>
      </div>
      <div class="flex flex-col gap-0.5">
        <dt class="text-xs font-semibold text-muted">Login terakhir</dt>
        <dd>{{ member.last_login_at ? formatTanggalJam(member.last_login_at) : 'Belum pernah' }}</dd>
      </div>
    </dl>
  </section>
</template>
