package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)
var reader = bufio.NewReader(os.Stdin)

var currentStudent Student
var subjectsMap = make(map[string]float64)
var average float64

type Student struct {
	name string
	numberOfSubjects int
}

func (s Student) String() string {
	return fmt.Sprintf("Name: %s, Number of subjects: %d", s.name, s.numberOfSubjects)
}

func getStudentInfo(){
	var name string
	var subjects int

	// console commands
	fmt.Print("Enter your name: ")
	name = readLine()

	fmt.Print("Enter the number of subjects you take: ")
	fmt.Scanln(&subjects)

	currentStudent = Student{name, subjects}
}

func getSubjectsAndGrades(){
	total_count := currentStudent.numberOfSubjects
	current_count := 0

	var title string
	var grade float64

	for {
		fmt.Printf("\033[32mSubject (%d/%d)\033[0m\n", current_count+1, total_count) // prints green colord text

		fmt.Print("Enter title of subject:")
		title = readLine()

		grade = readGrade()

		subjectsMap[title] = grade
		current_count ++
		fmt.Println("✅")

		if current_count == total_count{
			return
		}
	}
}

func validateGrade(grade float64) bool {
	return 0 <= grade && grade <= 100
}

func calculateAverage() float64{
	sum := 0.0
	for _, val := range subjectsMap{
		sum += val
	}
	
	average := sum / float64(len(subjectsMap))
	return average
}

func printReport() {
    fmt.Println("\033[36m--- Student Report ---\033[0m")
    fmt.Printf("Name: \033[33m%s\033[0m\n", currentStudent.name)
    fmt.Printf("Number of subjects: \033[33m%d\033[0m\n", currentStudent.numberOfSubjects)
    fmt.Println("\033[36mSubjects and Grades:\033[0m")
    for subject, grade := range subjectsMap {
        fmt.Printf("  \033[32m%s\033[0m: \033[35m%.2f\033[0m\n", subject, grade)
    }
    fmt.Printf("\033[36mAverage Grade: \033[34m%.2f\033[0m\n", calculateAverage())
}

func clearValues(){
	currentStudent = Student{}
	for k := range subjectsMap {
		delete(subjectsMap, k)
	}
	average = 0.0
}

func main() {
    for {
        fmt.Print("\033[H\033[2J") // clears console
        getStudentInfo()
        fmt.Print("\033[H\033[2J")
        getSubjectsAndGrades()
        fmt.Print("\033[H\033[2J")
        printReport()
        choice := showMenu()
        if choice == "Restart" {
            clearValues()
            continue
        } else {
            fmt.Println("\033[32mGoodbye!\033[0m")
            break
        }
    }
}

func readLine() string{
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func readGrade() float64 {
    for {
        fmt.Print("Enter grade: ")
        input := readLine()
        grade, err := strconv.ParseFloat(input, 64)
        if err != nil {
            fmt.Println("\033[31mInvalid input! Please enter a valid number.\033[0m")
            continue
        }
        if !validateGrade(grade) {
            fmt.Println("\033[31mGrade must be in the range 0-100!\033[0m")
            continue
        }
        return grade
    }
}

func showMenu() string {
    for {
        fmt.Println("\n\033[36mWhat would you like to do next?\033[0m")
        fmt.Println("1. Restart")
        fmt.Println("2. Close")
        fmt.Print("Enter your choice (1 or 2): ")
        choice := readLine()
        switch choice {
		case "1":
            return "Restart"
        case "2":
            return "Close"
        default:
            fmt.Println("\033[31mInvalid choice. Please enter 1 or 2.\033[0m")
        }
    }
}