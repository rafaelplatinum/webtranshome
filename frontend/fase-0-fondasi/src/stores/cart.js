import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'

const KUNCI = 'transhome.cart'

function bacaKeranjang() {
  try {
    const data = JSON.parse(localStorage.getItem(KUNCI) || '[]')
    return Array.isArray(data) ? data : []
  } catch {
    return []
  }
}

export const useCartStore = defineStore('cart', () => {
  const items = ref(bacaKeranjang())
  const jumlah = computed(() => items.value.reduce((total, item) => total + item.quantity, 0))

  watch(items, (data) => {
    try { localStorage.setItem(KUNCI, JSON.stringify(data)) } catch { /* Cart tetap tersedia selama halaman terbuka. */ }
  }, { deep: true })

  function tambah(product) {
    const item = items.value.find((entry) => entry.sku === product.sku)
    if (item) item.quantity += 1
    else items.value.push({
      sku: product.sku,
      name: product.name,
      slug: product.slug,
      image: product.primary_image_url || '',
      brand: product.brand?.name || '',
      price: product.price_general ?? null,
      showPrice: product.show_price !== false,
      unit: product.unit_sale || '',
      stockStatus: product.stock_status || 'TERSEDIA',
      stockLabel: product.stock_qty_label || '',
      quantity: 1,
    })
  }

  function ubahJumlah(sku, quantity) {
    if (quantity < 1) return hapus(sku)
    const item = items.value.find((entry) => entry.sku === sku)
    if (item) item.quantity = quantity
  }

  function hapus(sku) { items.value = items.value.filter((item) => item.sku !== sku) }
  function kosongkan() { items.value = [] }

  return { items, jumlah, tambah, ubahJumlah, hapus, kosongkan }
})
