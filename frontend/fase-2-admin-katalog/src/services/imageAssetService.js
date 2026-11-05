import http from '@/services/http'

const apiOrigin = String(import.meta.env.VITE_API_BASE_URL ?? '').replace(/\/api\/v1\/?$/, '')
const resolveMediaUrl = (item) => ({
  ...item,
  image_url: item.image_url?.startsWith('/media/') && apiOrigin ? `${apiOrigin}${item.image_url}` : item.image_url,
})

export async function listImageAssets({ signal } = {}) {
  const { data } = await http.get('/admin/image-assets', { signal })
  return data.data.map(resolveMediaUrl)
}

export async function createImageAsset({ file, slotKey, altText }) {
  const form = new FormData()
  form.append('file', file)
  form.append('slot_key', slotKey)
  form.append('alt_text', altText)
  const { data } = await http.post('/admin/image-assets', form, { tanpaToast: true })
  return resolveMediaUrl(data.data)
}

export async function updateImageAsset(id, changes) {
  const { data } = await http.patch(`/admin/image-assets/${id}`, changes)
  return resolveMediaUrl(data.data)
}

export async function deleteImageAsset(id) {
  await http.delete(`/admin/image-assets/${id}`)
}
