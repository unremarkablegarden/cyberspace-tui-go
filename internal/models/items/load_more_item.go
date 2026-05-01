package items

// ═══════════════════════════════════════════════════════════════════════════════
// LOAD MORE ITEM — sentinel item at the bottom of the list
// ═══════════════════════════════════════════════════════════════════════════════

// LoadMoreItem is a sentinel list item for triggering pagination.
type LoadMoreItem struct{}

func (l LoadMoreItem) FilterValue() string { return "" }
func (l LoadMoreItem) Title() string       { return "LOAD MORE" }
func (l LoadMoreItem) Description() string { return "" }
