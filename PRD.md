# PRD — go-osint-cli v2.0

## 1. Informasi Proyek

| Field | Nilai |
|---|---|
| Nama Proyek | go-osint-cli |
| Deskripsi | CLI untuk mengecek keberadaan username di 5 platform publik |
| Tech Stack | Go (Golang) 1.26+, standard library saja (tanpa dependensi) |
| Status | v1.1.0 — jalan, teruji |
| Penanggung Jawab | Nabil Aulia Dika (dezzorigin) |
| Versi Dokumen | 2.0 (18 Nov 2026 — ditulis ulang dari kode asli, PRD v1.0 adalah template generik dan tidak menggambarkan tool ini) |
| Live landing | https://go-osint-cli.pages.dev |

## 2. Latar Belakang

Saat menelusuri jejak digital seseorang (OSINT) atau memilih username baru, langkah pertama yang paling murah adalah mengecek apakah username itu sudah dipakai di platform-platform besar. Mengecek satu-satu lewat browser lambat dan berulang. Tool ini mengotomatisasi langkah itu: satu perintah, lima platform, hasilnya langsung terlihat.

## 3. Masalah yang Diselesaikan

- Mengecek username di banyak platform satu per satu itu lambat.
- Hasil cek manual sulit dibandingkan karena tidak ada format yang konsisten.
- Tool cek username yang ada sering mengklaim kepastian yang tidak bisa mereka buktikan (mis. menganggap “halaman merespons” = “akun ada”).

## 4. Tujuan & Sasaran

- Satu perintah mengecek 5 platform: GitHub, GitLab, Dev.to, Pinterest, Reddit.
- Hasil cepat (paralel, timeout 10 detik per platform) dan stabil (urutan output tetap).
- Hasil jujur: bedakan “ditemukan”, “tidak ditemukan”, dan “tidak bisa diverifikasi”.
- Tanpa dependensi, mudah di-build di Linux/macOS/Windows dengan Go 1.26+.

## 5. Ruang Lingkup

### In Scope
- Input username via argumen CLI (`go-osint-cli <username>`) atau prompt interaktif.
- Pengecekan 5 platform di atas via HTTP GET ke URL profil.
- Output berwarna per status + durasi total scan.
- Dokumentasi (README, PRD, arsitektur), CI, dan landing page.
- Unit test untuk pemetaan status → teks hasil.

### Out of Scope
- Login/autentikasi ke platform mana pun (tool ini anonim).
- Verifikasi isi halaman (mis. membedakan “akun ada” vs “halaman generik”) — di luar metode status-code.
- Platform di luar 5 yang terdaftar.
- Menyimpan/mengekspor hasil ke file atau database.
- GUI/web app — tool ini CLI; landing page hanya menjelaskan tool-nya.

## 6. Fitur & Persyaratan Fungsional

| # | Fitur | Kriteria |
|---|---|---|
| F1 | Input argumen | `go-osint-cli <username>` langsung memindai tanpa prompt |
| F2 | Input interaktif | Tanpa argumen, program meminta username lalu memindai |
| F3 | Validasi kosong | Username kosong → pesan error + exit code 1 |
| F4 | Cek 5 platform | GitHub, GitLab, Dev.to, Pinterest, Reddit — URL profil sesuai pola di `main.go` |
| F5 | Paralel | 5 permintaan berjalan bersamaan (goroutine), timeout 10 detik |
| F6 | Tanpa redirect | Redirect tidak diikuti; `3xx` dilaporkan apa adanya |
| F7 | User-Agent | Request memakai `User-Agent: go-osint-cli/<versi>` |
| F8 | Output berurutan | Hasil dicetak dalam urutan platform yang tetap, bukan urutan selesainya request |
| F9 | Warna status | Hijau = Found, merah = Not Found, kuning = Redirect |
| F10 | Durasi | Total waktu scan dicetak di akhir |

## 7. Persyaratan Non-Fungsional

- **Performa:** scan 5 platform selesai dalam hitungan detik (diuji ~2,8 dtk).
- **Keamanan:** tidak ada secret/credential; tidak ada file `.env`; tidak mengirim data selain request GET anonim.
- **Portabilitas:** satu binary Go, jalan di Linux/macOS/Windows.
- **Maintainability:** satu file `main.go` + `main_test.go`; fungsi pemetaan status murni dan teruji.
- **Kejujuran hasil:** batasan metode didokumentasikan di README (§ Batasan) dan landing page.

## 8. Arsitektur & Tech Stack

Go 1.26+, standard library saja (`net/http`, `sync`, `bufio`, `fmt`, `os`, `strings`, `time`). Tidak ada modul pihak ketiga (`go.mod` tanpa require).

Alur: input → 5 goroutine (satu HTTP GET per platform, tanpa follow redirect) → kumpulkan hasil per indeks → cetak berurutan → cetak durasi. Detail: `docs/architecture.md`, diagram: `docs/diagrams/flow.svg`.

## 9. Milestone / Roadmap

| Fase | Status |
|---|---|
| v1.0.0 — CLI awal (prompt interaktif, 5 platform, paralel) | ✅ |
| v1.1.0 — mode argumen, tanpa follow redirect, output berurutan, User-Agent, unit test | ✅ (18 Nov 2026) |
| Proper: README + PRD v2.0 + CI + docs + demo terminal + landing live | ✅ (18 Nov 2026) |
| Verifikasi isi halaman (bedakan akun asli vs halaman generik) | ⏳ ide — tidak dijanjikan |

## 10. Risiko & Mitigasi

| Risiko | Mitigasi |
|---|---|
| Platform mengubah perilaku status-code (mis. semua 200) | Dokumentasikan batasan; jangan klaim kepastian |
| Platform memblokir request anonim / rate limit | Timeout 10 dtk + pesan error per platform, bukan crash |
| Hasil disalahartikan sebagai bukti | Peringatan batasan di README dan landing |
| Go versi lama tidak bisa build | `go.mod` menetapkan `go 1.26.3`; CI memakai Go 1.27 |

## 11. Kriteria Keberhasilan

- `go build`, `go vet`, `go test`, `gofmt` semuanya bersih (lokal + CI).
- Dua skenario nyata teruji: username yang ada (`dezzorigin`) dan username yang jelas belum dipakai.
- README berisi screenshot + video dari sesi terminal asli.
- Landing live di https://go-osint-cli.pages.dev (HTTP 200).
- Tidak ada secret yang ter-commit.

## 12. Lampiran

- Kode: `main.go`, `main_test.go`
- Arsitektur: `docs/architecture.md`
- Demo: `docs/screenshots/`, `docs/demo/demo.mp4`
- Landing: `site/index.html`
- Lisensi: MIT (LICENSE)
