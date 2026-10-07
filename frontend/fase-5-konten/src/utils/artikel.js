// Artikel, promo, event, dan halaman (tabel articles). Dipakai halaman publik dan editor admin.

export const TIPE_ARTIKEL = {
  ARTIKEL: { label: 'Artikel', kelas: 'bg-subtle text-foreground' },
  PROMO: { label: 'Promo', kelas: 'bg-primary-soft text-primary-hover' },
  EVENT: { label: 'Event', kelas: 'bg-info-soft text-info-ink' },
  HALAMAN: { label: 'Halaman', kelas: 'bg-subtle text-muted' },
}

/**
 * Isi artikel → blok untuk ditampilkan (tanpa HTML mentah, jadi aman dari skrip sisipan).
 * Paragraf dipisah baris kosong; blok berawalan "## " jadi subjudul; blok yang semua barisnya "- " jadi daftar.
 */
export function blokIsi(isi) {
  return String(isi ?? '')
    .replace(/\r\n/g, '\n')
    .split(/\n\s*\n/)
    .map((blok) => blok.trim())
    .filter(Boolean)
    .map((blok) => {
      if (blok.startsWith('## ')) return { jenis: 'h2', teks: blok.slice(3).replace(/\s+/g, ' ').trim() }
      const baris = blok.split('\n').map((b) => b.trim())
      if (baris.every((b) => b.startsWith('- '))) return { jenis: 'ul', butir: baris.map((b) => b.slice(2).trim()) }
      return { jenis: 'p', teks: baris.join(' ') }
    })
}

/** Alamat publik artikel; token pratinjau hanya untuk draf (dibuka dari panel admin). */
export const tautanArtikel = (slug, token = null) => `/artikel/${slug}${token ? `?preview=${encodeURIComponent(token)}` : ''}`

/** Halaman daftar asal sebuah artikel: promo & event di /promo, sisanya di /artikel. */
export const daftarAsal = (tipe) => (tipe === 'PROMO' || tipe === 'EVENT' ? { label: 'Promo dan event', to: '/promo' } : { label: 'Artikel', to: '/artikel' })
