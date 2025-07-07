# 🟩 Rubik's Cube Timer (Terminal App) ⏱️

A beautiful terminal-based Rubik's Cube timer written in Go, featuring:
- 🎨 Large, colored ASCII-art timer display
- ⌨️ Start/stop with the space bar
- 📝 List of past solve times
- 🖥️ Responsive, modern TUI (Bubble Tea + Lipgloss)

## ✨ Features
- **Big, bold timer**: See your solve time 
- **Simple controls**: Press `[space]` to start/stop, `[q]` to quit
- **Past solves**: View your last 5 solve times right in the UI
- **Cross-platform**: Works on macOS, Linux, and Windows (with a compatible terminal)

## 🚀 Usage
1. **Install Go** (if you haven't already): https://go.dev/dl/
2. **Clone this repo** and install dependencies:
   ```sh
   git clone <your-repo-url>
   cd puzzle-timer
   go mod tidy
   ```
3. **Run the timer:**
   ```sh
   go run main.go
   ```
4. **Controls:**
   - Press `[space]` to start/stop the timer
   - Press `[q]` to quit

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

[space] start/stop  [q] quit

Past solves:
 1. 12.34567
 2. 13.12345
 3. 11.98765
```

## 📝 License
MIT
