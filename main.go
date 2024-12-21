package main

import (
	"fmt"
	"os"
	"strings"
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
		var s strings.Builder
		keyword := func(s string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
		}
		fmt.Fprintf(&s,
			"%s\n\nYou are clear to %s %s.",
			lipgloss.NewStyle().Bold(true).Render("REQUEST RECEIVED"),
			keyword(Comms.Request.Action),
			keyword(Comms.Request.Callsign),
		)

		Callsign := Comms.Request.Callsign
		if Callsign != "" {
			Callsign = ", " + Callsign
		}
		fmt.Fprintf(&s, "\n\nThank you! Please visit again!")

		fmt.Println(
			lipgloss.NewStyle().
				Width(40).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(s.String()),
		)
	}
}
