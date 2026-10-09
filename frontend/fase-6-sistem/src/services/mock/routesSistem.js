// Rute API tiruan sistem admin (Fase 6): pengguna admin, role & izin, log aktivitas.
// Bentuk mengikuti "Detail kontrak sistem" di docs/FRONTEND_BUILD_PLAN.md. TODO: ganti ke API.
import { activityLogs, menus } from '@/services/mock/data/admin'
import { rolePermissions, users } from '@/services/mock/data/akun'
import { permissions, roles } from '@/services/mock/data/sistem'
import { idBaru, simpanDb } from '@/services/mock/db'
import { catat, halaman } from '@/services/mock/routesAdmin'
import { hashPassword, kirimLinkReset } from '@/services/mock/routesAkun'
import { ROLE_ADMIN, cekAdmin, ditolak, namaPengguna, sekarang } from '@/services/mock/sesi'

const POLA_EMAIL = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const huruf = (teks) => String(teks ?? '').trim().toLowerCase()
const galatIsian = (errors, message = 'Periksa kembali isian yang ditandai.') => ({ status: 422, data: { message, errors } })
const tidakDitemukan = (pesan) => ({ status: 404, data: { message: pesan } })
const namaRole = (code) => roles.find((r) => r.code === code)?.name ?? code
const adalahAdmin = (u) => u.roles.some((r) => ROLE_ADMIN.includes(r))
const superAdminAktif = () => users.filter((u) => u.is_active && u.roles.includes('SUPER_ADMIN'))

// ---------- pengguna admin ----------

function bentukAdmin(u) {
  return {
    id: u.id, full_name: u.profile?.full_name ?? null, email: u.email, is_active: u.is_active,
    roles: u.roles.filter((r) => ROLE_ADMIN.includes(r)).map((code) => ({ code, name: namaRole(code) })),
    last_login_at: u.last_login_at ?? null, created_at: u.created_at ?? u.email_verified_at ?? null,
  }
}

function daftarAdmin({ headers, query }) {
  const { gagal } = cekAdmin(headers, 'user.view')
  if (gagal) return gagal
  const q = huruf(query.q)
  const hasil = users
    .filter(adalahAdmin)
    .filter((u) => !q || huruf(u.email).includes(q) || huruf(u.profile?.full_name).includes(q))
    .filter((u) => !query.role || u.roles.includes(query.role))
    .filter((u) => query.active == null || query.active === '' || u.is_active === (query.active === '1'))
    .sort((a, b) => Number(b.is_active) - Number(a.is_active) || (a.profile?.full_name ?? '').localeCompare(b.profile?.full_name ?? '', 'id'))
  const { data, meta } = halaman(hasil, query)
  // Pilihan role untuk filter dan form (tanpa perlu izin role.view).
  const pilihanRole = roles.filter((r) => ROLE_ADMIN.includes(r.code)).map(({ code, name, description }) => ({ code, name, description }))
  return { data: { data: data.map(bentukAdmin), meta: { ...meta, roles: pilihanRole } } }
}

async function tambahAdmin({ headers, body }) {
  const { u: pelaku, gagal } = cekAdmin(headers, 'user.create')
  if (gagal) return gagal
  const nama = String(body.full_name ?? '').trim()
  const email = huruf(body.email)
  const password = String(body.password ?? '')
  const errors = {}
  if (!nama) errors.full_name = ['Nama wajib diisi']
  if (!email) errors.email = ['Email wajib diisi']
  else if (!POLA_EMAIL.test(email)) errors.email = ['Format email belum benar']
  if (!ROLE_ADMIN.includes(body.role)) errors.role = ['Pilih role']
  if (password.length < 8) errors.password = ['Password sementara minimal 8 karakter']
  if (Object.keys(errors).length) return galatIsian(errors)
  if (users.some((x) => huruf(x.email) === email)) {
    return { status: 409, data: { message: 'Email sudah dipakai akun lain.', errors: { email: ['Email sudah dipakai akun lain'] } } }
  }
  const u = {
    id: idBaru(users), email, phone_number: null, is_active: true, roles: [body.role], password_hash: await hashPassword(password),
    email_verified_at: sekarang(), google_id: null, last_login_at: null, created_at: sekarang(), profile: { full_name: nama },
  }
  users.push(u)
  catat(pelaku, 'user.create', 'Pengguna admin', { sesudah: { full_name: nama, email, roles: [namaRole(body.role)] } })
  simpanDb()
  return { status: 201, data: { data: bentukAdmin(u) } }
}

