import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useMenuStore } from '@/stores/menu'
import { useUiStore } from '@/stores/ui'
import * as authService from '@/services/authService'

/**
 * Status login dan aksi masuk/keluar. Halaman memakai ini, bukan memanggil service auth langsung.
 * Navigasi setelah masuk/keluar diatur halaman pemanggil (redirect, /akun, /admin, /).
 */
export function useAuth() {
  const auth = useAuthStore()

  async function login(kredensial, { admin = false } = {}) {
    const sesi = admin ? await authService.loginAdmin(kredensial) : await authService.login(kredensial)
    auth.setSession(sesi)
    return sesi
  }

  function fetchMe() {
    return auth.restoreSession()
  }

  async function logout() {
    try {
      await authService.logout()
    } catch {
      // Gagal jaringan: tetap keluar secara lokal.
    }
    auth.logout()
    useMenuStore().reset()
    useUiStore().tampilkanToast({ pesan: 'Kamu sudah keluar.', jenis: 'success' })
  }

  return {
    user: computed(() => auth.user),
    isLoggedIn: computed(() => auth.isLoggedIn),
    isAdmin: computed(() => auth.isAdmin),
    login,
    logout,
    fetchMe,
  }
}
