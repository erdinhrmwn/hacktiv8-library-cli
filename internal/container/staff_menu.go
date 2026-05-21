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

type StaffMenu struct {
	authorController  *controller.AuthorController
	bookController    *controller.BookController
	userController    *controller.UserController
	activityController *controller.ActivityController

	currentUser *model.User
}

func (m *StaffMenu) Dashboard(ctx context.Context, user *model.User) {
	m.currentUser = user
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
			m.showActivityLogs(ctx)
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
			m.registerVisitor(ctx)
		case "Lihat Daftar Visitor":
			m.listVisitors(ctx)
		case "Kembali ke Dashboard":
			return
		}
	}
}

func (m *StaffMenu) registerVisitor(ctx context.Context) {
	name, err := utils.AskInput("Nama")
	if err != nil {
		return
	}

	email, err := utils.AskInput("Email")
	if err != nil {
		return
	}

	password, err := utils.AskInput("Password")
	if err != nil {
		return
	}

	err = m.userController.Create(ctx, controller.CreateUserInput{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "visitor",
	})
	if err != nil {
		fmt.Printf("\n❌ Gagal mendaftarkan visitor: %v\n\n", err)
		return
	}

	m.activityController.Log(ctx, "Register Visitor", fmt.Sprintf("%s mendaftarkan visitor %s", m.currentUser.Name, name))
	fmt.Printf("\n✅ Visitor berhasil didaftarkan\n\n")
}

func (m *StaffMenu) listVisitors(ctx context.Context) {
	users, err := m.userController.GetByRole(ctx, "visitor")
	if err != nil {
		fmt.Printf("\n❌ Gagal mengambil daftar visitor: %v\n\n", err)
		return
	}

	if len(users) == 0 {
		fmt.Printf("\n📭 Belum ada visitor terdaftar\n\n")
		return
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Nama", "Email", "Role"})
	for _, u := range users {
		t.Append([]any{u.ID, u.Name, u.Email, u.Role})
	}
	t.Render()
	fmt.Println()
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
			m.addBook(ctx)
		case "Update Stok Buku":
			m.updateBookStock(ctx)
		case "Hapus Buku":
			m.deleteBook(ctx)
		case "Tambah Penulis Baru":
			m.addAuthor(ctx)
		case "Lihat Daftar Penulis":
			m.listAuthors(ctx)
		case "Kembali ke Dashboard":
			return
		}
	}
}

func (m *StaffMenu) addBook(ctx context.Context) {
	isbn, err := utils.AskInput("ISBN")
	if err != nil {
		return
	}

	title, err := utils.AskInput("Title")
	if err != nil {
		return
	}

	authorIDStr, err := utils.AskInput("Author ID")
	if err != nil {
		return
	}
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		fmt.Printf("\n❌ Author ID harus berupa angka\n\n")
		return
	}

	_, genre, err := utils.SelectInput("Genre", []string{
		"Horror", "Sci-Fi", "Comedy", "Romance", "Slice of Life",
	})
	if err != nil {
		return
	}

	stockStr, err := utils.AskInput("Stock")
	if err != nil {
		return
	}
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		fmt.Printf("\n❌ Stock harus berupa angka\n\n")
		return
	}

	err = m.bookController.AddBook(ctx, controller.AddBookInput{
		ISBN:     isbn,
		Title:    title,
		Genre:    genre,
		AuthorID: authorID,
		Stock:    stock,
	})
	if err != nil {
		fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
		return
	}
	m.activityController.Log(ctx, "Add Book", fmt.Sprintf("%s menambahkan buku %s", m.currentUser.Name, title))
	fmt.Printf("\n✅ Buku berhasil ditambahkan\n\n")
}

