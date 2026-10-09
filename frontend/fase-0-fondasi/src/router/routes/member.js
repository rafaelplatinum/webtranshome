import MemberLayout from '@/layouts/MemberLayout.vue'

// Area member: wajib login dengan role CUSTOMER (CLAUDE.md §3).
export default [
  {
    path: '/akun',
    component: MemberLayout,
    meta: { area: 'member' },
    // Child '' wajib ada: tanpa itu /akun hanya mencocokkan induk dan RouterView kosong.
    children: [
      { path: '', name: 'akun', component: () => import('@/views/member/Ringkasan.vue'), meta: { title: 'Akun saya' } },
      { path: 'poin', name: 'akun-poin', component: () => import('@/views/member/RiwayatPoin.vue'), meta: { title: 'Riwayat poin' } },
      { path: 'profil', name: 'akun-profil', component: () => import('@/views/member/Profil.vue'), meta: { title: 'Profil' } },
      { path: 'keamanan', name: 'akun-keamanan', component: () => import('@/views/member/Keamanan.vue'), meta: { title: 'Keamanan' } },
      {
        path: ':bagian(.*)',
        name: 'akun-tidak-ditemukan',
        component: () => import('@/views/member/TidakDitemukan.vue'),
        meta: { title: 'Halaman tidak ditemukan' },
      },
    ],
  },
]
