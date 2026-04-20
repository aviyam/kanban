package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	boardState state = iota
	formState
	viewState
	deleteState
)

type model struct {
	state      state
	focused    status
	lists      []list.Model
	titleInput textinput.Model
	descInput  textarea.Model
	config     Config
	loaded     bool
	err        error
	editing    bool
}

func NewModel(config Config) model {
	ti := textinput.New()
	ti.Placeholder = "Task Title"
	ti.Focus()

	ta := textarea.New()
	ta.Placeholder = "Task Description (URLs will be clickable with configured open_link key)"

	return model{
		state:      boardState,
		titleInput: ti,
		descInput:  ta,
		config:     config,
	}
}

func (m *model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/3-4, height-6)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[todo].Title = "To Do"
	m.lists[inProgress].Title = "In Progress"
	m.lists[done].Title = "Done"

	tasks, err := LoadTasks()
	if err != nil {
		m.err = err
		return
	}

	for _, t := range tasks {
		m.lists[t.status].InsertItem(len(m.lists[t.status].Items()), t)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		} else {
			for i := range m.lists {
				m.lists[i].SetSize(msg.Width/3-4, msg.Height-6)
			}
		}
	case tea.KeyMsg:
		if m.state == formState {
			return m.updateForm(msg)
		}
		if m.state == viewState {
			return m.updateView(msg)
		}
		if m.state == deleteState {
			return m.updateDelete(msg)
		}
		// Handle board-level keys first
		newModel, boardCmd := m.updateBoard(msg)
		if boardCmd != nil || m.state != boardState {
			return newModel, boardCmd
		}
		m = newModel.(model)
	}

	if m.state == boardState {
		m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	} else if m.state == formState {
		if m.titleInput.Focused() {
			m.titleInput, cmd = m.titleInput.Update(msg)
		} else {
			m.descInput, cmd = m.descInput.Update(msg)
		}
	}
	return m, cmd
}

func (m model) updateBoard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	s := m.config.Shortcuts
	key := msg.String()

	switch key {
	case "ctrl+c", s.Quit:
		return m, tea.Quit
	case "enter":
		selectedItem := m.lists[m.focused].SelectedItem()
		if selectedItem != nil {
			m.state = viewState
		}
		return m, nil
	case "left", s.MoveLeft:
		m.focused = (m.focused + 2) % 3
		return m, nil
	case "right", s.MoveRight:
		m.focused = (m.focused + 1) % 3
		return m, nil
	case s.MoveTaskLeft:
		return m.moveTaskColumn(false)
	case s.MoveTaskRight:
		return m.moveTaskColumn(true)
	case s.MoveTaskDown:
		return m.moveTaskOrder(false)
	case s.MoveTaskUp:
		return m.moveTaskOrder(true)
	case s.OpenLink:
		selectedItem := m.lists[m.focused].SelectedItem()
		if selectedItem != nil {
			t := selectedItem.(Task)
			links := t.GetLinks()
			if len(links) > 0 {
				OpenURL(links[0])
			}
		}
		return m, nil
	case s.DeleteTask:
		selectedItem := m.lists[m.focused].SelectedItem()
		if selectedItem != nil {
			m.state = deleteState
		}
		return m, nil
	case s.AddTask:
		m.state = formState
		m.editing = false
		m.titleInput.SetValue("")
		m.descInput.SetValue("")
		m.titleInput.Focus()
		return m, tea.Batch(textinput.Blink)
	case s.EditTask:
		selectedItem := m.lists[m.focused].SelectedItem()
		if selectedItem != nil {
			t := selectedItem.(Task)
			m.state = formState
			m.editing = true
			m.titleInput.SetValue(t.title)
			m.descInput.SetValue(t.description)
			m.titleInput.Focus()
			return m, tea.Batch(textinput.Blink)
		}
		return m, nil
	}
	return m, nil
}

