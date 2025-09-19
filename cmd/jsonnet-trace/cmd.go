package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func usage(o io.Writer) {
	fmt.Fprintln(o)
	fmt.Fprintln(o, "jsonnet-trace {<option>} { <filename> }")
	fmt.Fprintln(o, "  Build a jsonnet file and collect trace information for debugging.")
	fmt.Fprintln(o, "  The built outcome and trace information will be served on localhost:8080.")
	fmt.Fprintln(o, "  Only a single root file is supported (but it can import other files).")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Available options:")
	fmt.Fprintln(o, "  -h / --help                This message")
	fmt.Fprintln(o)
}

func main() {
	showHelp := flag.Bool("help", false, "Show usage info")
	flag.Parse()

	if showHelp != nil && *showHelp {
		usage(os.Stderr)
		return
	}

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "jsonnet-trace can only take a single file")
		usage(os.Stderr)
	}

	filename := flag.Args()[0]

	fmt.Fprintf(os.Stdout, "Building jsonnet file %q...\n", filename)

	result, trace, err := buildWithTrace(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate trace for file %s: %s", filename, err.Error())
	}

	result = strings.Trim(result, " \n\r")
	lines := strings.Count(result, "\n") + 1
	p := tea.NewProgram(
		model{
			currentLine:     0,
			currentPosition: 0,
			lines:           lines,
			filename:        filename,
			json:            result,
			trace:           trace,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}

}

func buildWithTrace(filename string) (string, map[int]*ast.LocationRange, error) {
	vm := jsonnet.MakeTracingVM()
	result, trace, err := vm.EvaluateFileWithTrace(filename)
	if err != nil {
		return "", nil, fmt.Errorf("error generating trace: %w", err)
	}
	return result, trace, nil
}

func getOrigin(trace map[int]*ast.LocationRange, line int) string {
	loc, ok := trace[line]
	if !ok {
		return ""
	}
	filename, beginLine, endLine := loc.FileName, loc.Begin.Line, loc.End.Line
	if beginLine == endLine {
		return fmt.Sprintf("%s:%d", filename, beginLine)
	}
	return fmt.Sprintf("%s:%d-%d", filename, beginLine, endLine)
}

var (
	keyMap     = viewport.DefaultKeyMap()
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		// b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		// b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type model struct {
	currentLine     int
	currentPosition int
	lines           int
	json            string
	filename        string
	trace           map[int]*ast.LocationRange
	ready           bool
	viewport        viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c" || msg.String() == "q" || msg.String() == "esc":
			return m, tea.Quit

		case key.Matches(msg, keyMap.Down):
			m.down()
		case key.Matches(msg, keyMap.Up):
			m.up()
		case key.Matches(msg, keyMap.Left):
			m.viewport.ScrollLeft(1)
		case key.Matches(msg, keyMap.Right):
			m.viewport.ScrollRight(1)
		}

	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			if msg.Shift {
				m.viewport.ScrollLeft(1)
			} else {
				m.up()
			}
		case tea.MouseButtonWheelDown:
			if msg.Shift {
				m.viewport.ScrollRight(1)
			} else {
				m.down()
			}
		case tea.MouseButtonWheelLeft:
			m.viewport.ScrollLeft(1)
		case tea.MouseButtonWheelRight:
			m.viewport.ScrollRight(1)
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.json)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.contentView(), m.footerView())
}

func (m model) contentView() string {
	vpView := m.viewport.View()
	offset := m.viewport.YOffset

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	lineNumbers := []string{}
	for i := offset; i < m.viewport.Height+offset; i++ {
		if i == m.currentLine {
			lineNumbers = append(lineNumbers, style.Render(strconv.Itoa(i)))
		} else {
			lineNumbers = append(lineNumbers, strconv.Itoa(i))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, strings.Join(lineNumbers, "\n"), vpView)
}

func (m model) headerView() string {
	title := titleStyle.Render(fmt.Sprintf("%s (Current line: %d)", m.filename, m.currentLine))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%s", getOrigin(m.trace, m.currentLine)))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *model) up() {
	if m.currentLine > 0 {
		m.currentLine--
	}
	if m.currentPosition > 0 {
		m.currentPosition--
	} else {
		m.viewport.ScrollUp(1)
	}
}

func (m *model) down() {
	if m.currentLine < m.lines-1 {
		m.currentLine++
	}
	if m.currentPosition < m.viewport.Height-1 {
		m.currentPosition++
	} else {
		m.viewport.ScrollDown(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
