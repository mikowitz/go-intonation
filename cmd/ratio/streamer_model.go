package ratio

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/v2"
	intonation "github.com/mikowitz/intonation/pkg"
	"github.com/mikowitz/intonation/pkg/audio"
)

var (
	sineWave    = "⋅.˳˳.⋅ॱ˙˙ॱ⋅.˳˳.⋅ॱ˙˙ॱ⋅.˳˳.⋅"
	runes       = []rune(sineWave)
	sineSpinner = spinner.Spinner{
		Frames: []string{
			string(runes[4:20]),
			string(runes[5:21]),
			string(runes[6:22]),
			string(runes[7:23]),
			string(runes[8:24]),
			string(runes[9:25]),
			string(runes[0:16]),
			string(runes[1:17]),
			string(runes[2:18]),
			string(runes[3:19]),
		},
		FPS: time.Second / 10,
	}
)

type StreamerID string

const (
	RatioID StreamerID = "ratio"
	EdoID   StreamerID = "edo"
)

func (i StreamerID) pan() float64 {
	if i == RatioID {
		return -0.5
	}
	return 0.5
}

type StreamerModel struct {
	id          StreamerID
	playable    intonation.Playable
	streamer    *beep.Ctrl
	isStreaming bool
	trigger     key.Binding
	spinner     spinner.Model
	output      audio.AudioOutput
}

func (m *StreamerModel) setPlayable(playable intonation.Playable) {
	m.playable = playable
	if m.streamer != nil {
		m.isStreaming = false
		m.streamer.Streamer = nil
	}
}

func (m StreamerModel) Play() tea.Msg {
	return PlayMsg{id: m.id}
}

func (m StreamerModel) Pause() tea.Msg {
	return PauseMsg{id: m.id}
}

func NewStreamerModel(id StreamerID, trigger string, output audio.AudioOutput) StreamerModel {
	s := spinner.New()
	s.Spinner = sineSpinner
	return StreamerModel{
		id:      id,
		trigger: key.NewBinding(key.WithKeys(trigger)),
		spinner: s,
		output:  output,
	}
}

func (m StreamerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m StreamerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case ClearRatioMsg:
		m.isStreaming = false
		if m.streamer != nil {
			m.output.PauseStream(context.Background(), m.streamer)
		}
		m.playable = nil
	case PlayMsg:
		if msg.id == m.id {
			m.isStreaming = true
			m.streamer, _ = m.output.StreamChord(context.Background(), m.playable.Dyad(), m.id.pan())
			cmd = m.spinner.Tick
		}
	case PauseMsg:
		if msg.id == m.id {
			m.isStreaming = false
			if m.streamer != nil {
				m.output.PauseStream(context.Background(), m.streamer)
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.trigger):
			if m.isStreaming {
				m.isStreaming = false
				cmd = m.Pause
			} else {
				m.isStreaming = true
				cmds = append(cmds, m.Play)
				cmds = append(cmds, m.spinner.Tick)
			}
		}
	case spinner.TickMsg:
		if m.isStreaming {
			m.spinner, cmd = m.spinner.Update(msg)
		}
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m StreamerModel) View() string {
	if m.playable == nil {
		return ""
	}

	title := ""
	if m.id == RatioID {
		ratio := m.playable.(intonation.Ratio)
		offset := ratio.Approximate12EDOInterval().CentsOffset()
		sign := "+"
		if offset < 0 {
			sign = ""
		}
		title = fmt.Sprintf("%s (%s%.4f)", ratio.String(), sign, offset)
	} else {
		title = m.playable.(intonation.Interval).String()
	}
	spinner := fmt.Sprintf("%-10s", m.spinner.View())

	return streamerStyle.Render(fmt.Sprintf(
		"%s\n\n%s",
		title,
		spinner,
	))
}

var streamerStyle = lipgloss.NewStyle().
	Width(30).
	Align(lipgloss.Center).
	Border(lipgloss.RoundedBorder(), false, true, true, true).
	BorderForeground(lipgloss.Color("62"))
