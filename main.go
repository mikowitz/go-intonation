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

type model struct {
	input               textinput.Model
	ratio               intonation.Ratio
	edo                 intonation.TwelveEDOInterval
	output              audio.AudioOutput
	playing             bool
	playingEDO          bool
	currentlyPlaying    intonation.Ratio
	ratioStream         *beep.Ctrl
	currentlyPlayingEDO intonation.TwelveEDOInterval
	edoStream           *beep.Ctrl
	err                 error
}

func initialModel() model {
	output := internal.BeepAudioOutput{SampleRate: beep.SampleRate(48000)}
	ti := textinput.New()
	ti.Placeholder = "1/1"
	ti.Focus()
	ti.CharLimit = 11
	ti.Width = 20

	return model{
		input:  ti,
		ratio:  intonation.NewRatio(1, 1),
		output: output,
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
			log.Printf("m.playing? %t", m.playing)
			if !m.playing {
				m.playing = true
				m.PlayRatio()
			} else {
				m.playing = false
				m.PauseRatio()
				m.playingEDO = false
				m.PauseEDO()
			}
		case "P":
			log.Printf("m.playingEDO? %t", m.playingEDO)
			if !m.playingEDO {
				m.playingEDO = true
				m.PlayEDO()
			} else {
				m.playingEDO = false
				m.PauseEDO()
			}
		case "enter":
			m.SetRatio()
			m.ClearRatioStreamer()
			m.playingEDO = false
			m.CLearEDOStreamer()
			m.PlayRatio()
			m.input.CursorEnd()
		}
	}
	return m, cmd
}

func (m model) View() string {
	currentRatio := ""
	currentEDO := ""
	if m.playing {
		currentRatio = m.currentlyPlaying.String()
		currentEDO = m.currentlyPlaying.Approximate12EDOInterval().String()
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
	m.playing = true
	m.currentlyPlaying = m.ratio
	ratioCtrl, err := intonation.StreamChord(m.currentlyPlaying, context.TODO(), m.output)
	if err != nil {
		log.Printf("error playing ratio: %s (%v)", m.currentlyPlaying, err)
	}
	m.ratioStream = ratioCtrl
}

func (m *model) PauseRatio() {
	if m.ratioStream != nil {
		intonation.PauseStream(m.ratioStream, context.TODO(), m.output)
		m.ratioStream = nil
	}
}

func (m *model) ClearRatioStreamer() {
	if m.ratioStream != nil {
		m.ratioStream.Streamer = nil
	}
}

func (m *model) PlayEDO() {
	m.playingEDO = true
	m.currentlyPlayingEDO = m.edo
	edoCtrl, err := intonation.StreamChord(m.currentlyPlayingEDO, context.TODO(), m.output)
	if err != nil {
		log.Printf("error playing interval: %s (%v)", m.currentlyPlayingEDO, err)
	}
	m.edoStream = edoCtrl
}

func (m *model) PauseEDO() {
	if m.edoStream != nil {
		intonation.PauseStream(m.edoStream, context.TODO(), m.output)
		m.edoStream = nil
	}
}

func (m *model) CLearEDOStreamer() {
	if m.edoStream != nil {
		m.edoStream.Streamer = nil
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
