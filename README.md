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

## Roadmap

- [ ] Global search across all category files when no argument is passed (`v0.7.0`)
- [ ] Direct execution profile variables

## Installation & Usage (Development)

Clone the repository:

```bash
git clone [https://github.com/yellowey-com/snip-cli.git](https://github.com/yellowey-com/snip-cli.git)
cd snip-cli
```
