package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

var (
	styleHelp = lipgloss.NewStyle().Foreground(lipgloss.Color("28")).Align(lipgloss.Center) // green help text
	styleBg   = lipgloss.NewStyle().Background(lipgloss.Color("0"))                         // black background
)

type tickMsg time.Time

type solveResult struct {
	Time    float64 `json:"time"`
	Penalty string  `json:"penalty"`
}

type model struct {
	running           bool
	start             time.Time
	elapsed           time.Duration
	quitting          bool
	width             int
	height            int
	solves            []solveResult // store past solve times and penalties
	inspection        bool          // true if in inspection phase
	inspectionStart   time.Time     // when inspection started
	inspectionElapsed time.Duration
	penalty           string        // current penalty: "", "+2", or "DNF"
	viewingAllSolves  bool          // true if viewing all solves menu
	scrollOffset      int           // for scrolling in solves menu
}

const solvesFileName = ".solve_timer_solves.json"

func getSolvesFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return solvesFileName // fallback to current dir
	}
	return filepath.Join(home, solvesFileName)
}

func loadSolves() []solveResult {
	path := getSolvesFilePath()
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	var solves []solveResult
	if err := json.NewDecoder(f).Decode(&solves); err != nil {
		return nil
	}
	return solves
}

func saveSolves(solves []solveResult) {
	path := getSolvesFilePath()
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(solves)
}

func initialModel() model {
	return model{
		width:  80,
		height: 24,
		solves: loadSolves(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.viewingAllSolves {
			switch msg.String() {
			case "q", "esc":
				m.viewingAllSolves = false
				m.scrollOffset = 0
			case "up", "k":
				if m.scrollOffset > 0 {
					m.scrollOffset--
				}
			case "down", "j":
				if m.scrollOffset < len(m.solves)-1 {
					m.scrollOffset++
				}
			}
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			saveSolves(m.solves)
			return m, tea.Quit
		case "v":
			if len(m.solves) > 0 {
				m.viewingAllSolves = true
			}
		case " ":
			if m.inspection {
				// Start solve from inspection
				m.inspection = false
				m.running = true
				m.start = time.Now()
				m.elapsed = 0
				m.penalty = ""
				if m.inspectionElapsed > 15*time.Second && m.inspectionElapsed <= 17*time.Second {
					m.penalty = "+2"
				} else if m.inspectionElapsed > 17*time.Second {
					m.penalty = "DNF"
				}
				return m, tick()
			} else if !m.running {
				// Start inspection
				m.inspection = true
				m.inspectionStart = time.Now()
				m.inspectionElapsed = 0
				m.penalty = ""
				return m, tick()
			} else {
				// Stop solve
				m.running = false
				m.elapsed = time.Since(m.start)
				penalty := m.penalty
				if penalty == "DNF" {
					m.solves = append([]solveResult{{Time: m.elapsed.Seconds(), Penalty: "DNF"}}, m.solves...)
				} else if penalty == "+2" {
					m.solves = append([]solveResult{{Time: m.elapsed.Seconds() + 2, Penalty: "+2"}}, m.solves...)
				} else {
					m.solves = append([]solveResult{{Time: m.elapsed.Seconds(), Penalty: ""}}, m.solves...)
				}
				saveSolves(m.solves)
			}
		case "r":
			// Reset only the current timer/session, not past solves
			m.running = false
			m.inspection = false
			m.elapsed = 0
			m.penalty = ""
			m.inspectionElapsed = 0
			// Do not clear m.solves or call saveSolves
		}
	case tickMsg:
		if m.inspection {
			m.inspectionElapsed = time.Since(m.inspectionStart)
			if m.inspectionElapsed > 17*time.Second {
				m.penalty = "DNF"
			} else if m.inspectionElapsed > 15*time.Second {
				m.penalty = "+2"
			}
			return m, tick()
		}
		if m.running {
			m.elapsed = time.Since(m.start)
			return m, tick()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	if m.inspection || m.running {
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
	if m.viewingAllSolves {
		return m.viewAllSolvesView()
	}
	helptxt := styleHelp.Width(m.width).Render("[space] start/stop/next  [r] reset  [v] view all solves  [q] quit")
	var timerStr string
	var figStr string
	if m.inspection {
		secs := 15 - int(m.inspectionElapsed.Seconds())
		if secs < 0 {
			secs = 0
		}
		timerStr = fmt.Sprintf("Inspection: %2d", secs)
		if m.penalty == "+2" {
			timerStr += " (+2)"
		} else if m.penalty == "DNF" {
			timerStr += " (DNF)"
		}
		fig := figure.NewFigure(timerStr, "big", false)
		figStr = colorizeFiglet(fig.String())
	} else {
		timeVal := m.elapsed.Seconds()
		if !m.running && m.elapsed == 0 {
			timeVal = 0
		}
		display := fmt.Sprintf("%.5f", timeVal)
		if m.penalty == "+2" {
			display += " (+2)"
		} else if m.penalty == "DNF" {
			display = "DNF"
		}
		fig := figure.NewFigure(display, "big", false)
		figStr = colorizeFiglet(fig.String())
	}

	// Show up to 5 past solves
	past := ""
	if len(m.solves) > 0 {
		past = "Past solves:"
		for i, s := range m.solves {
			if i >= 5 {
				break
			}
			if s.Penalty == "DNF" {
				past += fmt.Sprintf("\n%2d. DNF", i+1)
			} else if s.Penalty == "+2" {
				past += fmt.Sprintf("\n%2d. %.5f (+2)", i+1, s.Time)
			} else {
				past += fmt.Sprintf("\n%2d. %.5f", i+1, s.Time)
			}
		}
	}

	content := fmt.Sprintf("%s\n%s\n\n%s", figStr, helptxt, past)
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		content,
		lipgloss.WithWhitespaceChars(" "),
	)
	return styleBg.Width(m.width).Height(m.height).Render(centered)
}

func (m model) viewAllSolvesView() string {
	header := styleHelp.Width(m.width).Render("All Past Solves [up/down] scroll  [q/esc] back")
	lines := []string{"All Past Solves:"}
	maxRows := m.height - 4 // header + padding
	start := m.scrollOffset
	end := start + maxRows
	if end > len(m.solves) {
		end = len(m.solves)
	}
	for i := start; i < end; i++ {
		s := m.solves[i]
		var line string
		if s.Penalty == "DNF" {
			line = fmt.Sprintf("%3d. DNF", i+1)
		} else if s.Penalty == "+2" {
			line = fmt.Sprintf("%3d. %.5f (+2)", i+1, s.Time)
		} else {
			line = fmt.Sprintf("%3d. %.5f", i+1, s.Time)
		}
		lines = append(lines, line)
	}
	if len(lines) == 1 {
		lines = append(lines, "No solves yet.")
	}
	content := strings.Join(lines, "\n")
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		content+"\n"+header,
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
