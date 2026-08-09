package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/unremarkablegarden/cyberspace-tui-go/internal/entities"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/api"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/external/cache"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/messages"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/items"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/models/keymaps"
	"github.com/unremarkablegarden/cyberspace-tui-go/internal/ui"
	"github.com/unremarkablegarden/cyberspace-tui-go/styles"
)

// NotesModel is the notes list screen
type NotesModel struct {
	list             list.Model
	loading          bool
	loadingMore      bool
	spinner          *spinner.Model
	err              error
	client           *api.Client
	cache            cache.ICache
	nextCursor       string
	hasMore          bool
	width            int
	height           int
	keys             keymaps.AppKeybinds
	help             help.Model
	confirmingDelete bool
	deletingNoteID   string
	deleting         bool
}

// NewNotesModel creates a new notes list screen
func NewNotesModel(client *api.Client, cache cache.ICache, keymap keymaps.AppKeybinds, sp *spinner.Model) NotesModel {
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
		cache:   cache,
		spinner: sp,
		loading: true,
		keys:    keymap,
		help:    h,
	}
}

func (m NotesModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetchNotes(false))
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

		switch msg.String() {
		case m.keys.GlobalKeybinds.Quit:
			return m, tea.Quit
		case m.keys.GlobalKeybinds.Help:
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case m.keys.GlobalKeybinds.Back:
			return m, func() tea.Msg { return messages.SwitchToFeed{} }
		case m.keys.GlobalKeybinds.Refresh:
			m.loading = true
			m.err = nil
			return m, tea.Batch(m.spinner.Tick, m.fetchNotes(true))
		case m.keys.GlobalKeybinds.Open:
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

		case m.keys.NotesKeybinds.New:
			return m, func() tea.Msg { return messages.SwitchToNoteCompose{} }
		case m.keys.NotesKeybinds.Edit:
			if ni, ok := m.list.SelectedItem().(items.NoteItem); ok {
				note := ni.Note
				return m, func() tea.Msg { return messages.SwitchToNoteCompose{Note: note, IsEdit: true} }
			}
		case m.keys.NotesKeybinds.Delete:
			if ni, ok := m.list.SelectedItem().(items.NoteItem); ok {
				m.confirmingDelete = true
				m.deletingNoteID = ni.Note.ID
				return m, nil
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
		return m, tea.Batch(m.spinner.Tick, m.fetchNotes(true))

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
	w, h := ui.SafeDimensions(m.width, m.height)

	if m.loading {
		return ui.RenderLoadingScreen(
			ui.LoadingScreenTexts{
				TitleText:    "ACCESSING PRIVATE NOTES",
				SubtitleText: " Loading notes...",
				BottomText:   "Retrieving encrypted data...",
			},
			*m.spinner,
			w, h,
		)
	}

	if m.err != nil {
		return ui.RenderErrorScreen(m.err, w, h)
	}

	var b strings.Builder
	b.WriteString(ui.RenderHeader("▓▒░ PRIVATE NOTES ░▒▓", w))

	noteCount := len(m.list.Items())
	label := fmt.Sprintf("  %d notes", noteCount)
	if m.hasMore {
		label = fmt.Sprintf("  %d+ notes", noteCount)
	}
	b.WriteString(styles.Dim.Render(label) + "\n")

	b.WriteString(m.list.View())
	b.WriteString("\n")
	b.WriteString(ui.RenderFooterWithList(
		m.help.View(m.keys.NotesHelpKeys()),
		m.list.Paginator.View(),
		w,
	))

	return b.String()
}

// SetSize updates the view dimensions
func (m *NotesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4)
}

func (m NotesModel) fetchNotes(isRefresh bool) tea.Cmd {
	return func() tea.Msg {
		if isRefresh {
			return m.notesFromAPI()
		}

		cacheNotes, notesFound := m.cache.Get(cache.DefaultNotesCacheKey + "notes")
		cacheCursor, cursorFound := m.cache.Get(cache.DefaultNotesCacheKey + "cursor")

		if !notesFound || !cursorFound {
			return m.notesFromAPI()
		}

		notes, ok := cacheNotes.([]entities.Note)
		if !ok {
			return m.notesFromAPI()
		}

		cursor, ok := cacheCursor.(string)
		if !ok {
			return m.notesFromAPI()
		}

		return messages.NotesLoadedMsg{
			Notes:  notes,
			Cursor: cursor,
		}
	}
}

func (m NotesModel) notesFromAPI() tea.Msg {
	n, c, err := m.syncNotesFromAPI()
	if err != nil {
		return messages.NotesLoadedErrMsg{Err: err}
	}
	return messages.NotesLoadedMsg{Notes: n, Cursor: c}
}

func (m NotesModel) syncNotesFromAPI() ([]entities.Note, string, error) {
	notes, cursor, err := m.client.FetchNotes(20)
	if err != nil {
		return []entities.Note{}, "", err
	}

	m.cache.Set(cache.DefaultNotesCacheKey+"notes", notes, 0)
	m.cache.Set(cache.DefaultNotesCacheKey+"cursor", cursor, 0)

	return notes, cursor, nil
}

func (m NotesModel) fetchMoreNotes() tea.Cmd {
	return func() tea.Msg {
		var notes []entities.Note
		if n, found := m.cache.Get(cache.DefaultNotesCacheKey + "notes"); !found {
			apiNotes, apiCursor, err := m.syncNotesFromAPI()
			if err != nil {
				return messages.NotesLoadedErrMsg{Err: err}
			}

			m.nextCursor = apiCursor
			notes = apiNotes
		} else {
			notes = n.([]entities.Note)
		}

		additionalNotes, cursor, err := m.client.FetchMoreNotes(20, m.nextCursor)
		if err != nil {
			return messages.NotesLoadedErrMsg{Err: err}
		}

		m.cache.Set(cache.DefaultNotesCacheKey+"notes", append(notes, additionalNotes...), 0)
		m.cache.Set(cache.DefaultNotesCacheKey+"cursor", cursor, 0)

		return messages.NotesLoadedMsg{
			Notes:        additionalNotes,
			Cursor:       cursor,
			IsAdditional: true,
		}
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
