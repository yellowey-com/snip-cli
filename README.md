# snip-cli

A fast, lightweight CLI snippet and cheat sheet manager written in Go. Keep your frequently used commands, configurations, and notes organized right in your terminal.

## Current Progress

- [x] Basic project scaffolding and modular structure (`v0.1.0`)
- [x] Abstract file storage interaction layer (`v0.1.0`)
- [x] Custom high-performance Markdown syntax parser using `strings.CutPrefix` (`v0.2.0`)
- [x] Transforming raw files into structured data slices in memory (`v0.2.0`)
- [x] Interactive terminal UI (Fuzzy Search) using Bubble Tea (`v0.3.0`)
- [x] System clipboard integration (copy on Enter) (`v0.4.0`)
- [x] Built-in command execution (`v0.5.0`)
- [x] Dynamic snippet management (`add`, `edit`, `remove`) (`v0.6.0`)
- [x] Global search across all category files when no argument is passed (`v0.7.0`)
- [x] Hide all default shortcuts and help footer from the UI (`v0.7.1`)
- [x] Interactive add/edit/delete snippets via in-picker modal forms (`v0.8.0`)
- [x] Centralized keybindings + contextual help footer (`v0.8.1`)
- [x] Fuzzy search directly from terminal arguments via `snip find "query"` (`v0.9.0`)
- [x] Execute snippets directly from CLI via `snip run <description>` (`v0.10.0`)
- [x] Full native Zsh autocompletion with multi-word description support (`v0.11.0`)
- [x] Inline placeholders to inject dynamic values before running a command (`v0.12.0`)
- [x] Clean up and refactor command line argument parsing (`v0.12.1`)
- [x] Redesign interactive TUI into a minimalist, boxed layout (`v0.12.2`)
- [x] Implement CLI command to rebuild and restart completion cache (`v0.13.0`)
- [x] Standardize project structure to `cmd/snip` for idiomatic Go distribution (`v0.14.0`)
- [x] Bash and Fish shell autocompletion support (`v0.15.0`)

## Roadmap

- [ ] Direct execution profile variables
- [ ] Auto sync snippets with a remote git repository
- [ ] Export snippets to JSON/YAML format and import from standard Cheat sheets
- [ ] Configurable keybindings and theme customization via `config.toml`
- [ ] Grouping and execution filtering by user-defined tags/labels
- [ ] Stat tracking (most used snippets sorted at the top)
- [ ] Custom validation hooks for inputs in the modal forms

## Installation & Usage (Development)

Clone the repository:

```bash
git clone [https://github.com/yellowey-com/snip-cli.git](https://github.com/yellowey-com/snip-cli.git)
cd snip-cli
```
