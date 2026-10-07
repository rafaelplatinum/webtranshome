import http from '@/services/http'

/** Menu sidebar yang boleh dilihat user (sudah difilter backend sesuai izin): [{ id, parent_id, name, code, route, icon, sort_order }]. */
export async function getMenus() {
  const { data } = await http.get('/admin/menus')
  return data.data
}

/** Ringkasan dashboard: { summary, attention, attention_total, logs, logs_scope }. */
export async function getDashboard() {
  const { data } = await http.get('/admin/dashboard')
  return data.data
}
