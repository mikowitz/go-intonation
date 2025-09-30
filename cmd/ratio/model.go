package ratio

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikowitz/intonation/pkg/audio"
)

type Model struct {
	input         InputModel
	ratioStreamer StreamerModel
	edoStreamer   StreamerModel
}

func NewRatioUI(output audio.AudioOutput) Model {
	return Model{
		input:         NewInputModel(),
		ratioStreamer: NewStreamerModel("ratio", "p", output),
		edoStreamer:   NewStreamerModel("edo", "P", output),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.ratioStreamer.Init(), m.edoStreamer.Init())
}

type PlayMsg struct {
	id StreamerID
}

type PauseMsg struct {
	id StreamerID
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case SetRatioMsg:
		m.ratioStreamer.setPlayable(msg.ratio)
		m.edoStreamer.setPlayable(msg.ratio.Approximate12EDOInterval().Interval())
		cmds = append(cmds, m.ratioStreamer.Play, m.edoStreamer.Pause)
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			if m.ratioStreamer.isStreaming && !m.edoStreamer.isStreaming {
				cmds = append(cmds, m.ratioStreamer.Pause, m.edoStreamer.Play)
			} else if m.edoStreamer.isStreaming && !m.ratioStreamer.isStreaming {
				cmds = append(cmds, m.edoStreamer.Pause, m.ratioStreamer.Play)
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	newInput, cmd := m.input.Update(msg)
	if in, ok := newInput.(InputModel); ok {
		m.input = in
	}
	cmds = append(cmds, cmd)

	newRatioStreamer, cmd := m.ratioStreamer.Update(msg)
	if rs, ok := newRatioStreamer.(StreamerModel); ok {
		m.ratioStreamer = rs
	}
	cmds = append(cmds, cmd)

	newEdoStreamer, cmd := m.edoStreamer.Update(msg)
	if es, ok := newEdoStreamer.(StreamerModel); ok {
		m.edoStreamer = es
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	streamers := lipgloss.JoinHorizontal(lipgloss.Top, m.ratioStreamer.View(), m.edoStreamer.View())
	return m.input.View() + "\n" + streamers + "\n"
}
