// Geometri rumah 3D (CSS 3D murni) dari sketsa "Landing3D". Setiap bagian = kumpulan sisi kotak (atas, depan,
// samping) yang diputar dengan transform. Warna lewat token --color-rumah-* (main.css); sisi depan/samping
// digelapkan dengan color-mix, jadi komponen tidak memuat kode warna langsung.

const W = 280 // lebar denah
const D = 200 // kedalaman denah
const H = 110 // tinggi dinding
const T = 10 // tebal dinding
const P = 'position:absolute;box-sizing:border-box;'

const w = (token) => `var(--color-${token})`
const gelap = (token, persen) => `color-mix(in srgb, var(--color-${token}) ${persen}%, black)`
const bening = (token, persen) => `color-mix(in srgb, var(--color-${token}) ${persen}%, transparent)`
const hitam = (persen) => `color-mix(in srgb, black ${persen}%, transparent)`

// Material: t = sisi atas, f = depan, l = samping kiri.
const M = {
  beton: { t: w('rumah-beton'), f: gelap('rumah-beton', 79), l: gelap('rumah-beton', 61) },
  pondasi: { t: w('rumah-pondasi'), f: gelap('rumah-pondasi', 77), l: gelap('rumah-pondasi', 58) },
  lantai: { t: w('rumah-lantai'), f: gelap('rumah-lantai', 86), l: gelap('rumah-lantai', 75) },
  teras: { t: w('rumah-teras'), f: gelap('rumah-teras', 86), l: gelap('rumah-teras', 74) },
  dindingLuar: { t: gelap('rumah-dinding', 89), f: w('rumah-dinding'), l: gelap('rumah-dinding', 91) },
  dindingDalam: { t: gelap('rumah-dinding-dalam', 85), f: w('rumah-dinding-dalam'), l: gelap('rumah-dinding-dalam', 97) },
  kolom: { t: w('rumah-kolom'), f: gelap('rumah-kolom', 83), l: gelap('rumah-kolom', 66) },
  kusen: { t: w('surface'), f: w('surface'), l: w('rumah-porselen') },
  kayu: { t: w('rumah-kayu'), f: gelap('rumah-kayu', 89), l: gelap('rumah-kayu', 66) },
  porselen: { t: w('surface'), f: w('rumah-porselen'), l: gelap('rumah-porselen', 90) },
  kabinet: { t: w('rumah-kabinet'), f: gelap('rumah-kabinet', 77), l: gelap('rumah-kabinet', 59) },
  pipa: { t: w('rumah-pipa'), f: gelap('rumah-pipa', 85), l: gelap('rumah-pipa', 70) },
  pipaPanas: { t: w('rumah-pipa-panas'), f: gelap('rumah-pipa-panas', 87), l: gelap('rumah-pipa-panas', 72) },
  pemanas: { t: w('surface'), f: w('rumah-porselen'), l: gelap('rumah-porselen', 68) },
  panel: { t: w('rumah-panel'), f: gelap('rumah-panel', 72), l: gelap('rumah-panel', 52) },
  lampu: { t: w('rumah-lampu'), f: gelap('rumah-lampu', 94), l: w('rumah-pipa-panas') },
  cerobong: { t: gelap('rumah-cerobong', 73), f: w('rumah-cerobong'), l: gelap('rumah-cerobong', 73) },
  bubungan: { t: w('rumah-bubungan'), f: gelap('rumah-bubungan', 74), l: gelap('rumah-bubungan', 55) },
  talang: { t: w('rumah-talang'), f: gelap('rumah-talang', 71), l: gelap('rumah-talang', 52) },
}

/** Tiga sisi tampak sebuah kotak. `o`: gaya tambahan (ls/fs/ts) dan isi (lk/fk/tk) per sisi. */
function kotak(daftar, x, y, z, lebar, dalam, tinggi, c, o = {}) {
  daftar.push({ gaya: `${P}left:${x - tinggi}px;top:${y}px;width:${tinggi}px;height:${dalam}px;background:${c.l};transform-origin:right center;transform:translateZ(${z}px) rotateY(90deg);${o.ls ?? ''}`, anak: o.lk ?? [] })
  daftar.push({ gaya: `${P}left:${x}px;top:${y + dalam}px;width:${lebar}px;height:${tinggi}px;background:${c.f};transform-origin:top center;transform:translateZ(${z}px) rotateX(90deg);${o.fs ?? ''}`, anak: o.fk ?? [] })
  daftar.push({ gaya: `${P}left:${x}px;top:${y}px;width:${lebar}px;height:${dalam}px;background:${c.t};transform:translateZ(${z + tinggi}px);${o.ts ?? ''}`, anak: o.tk ?? [] })
}

