package ratio

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Input    []key.Binding
	Playback []key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return k.Input
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.Input,
		k.Playback,
	}
}

func RatioKeymap() keyMap {
	return keyMap{
		Input: []key.Binding{
			allowedKeys,
			enterKeys,
			clearKeys,
			quitKeys,
		},
		Playback: []key.Binding{
			toggleRatio,
			toggleEdo,
			swap,
			muteAll,
		},
	}
}

var allowedKeys = key.NewBinding(
	key.WithKeys(
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "/", "backspace", "left", "right",
	),
	key.WithHelp("0-9, /", "enter JI ratio"),
)

var enterKeys = key.NewBinding(
	key.WithKeys("enter", "return"),
	key.WithHelp("Enter", "set ratio"),
)

var clearKeys = key.NewBinding(
	key.WithKeys("c"),
	key.WithHelp("c", "clear ratio"),
)

var quitKeys = key.NewBinding(
	key.WithKeys("q", "ctrl+c"),
	key.WithHelp("q/ctrl+c", "quit"),
)

var toggleRatio = key.NewBinding(
	key.WithKeys("p"),
	key.WithHelp("p", "play/pause ratio stream"),
)

var toggleEdo = key.NewBinding(
	key.WithKeys("P"),
	key.WithHelp("P", "play/pause EDO stream"),
)

var swap = key.NewBinding(
	key.WithKeys("s"),
	key.WithHelp("s", "swap playback"),
)

var muteAll = key.NewBinding(
	key.WithKeys("m"),
	key.WithHelp("m", "mute all streams"),
)
