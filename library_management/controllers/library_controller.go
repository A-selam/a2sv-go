package controllers

import (
	"fmt"
	"library_management/models"
	"library_management/services"
)
var libraryManager = services.GetLibrary()

func WelcomeText() {
	PrintBold("Welcome to the Library manager 🙌")
}

func App(){
	for {
        ClearScreen()
        WelcomeText()
        fmt.Println("1. Add Book")
        fmt.Println("2. Remove Book")
        fmt.Println("3. Add Member")
        fmt.Println("4. Borrow Book")
        fmt.Println("5. Return Book")
        fmt.Println("6. List Available Books")
        fmt.Println("7. List Borrowed Books")
        fmt.Println("q. Quit")
        fmt.Print("Select an option: ")
        choice := readLine()
        switch choice {
        case "1":
            AddBook()
        case "2":
            RemoveBook()
        case "3":
            AddMember()
        case "4":
            BorrowBook()
        case "5":
            ReturnBook()
        case "6":
            ListAvailableBooks()
        case "7":
            ListBorrowedBooks()
        case "q":
            return
        default:
            PrintError("Invalid option. Please try again.")
        }
        fmt.Println("Press Enter to continue...")
        readLine()
    }
}

func AddBook() {
    ClearScreen()
    PrintBlue("📚 Add a New Book")
    fmt.Println("Fill in the following fields (type 'b' to go back)")

    var Title string
    var Author string

    // Title input with validation and back option
	if Title = DetailsInputMethod("Enter the title of the book: "); Title == ""{
		return
	} 
	
	// Author input with validation and back option
	if Author = DetailsInputMethod("Enter the author of the book: "); Author == ""{
		return
	}

    book := models.Book{ID: 0, Title: Title, Author: Author, Status: "Available"}
    libraryManager.AddBook(book)
	PrintSuccess("✅ Book added successfully!")
}

func RemoveBook() {
    ClearScreen()
    PrintBlue("🗑️ Remove a Book")
    
    // List all books before prompting
    allBooks := DisplayBooks()

    var id int
    if id = BookIDInput("Enter the ID of the book you want to remove (type 'b' to go back): ", allBooks); id == -1{
		return 
	}

    libraryManager.RemoveBook(id)
    PrintSuccess(fmt.Sprintf("✅ Book with ID %d removed successfully!", id))
}

func AddMember() {
    ClearScreen()
    PrintBlue("👤 Add a New Member")
    fmt.Println("Fill in the following field (type 'b' to go back)")
	
	// additional validation for name of member cause ... why not
	var Name string
	for {
		Name = DetailsInputMethod("Enter your name: ")
		valid, msg := ValidateName(Name)
		if valid{
			break
		}
		PrintError(msg)
	}

    member := models.Member{ID: 0, Name: Name, BorrowedBooks: []models.Book{}}
    libraryManager.AddMember(member)
    PrintSuccess("✅ Member added successfully!")
}

func BorrowBook() {
	ClearScreen()
	PrintBlue("📖 Borrow a Book (1/2)")

	DisplayAllMembers()

	memberID := MemberIDInput("Enter the ID of the borrowing user (type 'b' to go back): ", false)
	if memberID == -1 {
		return
	}

	ClearScreen()
	PrintBlue("📖 Borrow a Book (2/2)")

	availableBooks := DisplayAvailableBooks() // shows and returns list
	if len(availableBooks) == 0 {
		PrintError("No books available for borrowing.")
		return
	}

	bookID := BookIDInput("Enter the ID of the book you want to borrow (type 'b' to go back):", availableBooks)
	if bookID == -1 {
		return
	}

	if err := libraryManager.BorrowBook(bookID, memberID); err != nil {
		PrintError(fmt.Sprintf("Operation was not successful: %s", err.Error()))
	} else {
		PrintSuccess("✅ Book borrowed successfully!")
	}
}

func ReturnBook() {
	ClearScreen()
	PrintBlue("📦 Return a Book (1/2)")

	DisplayAllMembers()

	memberID := MemberIDInput("Enter the ID of the user returning the book (type 'b' to go back): ", true)
	if memberID == -1 {
		return
	}

	ClearScreen()
	PrintBlue("📦 Return a Book (2/2)")
	PrintUnderline("📚 Borrowed Books by Member")

	borrowedBooks := DisplayBorrowedBooks(memberID)
	bookID := BookIDInput("Enter the ID of the book you want to return (type 'b' to go back):", borrowedBooks)
	if bookID == -1 {
		return
	}

	if err := libraryManager.ReturnBook(bookID, memberID); err != nil {
		PrintError(fmt.Sprintf("Operation was not successful: %s", err.Error()))
	} else {
		PrintSuccess("✅ Book returned successfully!")
	}
}

func ListAvailableBooks() {
    ClearScreen()
    PrintBlue("📚 List of Available Books")
    DisplayAvailableBooks()
}

func ListBorrowedBooks() {
	ClearScreen()
	PrintBlue("📚 List of Borrowed Books")
	DisplayAllMembers()

	memberID := MemberIDInput("Enter the ID of the user (type 'b' to go back): ", false)
	if memberID == -1 {
		return
	}

	DisplayBorrowedBooks(memberID)
}
