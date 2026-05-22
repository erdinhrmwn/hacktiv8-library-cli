# 📚 Library CLI

> Aplikasi **Command Line Interface (CLI)** untuk mengelola operasional perpustakaan — mencakup katalog buku, peminjaman, denda keterlambatan, pembayaran, dan audit log — semuanya dari terminal.

---

## ✨ Fitur

- 🔐 **Multi-Role Authentication** — Login aman dengan bcrypt hashing. Peran `staff` dan `visitor` dengan dashboard berbeda.
- 📦 **Manajemen Katalog** — CRUD `books` dan `authors`, validasi ISBN unique & stok.
- 🔄 **Peminjaman & Pengembalian** — Stok otomatis berkurang/tambah. Due date 7 hari.
- 💰 **Denda Otomatis** — Deteksi keterlambatan, auto-create `invoices` dengan denda Rp5.000/hari.
- 💳 **Pembayaran** — Bayar tagihan via `cash`/`transfer`, status otomatis `paid`.
- 📊 **Audit Trail** — Semua aksi penting tercatat di `activity_logs` via goroutine (non-blocking).

---

## 🏗️ Arsitektur

```
cmd/app/main.go
  └─ container (DI wiring)
       ├─ controller   ← input validation
       ├─ service      ← business logic
       └─ repository   ← database access
```

---

## 🛠️ Teknologi

- **Bahasa** — [Go](https://golang.org) 1.24+
- **Database** — [MySQL](https://dev.mysql.com/downloads/)
- **CLI** — [promptui](https://github.com/manifoldco/promptui) + [tablewriter](https://github.com/olekukonko/tablewriter)
- **Auth** — [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Config** — [godotenv](https://github.com/joho/godotenv)

---

## 🗄️ Database

7 tabel dengan foreign key constraint:

| # | Tabel | Deskripsi |
|---|---|---|
| 1 | `users` | Staff & visitor (autentikasi) |
| 2 | `authors` | Data penulis |
| 3 | `books` | Katalog buku |
| 4 | `loans` | Peminjaman |
| 5 | `invoices` | Denda keterlambatan |
| 6 | `payments` | Pembayaran invoice |
| 7 | `activity_logs` | Audit trail |

> ℹ️ Skema lengkap: [`sql/schema.sql`](sql/schema.sql) | Navigasi menu: [`REQUIREMENTS.md`](REQUIREMENTS.md)

---

## 🚀 Instalasi

### 1. Clone

```bash
git clone https://github.com/erdinhrmwn/hacktiv8-library-cli.git
cd hacktiv8-library-cli
```

### 2. Database

```bash
mysql -u root -p -e "CREATE DATABASE library_db;"
mysql -u root -p library_db < sql/schema.sql
```

### 3. Environment

```bash
cp .env.example .env
```

```ini
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=library_db
```

### 4. Dependensi

```bash
make tidy
```

### 5. Jalankan

```bash
make run
```

---

## 🎮 Penggunaan

Login dengan akun bawaan (dari seed data):

| Peran | Email | Password |
|---|---|---|
| Staff | `admin@library.com` | `password` |
| Visitor | `visitor@library.com` | `password` |

---

## 🧪 Testing

```bash
# Pastikan MySQL running
make test
```

---

## 📄 Lisensi

MIT © 2026 — Lihat [`LICENSE`](LICENSE).
