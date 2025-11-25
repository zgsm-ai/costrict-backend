package main

import (
	"fmt"
	"time"
)

// Library Management System
// This Go program implements a library management system with books, members, and loans

// Book represents a library book
type Book struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	ISBN            string    `json:"isbn"`
	Publisher       string    `json:"publisher"`
	PublishYear     int       `json:"publish_year"`
	Genre           string    `json:"genre"`
	TotalCopies     int       `json:"total_copies"`
	AvailableCopies int       `json:"available_copies"`
	Location        string    `json:"location"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Member represents a library member
type Member struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	MemberSince time.Time `json:"member_since"`
	Type        string    `json:"type"` // Student, Adult, Senior
	MaxBooks    int       `json:"max_books"`
}

// Loan represents a book loan
type Loan struct {
	ID         string     `json:"id"`
	BookID     string     `json:"book_id"`
	MemberID   string     `json:"member_id"`
	LoanDate   time.Time  `json:"loan_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	Status     string     `json:"status"` // Active, Returned, Overdue
	Fine       float64    `json:"fine"`
}

// Library represents the library system
type Library struct {
	Books   map[string]*Book   `json:"books"`
	Members map[string]*Member `json:"members"`
	Loans   map[string]*Loan   `json:"loans"`
	nextID  int
}

// NewLibrary creates a new library instance
func NewLibrary() *Library {
	return &Library{
		Books:   make(map[string]*Book),
		Members: make(map[string]*Member),
		Loans:   make(map[string]*Loan),
		nextID:  1,
	}
}

// getNextID generates the next available ID
func (lib *Library) getNextID() string {
	id := lib.nextID
	lib.nextID++
	return fmt.Sprintf("LIB%d", id)
}

// AddBook adds a new book to the library
func (lib *Library) AddBook(title, author, isbn, publisher, genre, location string, publishYear, totalCopies int) *Book {
	id := lib.getNextID()
	book := &Book{
		ID:              id,
		Title:           title,
		Author:          author,
		ISBN:            isbn,
		Publisher:       publisher,
		PublishYear:     publishYear,
		Genre:           genre,
		TotalCopies:     totalCopies,
		AvailableCopies: totalCopies,
		Location:        location,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	lib.Books[id] = book
	return book
}

// GetBook retrieves a book by ID
func (lib *Library) GetBook(id string) *Book {
	return lib.Books[id]
}

// UpdateBook updates book information
func (lib *Library) UpdateBook(id, title, author, isbn, publisher, genre, location string, publishYear, totalCopies int) bool {
	if book, ok := lib.Books[id]; ok {
		book.Title = title
		book.Author = author
		book.ISBN = isbn
		book.Publisher = publisher
		book.PublishYear = publishYear
		book.Genre = genre
		book.Location = location
		book.TotalCopies = totalCopies
		book.UpdatedAt = time.Now()
		return true
	}
	return false
}

// DeleteBook removes a book from the library
func (lib *Library) DeleteBook(id string) bool {
	if _, ok := lib.Books[id]; ok {
		delete(lib.Books, id)
		return true
	}
	return false
}

// AddMember adds a new member to the library
func (lib *Library) AddMember(name, email, phone, address, memberType string) *Member {
	id := lib.getNextID()

	// Set max books based on member type
	maxBooks := 5 // Default for Adult
	if memberType == "Student" {
		maxBooks = 3
	} else if memberType == "Senior" {
		maxBooks = 10
	}

	member := &Member{
		ID:          id,
		Name:        name,
		Email:       email,
		Phone:       phone,
		Address:     address,
		MemberSince: time.Now(),
		Type:        memberType,
		MaxBooks:    maxBooks,
	}
	lib.Members[id] = member
	return member
}

// GetMember retrieves a member by ID
func (lib *Library) GetMember(id string) *Member {
	return lib.Members[id]
}

// UpdateMember updates member information
func (lib *Library) UpdateMember(id, name, email, phone, address, memberType string) bool {
	if member, ok := lib.Members[id]; ok {
		member.Name = name
		member.Email = email
		member.Phone = phone
		member.Address = address
		member.Type = memberType

		// Update max books based on member type
		if memberType == "Student" {
			member.MaxBooks = 3
		} else if memberType == "Senior" {
			member.MaxBooks = 10
		} else {
			member.MaxBooks = 5
		}

		return true
	}
	return false
}

// DeleteMember removes a member from the library
func (lib *Library) DeleteMember(id string) bool {
	if _, ok := lib.Members[id]; ok {
		delete(lib.Members, id)
		return true
	}
	return false
}

// LoanBook loans a book to a member
func (lib *Library) LoanBook(bookID, memberID string) *Loan {
	// Check if book exists and has available copies
	book, bookOk := lib.Books[bookID]
	if !bookOk || book.AvailableCopies <= 0 {
		return nil
	}

	// Check if member exists
	member, memberOk := lib.Members[memberID]
	if !memberOk {
		return nil
	}

	// Check if member has reached their loan limit
	activeLoans := lib.GetActiveLoansByMember(memberID)
	if len(activeLoans) >= member.MaxBooks {
		return nil
	}

	// Create loan
	id := lib.getNextID()
	loanDate := time.Now()
	dueDate := loanDate.AddDate(0, 0, 14) // 2 weeks loan period

	loan := &Loan{
		ID:       id,
		BookID:   bookID,
		MemberID: memberID,
		LoanDate: loanDate,
		DueDate:  dueDate,
		Status:   "Active",
		Fine:     0,
	}

	lib.Loans[id] = loan

	// Update book availability
	book.AvailableCopies--
	book.UpdatedAt = time.Now()

	return loan
}

// ReturnBook processes the return of a loaned book
