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

		// Reply with the action acknowledgment
		Action := Comms.Request.Action
		if Action == "Take Off" {
			message = "You are clear to launch!\n\nThank you! Please visit again!"
		} else {
			message = "Please proceed to hangar " + fmt.Sprint(hangar) + Callsign + "."
		}

		// Insert the submission into the SQLite database
		db, err := sql.Open("sqlite3", "./data/users.sqlite")
		if err != nil {
			fmt.Println("Failed to connect to the database:", err)
			os.Exit(1)
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO users (callsign, hangar) VALUES (?, ?)", Comms.Request.Callsign, hangar)
		if err != nil {
			fmt.Println("Failed to insert data:", err)
			os.Exit(1)
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
