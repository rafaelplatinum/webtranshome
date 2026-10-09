import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { getMenus } from '@/services/adminService'

// Menu sidebar admin selalu dari API menus sesuai izin user (sudah difilter backend).
// JANGAN hardcode menu per role (CLAUDE.md §4).
export const useMenuStore = defineStore('menu', () => {
  // Bentuk item mengikuti tabel menus: { id, parent_id, name, code, route, icon, sort_order }
  const items = ref([])
  const status = ref('idle') // idle | loading | ready | error
  let milikUser = null

  const urut = (a, b) => (a.sort_order ?? 0) - (b.sort_order ?? 0)

  /** Grup (menu tanpa route) berisi anak-anaknya; menu ber-route tanpa induk tampil sebagai grup satu item. */
  const grup = computed(() => {
    const induk = items.value.filter((m) => m.parent_id == null).sort(urut)
    return induk
      .map((g) => ({
        ...g,
        anak: g.route ? [g] : items.value.filter((m) => m.parent_id === g.id && m.route).sort(urut),
      }))
      .filter((g) => g.anak.length)
  })

  /** Muat sekali per user; `paksa` untuk tombol Coba lagi. */
  async function muat(userId, { paksa = false } = {}) {
    if (!paksa && milikUser === userId && (status.value === 'ready' || status.value === 'loading')) return
    milikUser = userId
    status.value = 'loading'
    try {
      items.value = await getMenus()
      status.value = 'ready'
    } catch {
      status.value = 'error'
    }
  }

  function setItems(daftar) {
    items.value = [...daftar]
    status.value = 'ready'
  }

  function reset() {
    items.value = []
    status.value = 'idle'
    milikUser = null
  }

  return { items, status, grup, muat, setItems, reset }
})
