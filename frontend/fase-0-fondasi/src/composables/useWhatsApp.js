/** Creates prefilled context for the shared WhatsApp contact form; never returns wa.me links. */
export function useWhatsApp() {
  const generalLink = () => ({ topic: 'Lainnya' })

  function productLink(product) {
    return {
      topic: 'Tanya produk',
      product: product ? `${product.sku || ''}${product.sku ? ' — ' : ''}${product.name || ''}` : '',
      message: product?.name ? `Saya ingin bertanya tentang ${product.name}.` : '',
    }
  }

  function bulkLink(product) {
    return {
      topic: 'Proyek/renovasi',
      product: product ? `${product.sku || ''}${product.sku ? ' — ' : ''}${product.name || ''}` : '',
      message: 'Saya ingin menanyakan kebutuhan produk dalam jumlah besar.',
    }
  }

  function searchLink(topic) {
    return topic
      ? { topic: 'Tanya produk', message: `Saya sedang mencari ${topic} dan ingin menanyakan ketersediaannya.` }
      : generalLink()
  }

  function kontenLink(article) {
    return { topic: 'Lainnya', message: article?.title ? `Saya ingin bertanya tentang: ${article.title}` : '' }
  }

  function jasaLink(service, product = null) {
    return {
      topic: 'Jasa',
      product: service?.name || (product ? `${product.sku || ''}${product.sku ? ' — ' : ''}${product.name || ''}` : ''),
      message: service?.name ? `Saya ingin konsultasi ${service.name}.` : 'Saya ingin bertanya tentang layanan Transhome.',
    }
  }

  return { productLink, bulkLink, generalLink, searchLink, kontenLink, jasaLink }
}
