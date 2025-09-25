package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/mikowitz/intonation/internal"
	intonation "github.com/mikowitz/intonation/pkg"
	"github.com/mikowitz/intonation/pkg/audio"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type streamer struct {
	playable    intonation.Playable
	stream      *beep.Streamer
	isStreaming bool
	triggerKey  string
}

func (s streamer) Init() tea.Cmd {
	return nil
}

func (s streamer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return s, cmd
}

func (s streamer) View() string {
	if s.playable == nil {
		return ""
	}
	str := s.playable.String()

	if s.isStreaming {
		str += "\n\n(" + s.triggerKey + ")ause"
	} else {
		str += "\n\n(" + s.triggerKey + ")lay"
	}
	return str
}

func (s *streamer) Stream(output audio.AudioOutput) {
	s.isStreaming = true
	stream, err := intonation.StreamChord(s.playable, context.TODO(), output)
	if err != nil {
		log.Printf("error streaming: %s (%v)", s.playable, err)
	}
	s.stream = stream
}

func (s *streamer) Pause(output audio.AudioOutput) {
	if s.stream == nil {
		return
	}
	s.isStreaming = false
	intonation.PauseStream(s.stream, context.TODO(), output)
}

func (s *streamer) Clear() {
	s.isStreaming = false
	if s.stream != nil {
		s.stream.Streamer = nil
	}
}

type model struct {
	input        textinput.Model
	ratio        intonation.Ratio
	edo          intonation.TwelveEDOInterval
	output       audio.AudioOutput
	err          error
	muted        bool
	mutedStreams []string

	ratioStreamer streamer
	edoStreamer   streamer
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "1/1"
	ti.Focus()
	ti.CharLimit = 11
	ti.Width = 20

	output := internal.BeepAudioOutput{SampleRate: beep.SampleRate(48000)}

	return model{
		input:         ti,
		ratio:         intonation.NewRatio(1, 1),
		output:        output,
		ratioStreamer: &streamer{},
		edoStreamer:   &streamer{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.ratioStreamer.Init(), m.edoStreamer.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println(msg.String())
		switch msg.String() {
		case "ctrl+c", "q":
			speaker.Clear()
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "/", "backspace":
			m.input, cmd = m.input.Update(msg)
			m.SetRatio()
		case "up", "down", "left", "right":
			m.input, cmd = m.input.Update(msg)
			m.SetRatio()
		case "s":
			if m.ratioStreamer.isStreaming && !m.edoStreamer.isStreaming {
				m.PauseRatio()
				m.PlayEDO()
			} else if m.edoStreamer.isStreaming && !m.ratioStreamer.isStreaming {
				m.PauseEDO()
				m.PlayRatio()
			}
		case "S":
			m.PauseRatio()
			m.PauseEDO()
		case "p":
			if m.ratioStreamer.isStreaming {
				m.PauseRatio()
			} else {
				m.PlayRatio()
			}
		case "P":
			if m.edoStreamer.isStreaming {
				m.PauseEDO()
			} else {
				m.PlayEDO()
			}
		case "enter":
			m.SetRatio()
			m.ClearRatioStreamer()
			m.ClearEDOStreamer()
			m.PlayRatio()
			m.SetEDOStreamer()
			m.input.CursorEnd()
		}
	}

	m.ratioStreamer, cmd = m.ratioStreamer.Update(msg)
	m.ratioStreamer, cmd = m.ratioStreamer.Update(msg)
	return m, cmd
}

func (m model) View() string {
	streamerStyle := lipgloss.NewStyle().
		Width(20).
		Height(10).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	streamCommands := ""
	if m.ratioStreamer.isStreaming || m.edoStreamer.isStreaming {
		streamCommands = "(S)top all\n\n(s)wap"
	}
	ratioStreamer := streamerStyle.Render(m.ratioStreamer.View())
	edoStreamer := streamerStyle.Render(m.edoStreamer.View())
	commands := streamerStyle.Render(streamCommands)

	streamers := lipgloss.JoinHorizontal(lipgloss.Center, ratioStreamer, commands, edoStreamer)

	centerStyle := lipgloss.NewStyle().
		Width(64).
		Align(lipgloss.Center)

	inputStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	return fmt.Sprintf(
		"Enter a Just Intonation ratio (q to quit)\n\n%s\n%s\n",
		inputStyle.Render(
			fmt.Sprintf(
				"%s\n%s",
				centerStyle.Render(
					inputStyle.Render(
						m.input.View(),
					),
				),
				centerStyle.Render(
					fmt.Sprintf(
						"%-10s (%s)",
						m.ratio,
						m.ratio.Approximate12EDOInterval().String(),
					),
				),
			),
		),
		streamers,
	)
}

func (m *model) PlayRatio() {
	m.ratioStreamer = &streamer{
		playable:    m.ratio,
		stream:      &beep.Ctrl{},
		isStreaming: false,
		triggerKey:  "p",
	}
	m.ratioStreamer.Stream(m.output)
}

func (m *model) SetEDOStreamer() {
	m.edoStreamer = &streamer{
		playable:    m.edo,
		stream:      &beep.Ctrl{},
		isStreaming: false,
		triggerKey:  "P",
	}
}

func (m *model) PauseRatio() {
	if m.ratioStreamer != nil {
		m.ratioStreamer.Pause(m.output)
	}
}

func (m *model) ClearRatioStreamer() {
	if m.ratioStreamer != nil {
		m.ratioStreamer.Clear()
	}
}

func (m *model) PlayEDO() {
	m.edoStreamer.playable = m.edo
	m.edoStreamer.Stream(m.output)
}

func (m *model) PauseEDO() {
	if m.edoStreamer != nil {
		m.edoStreamer.Pause(m.output)
	}
}

func (m *model) ClearEDOStreamer() {
	if m.edoStreamer != nil {
		m.edoStreamer.Clear()
	}
}

func (m *model) SetRatio() {
	ratio, err := intonation.NewRatioFromString(m.input.Value())
	if err != nil {
		m.err = err
	}
	m.ratio = ratio
	m.edo = ratio.Approximate12EDOInterval().Interval()
	m.err = nil
}

func main() {
	// cmd.Run()
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
