package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

var (
	styleHelp = lipgloss.NewStyle().Foreground(lipgloss.Color("28")).Align(lipgloss.Center) // green help text
	styleBg   = lipgloss.NewStyle().Background(lipgloss.Color("0")) // black background
)

type tickMsg time.Time

type model struct {
	running   bool
	start     time.Time
	elapsed   time.Duration
	quitting  bool
	width     int
	height    int
	solves    []float64 // store past solve times in seconds
}

func initialModel() model {
	return model{width: 80, height: 24}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case " ":
			if !m.running {
				m.running = true
				m.start = time.Now()
				m.elapsed = 0
				return m, tick()
			} else {
				m.running = false
				m.elapsed = time.Since(m.start)
				m.solves = append([]float64{m.elapsed.Seconds()}, m.solves...)
			}
		}
	case tickMsg:
		if m.running {
			m.elapsed = time.Since(m.start)
			return m, tick()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	if m.running {
		return m, tick()
	}
	return m, nil
}

func colorizeFiglet(figStr string) string {
	c := color.New(color.FgHiGreen, color.Bold)
	lines := strings.Split(figStr, "\n")
	for i, line := range lines {
		lines[i] = c.Sprint(line)
	}
	return strings.Join(lines, "\n")
}

func (m model) View() string {
	helptxt := styleHelp.Width(m.width).Render("[space] start/stop  [q] quit")
	timerStr := fmt.Sprintf("%.5f", m.elapsed.Seconds())
	if !m.running && m.elapsed == 0 {
		timerStr = "0.00000"
	}
	fig := figure.NewFigure(timerStr, "big", false)
	figStr := colorizeFiglet(fig.String())

	// Show up to 5 past solves
	past := ""
	if len(m.solves) > 0 {
		past = "Past solves:"
		for i, s := range m.solves {
			if i >= 5 {
				break
			}
			past += fmt.Sprintf("\n%2d. %.5f", i+1, s)
		}
	}

	content := fmt.Sprintf("%s\n%s\n\n%s", figStr, helptxt, past)
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		content,
		lipgloss.WithWhitespaceChars(" "),
	)
	return styleBg.Width(m.width).Height(m.height).Render(centered)
}

func tick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Millisecond * 30)
		return tickMsg(time.Now())
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
