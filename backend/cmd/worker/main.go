// Command worker menjalankan proses background yang tidak boleh
// memperlambat request HTTP, misalnya:
//   - sync kontak member ke Mekari Qontak (modules/integration/qontak)
//   - retry sync yang gagal (maksimal 3 kali, dicatat di sync_logs)
//   - notifikasi WhatsApp setelah poin bertambah
//
// Entry point API ada di backend/transhome.go (dibuat goctl, lihat Makefile).
package main

import "log"

func main() {
	log.Println("transhome worker: belum ada job yang didaftarkan")
}
