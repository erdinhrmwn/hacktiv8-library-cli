package container

import (
	"context"
	"fmt"

	"github.com/manifoldco/promptui"
)

type VisitorMenu struct{}

func (m *VisitorMenu) Dashboard(ctx context.Context) {
	for {
		prompt := promptui.Select{
			Label: "VISITOR DASHBOARD",
			Items: []string{
				"Cari & Lihat Katalog Buku",
				"Buku yang Sedang Dipinjam (My Loans)",
				"Cek Tagihan Denda (My Invoices)",
				"Pengaturan Akun",
				"Logout",
			},
			Size: 10,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Cari & Lihat Katalog Buku":
			m.Catalog(ctx)
		case "Buku yang Sedang Dipinjam (My Loans)":
			fmt.Printf("\n🚧 My Loans — coming soon\n\n")
		case "Cek Tagihan Denda (My Invoices)":
			fmt.Printf("\n🚧 My Invoices — coming soon\n\n")
		case "Pengaturan Akun":
			m.Account(ctx)
		case "Logout":
			return
		}
	}
}

func (m *VisitorMenu) Catalog(ctx context.Context) {
	for {
		prompt := promptui.Select{
			Label: "Cari & Lihat Katalog Buku",
			Items: []string{
				"Tampilkan Semua Buku",
				"Cari Buku (Berdasarkan Judul)",
				"Lihat Detail Buku & Penulis (Berdasarkan ID)",
				"Kembali ke Dashboard",
			},
			Size: 10,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Tampilkan Semua Buku":
			fmt.Printf("\n🚧 Tampilkan Semua Buku — coming soon\n\n")
		case "Cari Buku (Berdasarkan Judul)":
			fmt.Printf("\n🚧 Cari Buku — coming soon\n\n")
		case "Lihat Detail Buku & Penulis (Berdasarkan ID)":
			fmt.Printf("\n🚧 Detail Buku — coming soon\n\n")
		case "Kembali ke Dashboard":
			return
		}
	}
}

func (m *VisitorMenu) Account(ctx context.Context) {
	for {
		prompt := promptui.Select{
			Label: "Pengaturan Akun",
			Items: []string{
				"Ubah Nama",
				"Ganti Password",
				"Kembali ke Dashboard",
			},
			Size: 10,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Ubah Nama":
			fmt.Printf("\n🚧 Ubah Nama — coming soon\n\n")
		case "Ganti Password":
			fmt.Printf("\n🚧 Ganti Password — coming soon\n\n")
		case "Kembali ke Dashboard":
			return
		}
	}
}
