package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

var (
	action string
	ack    bool
)

type Request struct {
	Action   string
	Callsign string
}

type Comms struct {
	Request      Request
	Confirmation bool
}

func main() {

	var request Request
	var Comms = Comms{Request: request}

	// Create a new form.

	form := huh.NewForm(

		huh.NewGroup(
			// Choose your action
			huh.NewSelect[string]().
				Options(huh.NewOptions("Land", "Take Off")...).
				Title("Choose your request").
				Value(&Comms.Request.Action),
			// Enter your callsign
			huh.NewInput().
				Value(&Comms.Request.Callsign).
				Title("Callsign").
				Placeholder("XYZ-1-F"),
			// Confirm your request
			huh.NewConfirm().
				Title("Submit?").
				Value(&Comms.Confirmation).
				Affirmative("Yes!").
				Negative("No."),
		),
	)
	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	// Create loading spinner
	sendComms := func() {
		time.Sleep(2 * time.Second)
	}

	_ = spinner.New().Title("Request in progress...").Action(sendComms).Run()

	{
		var message string
		hangar := rand.Intn(30)

		// Check if callsign is empty
		Callsign := Comms.Request.Callsign
		if Callsign != "" {
			Callsign = ", " + Callsign
		}

		// Insert the submission into the SQLite database

		db, err := sql.Open("sqlite3", "./data/users.sqlite")
		if err != nil {
			fmt.Println("Failed to connect to the database:", err)
			os.Exit(1)
		}
		defer db.Close()

		sqlStmt := ` CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, callsign VARCHAR(30), hangar INT); `
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println("Backend failure.")
		}

		// Reply with the action acknowledgment
		Action := Comms.Request.Action
		if Action == "Take Off" {
			var exists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE callsign = ?)", Comms.Request.Callsign).Scan(&exists)
			if err != nil {
				fmt.Println("Failed to check for existing callsign:", err)
				os.Exit(1)
			}

			if exists {
				_, err = db.Exec("DELETE FROM users WHERE (callsign) = ?", Comms.Request.Callsign)
				message = "You are clear to launch!\n\nThank you! Please visit again!"
				if err != nil {
					fmt.Println("Failed to delete data:", err)
					os.Exit(1)
				}
			} else {
				message = "You are not cleared for takeoff. Please land first."
			}
		} else {
			// Check if the callsign already exists
			var exists bool
			err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE callsign = ?)", Comms.Request.Callsign).Scan(&exists)
			if err != nil {
				fmt.Println("Failed to check for existing callsign:", err)
				os.Exit(1)
			}

			if !exists {
				_, err = db.Exec("INSERT INTO users (callsign, hangar) VALUES (?, ?)", Comms.Request.Callsign, hangar)
				message = "Please proceed to hangar " + fmt.Sprint(hangar) + Callsign + "."
				if err != nil {
					fmt.Println("Failed to insert data:", err)
					os.Exit(1)
				}
			} else {
				fmt.Println("Your request is already granted.")
			}

		}

		// Print the message
		fmt.Println(
			lipgloss.NewStyle().
				Width(40).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(message),
		)
	}
}
