package controller

import (
	"fmt"
	"strconv"
	"strings"

	"os"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
)

type LoanController struct {
	LoanService *service.LoanService
	loanRepo    *repository.LoanRepository // untuk list tanpa membuat service baru
}

func NewLoanController(ls *service.LoanService, lr *repository.LoanRepository) *LoanController {
	return &LoanController{LoanService: ls, loanRepo: lr}
}

// BorrowBook — alur peminjaman buku oleh staff.
// Staff memasukkan visitor ID dan memilih buku dari daftar aktif.
func (c *LoanController) BorrowBook(staff *model.User, visitors []model.User, books []model.Book) {
	fmt.Println("\n📚 Proses Peminjaman Buku")
	fmt.Println(strings.Repeat("─", 40))

	// 1. Pilih visitor
	visitorLabels := make([]string, len(visitors))
	for i, v := range visitors {
		visitorLabels[i] = fmt.Sprintf("#%d — %s (%s)", v.ID, v.Name, v.Email)
	}
	promptVisitor := promptui.Select{
		Label: "Pilih Visitor",
		Items: visitorLabels,
	}
	idxVisitor, _, err := promptVisitor.Run()
	if err != nil {
		fmt.Println("❌ Dibatalkan.")
		return
	}
	selectedVisitor := visitors[idxVisitor]

	// 2. Pilih buku (hanya tampilkan buku dengan stok > 0)
	availableBooks := []model.Book{}
	for _, b := range books {
		if b.Stock > 0 {
			availableBooks = append(availableBooks, b)
		}
	}
	if len(availableBooks) == 0 {
		fmt.Println("⚠️  Tidak ada buku yang tersedia untuk dipinjam.")
		return
	}
	bookLabels := make([]string, len(availableBooks))
	for i, b := range availableBooks {
		bookLabels[i] = fmt.Sprintf("#%d — %s (stok: %d)", b.ID, b.Title, b.Stock)
	}
	promptBook := promptui.Select{
		Label: "Pilih Buku",
		Items: bookLabels,
	}
	idxBook, _, err := promptBook.Run()
	if err != nil {
		fmt.Println("❌ Dibatalkan.")
		return
	}
	selectedBook := availableBooks[idxBook]

	// 3. Konfirmasi
	promptConfirm := promptui.Prompt{
		Label:     fmt.Sprintf("Pinjamkan '%s' ke %s? (y/n)", selectedBook.Title, selectedVisitor.Name),
		IsConfirm: true,
	}
	if _, err := promptConfirm.Run(); err != nil {
		fmt.Println("❌ Dibatalkan.")
		return
	}

	// 4. Proses ke service
	if err := c.LoanService.Borrow(selectedVisitor.ID, staff.ID, selectedBook.ID); err != nil {
		fmt.Println("❌ Gagal:", err)
		return
	}

	fmt.Printf("✅ Peminjaman berhasil! Due date: 7 hari dari sekarang.\n")
}

// ReturnBook — alur pengembalian buku oleh staff.
func (c *LoanController) ReturnBook() {
	fmt.Println("\n🔄 Proses Pengembalian Buku")
	fmt.Println(strings.Repeat("─", 40))

	// 1. Ambil semua loan aktif
	loans, err := c.loanRepo.FindAllActive()
	if err != nil {
		fmt.Println("❌ Gagal memuat data peminjaman:", err)
		return
	}
	if len(loans) == 0 {
		fmt.Println("⚠️  Tidak ada peminjaman aktif saat ini.")
		return
	}

	// 2. Tampilkan tabel loan aktif
	renderLoanTable(loans)

	// 3. Input loan ID
	promptID := promptui.Prompt{
		Label: "Masukkan Loan ID yang akan dikembalikan",
		Validate: func(s string) error {
			if _, err := strconv.Atoi(s); err != nil {
				return fmt.Errorf("harus berupa angka")
			}
			return nil
		},
	}
	idStr, err := promptID.Run()
	if err != nil {
		fmt.Println("❌ Dibatalkan.")
		return
	}
	loanID, _ := strconv.Atoi(idStr)

	// 4. Proses ke service
	err = c.LoanService.Return(loanID)
	if err != nil {
		// Return terlambat — error berisi info denda, bukan error fatal
		if strings.Contains(err.Error(), "terlambat") {
			fmt.Println("⚠️ ", err)
			return
		}
		fmt.Println("❌ Gagal:", err)
		return
	}

	fmt.Println("✅ Pengembalian berhasil, stok buku telah bertambah.")
}

// MyLoans — daftar peminjaman aktif milik visitor yang sedang login.
func (c *LoanController) MyLoans(visitorID int) {
	fmt.Println("\n📋 Buku yang Sedang Dipinjam")
	fmt.Println(strings.Repeat("─", 40))

	loans, err := c.loanRepo.FindActiveByVisitorID(visitorID)
	if err != nil {
		fmt.Println("❌ Gagal memuat data:", err)
		return
	}
	if len(loans) == 0 {
		fmt.Println("Kamu tidak sedang meminjam buku apapun.")
		return
	}
	renderLoanTable(loans)
}

// renderLoanTable menampilkan daftar loan dalam format tabel ASCII.
func renderLoanTable(loans []model.Loan) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Visitor", "Judul Buku", "Pinjam", "Due Date", "Status"})
	table.SetBorder(true)
	table.SetRowLine(false)

	for _, l := range loans {
		visitorName := "-"
		if l.Visitor != nil {
			visitorName = l.Visitor.Name
		}
		bookTitle := "-"
		if l.Book != nil {
			bookTitle = l.Book.Title
		}
		table.Append([]string{
			strconv.Itoa(l.ID),
			visitorName,
			bookTitle,
			l.BorrowDate.Format("2006-01-02"),
			l.DueDate.Format("2006-01-02"),
			l.Status,
		})
	}
	table.Render()
}
