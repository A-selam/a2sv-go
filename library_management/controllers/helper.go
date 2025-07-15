package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"os"
	"strconv"
	"strings"
)
var reader = bufio.NewReader(os.Stdin)

// Used by AddBook() controller to receive book details
func DetailsInputMethod(msg string) string{
	var input string
	for {
        fmt.Printf("%s", msg)
        input = readLine()
        if input == "b" {
            PrintBold("Returning to main menu...")
            return ""
        }
        if ValidateInput(input) {
            break
        }
        PrintError("Field cannot be empty. Please try again.")
    }
	return input
}

func BookIDInput(msg string, validBooks []models.Book) int {
	var id int
	for {
		fmt.Println(msg)
		input := readLine()
		if input == "b" {
			PrintBold("Returning to main menu...")
			return -1
		}
		var err error
		id, err = strconv.Atoi(input)
		if err != nil {
			PrintError("Invalid input. Please enter a valid integer ID.")
			continue
		}
		// Check if book ID exists in the valid list
		found := false
		for _, book := range validBooks {
			if book.ID == id {
				found = true
				break
			}
		}
		if !found {
			PrintError(fmt.Sprintf("Book with ID %d is not valid for this operation.", id))
			continue
		}
		break
	}
	return id
}

func MemberIDInput(prompt string, mustHaveBorrowed bool) int {
	var id int
	for {
		fmt.Print(prompt)
		input := readLine()
		if input == "b" {
			PrintBold("Returning to main menu...")
			return -1
		}
		var err error
		id, err = strconv.Atoi(input)
		if err != nil {
			PrintError("Invalid input. Please enter a valid integer ID.")
			continue
		}
		member, exists := libraryManager.Members[id]
		if !exists {
			PrintError(fmt.Sprintf("Member with ID %d does not exist.", id))
			continue
		}
		if mustHaveBorrowed && len(member.BorrowedBooks) == 0 {
			PrintError(fmt.Sprintf("Member with ID %d has not borrowed any books.", id))
			return -1
		}
		break
	}
	return id
}

func DisplayBooks() []models.Book{
	fmt.Println("Available books:")
    books := libraryManager.ListAllBooks()
    if len(books) == 0 {
        PrintError("No books available to remove.")
        return nil
    }
    for _, book := range libraryManager.Books {
        fmt.Printf("ID: %d | Title: %s | Author: %s | Status: %s\n", book.ID, book.Title, book.Author, book.Status)
    }
	return books
}

func DisplayAllMembers(){
	PrintUnderline("👥 Available Members")
	members := libraryManager.ListMembers()
    if len(members) == 0 {
        PrintError("No members available.")
        return
    }
    for _, member := range members {
        fmt.Printf("ID: %d | Name: %s\n", member.ID, member.Name)
    }
}

func DisplayAvailableBooks() []models.Book{
	availableBooks := libraryManager.ListAvailableBooks()
    if len(availableBooks) == 0 {
        PrintError("No books available to borrow.")
        return nil
    }
    PrintUnderline("Available books:")
    for _, book := range availableBooks {
        fmt.Printf("ID: %d | Title: %s | Author: %s\n", book.ID, book.Title, book.Author)
    }
	return availableBooks
}

func DisplayBorrowedBooks(memberID int) []models.Book{
	books := libraryManager.ListBorrowedBooks(memberID)
    if len(books) == 0 {
        PrintError("This member has not borrowed any books.")
        return nil
    }
    PrintUnderline(fmt.Sprintf("Books borrowed by %s:", libraryManager.Members[memberID].Name))
    for _, book := range books {
        fmt.Printf("ID: %d | Title: %s | Author: %s\n", book.ID, book.Title, book.Author)
    }
	return books
}

func readLine() string {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func ValidateName(input string) (bool, string) {
	name := strings.TrimSpace(input)
	if len(name) < 3 {
		return false, "Name must be at least 3 characters long."
	}
	for _, r := range name {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return false, "Name must contain only letters (no numbers or special characters)."
		}
	}
	return true, ""
}

func ValidateInput(input string) bool {
    return strings.TrimSpace(input) != ""
}

func PrintError(msg string) {
    fmt.Printf("\033[31m%s\033[0m\n", msg) // Red
}

func PrintSuccess(msg string) {
	fmt.Printf("\033[32m%s\033[0m\n", msg) // Green
}

func PrintBlue(msg string) {
    fmt.Printf("\033[36m%s\033[0m\n", msg) // Blue
}

func PrintBold(msg string) {
    fmt.Printf("\033[1m%s\033[0m\n", msg) // Blue
}

func PrintUnderline(msg string) {
	fmt.Printf("\033[4m%s\033[0m\n", msg) // Blue
}

func ClearScreen() {
    fmt.Print("\033[H\033[2J")
}