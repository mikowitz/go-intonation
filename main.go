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
)

type streamer struct {
	playable    intonation.Playable
	stream      *beep.Ctrl
	isStreaming bool
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
	input  textinput.Model
	ratio  intonation.Ratio
	edo    intonation.TwelveEDOInterval
	output audio.AudioOutput
	err    error

	ratioStreamer *streamer
	edoStreamer   *streamer
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
	return textinput.Blink
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
		case "p":
			if !m.ratioStreamer.isStreaming {
				m.PlayRatio()
			} else {
				m.PauseRatio()
			}
		case "P":
			if !m.edoStreamer.isStreaming {
				m.PlayEDO()
			} else {
				m.PauseEDO()
			}
		case "enter":
			m.SetRatio()
			m.ClearRatioStreamer()
			m.ClearEDOStreamer()
			m.PlayRatio()
			m.input.CursorEnd()
		}
	}
	return m, cmd
}

func (m model) View() string {
	currentRatio := ""
	currentEDO := ""
	if m.ratioStreamer.isStreaming {
		r, _ := m.ratioStreamer.playable.(intonation.Ratio)
		currentRatio = fmt.Sprintf("%s (%s)",
			m.ratioStreamer.playable.String(),
			r.Approximate12EDOInterval().String(),
		)
	}
	if m.edoStreamer.isStreaming {
		currentEDO = m.edoStreamer.playable.String()
	}
	return fmt.Sprintf(
		"Enter a Just Intonation ratio (q to quit)\n\nplaying: %-10s %s\n\n%s\n\n============\n\n%-10s (%s)",
		currentRatio,
		currentEDO,
		m.input.View(),
		m.ratio,
		m.ratio.Approximate12EDOInterval().String(),
	)
}

func (m *model) PlayRatio() {
	m.ratioStreamer = &streamer{
		playable:    m.ratio,
		stream:      &beep.Ctrl{},
		isStreaming: false,
	}
	m.ratioStreamer.Stream(m.output)
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
	m.edoStreamer = &streamer{
		playable: m.edo,
	}
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
