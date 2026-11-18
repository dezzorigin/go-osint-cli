# Changelog

Semua perubahan penting pada proyek ini dicatat di file ini.

Format mengikuti [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
dan proyek ini mengikuti [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> Catatan jujur: CHANGELOG sebelum v1.1.0 adalah template generik yang
> mengklaim fitur yang tidak pernah ada (domain lookup, IP geolocation,
> report generation, config file). Entri-entri itu dihapus — yang tercatat
> di bawah ini hanya yang benar-benar ada di kode.

## [1.1.0] - 2026-11-18

### Added
- Mode argumen: `go-osint-cli <username>` langsung memindai tanpa prompt
- Unit test (`main_test.go`): pemetaan status HTTP → teks hasil, dan validasi daftar platform
- CI GitHub Actions: `gofmt` + `go vet` + `go test` + `go build`
- Dokumentasi lengkap: README baru, PRD v2.0, `docs/architecture.md`, diagram alur
- Screenshot + video demo dari sesi terminal asli
- Landing page live di https://go-osint-cli.pages.dev
- `.gitignore` standar Go

### Changed
- Output kini dicetak dalam urutan platform yang tetap (sebelumnya acak, mengikuti goroutine yang selesai duluan)
- Redirect tidak lagi diikuti: `3xx` dilaporkan sebagai “Redirect — cannot verify” (sebelumnya request mengikuti redirect lalu melaporkan status akhirnya, mis. GitLab selalu “Unexpected Response (Status: 403)”)
- Request memakai header `User-Agent: go-osint-cli/<versi>`
- Username kosong kini keluar dengan exit code 1 (sebelumnya 0)
- README dan PRD ditulis ulang agar sesuai dengan kode yang sebenarnya (sebelumnya template generik: sempat mengklaim `npm install` dan diagram database)

### Fixed
- GitLab tidak lagi dilaporkan sebagai error 403 yang membingungkan — perilakunya yang sebenarnya (redirect ke halaman login) kini terlihat

## [1.0.0] - 2026-11-05

### Added
- CLI Go awal: prompt interaktif, cek username di GitHub, GitLab, Dev.to, Pinterest, dan Reddit secara paralel (timeout 10 detik), output berwarna, durasi total scan
- LICENSE MIT
