package container

import (
	"context"
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
	authorController   *controller.AuthorController
	bookController     *controller.BookController
	userController     *controller.UserController
	activityController *controller.ActivityController

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
			m.showAllBooks(ctx)
		case "Cari Buku (Berdasarkan Judul)":
			m.searchBooks(ctx)
		case "Lihat Detail Buku & Penulis (Berdasarkan ID)":
			m.showBookDetail(ctx)
		case "Kembali ke Dashboard":
			return
		}
	}
}

func (m *VisitorMenu) showAllBooks(ctx context.Context) {
	books, err := m.bookController.GetAllBooks(ctx)
	if err != nil {
		fmt.Printf("\n❌ Gagal menampilkan buku: %v\n\n", err)
		return
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "ISBN", "Title", "Genre", "Stock"})
	for _, b := range books {
		t.Append([]any{b.ID, b.ISBN, b.Title, b.Genre, b.Stock})
	}
	t.Render()
}

func (m *VisitorMenu) searchBooks(ctx context.Context) {
	query, err := utils.AskInput("Masukkan kata kunci pencarian")
	if err != nil {
		return
	}

	books, err := m.bookController.SearchBook(ctx, query)
	if err != nil {
		fmt.Printf("\n❌ Gagal mencari buku: %v\n\n", err)
		return
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "ISBN", "Title", "Genre", "Stock"})
	for _, b := range books {
		t.Append([]any{b.ID, b.ISBN, b.Title, b.Genre, b.Stock})
	}
	t.Render()
}

func (m *VisitorMenu) showBookDetail(ctx context.Context) {
	idStr, err := utils.AskInput("Masukkan ID buku")
	if err != nil {
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("\n❌ ID buku harus berupa angka\n\n")
		return
	}

	book, err := m.bookController.GetBookByID(ctx, id)
	if err != nil {
		fmt.Printf("\n❌ Gagal melihat detail buku: %v\n\n", err)
		return
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "ISBN", "Title", "Genre", "Stock", "Author Name", "Author Nationality"})
	t.Append([]any{book.ID, book.ISBN, book.Title, book.Genre, book.Stock, book.Author.Name, book.Author.Nationality})
	t.Render()
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
	go m.activityController.Log(context.Background(), "Update Profile", fmt.Sprintf("%s mengubah nama menjadi %s", m.currentUser.Email, newName))
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

	go m.activityController.Log(context.Background(), "Change Password", fmt.Sprintf("%s mengganti password", m.currentUser.Email))
	fmt.Printf("\n✅ Password berhasil diganti\n\n")
}
