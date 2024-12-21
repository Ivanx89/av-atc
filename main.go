package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
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
			huh.NewSelect[string]().
				Options(huh.NewOptions("Land", "Take Off")...).
				Title("Choose your request").
				Description("At Charm we truly have a burger for everyone.").
				Value(&Comms.Request.Action),

			huh.NewInput().
				Value(&Comms.Request.Callsign).
				Title("Callsign").
				Placeholder("XYZ-1-F"),

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

	sendComms := func() {
		time.Sleep(2 * time.Second)
	}

	_ = spinner.New().Title("Request in progress...").Action(sendComms).Run()

	{
		var message string

		Callsign := Comms.Request.Callsign
		if Callsign != "" {
			Callsign = ", " + Callsign
		}

		Action := Comms.Request.Action
		if Action == "Take Off" {
			message = "You are clear to launch!\n\nThank you! Please visit again!"
		} else {
			message = "LANDING"
		}

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
