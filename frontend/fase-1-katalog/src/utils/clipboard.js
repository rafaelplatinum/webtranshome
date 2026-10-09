/**
 * Salin teks ke clipboard. Mengembalikan true bila berhasil.
 * Browser yang menolak Clipboard API memakai cara lama (textarea + execCommand).
 */
export async function salinTeks(teks) {
  try {
    await navigator.clipboard.writeText(teks)
    return true
  } catch {
    const area = document.createElement('textarea')
    area.value = teks
    area.setAttribute('readonly', '')
    area.style.position = 'fixed'
    area.style.opacity = '0'
    document.body.append(area)
    area.select()
    let berhasil = false
    try {
      berhasil = document.execCommand('copy')
    } catch {
      berhasil = false
    }
    area.remove()
    return berhasil
  }
}
