# 📋 Dokumen Kebutuhan Sistem

> Spesifikasi lengkap kebutuhan fungsional dan non-fungsional **Aplikasi Manajemen Perpustakaan (CLI)** — mencakup modul autentikasi, katalog, peminjaman, denda, pembayaran, dan audit log.

---

## 1. Deskripsi Umum

Aplikasi berbasis **Command Line Interface (CLI)** untuk mengelola operasional perpustakaan secara end-to-end:

- Pencatatan buku dan penulis beserta relasi keduanya
- Transaksi peminjaman dengan sistem tenggat waktu (*due date*)
- Deteksi keterlambatan dan pembuatan tagihan denda otomatis (*invoicing*)
- Pencatatan pembayaran denda
- Pelacakan aktivitas sistem (*audit trail*)

---

## 2. Aktor

**Staff** — Petugas perpustakaan yang mengelola data master (buku, penulis, anggota), memproses transaksi peminjaman dan pengembalian, serta menangani pembayaran denda.

**Visitor** — Anggota perpustakaan yang dapat mencari dan melihat katalog buku, memantau status peminjaman pribadi, dan mengecek tagihan denda.

---

## 3. Kebutuhan Fungsional

### 3.1. Modul Autentikasi & Akun

- **`F-AUTH-01`** — Login menggunakan `email` dan `password`.
- **`F-AUTH-02`** — Password disimpan dalam bentuk hash (*bcrypt*), tidak boleh *plaintext*.
- **`F-AUTH-03`** — Staff dapat mendaftarkan akun baru untuk Visitor.
- **`F-AUTH-04`** — Sistem mendeteksi `role` saat login untuk menentukan menu yang ditampilkan.

### 3.2. Modul Katalog (Buku & Penulis)

- **`F-CAT-01`** — Staff dapat **CRUD** data Penulis (`authors`).
- **`F-CAT-02`** — Staff dapat **CRUD** data Buku (`books`).
- **`F-CAT-03`** — Validasi: `isbn` harus unik, `author_id` harus merujuk ke penulis yang valid.
- **`F-CAT-04`** — Visitor dan Staff dapat mencari dan melihat daftar buku beserta `stock` yang tersedia.

### 3.3. Modul Peminjaman

- **`F-LOAN-01`** — Staff mencatat peminjaman baru dengan mengaitkan `visitor_id`, `staff_id`, dan `book_id`.
- **`F-LOAN-02`** — Sistem otomatis mengisi `borrow_date` (hari ini) dan `due_date` (`borrow_date` + 7 hari).
- **`F-LOAN-03`** — Sistem **menolak** peminjaman jika `stock` buku = 0.
- **`F-LOAN-04`** — Setelah peminjaman berhasil, stok buku otomatis berkurang (`stock - 1`).

### 3.4. Modul Pengembalian & Denda

- **`F-RET-01`** — Staff memproses pengembalian → sistem mengisi `return_date` dan mengubah `status` menjadi `returned`.
- **`F-RET-02`** — Stok buku otomatis bertambah (`stock + 1`) setelah pengembalian.
- **`F-RET-03`** — Jika `return_date` > `due_date`, sistem **otomatis membuat tagihan** di `invoices` atas `visitor_id` dengan status `unpaid`.

### 3.5. Modul Pembayaran

- **`F-PAY-01`** — Visitor dapat melihat daftar tagihan (`invoices`) miliknya yang berstatus `unpaid`.
- **`F-PAY-02`** — Staff memproses pembayaran: memilih `invoice_id`, memasukkan `amount_paid` dan `method` (`cash` / `transfer`).
- **`F-PAY-03`** — Setelah pembayaran tersimpan di `payments`, status tagihan di `invoices` otomatis berubah menjadi `paid`.

### 3.6. Modul Log Aktivitas

- **`F-LOG-01`** — Sistem otomatis mencatat aksi penting (*Login, Add Book, Borrow Book, Return Book, Pay Invoice*) ke tabel `activity_logs`.
- **`F-LOG-02`** — Setiap log memiliki `key` (jenis aksi), `description` (detail), dan `date` (timestamp).

---

## 4. Kebutuhan Non-Fungsional

