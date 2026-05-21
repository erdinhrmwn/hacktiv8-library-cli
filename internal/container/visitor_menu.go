package container

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/controller"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/utils"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
)

type VisitorMenu struct {
	authorController *controller.AuthorController
	bookController   *controller.BookController
	userController   *controller.UserController

	currentUser *model.User
}

func (m *VisitorMenu) Dashboard(ctx context.Context, user *model.User) {
	m.currentUser = user

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
			Size:         10,
			HideSelected: true,
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
			Size:         10,
			HideSelected: true,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Tampilkan Semua Buku":
			books, err := m.bookController.GetAllBooks(ctx)
			if err != nil {
				fmt.Printf("\n🚧 Tampilkan Semua Buku — failed: %v\n\n", err)
				return
			}

			t := tablewriter.NewWriter(os.Stdout)
			t.Header([]string{"ID", "ISBN", "Title", "Genre", "Stock"})
			for _, book := range books {
				t.Append([]any{book.ID, book.ISBN, book.Title, book.Genre, book.Stock})
			}
			t.Render()
		case "Cari Buku (Berdasarkan Judul)":
			query := promptui.Prompt{
				Label: "Masukkan kata kunci pencarian",
			}
			queryText, err := query.Run()
			if err != nil {
				return
			}
			books, err := m.bookController.SearchBook(ctx, queryText)
			if err != nil {
				fmt.Printf("\n🚧 Cari Buku — failed: %v\n\n", err)
				return
			}

			t := tablewriter.NewWriter(os.Stdout)
			t.Header([]string{"ID", "ISBN", "Title", "Genre", "Stock"})
			for _, book := range books {
				t.Append([]any{book.ID, book.ISBN, book.Title, book.Genre, book.Stock})
			}
			t.Render()
		case "Lihat Detail Buku & Penulis (Berdasarkan ID)":
			prompt := promptui.Prompt{
				Label: "Masukkan ID buku",
				Validate: func(s string) error {
					_, err := strconv.Atoi(s)
					if err != nil {
						return errors.New("Invalid number")
					}

					return nil
				},
			}
			bookPrompt, err := prompt.Run()
			if err != nil {
				return
			}

			bookID, err := strconv.Atoi(bookPrompt)
			if err != nil {
				fmt.Printf("\n🚧 Lihat Detail Buku — failed: %v\n\n", err)
				return
			}

			book, err := m.bookController.GetBookByID(ctx, bookID)
			if err != nil {
				fmt.Printf("\n🚧 Lihat Detail Buku — failed: %v\n\n", err)
				return
			}

			t := tablewriter.NewWriter(os.Stdout)
			t.Header([]string{"ID", "ISBN", "Title", "Genre", "Stock", "Author ID", "Author Name", "Author Nationality"})
			t.Append([]any{book.ID, book.ISBN, book.Title, book.Genre, book.Stock, book.Author.ID, book.Author.Name, book.Author.Nationality})
			t.Render()
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
			Size:         10,
			HideSelected: true,
		}
		_, sel, err := prompt.Run()
		if err != nil {
			return
		}

		switch sel {
		case "Ubah Nama":
			m.changeName(ctx)
		case "Ganti Password":
			m.changePassword(ctx)
		case "Kembali ke Dashboard":
			return
		}
	}
}

func (m *VisitorMenu) changeName(ctx context.Context) {
	newName, err := utils.AskInput("Nama Baru")
	if err != nil {
		return
	}

	err = m.userController.Update(ctx, controller.UpdateUserInput{
		UserID: m.currentUser.ID,
		Name:   newName,
		Email:  m.currentUser.Email,
		Role:   m.currentUser.Role,
	})
	if err != nil {
		fmt.Printf("\n❌ Gagal mengubah nama: %v\n\n", err)
		return
	}

	m.currentUser.Name = newName
	fmt.Printf("\n✅ Nama berhasil diubah menjadi %s\n\n", newName)
}

func (m *VisitorMenu) changePassword(ctx context.Context) {
	oldPassword, err := utils.AskInput("Password Lama")
	if err != nil {
		return
	}

	newPassword, err := utils.AskInput("Password Baru")
	if err != nil {
		return
	}

	err = m.userController.ChangePassword(ctx, controller.ChangePasswordInput{
		UserID:      m.currentUser.ID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})
	if err != nil {
		fmt.Printf("\n❌ Gagal mengganti password: %v\n\n", err)
		return
	}

	fmt.Printf("\n✅ Password berhasil diganti\n\n")
}
