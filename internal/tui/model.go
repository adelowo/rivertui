package tui

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Tab int

const (
	JobsTab Tab = iota
	QueuesTab
	ClientsTab
)

type Model struct {
	activeTab  Tab
	queueModel QueueModel
	jobModel   JobModel
	keyMaps    *tuiKeyMaps
}

func New(client *river.Client[pgx.Tx]) Model {
	return Model{
		activeTab:  JobsTab,
		queueModel: NewQueueModel(client),
		jobModel:   NewJobModel(),
		keyMaps:    newListKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	tea.SetWindowTitle("rivertui")
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMaps.up):
			m.queueModel.table.MoveUp(1)

		case key.Matches(msg, m.keyMaps.down):
			m.queueModel.table.MoveDown(1)

		case key.Matches(msg, m.keyMaps.quit):
			return m, tea.Quit

		case key.Matches(msg, m.keyMaps.blurTable):

			if m.queueModel.table.Focused() {
				m.queueModel.table.Blur()
			} else {
				m.queueModel.table.Focus()
			}

		case key.Matches(msg, m.keyMaps.switchTabs):

			switch m.activeTab {
			case JobsTab:
				m.activeTab = QueuesTab

			case QueuesTab:
				m.activeTab = ClientsTab

			case ClientsTab:
				m.activeTab = JobsTab

			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var s string

	tabRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.tabContent("Jobs", JobsTab, activeTab, inactiveTab),
		m.tabContent("Queues", QueuesTab, activeTab, inactiveTab),
		m.tabContent("Clients", ClientsTab, activeTab, inactiveTab),
	)

	s += tabRow + "\n\n"

	switch m.activeTab {
	case QueuesTab:
		return s + m.queueModel.View()
	case JobsTab:
		return s + m.renderJobs()
	default:
		return s
	}
}

func (m Model) tabContent(title string, tab Tab, activeStyle, inactiveStyle lipgloss.Style) string {
	if m.activeTab == tab {
		return activeStyle.Render(title)
	}
	return inactiveStyle.Render(title)
}

func (m Model) renderJobs() string {
	if m.jobModel.showDetails && m.jobModel.selected >= 0 {
		return m.renderJobDetails()
	}

	s := ""
	if m.jobModel.showFilter {
		s += "Filter by status (press 1-3):\n"
		s += "1. Pending\n"
		s += "2. Running\n"
		s += "3. Completed\n\n"
	}

	s += "Jobs:\n"
	for i, job := range m.jobModel.Jobs {
		if m.jobModel.filter != "" && job.Status != m.jobModel.filter {
			continue
		}
		cursor := " "
		if m.jobModel.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s - %s (%s)\n",
			cursor,
			job.ID,
			job.Name,
			job.Status)
	}
	return s
}

func (m Model) renderJobDetails() string {
	job := m.jobModel.Jobs[m.jobModel.selected]

	s := fmt.Sprintf("Job Details - %s\n", job.Name)
	s += fmt.Sprintf("ID: %s\n\n", job.ID)

	s += fmt.Sprintf("State: %s\n", job.Status)
	s += fmt.Sprintf("Attempt: %d/%d\n", job.Attempt, job.MaxAttempts)
	s += fmt.Sprintf("Priority: %d\n\n", job.Priority)

	s += fmt.Sprintf("Queue: %s\n", job.Queue)
	s += fmt.Sprintf("Created: %s\n\n", job.CreatedAt.Format("15:04:05"))

	s += "Timeline:\n"
	for _, event := range job.Timeline {
		timeAgo := time.Since(event.Timestamp).Round(time.Second)
		s += fmt.Sprintf("• %s\n  %s ago\n", event.Status, timeAgo)
	}

	s += "\nArgs:\n"
	argsJSON, _ := json.MarshalIndent(job.Args, "", "  ")
	s += fmt.Sprintf("%s\n", string(argsJSON))

	s += "\nMetadata:\n"
	metaJSON, _ := json.MarshalIndent(job.Metadata, "", "  ")
	s += fmt.Sprintf("%s\n", string(metaJSON))

	s += fmt.Sprintf("\nAttempted By:\n%s\n", job.AttemptedBy)

	s += "\nPress 'enter' to go back"
	return s
}
