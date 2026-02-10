# Proxy TUI Controller

Go TUI application for managing Clash/Mihomo proxy services.

## Project
- Connects to Clash/Mihomo RESTful API (`http://127.0.0.1:9090`)
- Supports proxy group selection (Selector and URLTest types)
- Built with bubbletea and lipgloss (charmbracelet)

## Usage
```bash
# With Mihomo secret
MIHOMO_SECRET=YOUR_SECRET ./proxy-tui-controller

# Standard Clash
./proxy-tui-controller

# Mock mode for testing
MOCK_CLASH=1 ./proxy-tui-controller
```

## Controls
- `←/h` / `→/l`: Previous/Next group
- `↑/k` / `↓/j`: Previous/Next proxy
- `Enter`: Select proxy
- `r`: Refresh, `q`: Quit

## Agent Instructions
- **Read STATE.md** at the start of every conversation
- **Update STATE.md** after completing any meaningful work (features, bug fixes, etc.)
- Keep STATE.md as the single source of truth for project status
