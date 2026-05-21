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
	authorController *controller.AuthorController
	bookController   *controller.BookController
	userController   *controller.UserController
}

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
			var bookInput model.Book

			// ISBN Input
			isbnPrompt, err := utils.AskInput("ISBN")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}
			bookInput.ISBN = isbnPrompt

			// Title Input
			titlePrompt, err := utils.AskInput("Title")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}
			bookInput.Title = titlePrompt

			// Author ID Input
			authorIDPrompt, err := utils.AskInput("Author ID")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}
			bookInput.AuthorID, err = strconv.Atoi(authorIDPrompt)
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}

			// Genre Input
			_, genrePrompt, err := utils.SelectInput("Genre", []string{
				"Horror",
				"Sci-Fi",
				"Comedy",
				"Romance",
				"Slice of Life",
			})
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}
			bookInput.Genre = genrePrompt

			// Stock Input
			stockPrompt, err := utils.AskInput("Stock")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}
			bookInput.Stock, err = strconv.Atoi(stockPrompt)
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}

			// Save Book
			err = m.bookController.AddBook(ctx, controller.AddBookInput{
				ISBN:     bookInput.ISBN,
				Title:    bookInput.Title,
				Genre:    bookInput.Genre,
				AuthorID: bookInput.AuthorID,
				Stock:    bookInput.Stock,
			})
			if err != nil {
				fmt.Printf("\n❌ Gagal menambahkan buku: %v\n\n", err)
				return
			}
			fmt.Printf("\n✅ Buku berhasil ditambahkan\n\n")
		case "Update Stok Buku":
			books, err := m.bookController.GetAllBooks(ctx)
			if err != nil {
				fmt.Printf("\n❌ Gagal mengambil daftar buku: %v\n\n", err)
				return
			}

			bookTitles := make([]string, 0, len(books))
			for _, book := range books {
				bookTitles = append(bookTitles, book.Title)
			}

			selectIndex, _, err := utils.SelectInput("Pilih Buku", bookTitles)
			if err != nil {
				fmt.Printf("\n❌ Gagal memilih buku: %v\n\n", err)
				return
			}
			selectedBook := books[selectIndex]

			fmt.Printf("\n✅ Buku yang dipilih: %s\n\n", selectedBook.Title)

			newStock, err := utils.AskInput("Stok Baru")
			if err != nil {
				fmt.Printf("\n❌ Gagal mengubah stok buku: %v\n\n", err)
				return
			}
			selectedBook.Stock, err = strconv.Atoi(newStock)
			if err != nil {
				fmt.Printf("\n❌ Gagal mengubah stok buku: %v\n\n", err)
				return
			}

			err = m.bookController.UpdateBook(ctx, controller.UpdateBookInput{
				ID:       selectedBook.ID,
				Stock:    selectedBook.Stock,
				ISBN:     selectedBook.ISBN,
				Title:    selectedBook.Title,
				Genre:    selectedBook.Genre,
				AuthorID: selectedBook.AuthorID,
			})
			if err != nil {
				fmt.Printf("\n❌ Gagal mengubah stok buku: %v\n\n", err)
				return
			}

			fmt.Printf("\n✅ Stok buku berhasil diubah: %s (stok: %d)\n\n", selectedBook.Title, selectedBook.Stock)
		case "Hapus Buku":
			books, err := m.bookController.GetAllBooks(ctx)
			if err != nil {
				fmt.Printf("\n❌ Gagal mengambil daftar buku: %v\n\n", err)
				return
			}

			bookTitles := make([]string, 0, len(books))
			for _, book := range books {
				bookTitles = append(bookTitles, book.Title)
			}

			selectIndex, _, err := utils.SelectInput("Pilih Buku", bookTitles)
			if err != nil {
				fmt.Printf("\n❌ Gagal memilih buku: %v\n\n", err)
				return
			}
			selectedBook := books[selectIndex]

			err = m.bookController.DeleteBook(ctx, selectedBook.ID)
			if err != nil {
				fmt.Printf("\n❌ Gagal menghapus buku: %v\n\n", err)
				return
			}

			fmt.Printf("\n✅ Buku berhasil dihapus: %s\n\n", selectedBook.Title)
		case "Tambah Penulis Baru":
			var author model.Author
			namePrompt, err := utils.AskInput("Name")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambah penulis: %v\n\n", err)
				return
			}
			author.Name = namePrompt

			birthDatePrompt, err := utils.AskInput("BirthDate")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambah penulis: %v\n\n", err)
				return
			}
			author.BirthDate = birthDatePrompt

			nationalityPrompt, err := utils.AskInput("Nationality")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambah penulis: %v\n\n", err)
				return
			}
			author.Nationality = nationalityPrompt

			bioPrompt, err := utils.AskInput("Bio")
			if err != nil {
				fmt.Printf("\n❌ Gagal menambah penulis: %v\n\n", err)
				return
			}
			author.Bio = bioPrompt

			err = m.authorController.AddAuthor(ctx, controller.AddAuthorInput{
				Name:        author.Name,
				BirthDate:   author.BirthDate,
				Nationality: author.Nationality,
				Bio:         author.Bio,
			})
			if err != nil {
				fmt.Printf("\n❌ Gagal menambah penulis: %v\n\n", err)
				return
			}
			fmt.Printf("\n✅ Penulis berhasil ditambahkan\n\n")
		case "Lihat Daftar Penulis":
			authors, err := m.authorController.GetAllAuthors(ctx)
			if err != nil {
				fmt.Printf("\n❌ Gagal melihat daftar penulis: %v\n\n", err)
				return
			}

			t := tablewriter.NewWriter(os.Stdout)
			t.Header([]string{"ID", "Name", "BirthDate", "Nationality"})
			for _, author := range authors {
				t.Append([]any{author.ID, author.Name, author.BirthDate, author.Nationality})
			}
			t.Render()
		case "Kembali ke Dashboard":
			return
		}
	}
}
