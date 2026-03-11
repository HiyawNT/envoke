package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/HiyawNT/envoke/internal/models"
	"github.com/HiyawNT/envoke/internal/service"
	"github.com/catppuccin/go"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Catppuccin theme - using Mocha variant (dark theme)
var (
	theme = catppuccingo.Mocha

	// Base colors
	colorBase     = lipgloss.Color(theme.Base().Hex)
	colorSurface0 = lipgloss.Color(theme.Surface0().Hex)
	colorSurface1 = lipgloss.Color(theme.Surface1().Hex)
	colorSurface2 = lipgloss.Color(theme.Surface2().Hex)
	colorOverlay0 = lipgloss.Color(theme.Overlay0().Hex)
	colorOverlay1 = lipgloss.Color(theme.Overlay1().Hex)

	// Text colors
	colorText     = lipgloss.Color(theme.Text().Hex)
	colorSubtext0 = lipgloss.Color(theme.Subtext0().Hex)
	colorSubtext1 = lipgloss.Color(theme.Subtext1().Hex)

	// Accent colors
	colorMauve    = lipgloss.Color(theme.Mauve().Hex)
	colorRed      = lipgloss.Color(theme.Red().Hex)
	colorPeach    = lipgloss.Color(theme.Peach().Hex)
	colorYellow   = lipgloss.Color(theme.Yellow().Hex)
	colorGreen    = lipgloss.Color(theme.Green().Hex)
	colorTeal     = lipgloss.Color(theme.Teal().Hex)
	colorSky      = lipgloss.Color(theme.Sky().Hex)
	colorBlue     = lipgloss.Color(theme.Blue().Hex)
	colorLavender = lipgloss.Color(theme.Lavender().Hex)

	// Component Styles
	appTitleStyle = lipgloss.NewStyle().
			Foreground(colorMauve).
			Background(colorSurface0).
			Bold(true).
			Padding(0, 2).
			MarginBottom(0)

	headerStyle = lipgloss.NewStyle().
			Background(colorSurface0).
			Foreground(colorText).
			Padding(1, 2).
			Bold(true)

	panelStyle = lipgloss.NewStyle().
			Background(colorBase).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSurface2).
			Padding(1, 2).
			MarginRight(1)

	activePanelStyle = lipgloss.NewStyle().
				Background(colorBase).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorMauve).
				Padding(1, 2).
				MarginRight(1)

	panelHeaderStyle = lipgloss.NewStyle().
				Foreground(colorLavender).
				Bold(true).
				Underline(true)

	envItemStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Padding(0, 1)

	envSelectedStyle = lipgloss.NewStyle().
				Background(colorSurface1).
				Foreground(colorMauve).
				Bold(true).
				Padding(0, 1)

	envHoverStyle = lipgloss.NewStyle().
			Background(colorSurface0).
			Foreground(colorText).
			Padding(0, 1)

	activeTagStyle = lipgloss.NewStyle().
			Foreground(colorBase).
			Background(colorGreen).
			Padding(0, 1).
			Bold(true).
			MarginLeft(1)

	secretKeyStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	secretValueStyle = lipgloss.NewStyle().
				Foreground(colorText)

	maskedValueStyle = lipgloss.NewStyle().
				Foreground(colorOverlay0)

	statLabelStyle = lipgloss.NewStyle().
			Foreground(colorSubtext1)

	statValueStyle = lipgloss.NewStyle().
			Foreground(colorPeach).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorSubtext0).
			Background(colorSurface0).
			Padding(0, 1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Background(colorSurface1).
			Padding(0, 1).
			Bold(true).
			MarginRight(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorBase).
			Background(colorRed).
			Bold(true).
			Padding(0, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(colorBase).
			Background(colorGreen).
			Bold(true).
			Padding(0, 2)

	warningStyle = lipgloss.NewStyle().
			Foreground(colorBase).
			Background(colorYellow).
			Bold(true).
			Padding(0, 2)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(colorLavender).
			Bold(true).
			Padding(0, 0, 0, 2)

	modalStyle = lipgloss.NewStyle().
			Background(colorSurface0).
			Border(lipgloss.ThickBorder()).
			BorderForeground(colorMauve).
			Padding(2, 4).
			MarginTop(2).
			MarginBottom(2)

	modalTitleStyle = lipgloss.NewStyle().
			Foreground(colorMauve).
			Background(colorSurface1).
			Bold(true).
			Padding(0, 2).
			Align(lipgloss.Center)

	searchActiveStyle = lipgloss.NewStyle().
				Foreground(colorBase).
				Background(colorSky).
				Padding(0, 1).
				Bold(true)

	searchInactiveStyle = lipgloss.NewStyle().
				Foreground(colorSubtext1).
				Background(colorSurface1).
				Padding(0, 1)

	dividerStyle = lipgloss.NewStyle().
			Foreground(colorSurface2)

	scrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(colorOverlay1)

	// Special indicator styles
	indicatorActiveStyle = lipgloss.NewStyle().
				Foreground(colorMauve).
				Bold(true)

	indicatorInactiveStyle = lipgloss.NewStyle().
				Foreground(colorOverlay0)

	countBadgeStyle = lipgloss.NewStyle().
			Foreground(colorBase).
			Background(colorTeal).
			Padding(0, 1).
			Bold(true)

	emptyStateStyle = lipgloss.NewStyle().
			Foreground(colorSubtext0).
			Italic(true).
			Padding(2, 0)
)

type panel int

const (
	panelEnvironments panel = iota
	panelSecrets
)

type view int

const (
	viewMain view = iota
	viewAddSecret
	viewEditSecret
	viewDeleteConfirm
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Back     key.Binding
	Add      key.Binding
	Delete   key.Binding
	Toggle   key.Binding
	Edit     key.Binding
	Search   key.Binding
	Tab      key.Binding
	Quit     key.Binding
	Refresh  key.Binding
	PageUp   key.Binding
	PageDown key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left panel"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right panel"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch field"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("pgdn", "page down"),
		),
	}
}