func (m model) updateView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "enter":
		m.state = boardState
	case m.config.Shortcuts.OpenLink:
		selectedItem := m.lists[m.focused].SelectedItem()
		if selectedItem != nil {
			t := selectedItem.(Task)
			links := t.GetLinks()
			if len(links) > 0 {
				OpenURL(links[0])
			}
		}
	}
	return m, nil
}

func (m model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.state = boardState
		return m, nil
	case "tab":
		if m.titleInput.Focused() {
			m.titleInput.Blur()
			m.descInput.Focus()
		} else {
			m.descInput.Blur()
			m.titleInput.Focus()
		}
	case "enter":
		if m.titleInput.Focused() {
			m.titleInput.Blur()
			m.descInput.Focus()
			return m, nil
		}
		// Save task
		var newTask Task
		now := time.Now()
		if m.editing {
			oldTask := m.lists[m.focused].SelectedItem().(Task)
			newTask = oldTask
			newTask.title = m.titleInput.Value()
			newTask.description = m.descInput.Value()
			if newTask.status == done && newTask.completedAt == nil {
				newTask.completedAt = &now
			} else if newTask.status != done {
				newTask.completedAt = nil
			}
			SaveTask(newTask)
			m.lists[m.focused].SetItem(m.lists[m.focused].Index(), newTask)
		} else {
			newTask = Task{
				status:      m.focused,
				title:       m.titleInput.Value(),
				description: m.descInput.Value(),
				position:    len(m.lists[m.focused].Items()),
			}
			if newTask.status == done {
				newTask.completedAt = &now
			}
			id, _ := SaveTask(newTask)
			newTask.id = id
			m.lists[m.focused].InsertItem(len(m.lists[m.focused].Items()), newTask)
		}
		m.state = boardState
		return m, nil
	}

	var cmd tea.Cmd
	if m.titleInput.Focused() {
		m.titleInput, cmd = m.titleInput.Update(msg)
	} else {
		m.descInput, cmd = m.descInput.Update(msg)
	}
	return m, cmd
}

func (m model) moveTaskColumn(right bool) (tea.Model, tea.Cmd) {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return m, nil
	}

	t := selectedItem.(Task)
	m.lists[m.focused].RemoveItem(m.lists[m.focused].Index())
	m.syncPositions(m.focused)

	var nextStatus status
	if right {
		nextStatus = (m.focused + 1) % 3
	} else {
		nextStatus = (m.focused + 2) % 3
	}

	t.status = nextStatus
	t.position = len(m.lists[nextStatus].Items())

	if t.status == done {
		now := time.Now()
		t.completedAt = &now
	} else {
		t.completedAt = nil
	}

	SaveTask(t)
	m.lists[nextStatus].InsertItem(len(m.lists[nextStatus].Items()), t)
	m.focused = nextStatus

	return m, nil
}

func (m model) moveTaskOrder(up bool) (tea.Model, tea.Cmd) {
	idx := m.lists[m.focused].Index()
	items := m.lists[m.focused].Items()

	if up && idx > 0 {
		// Swap with previous
		item := items[idx]
		m.lists[m.focused].RemoveItem(idx)
		m.lists[m.focused].InsertItem(idx-1, item)
		m.lists[m.focused].Select(idx - 1)
	} else if !up && idx < len(items)-1 {
		// Swap with next
		item := items[idx]
		m.lists[m.focused].RemoveItem(idx)
		m.lists[m.focused].InsertItem(idx+1, item)
		m.lists[m.focused].Select(idx + 1)
	} else {
		return m, nil
	}

	m.syncPositions(m.focused)
	return m, nil
}

func (m model) syncPositions(s status) {
	for i, item := range m.lists[s].Items() {
		t := item.(Task)
		if t.position != i {
			t.position = i
			SaveTask(t)
			// Update the item in the list too so it's consistent
			m.lists[s].SetItem(i, t)
		}
	}
}

