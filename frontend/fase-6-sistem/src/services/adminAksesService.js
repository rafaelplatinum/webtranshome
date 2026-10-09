import http from '@/services/http'

// Pengguna admin, role, dan izin (Fase 6). Hanya Super Admin (izin user.* dan role.*) menurut usulan seed.

/**
 * Daftar admin (akun dengan role admin; member tidak ikut). Parameter: q (nama/email), role, active (1/0), page.
 * Item: { id, full_name, email, is_active, roles: [{ code, name }], last_login_at, created_at }.
 * Hasil: { data, meta } dengan `meta.roles` = pilihan role admin [{ code, name, description }].
 */
export async function listAdmins(params = {}, { signal } = {}) {
  const { data } = await http.get('/admin/users', { params, signal })
  return data
}

/** { full_name, email, role, password } → admin baru. 409 email sudah dipakai. */
export async function createAdmin(isian) {
  const { data } = await http.post('/admin/users', isian)
  return data.data
}

/** Ganti role admin: { roles: ['ADMIN_KATALOG'] }. 422 bila mencabut Super Admin dari diri sendiri. */
export async function updateAdminRoles(id, roles) {
  const { data } = await http.put(`/admin/users/${id}/roles`, { roles })
  return data.data
}

/** Aktif/nonaktif. Admin nonaktif tidak bisa masuk dan sesinya berakhir. */
export async function setAdminActive(id, isActive) {
  const { data } = await http.patch(`/admin/users/${id}`, { is_active: isActive })
  return data.data
}

/** Kirim link reset password ke email admin itu. Hasil: pesan dari server. */
export async function resetAdminPassword(id) {
  const { data } = await http.post(`/admin/users/${id}/reset-password`)
  return data.message
}

/** Role admin: [{ id, code, name, description, is_locked, user_count, permissions: [kode] }]. */
export async function listRoles({ signal } = {}) {
  const { data } = await http.get('/admin/roles', { signal })
  return data.data
}

/** Semua izin per menu: [{ menu_id, menu, group, menu_permission, permissions: [{ id, code, action, description }] }]. */
export async function listPermissions({ signal } = {}) {
  const { data } = await http.get('/admin/permissions', { signal })
  return data.data
}

/** Simpan izin satu role (daftar kode lengkap). Berlaku setelah admin terkait memuat ulang halaman. */
export async function saveRolePermissions(roleId, permissions) {
  const { data } = await http.put(`/admin/roles/${roleId}/permissions`, { permissions })
  return data.data
}