1. **Antarmuka** — CLI interaktif. Data ditampilkan dalam bentuk tabel ASCII agar mudah dibaca (*library: `tablewriter`*).
2. **Keamanan** — Password tidak pernah dikembalikan dalam bentuk *plaintext* ke antarmuka.
3. **Integritas Data** — Database wajib mengaktifkan *Foreign Key constraint* (`ON DELETE RESTRICT` / `CASCADE` sesuai kebutuhan).

---

## 5. Struktur Database

Database terdiri dari **7 tabel** utama:

| # | Tabel | Deskripsi |
|---|---|---|
| 1 | `users` | Data pengunjung dan staf (autentikasi) |
| 2 | `authors` | Data penulis buku |
| 3 | `books` | Katalog buku fisik |
| 4 | `loans` | Transaksi peminjaman buku |
| 5 | `invoices` | Tagihan denda keterlambatan |
| 6 | `payments` | Transaksi pembayaran tagihan |
| 7 | `activity_logs` | Log riwayat aktivitas sistem |

> ℹ️ Detail lengkap kolom, tipe data, dan relasi antar tabel tersedia di file `schema.sql`.

---

## 6. Navigasi Menu CLI

### 6.1. Menu Awal (Unauthenticated)

```
┌─────────────────────────────┐
│        LIBRARY CLI          │
├─────────────────────────────┤
│ 1. Login                    │
│ 0. Keluar Aplikasi          │
└─────────────────────────────┘
```

### 6.2. Dashboard Visitor

```
┌──── VISITOR DASHBOARD ──────────────────────────────────────┐
│                                                             │
│  1. Cari & Lihat Katalog Buku                               │
│     1.1. Tampilkan Semua Buku                               │
│     1.2. Cari Buku (Berdasarkan Judul)                      │
│     1.3. Lihat Detail Buku & Penulis (Berdasarkan ID)       │
│     0.  Kembali ke Dashboard                                │
│                                                             │
│  2. Buku yang Sedang Dipinjam (My Loans)                    │
│  3. Cek Tagihan Denda (My Invoices)                         │
│                                                             │
│  4. Pengaturan Akun                                         │
│     4.1. Ubah Nama                                          │
│     4.2. Ganti Password                                     │
│     0.  Kembali ke Dashboard                                │
│                                                             │
│  0. Logout                                                  │
└─────────────────────────────────────────────────────────────┘
```

### 6.3. Dashboard Staff

```
┌──── STAFF DASHBOARD ────────────────────────────────────────┐
│                                                             │
│  1. Kelola Anggota                                          │
│     1.1. Daftarkan Visitor Baru                             │
│     1.2. Lihat Daftar Visitor                               │
│     0.  Kembali ke Dashboard                                │
│                                                             │
│  2. Kelola Katalog (Buku & Penulis)                         │
│     2.1. Tambah Buku Baru                                   │
│     2.2. Update Stok Buku                                   │
│     2.3. Hapus Buku                                         │
│     2.4. Tambah Penulis Baru                                │
│     2.5. Lihat Daftar Penulis                               │
│     0.  Kembali ke Dashboard                                │
│                                                             │
│  3. Proses Peminjaman (Borrow Book)                         │
│  4. Proses Pengembalian (Return Book & Auto-Invoice)        │
│  5. Proses Pembayaran Denda (Payment)                       │
│  6. Pantau Log Aktivitas (Activity Logs)                    │
│                                                             │
│  0. Logout                                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Aturan Bisnis

- **Durasi peminjaman** — 7 hari sejak `borrow_date`.
- **Denda keterlambatan** — Dibuat otomatis saat `return_date` > `due_date`.
- **Stok buku** — Hanya berkurang saat dipinjam, bertambah saat dikembalikan.
- **Peminjaman stok habis** — Ditolak oleh sistem.
- **Foreign Key** — `ON DELETE RESTRICT` untuk menjaga integritas referensial.

---

## 8. Dependensi Teknis

| Library | Kegunaan |
|---|---|
| [promptui](https://github.com/manifoldco/promptui) | UI interaktif CLI (menu, input, seleksi) |
| [tablewriter](https://github.com/olekukonko/tablewriter) | Rendering data dalam bentuk tabel ASCII di terminal |
| [godotenv](https://github.com/joho/godotenv) | Memuat konfigurasi environment dari file `.env` |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | Hashing password (*bcrypt*) |
| MySQL Driver | Koneksi dan query ke database MySQL |
