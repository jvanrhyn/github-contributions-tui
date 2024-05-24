package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// model represents the state of the application.
type model struct {
	contributions [13][32]int     // Contributions data fetched from GitHub
	username      textinput.Model // GitHub username input by the user
	submittedName string          // GitHub username input by the user
	err           error           // Error encountered during data fetching
}

// fetchMsg is a message containing the fetched contributions data or an error.
type fetchMsg struct {
	contributions [13][32]int // Contributions data fetched from GitHub
	err           error       // Error encountered during data fetching
}

var githubToken string

// init loads the .env file and retrieves the GitHub token from environment variables.
func init() {
	// Load .env file from the same directory as the main program
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the GitHub token from environment variables
	githubToken = os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("No GitHub token provided. Set GITHUB_TOKEN in your .env file.")
	}
}

// main initializes and runs the Bubble Tea program.
func main() {
	p := bubbletea.NewProgram(initialModel(), bubbletea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
	}
	_ = result // Ignore the result if not needed
}

// initialModel returns the initial state of the model.
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter GitHub username"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		username: ti,
	}
}

// Init initializes the Bubble Tea program and enters the alternate screen mode.
func (m model) Init() bubbletea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model state accordingly.
func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	var cmd bubbletea.Cmd
	switch msg := msg.(type) {

	case bubbletea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, bubbletea.Quit

		case "enter":
			if m.username.Value() != "" {
				m.submittedName = m.username.Value()
				m.username.SetValue("")
				return m, fetchData(m.submittedName)
			}
		}

	case fetchMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.contributions = msg.contributions
	}

	m.username, cmd = m.username.Update(msg)
	return m, cmd
}

// View returns a string representation of the model's state.
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	// Define styles
	darkGrey := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	lightGrey := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	contributionColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#5AABE8"))

	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		"Github Contributuions\n\n%s\n\n%s",
		m.username.View(),
		"(ctrl+c to escape)\n\n",
	))
	b.WriteString("Contributions for : " + contributionColor.Render(fmt.Sprintf("%s\n\n", m.submittedName)) + "\n")

	for day := 0; day < 32; day++ {
		if day == 0 {
			b.WriteString(darkGrey.Render("[      ] "))
		} else {
			b.WriteString(lightGrey.Render(fmt.Sprintf("[%2d] ", day)))
		}
	}
	b.WriteString("\n")

	for month := 0; month < 13; month++ {
		for day := 0; day < 32; day++ {
			if day > 0 {
				if m.contributions[month][day] != 0 {
					b.WriteString(darkGrey.Render("[") + contributionColor.Render(
						fmt.Sprintf("%2d", m.contributions[month][day])) + darkGrey.Render("] "))
				} else {
					b.WriteString(darkGrey.Render("[ ✗] "))
				}
			} else {
				if m.contributions[month][day] != 0 {
					b.WriteString("[" + lightGrey.Render(
						fmt.Sprintf("%6d", m.contributions[month][day])) + "] ")
				} else {
					fmt.Fprintf(&b, "[      ] ")
				}
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// fetchData fetches the contributions data for the given GitHub username.
func fetchData(username string) bubbletea.Cmd {
	return func() bubbletea.Msg {

		currentDate := time.Now().Format("2006-01-02")
		fromDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		url := "https://api.github.com/graphql"
		query := fmt.Sprintf(`{ "query": "query { user(login: \"%s\") { contributionsCollection(from: \"%sT00:00:00Z\", to: \"%sT23:59:59Z\") { contributionCalendar { weeks { contributionDays { date, contributionCount } } } } } }" }`, username, fromDate, currentDate)
		payload := strings.NewReader(query)
		req, err := http.NewRequest("POST", url, payload)
		if err != nil {
			return fetchMsg{err: err}
		}

		req.Header.Add("Authorization", "bearer "+githubToken)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			return fetchMsg{err: err}
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fetchMsg{err: err}
		}

		var resp struct {
			Data struct {
				User struct {
					ContributionsCollection struct {
						ContributionCalendar struct {
							Weeks []struct {
								ContributionDays []struct {
									Date              string
									ContributionCount int
								}
							}
						}
					}
				}
			}
		}
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return fetchMsg{err: err}
		}

		var contributions [13][32]int
		startDate := time.Now().AddDate(-1, 0, 0)
		for _, week := range resp.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
			for _, day := range week.ContributionDays {
				date, err := time.Parse("2006-01-02", day.Date)
				if err != nil {
					return fetchMsg{err: err}
				}
				yearMonth, err := strconv.Atoi(date.Format("200601"))
				if err != nil {
					return fetchMsg{err: err}
				}
				monthIndex := (int(date.Year())-startDate.Year())*12 + int(date.Month()) - int(startDate.Month())
				if monthIndex < 0 || monthIndex >= 13 {
					continue // Skip out-of-bounds months
				}
				dayOfMonth := date.Day()
				contributions[monthIndex][dayOfMonth] = day.ContributionCount
				contributions[monthIndex][0] = yearMonth
			}
		}

		return fetchMsg{contributions: contributions}
	}
}
