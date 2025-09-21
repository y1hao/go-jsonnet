package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/pflag"
)

var (
	mainColor = lipgloss.Color("28")
	gray      = lipgloss.Color("245")
	keyMap    = viewport.DefaultKeyMap()

	headerStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			Background(mainColor).
			Padding(0, 1)
	}()
	dataStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(mainColor).
			Padding(0, 1)
	}()
	heighlightStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(mainColor).Bold(true)
	}()
	infoStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(gray)
	}()
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
	fmt.Fprintln(o, "  -J / --jpath <dir>         Specify an additional library search dir")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Environment variables:")
	fmt.Fprintln(o, "  JSONNET_PATH is a colon (semicolon on Windows) separated list of directories")
	fmt.Fprintln(o, "  added in reverse order before the paths specified by --jpath (i.e. left-most")
	fmt.Fprintln(o, "  wins). E.g. these are equivalent:")
	fmt.Fprintln(o, "    JSONNET_PATH=a:b jsonnet-trace -J c -J d")
	fmt.Fprintln(o, "    JSONNET_PATH=d:c:a:b jsonnet-trace")
	fmt.Fprintln(o, "    jsonnet-trace -J b -J a -J c -J d")
	fmt.Fprintln(o)
}

func main() {
	showHelp := pflag.Bool("help", false, "Show usage info")
	jpath := pflag.StringArrayP("jpath", "J", []string{}, "Additional library search dir")

	pflag.Parse()

	if showHelp != nil && *showHelp {
		usage(os.Stderr)
		return
	}

	if pflag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "jsonnet-trace can only take a single file")
		usage(os.Stderr)
		os.Exit(1)
	}

	filename := pflag.Args()[0]

	evalJpath := []string{}
	jsonnetPath := filepath.SplitList(os.Getenv("JSONNET_PATH"))
	for i := len(jsonnetPath) - 1; i >= 0; i-- {
		evalJpath = append(evalJpath, jsonnetPath[i])
	}
	if jpath != nil {
		evalJpath = append(evalJpath, *jpath...)
	}

	fmt.Fprintf(os.Stdout, "Building jsonnet file %q...\n", filename)

	result, trace, err := buildWithTrace(filename, evalJpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate trace for file %s: %s", filename, err.Error())
		os.Exit(1)
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

func buildWithTrace(filename string, evalJpath []string) (string, map[int]*ast.LocationRange, error) {
	vm := jsonnet.MakeTracingVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: evalJpath,
	})
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
		titleHeight := lipgloss.Height(m.inputView())
		infoHeight := lipgloss.Height(m.sourceView())
		verticalMarginHeight := titleHeight + infoHeight + 1

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = titleHeight
			m.viewport.SetContent(m.json)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - len(" > 1000 ")
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n\n%s\n%s", m.sourceView(), m.contentView(), m.inputView())
}

func (m model) contentView() string {
	vpView := m.viewport.View()
	offset := m.viewport.YOffset

	lineNumbers := []string{}
	marker := []string{}
	for i := offset; i < m.viewport.Height+offset; i++ {
		if i == m.currentLine {
			marker = append(marker, heighlightStyle.Render(" > "))
			lineNumbers = append(lineNumbers, heighlightStyle.Render(strconv.Itoa(i+1)))
		} else {
			marker = append(marker, "  ")
			lineNumbers = append(lineNumbers, infoStyle.Render(strconv.Itoa(i+1)))
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		strings.Join(marker, "\n"),
		strings.Join(lineNumbers, "\n"),
		vpView)
}

func (m model) inputView() string {
	info := dataStyle.Render(m.filename)
	spaces := strings.Repeat(" ", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, spaces, info)
}

func (m model) sourceView() string {
	return headerStyle.Render("SOURCE") + dataStyle.Render(getOrigin(m.trace, m.currentLine))
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