func (m model) updateDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		selectedItem := m.lists[m.focused].SelectedItem()
		if selectedItem != nil {
			t := selectedItem.(Task)
			DeleteTask(t.id)
			m.lists[m.focused].RemoveItem(m.lists[m.focused].Index())
			m.syncPositions(m.focused)
		}
		m.state = boardState
	case "n", "N", "esc":
		m.state = boardState
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}
	if !m.loaded {
		return "Loading..."
	}

	var content string
	if m.state == formState {
		content = m.formView()
	} else if m.state == viewState {
		content = m.taskView()
	} else if m.state == deleteState {
		content = m.deleteView()
	} else {
		var views []string
		columnColors := []string{
			m.config.Colors.TodoColumn,
			m.config.Colors.InProgressColumn,
			m.config.Colors.DoneColumn,
		}

		for i, l := range m.lists {
			style := columnStyle.Copy().BorderForeground(lipgloss.Color(columnColors[i]))
			if i == int(m.focused) {
				style = focusedColumnStyle.Copy().BorderForeground(lipgloss.Color(m.config.Colors.FocusedBorder))
			}
			// Apply column title color by wrapping the title string with style
			originalTitle := [...]string{"To Do", "In Progress", "Done"}[i]
			l.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(columnColors[i])).Bold(true).Render(originalTitle)
			views = append(views, style.Render(l.View()))
		}
		content = lipgloss.JoinHorizontal(lipgloss.Top, views...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, m.helpView())
}

func (m model) helpView() string {
	var help string
	s := m.config.Shortcuts
	if m.state == boardState {
		help = fmt.Sprintf(
			"%s view • %s move • %s move task col • %s move task order • %s add • %s edit • %s delete • %s open link • %s quit",
			helpKeyStyle.Render("Enter"),
			helpKeyStyle.Render(fmt.Sprintf("%s/%s", s.MoveLeft, s.MoveRight)),
			helpKeyStyle.Render(fmt.Sprintf("%s/%s", s.MoveTaskLeft, s.MoveTaskRight)),
			helpKeyStyle.Render(fmt.Sprintf("%s/%s", s.MoveTaskUp, s.MoveTaskDown)),
			helpKeyStyle.Render(s.AddTask),
			helpKeyStyle.Render(s.EditTask),
			helpKeyStyle.Render(s.DeleteTask),
			helpKeyStyle.Render(s.OpenLink),
			helpKeyStyle.Render(s.Quit),
		)
	} else if m.state == formState {
		help = fmt.Sprintf(
			"%s save • %s switch field • %s cancel",
			helpKeyStyle.Render("Enter"),
			helpKeyStyle.Render("Tab"),
			helpKeyStyle.Render("Esc"),
		)
	} else if m.state == deleteState {
		help = fmt.Sprintf(
			"%s confirm delete • %s cancel",
			helpKeyStyle.Render("y/Enter"),
			helpKeyStyle.Render("n/Esc"),
		)
	} else {
		help = fmt.Sprintf(
			"%s back • %s open link",
			helpKeyStyle.Render("Enter/Esc/q"),
			helpKeyStyle.Render(s.OpenLink),
		)
	}
	return helpStyle.Render(help)
}

func (m model) deleteView() string {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return "No task selected"
	}
	t := selectedItem.(Task)

	return lipgloss.NewStyle().Padding(1, 2).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(" Confirm Delete "),
			"",
			fmt.Sprintf("Are you sure you want to delete task: %s?", lipgloss.NewStyle().Bold(true).Render(t.title)),
			"",
			helpDescStyle.Render("This action cannot be undone."),
		),
	)
}

func (m model) taskView() string {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return "No task selected"
	}
	t := selectedItem.(Task)

	return lipgloss.NewStyle().Padding(1, 2).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(" Task Details "),
			"",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252")).Render(t.title),
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render(t.description),
		),
	)
}

func (m model) formView() string {
	title := " Add Task "
	if m.editing {
		title = " Edit Task "
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(title),
			"",
			"Title:",
			m.titleInput.View(),
			"",
			"Description:",
			m.descInput.View(),
		),
	)
}
