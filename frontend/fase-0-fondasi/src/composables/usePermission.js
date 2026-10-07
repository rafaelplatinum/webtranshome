import { useAuthStore } from '@/stores/auth'

/** Cek izin di script: `const { can } = usePermission(); can('product.update')` */
export function usePermission() {
  const auth = useAuthStore()
  return { can: (kode) => auth.can(kode) }
}

/**
 * Directive `v-can="'product.update'"`: elemen disembunyikan bila user tidak punya izin.
 * Hanya tampilan; backend tetap wajib menolak dengan 403 (CLAUDE.md §4).
 */
function terapkan(el, binding) {
  el.hidden = !useAuthStore().can(binding.value)
}

export const vCan = {
  mounted: terapkan,
  updated: terapkan,
}
