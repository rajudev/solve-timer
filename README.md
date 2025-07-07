# 🟩 Rubik's Cube Timer (Terminal App) ⏱️

A beautiful terminal-based Rubik's Cube timer written in Go, featuring:
- 🎨 Large, colored ASCII-art timer display
- ⌨️ Start/stop with the space bar
- 📝 Persistent list of all past solve times (across sessions)
- 🖥️ Responsive, modern TUI (Bubble Tea + Lipgloss)
- 🕒 WCA-style inspection and penalties (+2, DNF)
- 📜 Menu to view and scroll through all solves

## ✨ Features
- **Big, bold timer**: See your solve time in large ASCII art
- **WCA inspection**: 15s inspection phase before each solve, with automatic +2/DNF penalties
- **Simple controls**: Press `[space]` to start/stop, `[r]` to reset timer, `[v]` to view all solves, `[q]` to quit
- **Past solves**: View your last 5 solve times in the main UI, or browse all solves in a scrollable menu
- **Persistent history**: All solves are saved to `~/.solve_timer_solves.json` and loaded on startup
- **Cross-platform**: Works on macOS, Linux, and Windows (with a compatible terminal)

## 🚀 Usage
1. **Install Go** (if you haven't already): https://go.dev/dl/
2. **Clone this repo** and install dependencies:
   ```sh
   git clone <your-repo-url>
   cd solve-timer
   go mod tidy
   ```
3. **Run the timer:**
   ```sh
   go run main.go
   ```
4. **Controls:**
   - `[space]` start/stop/next (with inspection phase)
   - `[r]` reset current timer (does not clear history)
   - `[v]` view all solves (scroll with up/down, exit with q/esc)
   - `[q]` quit

## 📦 Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI framework)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) (styling)
- [go-figure](https://github.com/common-nighthawk/go-figure) (ASCII-art font)
- [fatih/color](https://github.com/fatih/color) (color output)

## 🖼️ Example
```
 _______  _______  _______  _______  _______  _______  _______ 
|       ||       ||       ||       ||       ||       ||       |
|   0   ||   .   ||   0   ||   0   ||   0   ||   0   ||   0   |
|_______||_______||_______||_______||_______||_______||_______|

[space] start/stop/next  [r] reset  [v] view all solves  [q] quit

Past solves:
 1. 12.34567
 2. 13.12345 (+2)
 3. DNF
```

## 📝 License
MIT
