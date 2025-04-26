package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Queue struct {
	Name      string
	CreatedAt time.Time
	Available int
	Running   int
	Status    string
}

type QueueList struct {
	Items  []Queue
	client *river.Client[pgx.Tx]
}

type QueueModel struct {
	queueList   QueueList
	cursor      int
	selected    int
	showDetails bool
}

func NewQueueModel(client *river.Client[pgx.Tx]) QueueModel {
	mockQueues := []Queue{
		{
			Name:      "default",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			Available: 10,
			Running:   5,
			Status:    "active",
		},
		{
			Name:      "high-priority",
			CreatedAt: time.Now().Add(-12 * time.Hour),
			Available: 15,
			Running:   8,
			Status:    "active",
		},
	}

	client.QueueList(context.Background(), river.NewQueueListParams())

	return QueueModel{
		queueList: QueueList{
			Items:  mockQueues,
			client: client,
		},
		cursor:      0,
		selected:    -1,
		showDetails: false,
	}
}

func (m QueueModel) View() string {
	if m.showDetails && m.selected >= 0 {
		return m.renderDetails()
	}
	return m.renderList()
}

func (m QueueModel) renderList() string {
	s := "Name\tCreated At\tAvailable\tRunning\tStatus\n"
	s += "────\t──────────\t─────────\t───────\t──────\n"

	for i, queue := range m.queueList.Items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\t%s\t%d\t%d\t%s\n",
			cursor,
			queue.Name,
			queue.CreatedAt.Format("2006-01-02 15:04"),
			queue.Available,
			queue.Running,
			queue.Status)
	}
	return s
}

func (m QueueModel) renderDetails() string {
	queue := m.queueList.Items[m.selected]
	return fmt.Sprintf("Queue Details - %s\n\n"+
		"Created At: %s\n"+
		"Available: %d\n"+
		"Running: %d\n"+
		"Status: %s\n\n"+
		"Press 'enter' to go back",
		queue.Name,
		queue.CreatedAt.Format("2006-01-02 15:04:05"),
		queue.Available,
		queue.Running,
		queue.Status)
}