type Model struct {
	service            *service.SecretService
	currentView        view
	activePanel        panel
	environments       []models.Environment
	currentEnvIndex    int
	envScrollOffset    int
	secrets            []secretItem
	selectedSecret     int
	secretScrollOffset int
	showValues         bool
	searchMode         bool
	searchQuery        string
	filteredSecrets    []secretItem
	width              int
	height             int
	keys               keyMap
	err                error
	successMsg         string
	statusTimeout      time.Time

	// Input fields
	keyInput     textinput.Model
	valueInput   textinput.Model
	searchInput  textinput.Model
	inputMode    string // "key" or "value"
	editingKey   string
	deleteTarget string
}

type secretItem struct {
	Key       string
	Value     string
	Masked    string
	CreatedAt string
}

type tickMsg time.Time

func NewModel(service *service.SecretService) Model {
	// Customize text input styles with Catppuccin
	keyInput := textinput.New()
	keyInput.Placeholder = "SECRET_KEY"
	keyInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorOverlay0)
	keyInput.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	keyInput.Cursor.Style = lipgloss.NewStyle().Foreground(colorMauve)
	keyInput.CharLimit = 128
	keyInput.Width = 40

	valueInput := textinput.New()
	valueInput.Placeholder = "secret value"
	valueInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorOverlay0)
	valueInput.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	valueInput.Cursor.Style = lipgloss.NewStyle().Foreground(colorMauve)
	valueInput.CharLimit = 1024
	valueInput.EchoMode = textinput.EchoPassword
	valueInput.EchoCharacter = '•'
	valueInput.Width = 40

	searchInput := textinput.New()
	searchInput.Placeholder = "Search secrets..."
	searchInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(colorOverlay0)
	searchInput.TextStyle = lipgloss.NewStyle().Foreground(colorText)
	searchInput.Cursor.Style = lipgloss.NewStyle().Foreground(colorSky)
	searchInput.CharLimit = 100
	searchInput.Width = 30

	return Model{
		service:     service,
		currentView: viewMain,
		activePanel: panelEnvironments,
		keys:        newKeyMap(),
		keyInput:    keyInput,
		valueInput:  valueInput,
		searchInput: searchInput,
		showValues:  false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadEnvironments,
		textinput.Blink,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) loadEnvironments() tea.Msg {
	envs, err := m.service.ListEnvironments()
	if err != nil {
		return errMsg{err}
	}
	return envsLoadedMsg{envs}
}

