```
 ██████╗██╗   ██╗██████╗ ███████╗██████╗ ███████╗██████╗  █████╗  ██████╗███████╗
██╔════╝╚██╗ ██╔╝██╔══██╗██╔════╝██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██║      ╚████╔╝ ██████╔╝█████╗  ██████╔╝███████╗██████╔╝███████║██║     █████╗
██║       ╚██╔╝  ██╔══██╗██╔══╝  ██╔══██╗╚════██║██╔═══╝ ██╔══██║██║     ██╔══╝
╚██████╗   ██║   ██████╔╝███████╗██║  ██║███████║██║     ██║  ██║╚██████╗███████╗
 ╚═════╝   ╚═╝   ╚═════╝ ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝
```
## [NEW OFFICAL TUI LIVES HERE](https://github.com/unremarkablegarden/cyberspace-opentui)

The OLD terminal client for [Cyberspace](https://cyberspace.online/) built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Browse the feed, read posts and replies, and switch themes.

## Features

- Browse the Cyberspace feed with cursor-based pagination
- View posts and threaded replies
- 5 built-in themes (Dark, Amber VT320, Matrix, Crypt, Brutalist)
- Live theme preview with instant switching
- Vim-style navigation throughout
- Mouse support (click to open posts, scroll)
- Responsive layout (adapts to terminal size)
- Token persistence with automatic refresh

## Building from source

Requires **Go 1.24+**.

```bash
go build -o cyberspace-cli .
./cyberspace-cli
```

Or use the included build script:

```bash
./build.sh
```

## Configuration

Auth tokens and theme preference are saved to `~/.cyberspace-cli.json`.

## Keybindings

### Feed

| Key | Action |
|-----|--------|
| `j` / `k` | Move down / up |
| `g` / `G` | Jump to top / bottom |
| `Enter` | Open post |
| `r` | Refresh feed |
| `t` | Switch theme |
| `L` | Log out |
| `?` | Toggle full help |
| `q` | Quit |

### Post detail

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll down / up |
| `f` / `pgdn` | Page down |
| `u` / `d` | Half-page up / down |
| `g` / `G` | Jump to top / bottom |
| `c` | Compose reply |
| `r` | Refresh |
| `esc` / `b` | Back to feed |
| `q` | Quit |

## Project structure

```
.
├── main.go              # App entrypoint, state machine, screen routing
├── config.go            # Auth token persistence and refresh logic
├── api/
│   ├── auth.go          # Login, token refresh, HTTP helpers
│   └── posts.go         # Posts and replies CRUD
├── models/
│   └── post.go          # Post, Reply, Attachment data types
├── styles/
│   └── theme.go         # Theme engine, color palette, style helpers
├── views/
│   ├── feed.go          # Feed screen (list + pagination)
│   ├── post_detail.go   # Post viewer + reply composer
│   ├── post_item.go     # Card rendering for feed items
│   ├── login.go         # Login screen
│   ├── theme_switcher.go # Theme selection modal
│   ├── keys.go          # Keybinding definitions
│   └── helpers.go       # Text utilities, layout helpers
└── themes/              # JSON theme definitions
    ├── dark.json
    ├── amber.json
    ├── matrix.json
    ├── red.json
    └── brutalist.json
```

## Requirements

- A [Cyberspace](https://cyberspace.online/) account
- Terminal with 256-color support
