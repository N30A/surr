package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type loadSinksMsg []listItem

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type listItem struct {
	sink sink
}

func (i listItem) Title() string { return i.sink.friendlyName }
func (i listItem) Description() string {
	if i.sink.isDefault {
		return "Selected"
	}
	return "Not selected"
}
func (i listItem) FilterValue() string { return i.sink.friendlyName }

type model struct {
	list list.Model
	err  error
}

func newModel() model {
	model := model{list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)}
	model.list.Title = "Output devices"
	return model
}

func (m model) Init() tea.Cmd {
	if err := ensurePactlOrExit(); err != nil {
		return func() tea.Msg {
			return errMsg{err}
		}
	}
	return loadSinks()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedItem, ok := m.list.SelectedItem().(listItem)
			if ok {
				return m, setDefaultSink(selectedItem.sink)
			}
			return m, loadSinks()
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case loadSinksMsg:
		var items []list.Item
		for _, item := range msg {
			items = append(items, item)
		}
		m.list.SetItems(items)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nOops! Something went wrong: %v\n\n[q, escape, ctrl-c] to exit", m.err)
	}
	return docStyle.Render(m.list.View())
}

func loadSinks() tea.Cmd {
	return func() tea.Msg {
		sinks, err := sinks()
		if err != nil {
			return errMsg{err}
		}

		items := make([]listItem, 0, len(sinks))
		for _, sink := range sinks {
			items = append(items, listItem{sink: sink})
		}
		return loadSinksMsg(items)
	}
}

func setDefaultSink(sink sink) tea.Cmd {
	return func() tea.Msg {
		if err := sink.setDefault(); err != nil {
			return errMsg{err}
		}
		return loadSinks()()
	}
}
