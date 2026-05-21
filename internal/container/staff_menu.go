package container

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
)

type StaffMenu struct{}

func (m *StaffMenu) Dashboard(ctx context.Context) {
	for {
		prompt := promptui.Select{
			Label: "STAFF DASHBOARD",
			Items: []string{
				"Kelola Anggota",
				"Kelola Katalog (Buku & Penulis)",
				"Proses Peminjaman (Borrow Book)",
				"Proses Pengembalian (Return Book & Auto-Invoice)",
				"Proses Pembayaran Denda (Payment)",
				"Pantau Log Aktivitas (Activity Logs)",
				"Logout",
			},
			Size:         10,
			HideSelected: true,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Kelola Anggota":
			m.Members(ctx)
		case "Kelola Katalog (Buku & Penulis)":
			m.Catalog(ctx)
		case "Proses Peminjaman (Borrow Book)":
			fmt.Printf("\n🚧 Proses Peminjaman — coming soon\n\n")
		case "Proses Pengembalian (Return Book & Auto-Invoice)":
			fmt.Printf("\n🚧 Proses Pengembalian — coming soon\n\n")
		case "Proses Pembayaran Denda (Payment)":
			fmt.Printf("\n🚧 Proses Pembayaran Denda — coming soon\n\n")
		case "Pantau Log Aktivitas (Activity Logs)":
			fmt.Printf("\n🚧 Pantau Log Aktivitas — coming soon\n\n")
		case "Logout":
			return
		}
	}
}

func (m *StaffMenu) Members(ctx context.Context) {
	for {
		prompt := promptui.Select{
			Label: "Kelola Anggota",
			Items: []string{
				"Daftarkan Visitor Baru",
				"Lihat Daftar Visitor",
				"Kembali ke Dashboard",
			},
			Size:         10,
			HideSelected: true,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Daftarkan Visitor Baru":
			fmt.Printf("\n🚧 Daftarkan Visitor Baru — coming soon\n\n")
		case "Lihat Daftar Visitor":
			fmt.Printf("\n🚧 Lihat Daftar Visitor — coming soon\n\n")
		case "Kembali ke Dashboard":
			return
		}
	}
}

func (m *StaffMenu) Catalog(ctx context.Context) {
	for {
		prompt := promptui.Select{
			Label: "Kelola Katalog (Buku & Penulis)",
			Items: []string{
				"Tambah Buku Baru",
				"Update Stok Buku",
				"Hapus Buku",
				"Tambah Penulis Baru",
				"Lihat Daftar Penulis",
				"Kembali ke Dashboard",
			},
			Size:         10,
			HideSelected: true,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Tambah Buku Baru":
			fmt.Printf("\n🚧 Tambah Buku Baru — coming soon\n\n")
		case "Update Stok Buku":
			fmt.Printf("\n🚧 Update Stok Buku — coming soon\n\n")
		case "Hapus Buku":
			fmt.Printf("\n🚧 Hapus Buku — coming soon\n\n")
		case "Tambah Penulis Baru":
			fmt.Printf("\n🚧 Tambah Penulis Baru — coming soon\n\n")
		case "Lihat Daftar Penulis":
			fmt.Printf("\n🚧 Lihat Daftar Penulis — coming soon\n\n")
		case "Kembali ke Dashboard":
			return
		}
	}
}
