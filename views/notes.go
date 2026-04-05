package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/models"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// OpenNotesMsg is sent to navigate to the notes screen
type OpenNotesMsg struct{}

// BackFromNotesMsg is sent when leaving the notes screen
type BackFromNotesMsg struct{}

// OpenNoteComposeMsg is sent to open the note compose/edit screen
type OpenNoteComposeMsg struct {
	Note   models.Note
	IsEdit bool
}

type notesLoadedMsg struct {
	Notes  []models.Note
	Cursor string
}

type moreNotesLoadedMsg struct {
	Notes  []models.Note
	Cursor string
}

type notesErrorMsg struct{ Err error }

type noteDeletedMsg struct{ NoteID string }
type noteDeleteErrMsg struct{ Err error }

// NoteItem implements list.Item for notes
type NoteItem struct{ Note models.Note }

func (n NoteItem) FilterValue() string { return n.Note.Content }
func (n NoteItem) Title() string       { return n.Note.Content }
func (n NoteItem) Description() string { return TimeAgo(n.Note.CreatedAt) }

// noteDelegate renders note items in the list
type noteDelegate struct{}

func (d noteDelegate) Height() int                               { return 3 }
func (d noteDelegate) Spacing() int                              { return 0 }
func (d noteDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d noteDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	switch it := item.(type) {
	case NoteItem:
		isSelected := index == m.Index()

		content := strings.TrimSpace(it.Note.Content)
		content = strings.ReplaceAll(content, "\n", " ")
		if len(content) > 72 {
			content = content[:72] + "…"
		}

		date := TimeAgo(it.Note.CreatedAt)
		var meta string
		if len(it.Note.Topics) > 0 {
			meta = date + "  [" + strings.Join(it.Note.Topics, "] [") + "]"
		} else {
			meta = date
		}

		var contentLine, metaLine string
		if isSelected {
			contentLine = styles.Bright.Render("▸ " + content)
			metaLine = styles.Dim.Render("  " + meta)
		} else {
			contentLine = styles.Normal.Render("  " + content)
			metaLine = styles.Dim.Render("  " + meta)
		}
		fmt.Fprintf(w, "%s\n%s\n", contentLine, metaLine)

	case LoadMoreItem:
		fmt.Fprint(w, zone.Mark("load-more-notes", styles.Dim.Render("  ▼ load more")))
	}
}

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
func NewNotesModel(baseURL, idToken string) NotesModel {
	l := list.New([]list.Item{}, noteDelegate{}, 0, 0)
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
		client:  api.NewClient(baseURL, idToken),
		spinner: NewSpinner(),
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
			return m, func() tea.Msg { return BackFromNotesMsg{} }
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchNotes())
		case key.Matches(msg, m.keys.New):
			return m, func() tea.Msg { return OpenNoteComposeMsg{} }
		case key.Matches(msg, m.keys.Edit):
			if ni, ok := m.list.SelectedItem().(NoteItem); ok {
				note := ni.Note
				return m, func() tea.Msg { return OpenNoteComposeMsg{Note: note, IsEdit: true} }
			}
		case key.Matches(msg, m.keys.Delete):
			if ni, ok := m.list.SelectedItem().(NoteItem); ok {
				m.confirmingDelete = true
				m.deletingNoteID = ni.Note.ID
				return m, nil
			}
		case key.Matches(msg, m.keys.Open):
			switch it := m.list.SelectedItem().(type) {
			case NoteItem:
				note := it.Note
				return m, func() tea.Msg { return OpenNoteComposeMsg{Note: note, IsEdit: true} }
			case LoadMoreItem:
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

	case notesLoadedMsg:
		m.loading = false
		m.deleting = false
		m.err = nil
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		items := notesToItems(msg.Notes)
		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case moreNotesLoadedMsg:
		m.loadingMore = false
		m.nextCursor = msg.Cursor
		m.hasMore = msg.Cursor != ""
		var items []list.Item
		for _, existing := range m.list.Items() {
			if _, ok := existing.(LoadMoreItem); ok {
				continue
			}
			items = append(items, existing)
		}
		for _, n := range msg.Notes {
			items = append(items, NoteItem{Note: n})
		}
		if m.hasMore {
			items = append(items, LoadMoreItem{})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case notesErrorMsg:
		m.loading = false
		m.loadingMore = false
		m.deleting = false
		m.err = msg.Err

	case noteDeletedMsg:
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, m.fetchNotes())

	case noteDeleteErrMsg:
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
	w, h := SafeDimensions(m.width, m.height)

	if m.loading {
		loadingBox := styles.DataBox("ACCESSING PRIVATE NOTES",
			"\n"+
				"  "+m.spinner.View()+styles.Normal.Render(" Loading notes...")+"\n"+
				"\n"+
				"  "+styles.Dim.Render("Retrieving encrypted data...")+"\n",
			50)
		return FullScreen(loadingBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	if m.err != nil {
		errorBox := styles.AlertBox(m.err.Error(), "error", 50) +
			"\n\n" +
			styles.Dim.Render("Press [esc] to go back, [r] to retry")
		return FullScreen(errorBox, w, h, lipgloss.Center, lipgloss.Center)
	}

	var b strings.Builder
	b.WriteString(RenderHeader("▓▒░ PRIVATE NOTES ░▒▓", w))

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
			return notesErrorMsg{Err: err}
		}
		return notesLoadedMsg{Notes: notes, Cursor: cursor}
	}
}

func (m NotesModel) fetchMoreNotes() tea.Cmd {
	return func() tea.Msg {
		notes, cursor, err := m.client.FetchMoreNotes(20, m.nextCursor)
		if err != nil {
			return notesErrorMsg{Err: err}
		}
		return moreNotesLoadedMsg{Notes: notes, Cursor: cursor}
	}
}

func (m NotesModel) deleteNote(noteID string) tea.Cmd {
	return func() tea.Msg {
		if err := m.client.DeleteNote(noteID); err != nil {
			return noteDeleteErrMsg{Err: err}
		}
		return noteDeletedMsg{NoteID: noteID}
	}
}

func notesToItems(notes []models.Note) []list.Item {
	items := make([]list.Item, len(notes))
	for i, n := range notes {
		items[i] = NoteItem{Note: n}
	}
	return items
}
