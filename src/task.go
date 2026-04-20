package main

import (
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"To Do", "In Progress", "Done"}[s]
}

type Task struct {
	id          int
	status      status
	title       string
	description string
	position    int
	completedAt *time.Time
}

func (t Task) Title() string       { return t.title }
func (t Task) Description() string { return t.description }
func (t Task) FilterValue() string { return t.title }

func (t Task) GetLinks() []string {
	re := regexp.MustCompile(`https?://[^\s/$.?#].[^\s]*`)
	return re.FindAllString(t.description, -1)
}

func (t *Task) Next() {
	if t.status < done {
		t.status++
	} else {
		t.status = todo
	}
}

func (t *Task) Prev() {
	if t.status > todo {
		t.status--
	} else {
		t.status = done
	}
}

// Ensure Task implements list.Item
var _ list.Item = Task{}
