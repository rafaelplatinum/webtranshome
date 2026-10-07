import { AxiosError } from 'axios'
import { cariRute } from './routes'

// Jeda buatan supaya skeleton dan state loading ikut teruji di mode mock.
const LATENSI_MS = 350

/**
 * Adapter axios untuk mode mock (VITE_USE_MOCK=true). Permintaan tidak keluar ke jaringan;
 * jawabannya dari services/mock/routes.js. Karena tetap lewat axios, interceptor di
 * http.js (token, 401, 403, toast) berjalan sama persis seperti dengan API asli.
 */
export async function mockAdapter(config) {
  await new Promise((selesai) => setTimeout(selesai, LATENSI_MS))

  const method = (config.method ?? 'get').toLowerCase()
  const url = new URL(config.url ?? '/', 'http://mock.local')
  const query = { ...Object.fromEntries(url.searchParams), ...(config.params ?? {}) }
  const body = typeof config.data === 'string' && config.data ? JSON.parse(config.data) : (config.data ?? {})

  const rute = cariRute(method, url.pathname)
  const hasil = rute
    ? await rute.handler({ params: rute.params, query, body, headers: config.headers })
    : { status: 404, data: { message: `Mock belum punya rute ${method.toUpperCase()} ${url.pathname}.` } }

  const status = hasil.status ?? 200
  const response = {
    data: hasil.data == null ? hasil.data : structuredClone(hasil.data),
    status,
    statusText: String(status),
    headers: {},
    config,
    request: {},
  }

  if (status >= 200 && status < 300) return response
  throw new AxiosError(
    response.data?.message ?? `Request failed with status code ${status}`,
    status >= 500 ? AxiosError.ERR_BAD_RESPONSE : AxiosError.ERR_BAD_REQUEST,
    config,
    response.request,
    response,
  )
}