const persegi = (kiri, atas, lebar, tinggi, tambahan) => `${P}left:${kiri}px;top:${atas}px;width:${lebar}px;height:${tinggi}px;${tambahan}`

function pondasi() {
  const s = []
  for (const [x, y] of [[16, 14], [128, 14], [240, 14], [16, 162], [128, 162], [240, 162]]) kotak(s, x, y, -26, 26, 26, 26, M.pondasi)
  kotak(s, 0, 0, 0, W, D + 40, 16, M.beton, {
    ts: `background-image:linear-gradient(90deg,${hitam(7)} 1px,transparent 1px),linear-gradient(${hitam(7)} 1px,transparent 1px);background-size:40px 40px`,
  })
  return s
}

function lantai() {
  const s = []
  const nat = `background-image:linear-gradient(${w('rumah-nat')} 1px,transparent 1px),linear-gradient(90deg,${w('rumah-nat')} 1px,transparent 1px);background-size:24px 24px`
  kotak(s, 6, 6, 0, W - 12, D - 12, 6, M.lantai, { ts: nat })
  kotak(s, 90, D + 2, 0, 100, 34, 6, M.teras, { ts: nat })
  return s
}

function dinding() {
  const s = []
  const lubang = `background:${w('rumah-lubang')};`
  kotak(s, 0, 0, 0, W, T, H, M.dindingDalam)
  kotak(s, W - T, 0, 0, T, D, H, M.dindingDalam)
  kotak(s, 140, T, 0, T, 96, H, M.dindingDalam)
  kotak(s, 0, 0, 0, T, D, H, M.dindingLuar, { lk: [persegi(H - 92, 34, 44, 54, lubang), persegi(H - 92, 112, 44, 54, lubang)] })
  kotak(s, 0, D - T, 0, W, T, H, M.dindingLuar, { fk: [persegi(26, 44, 70, 46, lubang), persegi(120, 0, 42, 82, lubang), persegi(186, 44, 70, 46, lubang)] })
  for (const [x, y] of [[-4, -4], [W - 10, -4], [-4, D - 10], [W - 10, D - 10]]) kotak(s, x, y, 0, 14, 14, H + 4, M.kolom)
  return s
}

function bukaan() {
  const s = []
  const kaca = `background:${w('rumah-kaca')};border:5px solid ${w('surface')};box-shadow:inset 0 0 0 1px ${w('rumah-kaca-bingkai')};`
  const garisKayu = gelap('rumah-kayu', 66)
  const kusenTengah = [persegi(31, 0, 3, 36, `background:${w('surface')};`)]
  kotak(s, 120, D - 4, 0, 42, 6, 82, M.kayu, {
    fs: `border:3px solid ${garisKayu}`,
    fk: [persegi(8, 10, 26, 28, `border:2px solid ${garisKayu};`), persegi(8, 44, 26, 28, `border:2px solid ${garisKayu};`), persegi(30, 40, 6, 6, `background:${w('accent')};border-radius:3px;`)],
  })
  kotak(s, 26, D - 4, 44, 70, 6, 46, M.kusen, { fs: kaca, fk: kusenTengah })
  kotak(s, 186, D - 4, 44, 70, 6, 46, M.kusen, { fs: kaca, fk: kusenTengah })
  kotak(s, -2, 34, 48, 6, 54, 44, M.kusen, { ls: kaca })
  kotak(s, -2, 112, 48, 6, 54, 44, M.kusen, { ls: kaca })
  return s
}

