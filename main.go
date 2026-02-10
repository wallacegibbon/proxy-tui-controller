package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error
type proxiesLoadedMsg struct {
	proxies map[string]Proxy
	groups  []string
}

const (
	maxVisibleProxies = 20
)

var (
	headerStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true)
	selectedGroupStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("255")).Bold(true)
	normalStyle        = lipgloss.NewStyle()
	helpStyle          = lipgloss.NewStyle().Faint(true)
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type model struct {
	client         *ClashClient
	proxies        map[string]Proxy
	groups         []string
	currentIdx     int
	cursor         int
	loading        bool
	err            error
	viewportOffset int
}

func initialModel() model {
	client := NewClashClient("")
	return model{
		client:         client,
		proxies:        make(map[string]Proxy),
		groups:         make([]string, 0),
		currentIdx:     0,
		cursor:         0,
		loading:        true,
		err:            nil,
		viewportOffset: 0,
	}
}

func loadProxiesCmd(client *ClashClient) tea.Cmd {
	return func() tea.Msg {
		proxies, err := client.GetProxies()
		if err != nil {
			return errMsg(err)
		}

		groups := make([]string, 0)
		for name, proxy := range proxies.Proxies {
			if proxy.Type == "Selector" || proxy.Type == "URLTest" {
				groups = append(groups, name)
			}
		}

		return proxiesLoadedMsg{
			proxies: proxies.Proxies,
			groups:  groups,
		}
	}
}

func (m model) Init() tea.Cmd {
	return loadProxiesCmd(m.client)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.loading = false
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		// Adjust viewport offset if needed when window resizes
		m.adjustViewport()
		return m, nil

	case proxiesLoadedMsg:
		m.loading = false
		m.proxies = msg.proxies
		m.groups = msg.groups
		m.currentIdx = 0
		// Set cursor to the active proxy in the first group
		if len(m.groups) > 0 {
			if proxy, ok := m.proxies[m.groups[0]]; ok {
				cursorFound := false
				for i, p := range proxy.All {
					if p == proxy.Now {
						m.cursor = i
						cursorFound = true
						break
					}
				}
				if !cursorFound && len(proxy.All) > 0 {
					m.cursor = 0
				}
			}
		}
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.Type {
		case tea.KeyUp, tea.KeyCtrlK:
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor > 0 {
						m.cursor--
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case tea.KeyDown, tea.KeyCtrlJ:
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor < len(proxy.All)-1 {
						m.cursor++
						m.adjustViewport()
					}
				}
			}
			return m, nil

		case tea.KeyLeft:
			if m.currentIdx > 0 {
				m.currentIdx--
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil

		case tea.KeyRight:
			if m.currentIdx < len(m.groups)-1 {
				m.currentIdx++
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil

		case tea.KeyEnter:
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && m.cursor < len(proxy.All) {
					selectedProxy := proxy.All[m.cursor]
					if err := m.client.SelectProxy(group, selectedProxy); err != nil {
						m.err = err
						return m, nil
					}
					return m, loadProxiesCmd(m.client)
				}
			}
			return m, nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}

		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, loadProxiesCmd(m.client)
		case "h":
			if m.currentIdx > 0 {
				m.currentIdx--
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil
		case "l":
			if m.currentIdx < len(m.groups)-1 {
				m.currentIdx++
				// Set cursor to active proxy in the new group
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok {
					for i, p := range proxy.All {
						if p == proxy.Now {
							m.cursor = i
							break
						}
					}
					m.viewportOffset = 0
				} else {
					m.cursor = 0
				}
			}
			return m, nil
		case "k":
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor > 0 {
						m.cursor--
						m.adjustViewport()
					}
				}
			}
			return m, nil
		case "j":
			if m.currentIdx < len(m.groups) {
				group := m.groups[m.currentIdx]
				if proxy, ok := m.proxies[group]; ok && len(proxy.All) > 0 {
					if m.cursor < len(proxy.All)-1 {
						m.cursor++
						m.adjustViewport()
					}
				}
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *model) adjustViewport() {
	if len(m.groups) == 0 {
		return
	}
	group := m.groups[m.currentIdx]
	proxy, ok := m.proxies[group]
	if !ok {
		return
	}

	// Ensure cursor is visible within viewport
	if m.cursor < m.viewportOffset {
		m.viewportOffset = m.cursor
	} else if m.cursor >= m.viewportOffset+maxVisibleProxies {
		m.viewportOffset = m.cursor - maxVisibleProxies + 1
	}

	// Ensure offset doesn't go negative
	if m.viewportOffset < 0 {
		m.viewportOffset = 0
	}

	// Ensure offset doesn't exceed list length
	maxOffset := len(proxy.All) - maxVisibleProxies
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.viewportOffset > maxOffset {
		m.viewportOffset = maxOffset
	}
}

func (m model) View() string {

	if m.loading {
		return "Loading proxies from Clash...\n"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'r' to retry, 'q' to quit\n", m.err)
	}

	if len(m.groups) == 0 {
		return "No proxy groups found. Press 'r' to refresh, 'q' to quit\n"
	}

	var s string
	s += headerStyle.Render("Clash Proxy Controller") + "\n\n"

	for i, group := range m.groups {
		proxy, ok := m.proxies[group]
		if !ok {
			continue
		}

		if i == m.currentIdx {
			s += selectedGroupStyle.Render(fmt.Sprintf("▸ %s", group)) + "\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("  %s", group)) + "\n"
		}

		if i == m.currentIdx {
			// Only render visible proxies based on viewport offset
			visibleProxies := proxy.All
			if len(proxy.All) > maxVisibleProxies {
				// Calculate visible range
				startIdx := m.viewportOffset
				if startIdx < 0 {
					startIdx = 0
				}
				endIdx := startIdx + maxVisibleProxies
				if endIdx > len(proxy.All) {
					endIdx = len(proxy.All)
				}
				visibleProxies = proxy.All[startIdx:endIdx]
			}

			for j, p := range visibleProxies {
				actualIdx := j + m.viewportOffset
				if actualIdx == m.cursor && p == proxy.Now {
					s += "  [▶ ◆] " + p + "\n"
				} else if actualIdx == m.cursor {
					s += "  [▶] " + p + "\n"
				} else if p == proxy.Now {
					s += "     ◆ " + p + "\n"
				} else {
					s += "        " + p + "\n"
				}
			}

			// Show scroll indicator if there are more proxies
			if len(proxy.All) > maxVisibleProxies {
				s += helpStyle.Render(fmt.Sprintf("    (%d-%d / %d)", m.viewportOffset+1,
					min(m.viewportOffset+maxVisibleProxies, len(proxy.All)), len(proxy.All))) + "\n"
			}
		}

		if i < len(m.groups)-1 {
			s += "\n"
		}
	}

	s += "\n"
	s += helpStyle.Render("Controls:\n")
	s += helpStyle.Render("  ←/h : 上一个组     →/l : 下一个组\n")
	s += helpStyle.Render("  ↑/k : 上一个代理   ↓/j : 下一个代理\n")
	s += helpStyle.Render("  Enter: 选择代理      r : 刷新\n")
	s += helpStyle.Render("  q : 退出\n")

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
