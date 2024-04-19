package main

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Create","Read Data","Update data", "Delete Data"}

type model struct {
	cursor int
	choice string
}

type day struct{
	Date string `csv:"Date"`
	Studying int `csv:"Studying"`
	WorkOut int `csv:"WorkOut"`
	Sleep int `csv:"Sleep"`
	Hobbies int `csv:"Hobbies"`
	Loss int `csv:"Loss"`
	Eating string `csv:"Eating"`
	

}
func createCSV() error {
    // Data to write into the CSV file	
	currentTime := time.Now()

    data := [][]string{
		{"Date", "Studying", "WorkOut","Sleep","Hobbies","Loss","Eating"},
        {currentTime.Format("01-02-2006"), "8", "1","6","1","1","Good"},
    }

    // Create a new CSV file
    file, err := os.Create("data.csv")
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, row := range data {
        err := writer.Write(row)
        if err != nil {
            return err
        }
    }

    return nil
}


func updateCSV(){
	file, err := os.OpenFile("data.csv", os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        panic( err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()


	rows := [][]string{
    {"Number", "Number", "Number","Number","Number","Text"},
}
    t := table.New().Headers("Studying", "Work out", "Sleep","Hobbies","Loss","Eating").
    Rows(rows...)

	var style=lipgloss.NewStyle().Bold(true).BorderStyle(lipgloss.RoundedBorder())
	var data string 
	currentTime := time.Now()
	fmt.Println(style.Render("Give me the Data with commas between them using this kind of data: "))
	fmt.Println(t)

	fmt.Scan(&data)
	fmt.Println(data)

	values := strings.Split(data, ",")

	fmt.Println(values)
	row:=[]string{ currentTime.Format("01-02-2006"),values[0],values[1],values[2],values[3],values[4],values[5]}
	err = writer.Write(row)
    if err != nil {
        panic(err)
	}

}

func newData() {

	 _, err := os.Stat("data.csv")
    if os.IsNotExist(err) {
        err := createCSV()
        if err != nil {
            fmt.Println("Error creating CSV:", err)
            return
        }
        fmt.Println("CSV file created with data.")
    } else if err == nil {
        fmt.Println("CSV file already exists.")
    } else {
        fmt.Println("Error checking CSV file:", err)
    }
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	var style = lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    Foreground(lipgloss.Color("#FCAE1E")).MarginTop(2).MarginLeft(2).Blink(true)
	
	var style2= lipgloss.NewStyle().Foreground(lipgloss.Color("#FCAE1E")).Blink(true)

	s.WriteString(style.Render("What do you want to do?\n\n"))
	s.WriteString("\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString(style2.Render("	(^_^) "))
		} else {
			s.WriteString("	( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func main() {
	p := tea.NewProgram(model{})
	var style = lipgloss.NewStyle().
    Bold(true).
    Background(lipgloss.Color("#FDA172")).
    MarginLeft(10).
    BorderStyle(lipgloss.DoubleBorder()).BorderBottom(true).Italic(true)
	
	var menu = lipgloss.NewStyle().MarginLeft(2)

	fmt.Println(style.Render("Gorest"))
	fmt.Println(menu.Render("Lets see how was your week!"))
	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok && m.choice != "" {
		fmt.Println("You choose ",m.choice)
		if(m.choice=="Update data"){
			updateCSV()
		} else if(m.choice=="Create"){
			createCSV()
		}

	}
}
