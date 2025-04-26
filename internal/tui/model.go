package tui

import (
	"encoding/json"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Tab int

const (
	QueuesTab Tab = iota
	JobsTab
	ClientsTab
)

type Model struct {
	activeTab   Tab
	queueModel  QueueModel
	jobModel    JobModel
	clientModel ClientModel
}

func New(client *river.Client[pgx.Tx]) Model {
	return Model{
		activeTab:   QueuesTab,
		queueModel:  NewQueueModel(client),
		jobModel:    NewJobModel(),
		clientModel: NewClientModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			switch m.activeTab {
			case QueuesTab:
				m.activeTab = JobsTab
			case JobsTab:
				m.activeTab = ClientsTab
			case ClientsTab:
				m.activeTab = QueuesTab
			}
		case "up", "k":
			switch m.activeTab {
			case QueuesTab:
				if m.queueModel.cursor > 0 {
					m.queueModel.cursor--
				}
			case JobsTab:
				if m.jobModel.cursor > 0 {
					m.jobModel.cursor--
				}
			case ClientsTab:
				if m.clientModel.cursor > 0 {
					m.clientModel.cursor--
				}
			}
		case "down", "j":
			switch m.activeTab {
			case QueuesTab:
				if m.queueModel.cursor < len(m.queueModel.queueList.Items)-1 {
					m.queueModel.cursor++
				}
			case JobsTab:
				if m.jobModel.cursor < len(m.jobModel.Jobs)-1 {
					m.jobModel.cursor++
				}
			case ClientsTab:
				if m.clientModel.cursor < len(m.clientModel.Clients)-1 {
					m.clientModel.cursor++
				}
			}
		case "enter":
			switch m.activeTab {
			case QueuesTab:
				m.queueModel.showDetails = !m.queueModel.showDetails
				m.queueModel.selected = m.queueModel.cursor
			case JobsTab:
				m.jobModel.showDetails = !m.jobModel.showDetails
				m.jobModel.selected = m.jobModel.cursor
			case ClientsTab:
				m.clientModel.showDetails = !m.clientModel.showDetails
				m.clientModel.selected = m.clientModel.cursor
			}
		case "f":
			if m.activeTab == JobsTab {
				m.jobModel.showFilter = !m.jobModel.showFilter
			}
		case "1":
			if m.activeTab == JobsTab && m.jobModel.showFilter {
				m.jobModel.filter = StatusPending
			}
		case "2":
			if m.activeTab == JobsTab && m.jobModel.showFilter {
				m.jobModel.filter = StatusRunning
			}
		case "3":
			if m.activeTab == JobsTab && m.jobModel.showFilter {
				m.jobModel.filter = StatusCompleted
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s string

	// Tab styling
	activeTab := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#666"))
	inactiveTab := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))

	// Render tabs
	queuesTab := "Queues"
	jobsTab := "Jobs"
	clientsTab := "Clients"

	switch m.activeTab {
	case QueuesTab:
		queuesTab = activeTab.Render(queuesTab)
		jobsTab = inactiveTab.Render(jobsTab)
		clientsTab = inactiveTab.Render(clientsTab)
	case JobsTab:
		queuesTab = inactiveTab.Render(queuesTab)
		jobsTab = activeTab.Render(jobsTab)
		clientsTab = inactiveTab.Render(clientsTab)
	case ClientsTab:
		queuesTab = inactiveTab.Render(queuesTab)
		jobsTab = inactiveTab.Render(jobsTab)
		clientsTab = activeTab.Render(clientsTab)
	}

	s += fmt.Sprintf("%s | %s | %s\n\n", queuesTab, jobsTab, clientsTab)

	// Render content based on active tab
	switch m.activeTab {
	case QueuesTab:
		s += m.queueModel.View()
	case JobsTab:
		s += m.renderJobs()
	case ClientsTab:
		s += m.renderClients()
	}

	s += "\nPress 'tab' to switch tabs, 'q' to quit\n"
	return s
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

func (m Model) renderClients() string {
	if m.clientModel.showDetails && m.clientModel.selected >= 0 {
		return m.renderClientDetails()
	}

	s := "ID\tCreated\tRunning\tStatus\n"
	s += "──\t───────\t───────\t──────\n"

	for i, client := range m.clientModel.Clients {
		cursor := " "
		if m.clientModel.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\t%s\t%d\t%s\n",
			cursor,
			client.ID,
			client.CreatedAt.Format("2006-01-02 15:04"),
			client.Running,
			client.Status)
	}
	return s
}

func (m Model) renderClientDetails() string {
	client := m.clientModel.Clients[m.clientModel.selected]
	return fmt.Sprintf("Client Details\n\n"+
		"ID: %s\n"+
		"Created At: %s\n"+
		"Running Jobs: %d\n"+
		"Status: %s\n\n"+
		"Concurrency Settings:\n"+
		"Global Limit: unlimited\n"+
		"Local Limit: 0\n\n"+
		"Partitioning:\n"+
		"Partition by kind: disabled\n"+
		"Partition by args: disabled\n\n"+
		"Press 'enter' to go back",
		client.ID,
		client.CreatedAt.Format("2006-01-02 15:04:05"),
		client.Running,
		client.Status)
}