func (m Model) loadSecrets() tea.Msg {
	if len(m.environments) == 0 {
		return secretsLoadedMsg{nil}
	}

	env := m.environments[m.currentEnvIndex]
	secrets, err := m.service.ListSecretsWithValues(env.Name)
	if err != nil {
		return errMsg{err}
	}

	items := make([]secretItem, 0, len(secrets))
	for key, value := range secrets {
		masked := strings.Repeat("•", min(len(value), 20))
		if len(masked) == 0 {
			masked = "••••••••"
		}
		items = append(items, secretItem{
			Key:    key,
			Value:  value,
			Masked: masked,
		})
	}

	return secretsLoadedMsg{items}
}

type envsLoadedMsg struct {
	envs []models.Environment
}

type secretsLoadedMsg struct {
	secrets []secretItem
}

type errMsg struct {
	err error
}

type successMsg struct {
	message string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Clear status messages after timeout
	if !m.statusTimeout.IsZero() && time.Now().After(m.statusTimeout) {
		m.err = nil
		m.successMsg = ""
		m.statusTimeout = time.Time{}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		return m, tickCmd()

	case envsLoadedMsg:
		m.environments = msg.envs
		if len(m.environments) > 0 {
			return m, m.loadSecrets
		}
		return m, nil

	case secretsLoadedMsg:
		m.secrets = msg.secrets
		m.selectedSecret = 0
		m.secretScrollOffset = 0
		m.applySearchFilter()
		return m, nil

	case errMsg:
		m.err = msg.err
		m.statusTimeout = time.Now().Add(5 * time.Second)
		return m, nil

	case successMsg:
		m.successMsg = msg.message
		m.statusTimeout = time.Now().Add(3 * time.Second)
		return m, nil

	case tea.KeyMsg:
		// Handle modal views
		if m.currentView == viewAddSecret || m.currentView == viewEditSecret {
			return m.handleInputView(msg)
		}

		if m.currentView == viewDeleteConfirm {
			return m.handleDeleteConfirm(msg)
		}

		// Handle search mode
		if m.searchMode {
			return m.handleSearchInput(msg)
		}

		// Main view navigation
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Refresh):
			return m, tea.Batch(m.loadEnvironments, func() tea.Msg {
				return successMsg{"Refreshed"}
			})

		case key.Matches(msg, m.keys.Left):
			m.activePanel = panelEnvironments

		case key.Matches(msg, m.keys.Right):
			if len(m.environments) > 0 {
				m.activePanel = panelSecrets
			}

		case key.Matches(msg, m.keys.Search):
			if m.activePanel == panelSecrets {
				m.searchMode = true
				m.searchInput.Focus()
				return m, textinput.Blink
			}

		case key.Matches(msg, m.keys.Up):
			m.navigateUp()

		case key.Matches(msg, m.keys.Down):
			m.navigateDown()

		case key.Matches(msg, m.keys.PageUp):
			m.pageUp()

		case key.Matches(msg, m.keys.PageDown):
			m.pageDown()

		case key.Matches(msg, m.keys.Enter):
			if m.activePanel == panelEnvironments && len(m.environments) > 0 {
				m.activePanel = panelSecrets
				return m, m.loadSecrets
			}

		case key.Matches(msg, m.keys.Toggle):
			if m.activePanel == panelSecrets {
				m.showValues = !m.showValues
			}

		case key.Matches(msg, m.keys.Add):
			if m.activePanel == panelSecrets && len(m.environments) > 0 {
				m.currentView = viewAddSecret
				m.inputMode = "key"
				m.keyInput.Focus()
				m.keyInput.SetValue("")
				m.valueInput.SetValue("")
				return m, textinput.Blink
			}

		case key.Matches(msg, m.keys.Edit):
			if m.activePanel == panelSecrets && len(m.filteredSecrets) > 0 {
				secret := m.filteredSecrets[m.selectedSecret]
				m.currentView = viewEditSecret
				m.inputMode = "value"
				m.editingKey = secret.Key
				m.keyInput.SetValue(secret.Key)
				m.keyInput.Blur()
				m.valueInput.SetValue(secret.Value)
				m.valueInput.Focus()
				return m, textinput.Blink
			}

		case key.Matches(msg, m.keys.Delete):
			if m.activePanel == panelSecrets && len(m.filteredSecrets) > 0 {
				secret := m.filteredSecrets[m.selectedSecret]
				m.deleteTarget = secret.Key
				m.currentView = viewDeleteConfirm
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) navigateUp() {
	if m.activePanel == panelEnvironments {
		if m.currentEnvIndex > 0 {
			m.currentEnvIndex--
			m.adjustEnvScroll()
		}
	} else if m.activePanel == panelSecrets {
		if m.selectedSecret > 0 {
			m.selectedSecret--
			m.adjustSecretScroll()
		}
	}
}

func (m *Model) navigateDown() {
	if m.activePanel == panelEnvironments {
		if m.currentEnvIndex < len(m.environments)-1 {
			m.currentEnvIndex++
			m.adjustEnvScroll()
		}
	} else if m.activePanel == panelSecrets {
		if m.selectedSecret < len(m.filteredSecrets)-1 {
			m.selectedSecret++
			m.adjustSecretScroll()
		}
	}
}

func (m *Model) pageUp() {
	if m.activePanel == panelEnvironments {
		m.currentEnvIndex = max(0, m.currentEnvIndex-5)
		m.adjustEnvScroll()
	} else if m.activePanel == panelSecrets {
		m.selectedSecret = max(0, m.selectedSecret-5)
		m.adjustSecretScroll()
	}
}

func (m *Model) pageDown() {
	if m.activePanel == panelEnvironments {
		m.currentEnvIndex = min(len(m.environments)-1, m.currentEnvIndex+5)
		m.adjustEnvScroll()
	} else if m.activePanel == panelSecrets {
		m.selectedSecret = min(len(m.filteredSecrets)-1, m.selectedSecret+5)
		m.adjustSecretScroll()
	}
}

func (m *Model) adjustEnvScroll() {
	visibleHeight := m.height - 10
	if m.currentEnvIndex < m.envScrollOffset {
		m.envScrollOffset = m.currentEnvIndex
	} else if m.currentEnvIndex >= m.envScrollOffset+visibleHeight {
		m.envScrollOffset = m.currentEnvIndex - visibleHeight + 1
	}
}

func (m *Model) adjustSecretScroll() {
	visibleHeight := m.height - 12
	if m.selectedSecret < m.secretScrollOffset {
		m.secretScrollOffset = m.selectedSecret
	} else if m.selectedSecret >= m.secretScrollOffset+visibleHeight {
		m.secretScrollOffset = m.selectedSecret - visibleHeight + 1
	}
}

func (m *Model) applySearchFilter() {
	if m.searchQuery == "" {
		m.filteredSecrets = m.secrets
		return
	}

	query := strings.ToLower(m.searchQuery)
	filtered := []secretItem{}
	for _, secret := range m.secrets {
		if strings.Contains(strings.ToLower(secret.Key), query) ||
			strings.Contains(strings.ToLower(secret.Value), query) {
			filtered = append(filtered, secret)
		}
	}
	m.filteredSecrets = filtered

	// Adjust selection if needed
	if m.selectedSecret >= len(m.filteredSecrets) {
		m.selectedSecret = max(0, len(m.filteredSecrets)-1)
	}
}

func (m Model) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.searchMode = false
		m.searchInput.Blur()
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m.applySearchFilter()
		return m, nil

	case "enter":
		m.searchMode = false
		m.searchInput.Blur()
		m.searchQuery = m.searchInput.Value()
		m.applySearchFilter()
		return m, nil
	}

	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m Model) handleDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		env := m.environments[m.currentEnvIndex]
		err := m.service.DeleteSecret(env.Name, m.deleteTarget)
		if err != nil {
			m.err = err
			m.statusTimeout = time.Now().Add(5 * time.Second)
		} else {
			m.successMsg = fmt.Sprintf("Deleted secret: %s", m.deleteTarget)
			m.statusTimeout = time.Now().Add(3 * time.Second)
		}
		m.currentView = viewMain
		return m, m.loadSecrets

	case "n", "N", "esc":
		m.currentView = viewMain
		return m, nil
	}

	return m, nil
}