function sanitasi() {
  const s = []
  kotak(s, 40, 96, 0, 200, 5, 5, M.pipa)
  kotak(s, 40, 72, 0, 5, 26, 5, M.pipa)
  kotak(s, 236, 96, 0, 5, 64, 5, M.pipaPanas)
  kotak(s, 24, 22, 0, 96, 50, 28, M.porselen, { ts: `background:${w('rumah-air')};box-shadow:inset 0 0 0 7px ${w('surface')}` })
  kotak(s, 154, 18, 0, 34, 12, 42, M.porselen)
  kotak(s, 158, 30, 0, 26, 30, 22, M.porselen, { ts: 'border-radius:40% 40% 50% 50%' })
  kotak(s, 206, 20, 0, 52, 30, 40, M.kabinet, {
    ts: `background:${w('surface')};box-shadow:inset 0 0 0 6px ${w('rumah-porselen')},inset 0 0 0 12px ${gelap('rumah-porselen', 90)}`,
  })
  kotak(s, 28, 140, 0, 26, 26, 62, M.pemanas)
  kotak(s, 230, 160, 0, 26, 8, 44, M.panel, { fk: [persegi(6, 26, 14, 12, `background:${w('rumah-indikator')};`)] })
  kotak(s, 128, 120, 92, 22, 22, 8, M.lampu, { ts: `box-shadow:0 0 0 8px ${bening('rumah-lampu', 35)}` })
  return s
}

const RH = 84 // tinggi atap

function atap() {
  const s = []
  const [ox, oy, rw, rd] = [-16, -16, W + 32, D + 32]
  const L = Math.sqrt((rd / 2) ** 2 + RH ** 2)
  const sudut = (Math.atan(RH / (rd / 2)) * 180) / Math.PI
  const genteng = (dasar) =>
    `background:${dasar};background-image:repeating-linear-gradient(180deg,transparent 0 12px,${bening('rumah-bubungan', 55)} 12px 14px),repeating-linear-gradient(90deg,transparent 0 20px,${hitam(14)} 20px 21px);`
  const tinggiGevel = (RH * D) / rd
  s.push({
    gaya: `${P}left:${-tinggiGevel}px;top:0;width:${tinggiGevel}px;height:${D}px;background:${w('rumah-gevel')};clip-path:polygon(0 50%,100% 0,100% 100%);transform-origin:right center;transform:translateZ(0px) rotateY(90deg)`,
    anak: [persegi(tinggiGevel * 0.35, D / 2 - 14, 22, 28, `background:${w('rumah-kaca')};border:3px solid ${w('surface')};`)],
  })
  s.push({ gaya: `${P}left:${ox}px;top:${oy + rd / 2 - L}px;width:${rw}px;height:${L}px;${genteng(w('rumah-genteng'))}transform-origin:bottom center;transform:translateZ(${RH}px) rotateX(${sudut}deg)`, anak: [] })
  kotak(s, 196, 14, 30, 24, 24, 74, M.cerobong)
  s.push({ gaya: `${P}left:${ox}px;top:${oy + rd / 2}px;width:${rw}px;height:${L}px;${genteng(w('primary'))}transform-origin:top center;transform:translateZ(${RH}px) rotateX(-${sudut}deg)`, anak: [] })
  kotak(s, ox, oy + rd / 2 - 4, RH - 2, rw, 8, 7, M.bubungan)
  kotak(s, ox, oy + rd - 2, -8, rw, 8, 7, M.talang)
  return s
}

export const UKURAN_DENAH = { lebar: W, dalam: D }

/**
 * Enam bagian, urut dari bawah ke atas. z0 = posisi saat dirakit, z1 = saat dipisah; dx/dy = geser saat dipisah;
 * h = tinggi bagian (letak lencana nomor); bx/by = letak lencana di denah.
 */
export function bangunBagianRumah() {
  return [
    { key: 'pondasi', sisi: pondasi(), z0: 0, z1: 0, dx: 0, dy: 0, h: 16, bx: -64, by: 100 },
    { key: 'lantai', sisi: lantai(), z0: 16, z1: 76, dx: 0, dy: 0, h: 6, bx: -64, by: 100 },
    { key: 'dinding', sisi: dinding(), z0: 22, z1: 146, dx: 0, dy: 0, h: H, bx: -64, by: 100 },
    { key: 'bukaan', sisi: bukaan(), z0: 22, z1: 146, dx: -56, dy: 70, h: 92, bx: 141, by: D + 30 },
    { key: 'sanitasi', sisi: sanitasi(), z0: 26, z1: 300, dx: 0, dy: 0, h: 62, bx: -64, by: 100 },
    { key: 'atap', sisi: atap(), z0: 132, z1: 430, dx: 0, dy: 0, h: RH, bx: -64, by: 100 },
  ]
}
