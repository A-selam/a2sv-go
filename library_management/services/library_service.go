package services

import (
	"fmt"
	"library_management/models"
)

type LibraryManager interface{
	// methods specified by the task doc
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book

	// methods added by me
	AddMember(member models.Member)
	ListAllBooks() []models.Book
}

type Library struct {
	Books map[int]models.Book
	Members map[int]models.Member
}

var bookIDCounter int = 0
var memberIDCounter int = 0
var currentLibrary Library

func Init(){
	currentLibrary = Library{
		Books: make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func GetLibrary() *Library{
	return &currentLibrary
}

// this function generates an ID for each new book entries
func GenerateBookID() int {
	bookIDCounter++
	return bookIDCounter
}

// this function generates an ID for each new member entries
func GenerateMemberID() int {
	memberIDCounter++
	return memberIDCounter
}

func (l *Library) AddBook(book models.Book){
	generatedID := GenerateBookID()
	book.ID = generatedID
	l.Books[generatedID] = book
	// fmt.Println(l.Books)
}

func (l *Library) RemoveBook(bookID int){
	delete(l.Books, bookID)
	for memberID, member := range l.Members{
		newBorrowed := []models.Book{}
		for _, b := range member.BorrowedBooks{
			if b.ID != bookID {
				newBorrowed = append(newBorrowed, b)
			}
        }
        member.BorrowedBooks = newBorrowed
        l.Members[memberID] = member
    }
}

func (l *Library) AddMember(member models.Member){
	generatedID := GenerateMemberID()
	member.ID = generatedID
	l.Members[generatedID] = member
	// fmt.Println(l.Members)
}

func (l *Library) BorrowBook(bookID, memberID int) error {
    book, bookExists := l.Books[bookID]
    if !bookExists {
        return fmt.Errorf("book with ID %d does not exist", bookID)
    }
    if book.Status == "Borrowed" {
        return fmt.Errorf("book with ID %d is already borrowed", bookID)
    }
    member, memberExists := l.Members[memberID]
    if !memberExists {
        return fmt.Errorf("member with ID %d does not exist", memberID)
    }

    // Update book status and member's borrowed books
    book.Status = "Borrowed"
    l.Books[bookID] = book
    member.BorrowedBooks = append(member.BorrowedBooks, book)
    l.Members[memberID] = member

    return nil
}

func (l *Library) ReturnBook(bookID, memberID int) error {
    book, bookExists := l.Books[bookID]
    if !bookExists {
        return fmt.Errorf("book with ID %d does not exist", bookID)
    }
    member, memberExists := l.Members[memberID]
    if !memberExists {
        return fmt.Errorf("member with ID %d does not exist", memberID)
    }

    // Update book status and member's borrowed books
    book.Status = "Available"
    l.Books[bookID] = book
    member.BorrowedBooks = append(member.BorrowedBooks, book)
    l.Members[memberID] = member

    return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	availableBooks := []models.Book{}

	for _, book := range l.Books{
		if book.Status == "Available"{
			availableBooks = append(availableBooks, book)
		}
	}

	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberId int) []models.Book {
	return l.Members[memberId].BorrowedBooks
}

func (l *Library) ListAllBooks() []models.Book {
	allBooks := []models.Book{}
	for _, book := range l.Books{
		allBooks = append(allBooks, book)
	} 

	return allBooks
}

func (l *Library) ListMembers() []models.Member {
	allMembers := []models.Member{}
	for _, member := range l.Members{
		allMembers = append(allMembers, member)
	} 

	return allMembers
}