func (m Model) handleInputView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.currentView = viewMain
		m.keyInput.Blur()
		m.valueInput.Blur()
		return m, nil

	case "tab":
		if m.currentView == viewEditSecret {
			return m, nil // Can't change key when editing
		}
		if m.inputMode == "key" {
			m.inputMode = "value"
			m.keyInput.Blur()
			m.valueInput.Focus()
		} else {
			m.inputMode = "key"
			m.valueInput.Blur()
			m.keyInput.Focus()
		}
		return m, nil

	case "enter":
		key := m.keyInput.Value()
		value := m.valueInput.Value()

		if key != "" && value != "" {
			env := m.environments[m.currentEnvIndex]
			err := m.service.SetSecret(env.Name, key, value)
			if err != nil {
				m.err = err
				m.statusTimeout = time.Now().Add(5 * time.Second)
			} else {
				if m.currentView == viewEditSecret {
					m.successMsg = fmt.Sprintf("Updated secret: %s", key)
				} else {
					m.successMsg = fmt.Sprintf("Added secret: %s", key)
				}
				m.statusTimeout = time.Now().Add(3 * time.Second)
				m.currentView = viewMain
				m.keyInput.Blur()
				m.valueInput.Blur()
				return m, m.loadSecrets
			}
		}
		return m, nil
	}

	if m.inputMode == "key" && m.currentView == viewAddSecret {
		m.keyInput, cmd = m.keyInput.Update(msg)
	} else {
		m.valueInput, cmd = m.valueInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorBase).
			Render("Loading...")
	}

	switch m.currentView {
	case viewAddSecret:
		return m.viewAddSecretModal()
	case viewEditSecret:
		return m.viewEditSecretModal()
	case viewDeleteConfirm:
		return m.viewDeleteConfirmModal()
	default:
		return m.viewMain()
	}
}

