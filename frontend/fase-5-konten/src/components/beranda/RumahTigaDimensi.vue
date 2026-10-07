<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { BAGIAN_RUMAH } from './bagianRumah'
import { bangunBagianRumah, UKURAN_DENAH } from './geometriRumah'

// Rumah 3D CSS murni (HouseExploded di rencana build). Hanya visual: pilih bagian, rakit/pisah, dan putar
// dikendalikan tombol di HeroRumah. Transisi otomatis mati bila pengguna memilih "kurangi animasi" (main.css).
const props = defineProps({
  // true: bagian dipisah (exploded); false: dirakit jadi rumah.
  terpisah: { type: Boolean, default: true },
  // Sudut putar dalam derajat, dari slider "Putar rumah".
  sudut: { type: Number, default: -40 },
  // Indeks bagian yang disorot (0 = pondasi … 5 = atap), atau null.
  pilihan: { type: Number, default: null },
})

// Geometri dibangun sekali; yang berubah hanya transform per bagian.
const bagian = bangunBagianRumah()
const { lebar: W, dalam: D } = UKURAN_DENAH

// Ukuran sketsa: skala 0,84 di wadah setinggi 660 px. Layar sempit memperkecil adegan sesuai lebar wadah.
const wadah = ref(null)
const lebarWadah = ref(typeof window === 'undefined' ? 560 : Math.min(window.innerWidth - 32, 560))
let pengamat = null
onMounted(() => {
  pengamat = new ResizeObserver(([entri]) => (lebarWadah.value = entri.contentRect.width))
  pengamat.observe(wadah.value)
})
onBeforeUnmount(() => pengamat?.disconnect())

const skala = computed(() => Math.min(0.84, Math.max(0.48, lebarWadah.value / 560)))
const tinggi = computed(() => Math.round((660 * skala.value) / 0.84))

const gayaAdegan = computed(
  () =>
    `position:relative;width:${W}px;height:${D}px;transform-style:preserve-3d;` +
    `transform:translateY(${(190 * skala.value) / 0.84}px) scale(${skala.value}) rotateX(58deg) rotateZ(${props.sudut}deg)`,
)

const daftar = computed(() =>
  bagian.map((p, i) => {
    const dipilih = props.pilihan === i
    const redup = props.pilihan !== null && !dipilih
    const z = (props.terpisah ? p.z1 : p.z0) + (props.terpisah && dipilih ? 28 : 0)
    const [dx, dy] = props.terpisah ? [p.dx, p.dy] : [0, 0]
    const info = BAGIAN_RUMAH[i]
    return {
      key: p.key,
      nomor: info.nomor,
      sisi: p.sisi,
      redup,
      gaya: `position:absolute;left:0;top:0;width:${W}px;height:${D}px;transform-style:preserve-3d;transform:translate3d(${dx}px,${dy}px,${z}px);transition:transform .8s cubic-bezier(.2,.8,.2,1)`,
      // Lencana nomor selalu menghadap ke depan (rotasi adegan dibalik).
      lencana:
        `position:absolute;left:${p.bx}px;top:${p.by}px;width:30px;height:30px;margin:-15px 0 0 -15px;border-radius:15px;` +
        `border:2px solid var(--color-surface);box-sizing:border-box;display:${props.terpisah ? 'flex' : 'none'};align-items:center;justify-content:center;` +
        `font:800 13px var(--font-sans);background:var(--color-${info.titik});color:var(--color-${info.tinta});` +
        `transform:translateZ(${p.h + 14}px) rotateZ(${-props.sudut}deg) rotateX(-58deg);${redup ? 'opacity:.35' : ''}`,
    }
  }),
)
</script>

<template>
  <!-- Opacity dipasang per sisi (bukan di bagian) supaya efek 3D tidak menjadi datar. -->
  <div
    ref="wadah"
    class="relative flex w-full items-center justify-center overflow-hidden"
    :style="{ height: `${tinggi}px`, perspective: '1700px' }"
    aria-hidden="true"
    data-rumah-3d
  >
    <div :style="gayaAdegan" data-adegan-rumah>
      <div class="absolute rounded-full bg-rumah-bayangan" style="left: -70px; top: -60px; width: 420px; height: 340px; transform: translateZ(-40px)" />
      <div v-for="p in daftar" :key="p.key" :style="p.gaya" :data-bagian-rumah="p.key">
        <div v-for="(s, j) in p.sisi" :key="j" :style="p.redup ? `${s.gaya};opacity:.3` : s.gaya">
          <div v-for="(a, k) in s.anak" :key="k" :style="a" />
        </div>
        <div :style="p.lencana">{{ p.nomor }}</div>
      </div>
    </div>
  </div>
</template>
