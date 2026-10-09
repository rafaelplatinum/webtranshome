import { toValue, watchEffect } from 'vue'

/**
 * Title dan meta description per halaman (pakai meta_title / meta_description jika ada).
 * Menerima nilai biasa, ref, atau getter: useSeo({ title: () => produk.value?.name }).
 * Catatan SPA: pratinjau link di WhatsApp/Facebook tetap memakai meta default di index.html.
 */
export function useSeo({ title, description } = {}) {
  watchEffect(() => {
    const judul = toValue(title)
    const deskripsi = toValue(description)
    if (judul) document.title = `${judul} | Transhome`
    if (deskripsi) document.querySelector('meta[name="description"]')?.setAttribute('content', deskripsi)
  })
}
