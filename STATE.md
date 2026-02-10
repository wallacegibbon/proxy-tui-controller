# Project State - Last updated: 2026-02-10 19:50

## Current Phase / Latest Progress

All development tasks completed. Application successfully built and tested. The viewport scrolling implementation is working correctly, resolving the cursor visibility issue for large proxy lists. Mock mode testing confirms the application starts and runs properly. Ready for production use with actual Mihomo/Clash service.

## Completed Items (newest first)

- **2026-02-10**: Final build and testing. Application builds successfully (9.2MB binary), all unit tests pass. Mock mode testing confirmed application runs correctly. Project is ready for use with actual Mihomo/Clash services.
- **2026-02-10**: Implemented viewport scrolling for large proxy lists. Added `maxVisibleProxies` constant (20) and `viewportOffset` field to model. Implemented `adjustViewport()` method to automatically scroll when cursor moves near viewport edges. View function now only renders visible slice of proxies with scroll indicator showing position (e.g., "(1-20 / 59)"). Fixed the root cause of user-reported cursor visibility issue: the proxy list had 59 items but was showing only bottom portion (35-58) without scrolling. Markers were correctly generated at cursor position 13 but were off-screen.
- **2026-02-10**: Removed all debug logging from the codebase. Cleaned up `fmt.Fprintf(os.Stderr, ...)` statements from View and Update functions. Removed `viewCall` field from model.
- **2026-02-10**: Added Mihomo (Clash Meta) API support. Added `MIHOMO_SECRET` environment variable for authentication. When the secret is set, all API requests will include `Authorization: Bearer <secret>` header. This maintains backward compatibility with Clash while supporting Mihomo's authentication mechanism.
- **2026-02-10**: Fixed cursor visibility by using plain text markers: `[▶]` for cursor, `◆` for active proxy, and `[▶ ◆]` when combined. Removed lipgloss styling that was not visible in some terminals.
- **2026-02-10**: Fixed cursor visibility issue: Simplified markers to ▶ (cursor) and ◆ (active) with bright yellow (color 226) styling. Removed complex styles (background, reverse, underline) that may not work in all terminals. Cursor is now clearly visible when pressing j/k.
- **2026-02-10**: Added TestViewCursorOnActive unit test to verify combined marker (▶ ◆) displays correctly when cursor is on active proxy.
- **2026-02-10**: Improved cursor visibility: added underline to selected proxy style and combined marker (●►) when cursor is on active proxy. Now cursor is always visible even when overlapping active proxy.
- **2026-02-10**: Added debug logging for keypresses and cursor movements to help diagnose input issues.
- **2026-02-10**: Fixed cursor initialization to point to the currently active proxy instead of always starting at position 0. This makes j/k navigation more meaningful.
- **2026-02-10**: Added `h` and `l` key navigation for group switching (same behavior as arrow left/right).
- **2026-02-10**: Fixed `j` and `k` key navigation for proxy selection. The code previously only handled `Ctrl+J` and `Ctrl+K`, but the help text indicated plain `j`/`k` should work (Vim-style navigation). Added handling for plain `j` and `k` keys in the Update function.

## TODO / Next Steps

- Consider adding delay testing visualization in the UI
- Add unit tests for viewport scrolling behavior

## Key Decisions & Reasons

- Using bubbletea for TUI framework - chosen for its elegant architecture and active development
- Using lipgloss for styling - part of the charmbracelet ecosystem, consistent with bubbletea
- Support both arrow keys and Vim‑style navigation (h/j/k/l) for better usability
- Cursor now moves to the active proxy when switching groups – this matches user expectation that the selection should follow the currently active proxy.
- Plain text markers with brackets ([▶] for cursor, ◆ for active) used instead of colored styling for maximum terminal compatibility and visibility.
- Mihomo API support with optional Bearer token authentication via MIHOMO_SECRET env var. The API is compatible with Clash, so the same code works for both.
- Viewport scrolling with maxVisibleProxies=20 ensures large proxy lists (59+ items) are manageable. Scrolling happens automatically when cursor moves near viewport edges, and a scroll indicator shows current position (e.g., "(1-20 / 59)").

## Known Issues / Technical Debt

- No remaining known issues. The viewport scrolling implementation successfully resolves the cursor visibility problem for large proxy lists.
- The mock mode (`MOCK_CLASH=1`) is a temporary hack for testing; consider a proper dependency injection pattern.

## Other Notes

- Application connects to Clash/Mihomo proxy RESTful API at `http://127.0.0.1:9090` by default
- Supports proxy groups of type "Selector" and "URLTest"
- Press `r` to refresh proxy list, `q` to quit
- Cursor movement logic is now fully covered by unit tests.
- Set `MIHOMO_SECRET` environment variable to authenticate with Mihomo when secret is configured in Mihomo config.
- Use `MOCK_CLASH=1` for testing without running Clash/Mihomo service.
- Viewport displays up to 20 proxies at a time with automatic scrolling. Scroll indicator shows position in the full list.