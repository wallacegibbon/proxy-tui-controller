Output demo.gif

Set FontSize 18
Set Width 1200
Set Height 800
Set TypingSpeed 50ms
Set Shell "bash"

# Start with a clean terminal
Type "clear" Enter
Sleep 1s

# Show the command being run
Type "MOCK_CLASH=1 ./proxy-controller-tui" Enter
Sleep 2s

# Navigate through groups (use right arrow)
Type "→"
Sleep 500ms

# Navigate down through proxies (use down arrow)
Type "↓"
Sleep 500ms

# Navigate down again
Type "↓"
Sleep 500ms

# Navigate up
Type "↑"
Sleep 500ms

Sleep 1s

# Go to next group
Type "→"
Sleep 500ms

# Navigate proxies in this group
Type "↓"
Sleep 500ms
Type "↓"
Sleep 500ms

Sleep 1s

# Quit the application
Type "q"

Sleep 500ms
