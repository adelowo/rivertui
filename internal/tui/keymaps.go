package tui

import "github.com/charmbracelet/bubbles/key"

type tuiKeyMaps struct {
	down       key.Binding
	up         key.Binding
	quit       key.Binding
	switchTabs key.Binding
	blurTable  key.Binding
}

func newListKeyMap() *tuiKeyMaps {
	return &tuiKeyMaps{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("k", "go up"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("j", "go down"),
		),
		switchTabs: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "Switch Tabs"),
		),
		blurTable: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Blur or Unblur"),
		),
	}
}
