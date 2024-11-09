package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Parse the response and extract roadmap details
func parseResponse(response string) (goals string, deadlines string, descriptions string) {
	goalRegex := regexp.MustCompile(`(?i)Goal\s*\d+:\s*(.+)`)
	deadlineRegex := regexp.MustCompile(`(?i)Deadline:\s*([^\n]+)`)

	var goalList, deadlineList, descriptionList []string

	// Extract all goals and deadlines
	goalMatches := goalRegex.FindAllStringSubmatch(response, -1)
	deadlineMatches := deadlineRegex.FindAllStringSubmatch(response, -1)

	for _, match := range goalMatches {
		if len(match) > 1 {
			goalList = append(goalList, strings.TrimSpace(match[1]))
		}
	}

	for _, match := range deadlineMatches {
		if len(match) > 1 {
			deadlineList = append(deadlineList, strings.TrimSpace(match[1]))
		}
	}

	// Extract descriptions after each deadline
	goalSections := strings.Split(response, "Goal ")
	for i := 1; i < len(goalSections); i++ {
		deadlineIndex := strings.Index(goalSections[i], "Deadline:")
		if deadlineIndex != -1 {
			newlineIndex := strings.Index(goalSections[i][deadlineIndex:], "\n")
			if newlineIndex != -1 {
				descriptionStart := deadlineIndex + newlineIndex + 1
				description := strings.TrimSpace(goalSections[i][descriptionStart:])
				descriptionList = append(descriptionList, description)
			}
		}
	}

	// Join extracted data into single strings
	goals = strings.Join(goalList, "; ")
	deadlines = strings.Join(deadlineList, "; ")
	descriptions = strings.Join(descriptionList, "; ")

	return goals, deadlines, descriptions
}

// Store roadmap data in the MySQL database
func storeInDatabase(userId int, roadmapId int, goals string, deadlines string, descriptions string) {
	dsn := "root:@tcp(127.0.0.1:3306)/roadmap_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	// Prepare the insert query
	query := `INSERT INTO roadmap (userId, roadmapId, goals, deadlines, description) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Error preparing query: %v", err)
	}
	defer stmt.Close()

	// Execute the query
	_, err = stmt.Exec(userId, roadmapId, goals, deadlines, descriptions)
	if err != nil {
		log.Fatalf("Error inserting data into database: %v", err)
	} else {
		fmt.Println("Data inserted successfully!")
	}
}

func main() {
	var userInput string
	fmt.Println("Enter your input: ")
	fmt.Scanln(&userInput)
	fmt.Println(userInput)
	response := GetGeminiResponse(userInput)
	cleanText := strings.ReplaceAll(response, "*", "")
	cleanText = strings.ReplaceAll(cleanText, "#", "")
	goals, deadlines, descriptions := parseResponse(cleanText)

	// Example userId and roadmapId
	userId := 1
	roadmapId := 101

	// Store data in the database
	storeInDatabase(userId, roadmapId, goals, deadlines, descriptions)
}