func (m Model) viewMain() string {
	// Calculate panel widths
	leftWidth := max(30, m.width/3)
	rightWidth := m.width - leftWidth - 6

	// Top bar with stats
	topBar := m.renderTopBar()

	// Bottom help bar (render first to calculate height)
	helpBar := m.renderHelpBar()

	// Status bar
	statusBar := m.renderStatusBar()

	// Calculate reserved height
	topBarHeight := lipgloss.Height(topBar)
	statusBarHeight := lipgloss.Height(statusBar)
	helpBarHeight := lipgloss.Height(helpBar)
	reservedHeight := topBarHeight + statusBarHeight + helpBarHeight + 2

	// Available height for panels
	panelHeight := max(10, m.height-reservedHeight)

	// Left panel - Environments
	leftPanel := m.renderEnvironmentsPanel(leftWidth, panelHeight)

	// Right panel - Secrets
	rightPanel := m.renderSecretsPanel(rightWidth, panelHeight)

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	// Full view with background
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		topBar,
		panels,
		statusBar,
		helpBar,
	)

	return lipgloss.NewStyle().
		Background(colorBase).
		Width(m.width).
		Height(m.height).
		Render(mainView)
}

func (m Model) renderTopBar() string {
	// Title with icon
	title := appTitleStyle.Render("🔐 ENVOKE")

	totalEnvs := len(m.environments)
	totalSecrets := len(m.secrets)
	activeEnv := "N/A"

	if len(m.environments) > 0 {
		env := m.environments[m.currentEnvIndex]
		activeEnv = env.Name
	}

	// Stats with badges
	envBadge := countBadgeStyle.Render(fmt.Sprintf("%d", totalEnvs))
	secretBadge := countBadgeStyle.Render(fmt.Sprintf("%d", totalSecrets))

	stats := lipgloss.JoinHorizontal(
		lipgloss.Left,
		statLabelStyle.Render("Environments: "),
		envBadge,
		"  ",
		statLabelStyle.Render("Secrets: "),
		secretBadge,
		"  ",
		statLabelStyle.Render("Current: "),
		statValueStyle.Render(activeEnv),
	)

	leftSide := lipgloss.JoinHorizontal(lipgloss.Left, title, "  ", stats)

	bar := headerStyle.Width(m.width - 4).Render(leftSide)

	return bar
}

