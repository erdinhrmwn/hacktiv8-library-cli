# 📚 Library CLI

> Aplikasi **Command Line Interface (CLI)** untuk mengelola operasional perpustakaan — mencakup katalog buku, peminjaman, denda keterlambatan, pembayaran, dan audit log — semuanya dari terminal.

---

## ✨ Fitur

- 🔐 **Multi-Role Authentication** — Login aman dengan hashing password. Peran `staff` (admin) dan `visitor` (pengunjung).
- 📦 **Manajemen Katalog** — Kelola data `books` dan `authors` beserta relasi antar keduanya.
- 🔄 **Peminjaman & Pengembalian** — Stok buku otomatis berkurang/tambah saat transaksi pinjam-kembali.
- 💰 **Denda Otomatis (Invoicing)** — Deteksi keterlambatan dan buat tagihan (`invoices`) secara otomatis untuk pengunjung.
- 💳 **Pembayaran** — Catat pembayaran denda lengkap dengan struk.
- 📊 **Audit Trail** — Setiap aktivitas penting (login, CRUD buku, dsb.) tercatat di tabel `activity_logs`.

---

## 🛠️ Teknologi

- **Bahasa** — [Go](https://golang.org) 1.24+
- **Database** — [MySQL](https://dev.mysql.com/downloads/)
- **CLI Interaktif** — [promptui](https://github.com/manifoldco/promptui)
- **Tabel Terminal** — [tablewriter](https://github.com/olekukonko/tablewriter)
- **Env Loader** — [godotenv](https://github.com/joho/godotenv)

---

## 🗄️ Skema Database

Aplikasi menggunakan **7 tabel** utama:

| # | Tabel | Deskripsi |
|---|---|---|
| 1 | `users` | Data pengunjung dan staf |
| 2 | `authors` | Data penulis buku |
| 3 | `books` | Katalog buku fisik |
| 4 | `loans` | Transaksi peminjaman buku |
| 5 | `invoices` | Tagihan / denda keterlambatan |
| 6 | `payments` | Transaksi pembayaran tagihan |
| 7 | `activity_logs` | Log riwayat aktivitas sistem |

> ℹ️ Rancangan lengkap ERD dan relasi antar tabel tersedia di [`REQUIREMENTS.md`](REQUIREMENTS.md).

---

## 🚀 Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/erdinhrmwn/hacktiv8-library-cli.git
cd hacktiv8-library-cli
```

### 2. Setup Database

Buat database baru di MySQL, lalu jalankan skrip DDL:

```bash
# Buat database (via MySQL CLI)
mysql -u root -p -e "CREATE DATABASE library_db;"

# Jalankan skrip schema
mysql -u root -p library_db < schema.sql
```

### 3. Konfigurasi Environment

```bash
cp .env.example .env
```

Sesuaikan kredensial database di dalam file `.env`:

```ini
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=library_db
```

### 4. Install Dependensi

```bash
go mod tidy
```

### 5. Jalankan Aplikasi

```bash
go run main.go
```

---

## 🎮 Penggunaan

Saat aplikasi dijalankan, kamu akan disambut menu **login**. Gunakan akun bawaan berikut (jika sudah melakukan *seeding*):

- **Staff** — `admin@library.local` / `admin123`

Navigasi dilakukan dengan **mengetik angka** pilihan menu, lalu tekan `Enter`.

---

## 📦 Dependensi

- [**promptui**](https://github.com/manifoldco/promptui) — Membangun UI interaktif: menu, form input, dan seleksi di terminal.
- [**tablewriter**](https://github.com/olekukonko/tablewriter) — Merender data dalam bentuk tabel yang rapi di terminal.
- [**godotenv**](https://github.com/joho/godotenv) — Memuat konfigurasi dari file `.env` ke dalam environment variabel.

---

## 📄 Lisensi

MIT © 2026 — Lihat [`LICENSE`](LICENSE) untuk detail selengkapnya.
