import { nextTick } from 'vue'

// Animasi situs publik tanpa library: View Transitions API (pindah halaman) dan IntersectionObserver (scroll reveal).
// Panel admin tidak dianimasikan. Semua mati bila pengguna meminta animasi dikurangi.
// Gaya & durasi ada di main.css (bagian "Animasi situs publik").

const kurangiGerak = () => window.matchMedia?.('(prefers-reduced-motion: reduce)').matches === true
const halamanAdmin = (route) => route.path.startsWith('/admin')

/**
 * Pasang di router: menandai area di <html data-area="publik|admin"> (dipakai CSS efek tekan tombol & modal),
 * lalu crossfade saat pindah halaman publik. Ganti filter/query di halaman yang sama tidak dianimasikan.
 */
export function pasangTransisiHalaman(router) {
  let selesai = null

  router.beforeResolve((to, from) => {
    const pindahHalaman = from.matched.length > 0 && to.path !== from.path
    if (!document.startViewTransition || !pindahHalaman || halamanAdmin(to) || halamanAdmin(from) || kurangiGerak()) return
    // Tampilan lama dipotret dulu, lalu navigasi dilanjutkan; tampilan baru dipotret setelah halaman baru dirender.
    return new Promise((lanjut) => {
      document.startViewTransition(
        () =>
          new Promise((beres) => {
            selesai = beres
            lanjut()
          }),
      )
    })
  })

  router.afterEach(async (to) => {
    document.documentElement.dataset.area = halamanAdmin(to) ? 'admin' : 'publik'
    if (!selesai) return
    const beres = selesai
    selesai = null
    await nextTick()
    beres()
  })

  router.onError(() => {
    selesai?.()
    selesai = null
  })
}

// ---------- scroll reveal ----------

const JEDA_STAGGER = 60
const MAKS_URUTAN = 5

let pengamat = null
function amati() {
  pengamat ??= new IntersectionObserver(
    (entri) => {
      for (const e of entri) {
        if (!e.isIntersecting) continue
        e.target.classList.add('muncul-tampil')
        pengamat.unobserve(e.target)
      }
    },
    { rootMargin: '0px 0px -8% 0px' },
  )
  return pengamat
}

// Setelah selesai, kelas dilepas supaya transisi milik elemen itu sendiri (hover, dll.) tidak tertimpa.
function bersihkan(event) {
  if (event.target !== event.currentTarget || !event.currentTarget.classList.contains('muncul-tampil')) return
  event.currentTarget.classList.remove('muncul', 'muncul-tampil')
  event.currentTarget.style.removeProperty('--jeda-muncul')
  event.currentTarget.removeEventListener('transitionend', bersihkan)
}

/**
 * `v-reveal` (atau `v-reveal="urutan"` untuk muncul bergantian): elemen memudar dan naik 20 px saat masuk layar.
 * Elemen yang sudah terlihat saat halaman dibuka langsung tampil, supaya konten pertama tidak tertunda.
 */
export const vReveal = {
  mounted(el, binding) {
    if (!(el instanceof Element) || typeof IntersectionObserver === 'undefined' || kurangiGerak()) return
    if (document.documentElement.dataset.area === 'admin') return
    if (el.getBoundingClientRect().top < window.innerHeight) return
    const urutan = Math.min(Math.max(Number(binding.value) || 0, 0), MAKS_URUTAN)
    el.style.setProperty('--jeda-muncul', `${urutan * JEDA_STAGGER}ms`)
    el.classList.add('muncul')
    el.addEventListener('transitionend', bersihkan)
    amati().observe(el)
  },
  beforeUnmount(el) {
    pengamat?.unobserve(el)
  },
}