func (m *StaffMenu) updateBookStock(ctx context.Context) {
	books, err := m.bookController.GetAllBooks(ctx)
	if err != nil {
		fmt.Printf("\n❌ Gagal mengambil daftar buku: %v\n\n", err)
		return
	}

	titles := make([]string, len(books))
	for i, b := range books {
		titles[i] = b.Title
	}

	idx, _, err := utils.SelectInput("Pilih Buku", titles)
	if err != nil {
		return
	}
	selected := books[idx]

	stockStr, err := utils.AskInput("Masukkan Stok Baru")
	if err != nil {
		return
	}
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		fmt.Printf("\n❌ Stok harus berupa angka\n\n")
		return
	}

	err = m.bookController.UpdateBook(ctx, controller.UpdateBookInput{
		ID:       selected.ID,
		ISBN:     selected.ISBN,
		Title:    selected.Title,
		Genre:    selected.Genre,
		AuthorID: selected.AuthorID,
		Stock:    stock,
	})
	if err != nil {
		fmt.Printf("\n❌ Gagal mengubah stok: %v\n\n", err)
		return
	}
	m.activityController.Log(ctx, "Update Book", fmt.Sprintf("%s mengubah stok %s menjadi %d", m.currentUser.Name, selected.Title, stock))
	fmt.Printf("\n✅ Stok %s berhasil diubah menjadi %d\n\n", selected.Title, stock)
}

func (m *StaffMenu) deleteBook(ctx context.Context) {
	books, err := m.bookController.GetAllBooks(ctx)
	if err != nil {
		fmt.Printf("\n❌ Gagal mengambil daftar buku: %v\n\n", err)
		return
	}

	titles := make([]string, len(books))
	for i, b := range books {
		titles[i] = b.Title
	}

	idx, _, err := utils.SelectInput("Pilih Buku", titles)
	if err != nil {
		return
	}
	selected := books[idx]

	err = m.bookController.DeleteBook(ctx, selected.ID)
	if err != nil {
		fmt.Printf("\n❌ Gagal menghapus buku: %v\n\n", err)
		return
	}
	m.activityController.Log(ctx, "Delete Book", fmt.Sprintf("%s menghapus buku %s", m.currentUser.Name, selected.Title))
	fmt.Printf("\n✅ %s berhasil dihapus\n\n", selected.Title)
}

func (m *StaffMenu) addAuthor(ctx context.Context) {
	name, err := utils.AskInput("Name")
	if err != nil {
		return
	}

	birthDate, err := utils.AskInput("BirthDate")
	if err != nil {
		return
	}

	nationality, err := utils.AskInput("Nationality")
	if err != nil {
		return
	}

	bio, err := utils.AskInput("Bio")
	if err != nil {
		return
	}

	err = m.authorController.AddAuthor(ctx, controller.AddAuthorInput{
		Name:        name,
		BirthDate:   birthDate,
		Nationality: nationality,
		Bio:         bio,
	})
	if err != nil {
		fmt.Printf("\n❌ Gagal menambah penulis: %v\n\n", err)
		return
	}
	m.activityController.Log(ctx, "Add Author", fmt.Sprintf("%s menambahkan penulis %s", m.currentUser.Name, name))
	fmt.Printf("\n✅ Penulis berhasil ditambahkan\n\n")
}

func (m *StaffMenu) listAuthors(ctx context.Context) {
	authors, err := m.authorController.GetAllAuthors(ctx)
	if err != nil {
		fmt.Printf("\n❌ Gagal melihat daftar penulis: %v\n\n", err)
		return
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Name", "BirthDate", "Nationality"})
	for _, a := range authors {
		t.Append([]any{a.ID, a.Name, a.BirthDate, a.Nationality})
	}
	t.Render()
	fmt.Println()
}

func (m *StaffMenu) showActivityLogs(ctx context.Context) {
	logs, err := m.activityController.GetAll(ctx)
	if err != nil {
		fmt.Printf("\n❌ Gagal mengambil log aktivitas: %v\n\n", err)
		return
	}

	if len(logs) == 0 {
		fmt.Printf("\n📭 Belum ada aktivitas tercatat\n\n")
		return
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Key", "Description", "Date"})
	for _, l := range logs {
		t.Append([]any{l.ID, l.Key, l.Description, l.Date.Format("2006-01-02 15:04:05")})
	}
	t.Render()
	fmt.Println()
}
