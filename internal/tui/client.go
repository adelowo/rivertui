package tui

import "time"

type Client struct {
	ID        string
	CreatedAt time.Time
	Running   int
	Status    string
}

type ClientModel struct {
	Clients     []Client
	cursor      int
	selected    int
	showDetails bool
}

func NewClientModel() ClientModel {
	mockClients := []Client{
		{
			ID:        "d8981704c63568_2025_0a_10T15_22_21_313530",
			CreatedAt: time.Now().Add(-16 * 24 * time.Hour),
			Running:   1,
			Status:    "Active",
		},
		{
			ID:        "91853399454615_2025_0a_10T15_31_36_096670",
			CreatedAt: time.Now().Add(-16 * 24 * time.Hour),
			Running:   1,
			Status:    "Active",
		},
	}

	return ClientModel{
		Clients:     mockClients,
		cursor:      0,
		selected:    -1,
		showDetails: false,
	}
}
