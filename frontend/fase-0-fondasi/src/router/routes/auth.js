import AuthLayout from '@/layouts/AuthLayout.vue'

// Setiap halaman berdiri sendiri (bukan di bawah induk "/") agar tidak bentrok dengan Beranda.
// area `guest`: yang sudah masuk diarahkan ke /akun atau /admin. area `public`: link dari email
// dan callback Google tetap bisa dibuka walau sedang masuk.
function halaman(path, name, title, component, area = 'guest') {
  return {
    path,
    component: AuthLayout,
    meta: { area },
    children: [{ path: '', name, component, meta: { title } }],
  }
}

export default [
  halaman('/masuk', 'masuk', 'Masuk', () => import('@/views/auth/Masuk.vue')),
  halaman('/daftar', 'daftar', 'Daftar member', () => import('@/views/auth/Daftar.vue')),
  halaman('/lupa-password', 'lupa-password', 'Lupa password', () => import('@/views/auth/LupaPassword.vue')),
  halaman('/verifikasi', 'verifikasi', 'Verifikasi email', () => import('@/views/auth/Verifikasi.vue'), 'public'),
  halaman('/reset-password', 'reset-password', 'Buat password baru', () => import('@/views/auth/ResetPassword.vue'), 'public'),
  halaman('/masuk/google', 'masuk-google', 'Masuk dengan Google', () => import('@/views/auth/MasukGoogle.vue'), 'public'),
  halaman('/admin/masuk', 'admin-masuk', 'Masuk admin', () => import('@/views/admin/Masuk.vue')),
]
