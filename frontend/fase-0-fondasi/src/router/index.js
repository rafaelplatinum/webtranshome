import { createRouter, createWebHistory } from 'vue-router'
import publicRoutes from './routes/public'
import authRoutes from './routes/auth'
import memberRoutes from './routes/member'
import adminRoutes from './routes/admin'
import { pasangGuard } from './guard'
import { pasangTransisiHalaman } from '@/composables/useAnimasi'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  // Urutan tidak menentukan prioritas: Vue Router memilih path paling spesifik.
  routes: [...publicRoutes, ...authRoutes, ...memberRoutes, ...adminRoutes],
  scrollBehavior(to, from, savedPosition) {
    // Back/Forward: kembali ke posisi semula (daftar katalog sudah tampil dari cache useKatalog).
    if (savedPosition) return savedPosition
    // Hanya query yang berubah (filter, urutan, halaman katalog): posisi tetap; paginasi menggulir sendiri.
    if (to.path === from.path) return false
    return { top: 0 }
  },
})

pasangGuard(router)
pasangTransisiHalaman(router)

export default router
