/** true bila pengguna meminta animasi dikurangi (prefers-reduced-motion). */
export function kurangiGerak() {
  return window.matchMedia?.('(prefers-reduced-motion: reduce)').matches === true
}

/** Nilai `behavior` untuk scrollTo / scrollBy / scrollIntoView yang menghormati reduced motion. */
export function perilakuGulir() {
  return kurangiGerak() ? 'auto' : 'smooth'
}
