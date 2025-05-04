package tui

import "github.com/charmbracelet/bubbles/key"

type tuiKeyMaps struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding

	down       key.Binding
	up         key.Binding
	quit       key.Binding
	switchTabs key.Binding
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
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}