function cariAdmin(params) {
  return users.find((x) => String(x.id) === String(params.id) && adalahAdmin(x))
}

function ubahRole({ headers, params, body }) {
  const { u: pelaku, gagal } = cekAdmin(headers, 'user.update')
  if (gagal) return gagal
  const u = cariAdmin(params)
  if (!u) return tidakDitemukan('Admin tidak ditemukan.')
  const baru = [...new Set(Array.isArray(body.roles) ? body.roles : [])]
  if (!baru.length || baru.some((r) => !ROLE_ADMIN.includes(r))) return galatIsian({ roles: ['Pilih role'] })
  const lepasSuper = u.roles.includes('SUPER_ADMIN') && !baru.includes('SUPER_ADMIN')
  if (lepasSuper && u.id === pelaku.id) return galatIsian({ roles: ['Kamu tidak bisa mencabut role Super Admin dari akunmu sendiri'] }, 'Kamu tidak bisa mencabut role Super Admin dari akunmu sendiri.')
  if (lepasSuper && u.is_active && superAdminAktif().length === 1) return galatIsian({ roles: ['Harus ada minimal satu Super Admin aktif'] }, 'Harus ada minimal satu Super Admin aktif.')
  const sebelum = { email: u.email, roles: u.roles.map(namaRole) }
  u.roles = [...u.roles.filter((r) => !ROLE_ADMIN.includes(r)), ...baru]
  catat(pelaku, 'user.update', 'Pengguna admin', { sebelum, sesudah: { email: u.email, roles: u.roles.map(namaRole) } })
  simpanDb()
  return { data: { data: bentukAdmin(u) } }
}

function ubahStatus({ headers, params, body }) {
  const { u: pelaku, gagal } = cekAdmin(headers, 'user.update')
  if (gagal) return gagal
  const u = cariAdmin(params)
  if (!u) return tidakDitemukan('Admin tidak ditemukan.')
  const aktif = Boolean(body.is_active)
  if (!aktif && u.id === pelaku.id) return galatIsian({ is_active: ['Kamu tidak bisa menonaktifkan akunmu sendiri'] }, 'Kamu tidak bisa menonaktifkan akunmu sendiri.')
  if (!aktif && u.roles.includes('SUPER_ADMIN') && u.is_active && superAdminAktif().length === 1) {
    return galatIsian({ is_active: ['Harus ada minimal satu Super Admin aktif'] }, 'Harus ada minimal satu Super Admin aktif.')
  }
  const sebelum = { email: u.email, is_active: u.is_active }
  u.is_active = aktif
  catat(pelaku, 'user.update', 'Pengguna admin', { sebelum, sesudah: { email: u.email, is_active: aktif } })
  simpanDb()
  return { data: { data: bentukAdmin(u) } }
}

function resetPasswordAdmin({ headers, params }) {
  const { u: pelaku, gagal } = cekAdmin(headers, 'user.update')
  if (gagal) return gagal
  const u = cariAdmin(params)
  if (!u) return tidakDitemukan('Admin tidak ditemukan.')
  if (!u.is_active) return galatIsian({}, 'Aktifkan akun ini dulu sebelum mengirim link reset password.')
  kirimLinkReset(u)
  catat(pelaku, 'user.reset_password', 'Pengguna admin', { sesudah: { email: u.email } })
  simpanDb()
  return { data: { message: `Link reset password dikirim ke ${u.email}.` } }
}

// ---------- role & izin ----------

const urutanIzin = (kode) => permissions.findIndex((p) => p.code === kode)

function bentukRole(r) {
  return {
    id: r.id, code: r.code, name: r.name, description: r.description, is_locked: r.code === 'SUPER_ADMIN',
    user_count: users.filter((u) => u.is_active && u.roles.includes(r.code)).length,
    permissions: [...(rolePermissions[r.code] ?? [])].sort((a, b) => urutanIzin(a) - urutanIzin(b)),
  }
}

function daftarRole({ headers }) {
  const { gagal } = cekAdmin(headers, 'role.view')
  if (gagal) return gagal
  return { data: { data: roles.filter((r) => ROLE_ADMIN.includes(r.code)).map(bentukRole) } }
}

