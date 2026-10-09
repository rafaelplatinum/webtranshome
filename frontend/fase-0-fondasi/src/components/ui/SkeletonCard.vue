<script setup>
import Skeleton from './Skeleton.vue'
</script>

<template>
  <!-- Bentuk sama dengan ProductCard agar tidak ada layout shift saat data datang. -->
  <div class="flex flex-col overflow-hidden rounded-xl border border-border bg-surface" aria-hidden="true">
    <Skeleton class="aspect-square rounded-none" />
    <div class="flex flex-col gap-2 p-4">
      <Skeleton class="h-3 w-16" />
      <Skeleton class="h-4 w-full" />
      <Skeleton class="h-4 w-2/3" />
      <Skeleton class="mt-1 h-6 w-28" />
      <Skeleton class="h-5 w-20" />
      <Skeleton class="mt-2 h-11 w-full" />
    </div>
  </div>
</template>
