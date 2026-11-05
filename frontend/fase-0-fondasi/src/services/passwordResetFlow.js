// Data reset hanya hidup di memory tab ini. Refresh atau tab baru akan memulai ulang alur.
let challenge = null
let resetToken = ''

export function mulaiPasswordReset(channel, destination) {
  challenge = { channel, destination }
  resetToken = ''
}

export function dataPasswordReset() {
  return challenge ? { ...challenge } : null
}

export function simpanTokenReset(token) {
  resetToken = String(token ?? '')
}

export function ambilTokenReset() {
  return resetToken
}

export function bersihkanPasswordReset() {
  challenge = null
  resetToken = ''
}