func (m Model) renderEnvironmentsPanel(width, height int) string {
	var content strings.Builder

	// Panel header with decorative border
	countBadge := countBadgeStyle.Render(fmt.Sprintf("%d", len(m.environments)))
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		panelHeaderStyle.Render("ENVIRONMENTS "),
		countBadge,
	)

	content.WriteString(header)
	content.WriteString("\n")
	content.WriteString(dividerStyle.Render(strings.Repeat("─", width-4)))
	content.WriteString("\n\n")

	if len(m.environments) == 0 {
		empty := emptyStateStyle.Render("📭 No environments found")
		content.WriteString(empty)
		content.WriteString("\n")
	} else {
		visibleHeight := height - 8
		start := m.envScrollOffset
		end := min(start+visibleHeight, len(m.environments))

		for i := start; i < end; i++ {
			env := m.environments[i]

			indicator := indicatorInactiveStyle.Render("  ")
			style := envItemStyle

			if i == m.currentEnvIndex {
				indicator = indicatorActiveStyle.Render("▶ ")
				style = envSelectedStyle
			}

			status := ""
			if env.IsActive {
				status = activeTagStyle.Render(" ACTIVE ")
			}

			line := fmt.Sprintf("%s%s%s", indicator, env.Name, status)
			rendered := style.Width(width - 6).Render(line)

			content.WriteString(rendered)
			content.WriteString("\n")
		}

		// Scroll indicator
		if len(m.environments) > visibleHeight {
			content.WriteString("\n")
			scrollInfo := scrollIndicatorStyle.Render(
				fmt.Sprintf("↕ %d-%d of %d", start+1, end, len(m.environments)),
			)
			content.WriteString(scrollInfo)
		}
	}

	panelStr := content.String()

	if m.activePanel == panelEnvironments {
		return activePanelStyle.Width(width).Height(height).Render(panelStr)
	}
	return panelStyle.Width(width).Height(height).Render(panelStr)
}

