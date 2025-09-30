package ratio

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	intonation "github.com/mikowitz/intonation/pkg"
)

type SetRatioMsg struct {
	ratio intonation.Ratio
}

func (r SetRatioMsg) String() string {
	return "set current ratio on streamers"
}

type InputModel struct {
	input textinput.Model
	ratio intonation.Ratio
	err   error
}

func NewInputModel() InputModel {
	ti := textinput.New()
	ti.Placeholder = "1/1"
	ti.Focus()
	ti.Width = 20

	return InputModel{
		input: ti,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

var allowedKeys = key.NewBinding(
	key.WithKeys(
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "/", "backspace", "left", "right",
	),
	key.WithHelp("0-9, /", "enter JI ratio"),
)

var enterKeys = key.NewBinding(
	key.WithKeys("enter", "return"),
	key.WithHelp("Enter", "set the ratio"),
)

func setRatios(ratio intonation.Ratio) tea.Cmd {
	return func() tea.Msg {
		return SetRatioMsg{ratio: ratio}
	}
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, enterKeys):
			m.UpdateRatio()
			if m.ratio.Denom > 0 {
				cmd = setRatios(m.ratio)
			}
		case key.Matches(msg, allowedKeys):
			m.input, cmd = m.input.Update(msg)
			m.UpdateRatio()
		}
	}
	return m, cmd
}

func (m *InputModel) UpdateRatio() {
	r, err := intonation.NewRatioFromString(m.input.Value())
	if err != nil {
		m.ratio.Numer = 0
		m.ratio.Denom = 0
		m.err = err
	} else {
		m.ratio = r
		m.err = nil
	}
}

func (m InputModel) View() string {
	outputStr := ""
	if m.ratio.Denom != 0 {
		outputStr = fmt.Sprintf("%s\t%s", m.ratio.String(), m.ratio.Approximate12EDOInterval())
	}

	return inputStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Center,
			textInputStyle.Render(m.input.View()),
			lightStyle.Render(outputStr),
		),
	)
}

var textInputStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.InnerHalfBlockBorder()).
	BorderRight(true).
	Width(20)

var inputStyle = lipgloss.NewStyle().
	Width(62).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63"))

var lightStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("12")).
	AlignHorizontal(lipgloss.Right).
	Width(40)
