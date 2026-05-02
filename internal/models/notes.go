package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// NotesModel is the notes list screen
type NotesModel struct {
	list             list.Model
	loading          bool
	loadingMore      bool
	spinner          spinner.Model
	err              error
	client           *api.Client
	nextCursor       string
	hasMore          bool
	width            int
	height           int
	keys             NotesKeyMap
	help             help.Model
	confirmingDelete bool
	deletingNoteID   string
	deleting         bool
}

// NewNotesModel creates a new notes list screen
func NewNotesModel(client *api.Client) NotesModel {
	l := list.New([]list.Item{}, items.NoteDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.Styles = styles.ListStyles()
	l.Paginator.ActiveDot = styles.Bright.Render("▄")
	l.Paginator.InactiveDot = styles.Dark.Render("▄")
	l.KeyMap.Quit.SetEnabled(false)
	l.KeyMap.ForceQuit.SetEnabled(false)

	h := help.New()
	h.Styles = styles.HelpStyles()

	return NotesModel{
		list:    l,
		client:  client,
		spinner: items.NewSpinner(),
		loading: true,
		keys:    NewNotesKeyMap(),
		help:    h,
	}
}

func (m NotesModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchNotes())
}

func (m NotesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				m.confirmingDelete = false
				m.deleting = true
				return m, tea.Batch(m.spinner.Tick, m.deleteNote(m.deletingNoteID))
			default:
				m.confirmingDelete = false
				m.deletingNoteID = ""
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg { return messages.SwitchToFeed{} }
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchNotes())
		case key.Matches(msg, m.keys.New):
			return m, func() tea.Msg { return messages.SwitchToNoteCompose{} }
		case key.Matches(msg, m.keys.Edit):
			if ni, ok := m.list.SelectedItem().(items.NoteItem); ok {
				note := ni.Note
				return m, func() tea.Msg { return messages.SwitchToNoteCompose{Note: note, IsEdit: true} }
			}
		case key.Matches(msg, m.keys.Delete):
			if ni, ok := m.list.SelectedItem().(items.NoteItem); ok {
				m.confirmingDelete = true
				m.deletingNoteID = ni.Note.ID
				return m, nil
			}
		case key.Matches(msg, m.keys.Open):
			switch it := m.list.SelectedItem().(type) {
			case items.NoteItem:
				note := it.Note
				return m, func() tea.Msg { return messages.SwitchToNoteCompose{Note: note, IsEdit: true} }
			case items.LoadMoreItem:
				if !m.loadingMore {
					m.loadingMore = true
					return m, tea.Batch(m.spinner.Tick, m.fetchMoreNotes())
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && !m.loading {
			if zone.Get("load-more-notes").InBounds(msg) && m.hasMore && !m.loadingMore {
				m.loadingMore = true
				return m, tea.Batch(m.spinner.Tick, m.fetchMoreNotes())
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case messages.NotesLoadedMsg:
		m.loading = false
		m.loadingMore = false
		m.deleting = false
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		m.err = nil

		localItems := items.BuildListItems(
			&m.list, msg.IsAdditional, m.hasMore,
			items.NotesToItems(msg.Notes),
		)
		return m, m.list.SetItems(localItems)

	case messages.NotesLoadedErrMsg:
		m.loading = false
		m.loadingMore = false
		m.deleting = false
		m.err = msg.Err

	case messages.NoteDeleteMsg:
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, m.fetchNotes())

	case messages.NoteDeleteErrMsg:
		m.deleting = false
		m.err = msg.Err

	case ThemeChangedMsg:
		m.spinner.Style = styles.Spinner
		m.list.Styles = styles.ListStyles()
		m.help.Styles = styles.HelpStyles()
		m.list.Paginator.ActiveDot = styles.Bright.Render("▄")
		m.list.Paginator.InactiveDot = styles.Dark.Render("▄")
	}

	if !m.loading {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m NotesModel) View() string {
	w, h := items.SafeDimensions(m.width, m.height)

	if m.loading {
		loadingBox := styles.DataBox("ACCESSING PRIVATE NOTES",
			"\n"+
				"  "+m.spinner.View()+styles.Normal.Render(" Loading notes...")+"\n"+
				"\n"+
				"  "+styles.Dim.Render("Retrieving encrypted data...")+"\n",
			50)
		return items.FullScreen(loadingBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	if m.err != nil {
		errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
			"\n\n" +
			styles.Dim.Render("Press [esc] to go back, [r] to retry")
		return items.FullScreen(errorBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	var b strings.Builder
	b.WriteString(items.RenderHeader("▓▒░ PRIVATE NOTES ░▒▓", w))

	noteCount := len(m.list.Items())
	label := fmt.Sprintf("  %d notes", noteCount)
	if m.hasMore {
		label = fmt.Sprintf("  %d+ notes", noteCount)
	}
	b.WriteString(styles.Dim.Render(label) + "\n")

	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(m.renderFooter(w))

	return b.String()
}

func (m NotesModel) renderFooter(width int) string {
	helpView := m.help.View(m.keys)
	helpWidth := lipgloss.Width(helpView)

	var status string
	if m.confirmingDelete {
		status = styles.Error.Render(" delete note? [y/n]")
	} else if m.deleting {
		status = styles.Dim.Render(" [deleting...]")
	}
	statusWidth := lipgloss.Width(status)

	paginatorView := m.list.Paginator.View()
	paginatorWidth := lipgloss.Width(paginatorView)

	dividerWidth := width - helpWidth - statusWidth - paginatorWidth - 2
	if dividerWidth < 1 {
		dividerWidth = 1
	}

	return helpView + status + " " + styles.Divider(dividerWidth) + " " + paginatorView
}

// SetSize updates the view dimensions
func (m *NotesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

func (m NotesModel) fetchNotes() tea.Cmd {
	return func() tea.Msg {
		notes, cursor, err := m.client.FetchNotes(20)
		if err != nil {
			return messages.NotesLoadedErrMsg{Err: err}
		}
		return messages.NotesLoadedMsg{Notes: notes, Cursor: cursor}
	}
}

func (m NotesModel) fetchMoreNotes() tea.Cmd {
	return func() tea.Msg {
		notes, cursor, err := m.client.FetchMoreNotes(20, m.nextCursor)
		if err != nil {
			return messages.NotesLoadedErrMsg{Err: err}
		}
		return messages.NotesLoadedMsg{Notes: notes, Cursor: cursor, IsAdditional: true}
	}
}

func (m NotesModel) deleteNote(noteID string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.DeleteNote(noteID); err != nil {
			return messages.NoteDeleteErrMsg{Err: err}
		}
		return messages.NoteDeleteMsg{NoteID: noteID}
	}
}
