package tui

import "time"

type JobStatus string

const (
	StatusCreated   JobStatus = "created"
	StatusScheduled JobStatus = "scheduled"
	StatusWait      JobStatus = "wait"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusPending   JobStatus = "pending"
)

type StatusEvent struct {
	Status    JobStatus
	Timestamp time.Time
}

type Job struct {
	ID          string
	Name        string
	Status      JobStatus
	Attempt     int
	MaxAttempts int
	Priority    int
	Queue       string
	CreatedAt   time.Time
	Timeline    []StatusEvent
	Args        map[string]interface{}
	Metadata    map[string]interface{}
	AttemptedBy string
}

type JobModel struct {
	Jobs        []Job
	cursor      int
	filter      JobStatus
	showFilter  bool
	showDetails bool
	selected    int
}

func NewJobModel() JobModel {
	now := time.Now()
	mockJob := Job{
		ID:          "244894479",
		Name:        "AITrainingBatch",
		Status:      StatusCompleted,
		Attempt:     1,
		MaxAttempts: 25,
		Priority:    1,
		Queue:       "long_running",
		CreatedAt:   now.Add(-8 * time.Minute),
		Timeline: []StatusEvent{
			{Status: StatusCreated, Timestamp: now.Add(-8 * time.Minute)},
			{Status: StatusScheduled, Timestamp: now.Add(-8 * time.Minute)},
			{Status: StatusWait, Timestamp: now.Add(-7 * time.Minute)},
			{Status: StatusRunning, Timestamp: now.Add(-6 * time.Minute)},
			{Status: StatusCompleted, Timestamp: now.Add(-21 * time.Second)},
		},
		Args: map[string]interface{}{
			"LLMversion": "GPT-4",
		},
		Metadata: map[string]interface{}{
			"periodic": true,
		},
		AttemptedBy: "d8901786443558_2025_04_10T15_22_21_313530",
	}

	return JobModel{
		Jobs:        []Job{mockJob},
		cursor:      0,
		filter:      "",
		showFilter:  false,
		showDetails: false,
		selected:    -1,
	}
}
