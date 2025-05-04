package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// ENUM(paused,active)
type QueueStatus string

type Queue struct {
	Name      string
	CreatedAt string
	Status    QueueStatus
}

func (i Queue) Title() string { return "" }
func (i Queue) Description() string {
	return ""
}

func (i Queue) FilterValue() string { return i.Name }

type QueueModel struct {
	client      *river.Client[pgx.Tx]
	items       []Queue
	cursor      int
	selected    int
	showDetails bool
	table       table.Model
}

func NewQueueModel(client *river.Client[pgx.Tx]) QueueModel {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Created at", Width: 10},
		{Title: "Status", Width: 6},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	queues, err := client.QueueList(context.Background(), river.NewQueueListParams())
	if err != nil {
		panic(err.Error())
	}

	for _, queue := range queues.Queues {

		var status = QueueStatusActive

		if queue.PausedAt != nil {
			status = QueueStatusPaused
		}

		rows = append(rows, table.Row{
			queue.Name,
			queue.CreatedAt.String(),
			status.String(),
		})
	}

	t.SetRows(rows)

	return QueueModel{
		cursor:      0,
		table:       t,
		selected:    -1,
		showDetails: false,
		client:      client,
	}
}

func (m QueueModel) View() string {
	if m.showDetails && m.selected >= 0 {
		return m.renderDetails()
	}
	return m.renderList()
}

func (m QueueModel) renderList() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func (m QueueModel) renderDetails() string {
	queue := m.items[m.selected]
	return fmt.Sprintf("Queue Details - %s\n\n"+
		"Created At: %s\n"+
		"Press 'enter' to go back",
		queue.Name,
		queue.Status)
}