func (m Model) renderSecretsPanel(width, height int) string {
	var content strings.Builder

	// Panel header
	secretCount := len(m.filteredSecrets)
	totalCount := len(m.secrets)

	var countDisplay string
	if m.searchQuery != "" {
		countDisplay = fmt.Sprintf("%d/%d", secretCount, totalCount)
	} else {
		countDisplay = fmt.Sprintf("%d", secretCount)
	}

	countBadge := countBadgeStyle.Render(countDisplay)
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		panelHeaderStyle.Render("SECRETS "),
		countBadge,
	)

	content.WriteString(header)
	content.WriteString("\n")
	content.WriteString(dividerStyle.Render(strings.Repeat("─", width-4)))
	content.WriteString("\n\n")

	if len(m.environments) == 0 {
		empty := emptyStateStyle.Render("← Select an environment first")
		content.WriteString(empty)
		content.WriteString("\n")
	} else if len(m.filteredSecrets) == 0 {
		if m.searchQuery != "" {
			empty := emptyStateStyle.Render("🔍 No secrets match your search")
			content.WriteString(empty)
		} else {
			empty := emptyStateStyle.Render("📝 No secrets yet. Press 'a' to add one.")
			content.WriteString(empty)
		}
		content.WriteString("\n")
	} else {
		// Show search indicator
		if m.searchQuery != "" {
			searchBadge := searchActiveStyle.Render(fmt.Sprintf("🔍 %s", m.searchQuery))
			content.WriteString(searchBadge)
			content.WriteString("\n\n")
		}

		visibleHeight := height - 12
		if m.searchQuery != "" {
			visibleHeight -= 2
		}

		start := m.secretScrollOffset
		end := min(start+visibleHeight, len(m.filteredSecrets))

		// Calculate max key length for alignment
		maxKeyLen := 0
		for i := start; i < end; i++ {
			if len(m.filteredSecrets[i].Key) > maxKeyLen {
				maxKeyLen = len(m.filteredSecrets[i].Key)
			}
		}
		maxKeyLen = min(maxKeyLen, 30)

		for i := start; i < end; i++ {
			secret := m.filteredSecrets[i]

			indicator := indicatorInactiveStyle.Render("  ")
			bgStyle := lipgloss.NewStyle()

			if i == m.selectedSecret {
				indicator = indicatorActiveStyle.Render("▶ ")
				bgStyle = envSelectedStyle
			}

			displayValue := secret.Masked
			valueStyle := maskedValueStyle
			if m.showValues {
				displayValue = secret.Value
				valueStyle = secretValueStyle
			}

			// Truncate if too long
			maxValueLen := width - maxKeyLen - 15
			if len(displayValue) > maxValueLen {
				displayValue = displayValue[:maxValueLen-3] + "..."
			}

			keyStyled := secretKeyStyle.Render(fmt.Sprintf("%-*s", maxKeyLen, secret.Key))
			valueStyled := valueStyle.Render(displayValue)

			line := fmt.Sprintf("%s%s = %s", indicator, keyStyled, valueStyled)

			if i == m.selectedSecret {
				line = bgStyle.Width(width - 6).Render(line)
			}

			content.WriteString(line)
			content.WriteString("\n")
		}

		// Scroll indicator
		if len(m.filteredSecrets) > visibleHeight {
			content.WriteString("\n")
			scrollInfo := scrollIndicatorStyle.Render(
				fmt.Sprintf("↕ %d-%d of %d", start+1, end, len(m.filteredSecrets)),
			)
			content.WriteString(scrollInfo)
			content.WriteString("\n")
		}

		// Visibility status
		content.WriteString("\n")
		visStatus := "🔒 hidden"
		visStyle := lipgloss.NewStyle().Foreground(colorOverlay1)
		if m.showValues {
			visStatus = "👁️  visible"
			visStyle = lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
		}
		content.WriteString(visStyle.Render(fmt.Sprintf("Values: %s", visStatus)))
	}

	panelStr := content.String()

	if m.activePanel == panelSecrets {
		return activePanelStyle.Width(width).Height(height).Render(panelStr)
	}
	return panelStyle.Width(width).Height(height).Render(panelStr)
}

func (m Model) renderStatusBar() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf(" ✗ Error: %s ", m.err.Error()))
	}
	if m.successMsg != "" {
		return successStyle.Render(fmt.Sprintf(" ✓ %s ", m.successMsg))
	}
	return ""
}

func (m Model) renderHelpBar() string {
	var helps []string

	if m.activePanel == panelEnvironments {
		helps = []string{
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("↑↓"), helpStyle.Render("navigate")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("→"), helpStyle.Render("secrets")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("r"), helpStyle.Render("refresh")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("q"), helpStyle.Render("quit")),
		}
	} else {
		helps = []string{
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("↑↓"), helpStyle.Render("nav")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("←"), helpStyle.Render("envs")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("a"), helpStyle.Render("add")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("e"), helpStyle.Render("edit")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("d"), helpStyle.Render("delete")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("t"), helpStyle.Render("toggle")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("/"), helpStyle.Render("search")),
			lipgloss.JoinHorizontal(lipgloss.Left, helpKeyStyle.Render("q"), helpStyle.Render("quit")),
		}
	}

	helpText := lipgloss.JoinHorizontal(lipgloss.Left, helps...)
	return lipgloss.NewStyle().
		Background(colorSurface0).
		Foreground(colorText).
		Padding(0, 2).
		Width(m.width - 4).
		Render(helpText)
}

