# Arsitektur — go-osint-cli

## Ringkasan

CLI Go satu file (`main.go`, ~150 baris) tanpa dependensi di luar standard
library. Satu perintah mengecek satu username di lima platform publik
dengan lima HTTP GET paralel, lalu mencetak hasilnya dalam urutan yang
tetap.

```
pengguna ──► input (argumen / prompt) ──► 5× HTTP GET paralel
                                                    │
                                        ┌───────────┼───────────┐
                                        ▼           ▼           ▼
                                   GitHub  GitLab  Dev.to  Pinterest  Reddit
                                        │           │           │
                                        └───────────┼───────────┘
                                                    ▼
                                   status → teks hasil (formatResult)
                                                    │
                                                    ▼
                                   cetak berurutan + durasi total
```

Diagram visual: [`docs/diagrams/flow.svg`](diagrams/flow.svg).

## Lapisan

| Lapisan | Isi | File |
|---|---|---|
| Entry & input | Argumen CLI (`os.Args`), fallback prompt interaktif (`bufio`), validasi kosong → error + exit 1 | `main.go: main()` |
| Platform | Daftar 5 platform: nama + pola URL profil (`%s` = username) | `main.go: platforms` |
| Pemeriksaan | Satu goroutine per platform; `http.Client` timeout 10 dtk, **tanpa follow redirect** (`CheckRedirect` → `ErrUseLastResponse`), header `User-Agent: go-osint-cli/<versi>` | `main.go: check()` |
| Pemetaan hasil | Fungsi murni: status HTTP → satu baris teks (berwarna) | `main.go: formatResult()` |
| Output | Hasil dikumpulkan per indeks lalu dicetak berurutan + durasi total | `main.go: main()` |
| Test | Table test pemetaan status; validasi daftar platform | `main_test.go` |

## Keputusan desain

1. **Paralel, bukan berurutan.** Lima request jalan bersamaan (goroutine +
   `sync.WaitGroup`) supaya scan selesai dalam ~3 detik, bukan 5×
   latency.
2. **Hasil dikumpulkan, baru dicetak.** Output goroutine yang dicetak
   langsung akan acak urutannya. Hasil disimpan per indeks platform
   (`results[i]`) lalu dicetak berurutan — output stabil dan enak dibaca.
3. **Redirect tidak diikuti.** Platform yang mengarahkan (GitLab →
   `/users/sign_in`) akan tampak sebagai `302 — cannot verify`, bukan
   status membingungkan dari halaman tujuan redirect (dulu: `403`).
4. **User-Agent yang mengidentifikasi tool.** Request membawa
   `go-osint-cli/<versi> (+URL repo)` — sopan dan mudah dikenali di
   log server, bukan UA default Go yang sering diblokir.
5. **Pemetaan status dipisah jadi fungsi murni** (`formatResult`) supaya
   bisa di-unit-test tanpa jaringan.

## Metode & batasannya

Tool ini membaca **kode status HTTP saja**:

- `200` → Found — profil merespons di URL itu.
- `404` → Not Found (Available).
- `3xx` → Redirect — tidak bisa diverifikasi (biasanya login wall).
- Lainnya → Unexpected Response.

Batasan yang sudah diverifikasi dengan pengujian langsung (18 Nov 2026):

- **Reddit dan Pinterest mengembalikan `200` untuk username apa pun**,
  termasuk username yang jelas tidak ada — “Found” di sana hanya
  berarti halamannya merespons.
- **GitLab mengarahkan semua permintaan anonim ke `/users/sign_in`**
  (`302`), jadi selalu “cannot verify” tanpa login.

Batasan ini juga tertulis di README (§ Batasan) dan di landing page —
sengaja, supaya output tidak disalahartikan sebagai bukti.

## Build & test lokal

```bash
go test ./...   # unit test
go vet ./...    # static check
gofmt -l .      # harus kosong
go build -o go-osint-cli .
```

CI (`.github/workflows/ci.yml`) menjalankan keempat langkah itu di setiap
push dan pull request ke `main` (Go 1.27).

## File di repo

```
.
├── main.go                  # CLI: input, 5× GET paralel, output
├── main_test.go             # unit test formatResult + daftar platform
├── go.mod                   # module go-osint-cli, go 1.26.3, tanpa require
├── README.md                # cara pakai, contoh asli, batasan
├── PRD.md                   # PRD v2.0
├── CHANGELOG.md             # 1.0.0, 1.1.0
├── LICENSE                  # MIT
├── CODE_OF_CONDUCT.md
├── .gitignore
├── .github/workflows/ci.yml
├── docs/
│   ├── architecture.md      # file ini
│   ├── diagrams/flow.svg    # diagram alur
│   ├── screenshots/         # 3 tangkapan sesi terminal asli
│   └── demo/demo.mp4        # video demo 9 detik (tanpa suara)
└── site/
    └── index.html           # landing page (deploy: go-osint-cli.pages.dev)
```
