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
	response := `## Roadmap name: Physics Odyssey

**Goal #1: Foundational Physics**

**Deadline:** 01.12 - 31.01

**Mini-goals:**

* **Master algebra and trigonometry:** Solid understanding of equations, manipulating expressions, and solving for unknowns. (01.12)
* **Mechanics 101:** Explore concepts of motion, forces, energy, and momentum. (15.12)
* **Electricity & Magnetism Basics:** Understand electric charge, current, fields, and basic magnetic interactions. (31.12)
* **Waves & Optics:** Learn about wave behavior, light, and simple optical phenomena. (15.01)
* **Thermodynamics & Heat:** Grasp the concepts of temperature, heat transfer, and the laws of thermodynamics. (31.01)

**Goal #2: Deep Dive into Core Physics**

**Deadline:** 01.02 - 31.03

**Mini-goals:**

* **Classical Mechanics:** Delve deeper into mechanics with advanced topics like Lagrangian and Hamiltonian mechanics. (15.02)
* **Electrodynamics:**  Uncover the secrets of electric and magnetic fields, including Maxwell's equations. (28.02)
* **Quantum Mechanics:**  Explore the strange world of the very small, with wave-particle duality, superposition, and uncertainty. (15.03)
* **Statistical Mechanics:**  Bridge the gap between microscopic and macroscopic systems with probability and statistics. (31.03)

**Goal #3:  Expanding your Horizons**

**Deadline:** 01.04 - Ongoing

**Mini-goals:**

* **Choose your focus:**  Pick a subfield of physics that excites you (e.g., astrophysics, particle physics, condensed matter physics). (15.04)
* **Advanced coursework/independent study:**  Dive deeper into your chosen field through specialized courses, online resources, or self-study. (Ongoing)
* **Explore the history of physics:**  Learn about the fascinating journey of physics from early discoveries to modern theories. (Ongoing)
* **Connect with the physics community:**  Engage with online forums, attend physics lectures, or join physics clubs to learn from others. (Ongoing)

**Note:** This roadmap is a suggestion, adjust the pace and complexity of the mini-goals to suit your learning style and goals.
`

	cleanText := strings.ReplaceAll(response, "*", "")
	cleanText = strings.ReplaceAll(cleanText, "#", "")
	goals, deadlines, descriptions := parseResponse(cleanText)

	// Example userId and roadmapId
	userId := 1
	roadmapId := 101

	// Store data in the database
	storeInDatabase(userId, roadmapId, goals, deadlines, descriptions)
}