func (m Model) viewAddSecretModal() string {
	// Modal title
	title := modalTitleStyle.
		Width(50).
		Render("✨ Add New Secret")

	// Form content
	var form strings.Builder
	form.WriteString("\n\n")

	// Key input
	keyLabel := inputLabelStyle.Render("🔑 Key:")
	form.WriteString(keyLabel)
	form.WriteString("\n  ")

	keyInputBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSurface2).
		Padding(0, 1).
		Render(m.keyInput.View())

	if m.inputMode == "key" {
		keyInputBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMauve).
			Padding(0, 1).
			Render(m.keyInput.View())
	}

	form.WriteString(keyInputBorder)
	form.WriteString("\n\n")

	// Value input
	valueLabel := inputLabelStyle.Render("🔒 Value:")
	form.WriteString(valueLabel)
	form.WriteString("\n  ")

	valueInputBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSurface2).
		Padding(0, 1).
		Render(m.valueInput.View())

	if m.inputMode == "value" {
		valueInputBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMauve).
			Padding(0, 1).
			Render(m.valueInput.View())
	}

	form.WriteString(valueInputBorder)
	form.WriteString("\n\n")

	// Help text
	help := lipgloss.NewStyle().
		Foreground(colorSubtext0).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				helpKeyStyle.Render("tab"),
				" switch  ",
				helpKeyStyle.Render("enter"),
				" save  ",
				helpKeyStyle.Render("esc"),
				" cancel",
			),
		)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		form.String(),
		help,
	)

	modal := modalStyle.Render(content)

	// Center the modal on colored background
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)

	return lipgloss.NewStyle().
		Background(colorBase).
		Render(centered)
}

func (m Model) viewEditSecretModal() string {
	// Modal title
	title := modalTitleStyle.
		Width(50).
		Render(fmt.Sprintf("✏️  Edit Secret"))

	// Form content
	var form strings.Builder
	form.WriteString("\n\n")

	// Key display (read-only)
	keyLabel := inputLabelStyle.Render("🔑 Key (read-only):")
	form.WriteString(keyLabel)
	form.WriteString("\n  ")

	keyDisplay := lipgloss.NewStyle().
		Foreground(colorSubtext1).
		Background(colorSurface1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSurface2).
		Padding(0, 1).
		Width(40).
		Render(m.editingKey)

	form.WriteString(keyDisplay)
	form.WriteString("\n\n")

	// Value input
	valueLabel := inputLabelStyle.Render("🔒 New Value:")
	form.WriteString(valueLabel)
	form.WriteString("\n  ")

	valueInputBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMauve).
		Padding(0, 1).
		Render(m.valueInput.View())

	form.WriteString(valueInputBorder)
	form.WriteString("\n\n")

	// Help text
	help := lipgloss.NewStyle().
		Foreground(colorSubtext0).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				helpKeyStyle.Render("enter"),
				" save  ",
				helpKeyStyle.Render("esc"),
				" cancel",
			),
		)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		form.String(),
		help,
	)

	modal := modalStyle.Render(content)

	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)

	return lipgloss.NewStyle().
		Background(colorBase).
		Render(centered)
}

func (m Model) viewDeleteConfirmModal() string {
	// Modal title with warning style
	title := lipgloss.NewStyle().
		Foreground(colorBase).
		Background(colorRed).
		Bold(true).
		Padding(0, 2).
		Width(50).
		Align(lipgloss.Center).
		Render("⚠️  Confirm Deletion")

	// Warning message
	message := lipgloss.NewStyle().
		Foreground(colorText).
		Padding(2, 2).
		Render(
			fmt.Sprintf(
				"Are you sure you want to delete this secret?\n\n"+
					"  %s\n\n"+
					"%s",
				secretKeyStyle.Render(m.deleteTarget),
				lipgloss.NewStyle().
					Foreground(colorSubtext0).
					Italic(true).
					Render("This action cannot be undone."),
			),
		)

	// Help text
	help := lipgloss.NewStyle().
		Foreground(colorSubtext0).
		Padding(1, 0).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				helpKeyStyle.Render("y"),
				" confirm  ",
				helpKeyStyle.Render("n"),
				" cancel",
			),
		)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		message,
		help,
	)

	modal := modalStyle.Render(content)

	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)

	return lipgloss.NewStyle().
		Background(colorBase).
		Render(centered)
}

// Run starts the TUI
func Run(service *service.SecretService) error {
	p := tea.NewProgram(
		NewModel(service),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