// Izin dikelompokkan per menu (urut seperti sidebar). `menu_permission` = izin yang membuat menu terlihat.
function katalogIzin() {
  const induk = (m) => menus.find((x) => x.id === m.parent_id)
  return menus
    .filter((m) => permissions.some((p) => p.menu_id === m.id))
    .sort((a, b) => induk(a).sort_order - induk(b).sort_order || a.sort_order - b.sort_order)
    .map((m) => ({
      menu_id: m.id, menu: m.name, group: induk(m).name, menu_permission: m.izin,
      permissions: permissions.filter((p) => p.menu_id === m.id).map(({ id, code, action, description }) => ({ id, code, action, description })),
    }))
}

function daftarIzin({ headers }) {
  const { gagal } = cekAdmin(headers, 'role.view')
  if (gagal) return gagal
  return { data: { data: katalogIzin() } }
}

function simpanIzinRole({ headers, params, body }) {
  const { u: pelaku, gagal } = cekAdmin(headers, 'role.update')
  if (gagal) return gagal
  const r = roles.find((x) => String(x.id) === String(params.id) && ROLE_ADMIN.includes(x.code))
  if (!r) return tidakDitemukan('Role tidak ditemukan.')
  if (r.code === 'SUPER_ADMIN') return ditolak('Izin Super Admin tidak bisa diubah.')
  // Daftar wajib dikirim lengkap: body tanpa `permissions` tidak boleh menghapus semua izin.
  if (!Array.isArray(body.permissions)) return galatIsian({ permissions: ['Kirim daftar izin lengkap'] })
  const baru = [...new Set(body.permissions)]
  if (baru.some((k) => urutanIzin(k) < 0)) return galatIsian({ permissions: ['Ada kode izin yang tidak dikenal'] })
  // Izin lain di satu menu butuh izin yang membuat menunya terlihat.
  const kurang = katalogIzin().find((g) => g.permissions.some((p) => baru.includes(p.code)) && !baru.includes(g.menu_permission))
  if (kurang) return galatIsian({ permissions: [`Centang juga "${kurang.menu_permission}" supaya menu ${kurang.menu} terlihat`] })
  const sebelum = { role: r.name, permissions: bentukRole(r).permissions }
  rolePermissions[r.code] = baru.sort((a, b) => urutanIzin(a) - urutanIzin(b))
  catat(pelaku, 'role.update', 'Role dan izin', { sebelum, sesudah: { role: r.name, permissions: rolePermissions[r.code] } })
  simpanDb()
  return { data: { data: bentukRole(r) } }
}

// ---------- log aktivitas ----------

function daftarLog({ headers, query }) {
  const { gagal } = cekAdmin(headers, 'log.view')
  if (gagal) return gagal
  const hasil = activityLogs
    .filter((l) => !query.admin_id || String(l.admin_id) === String(query.admin_id))
    .filter((l) => !query.module || l.module === query.module)
    .filter((l) => !query.from || l.created_at.slice(0, 10) >= query.from)
    .filter((l) => !query.to || l.created_at.slice(0, 10) <= query.to)
    .sort((a, b) => b.created_at.localeCompare(a.created_at) || b.id - a.id)
  const { data, meta } = halaman(hasil, query)
  return { data: { data: data.map((l) => ({ ...l, admin_name: namaPengguna(l.admin_id) })), meta } }
}

function opsiLog({ headers }) {
  const { gagal } = cekAdmin(headers, 'log.view')
  if (gagal) return gagal
  const idAdmin = new Set([...users.filter(adalahAdmin).map((u) => u.id), ...activityLogs.map((l) => l.admin_id)])
  const admins = [...idAdmin].map((id) => ({ id, name: namaPengguna(id) })).sort((a, b) => a.name.localeCompare(b.name, 'id'))
  const modules = [...new Set(activityLogs.map((l) => l.module))].sort((a, b) => a.localeCompare(b, 'id'))
  return { data: { data: { admins, modules } } }
}

export const ruteSistem = [
  ['get', '/admin/users', daftarAdmin],
  ['post', '/admin/users', tambahAdmin],
  ['put', '/admin/users/:id/roles', ubahRole],
  ['post', '/admin/users/:id/reset-password', resetPasswordAdmin],
  ['patch', '/admin/users/:id', ubahStatus],
  ['get', '/admin/roles', daftarRole],
  ['get', '/admin/permissions', daftarIzin],
  ['put', '/admin/roles/:id/permissions', simpanIzinRole],
  ['get', '/admin/activity-logs/options', opsiLog],
  ['get', '/admin/activity-logs', daftarLog],
]
