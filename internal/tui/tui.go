package tui

import (
	"fmt"
	"strings"

	"github.com/HiyawNT/envoke/internal/models"
	"github.com/HiyawNT/envoke/internal/service"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Distinctive color scheme inspired by terminal aesthetics
// Using a sophisticated dark cyberpunk palette with electric accents
var (
	// Primary colors - deep space theme with electric highlights
	colorBg          = lipgloss.Color("#0a0e14")
	colorBgAlt       = lipgloss.Color("#1a1f29")
	colorBorder      = lipgloss.Color("#2d3748")
	colorBorderLight = lipgloss.Color("#4a5568")

	// Accent colors - electric cyan and warm amber
	colorPrimary   = lipgloss.Color("#00d9ff") // Electric cyan
	colorSecondary = lipgloss.Color("#ffb454") // Warm amber
	colorSuccess   = lipgloss.Color("#7ee787") // Soft green
	colorDanger    = lipgloss.Color("#ff6b9d") // Soft pink/red
	colorWarning   = lipgloss.Color("#ffd60a") // Bright yellow

	// Text colors
	colorText     = lipgloss.Color("#e6edf3")
	colorTextDim  = lipgloss.Color("#8b949e")
	colorTextDark = lipgloss.Color("#6e7681")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorTextDim).
			Italic(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	activeBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(1, 2)

	statusBarStyle = lipgloss.NewStyle().
			Background(colorBgAlt).
			Foreground(colorTextDim).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorTextDark).
			Padding(0, 1)

	keyStyle = lipgloss.NewStyle().
			Foreground(colorSecondary).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(colorText)

	maskedStyle = lipgloss.NewStyle().
			Foreground(colorTextDark)

	itemStyle = lipgloss.NewStyle().
			Foreground(colorText).
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true).
				PaddingLeft(1)

	envActiveStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)
)

type view int

const (
	viewEnvironments view = iota
	viewSecrets
	viewAddSecret
	viewEditSecret
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Back   key.Binding
	Add    key.Binding
	Delete key.Binding
	Toggle key.Binding
	Edit   key.Binding
	Quit   key.Binding
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
			key.WithHelp("←/h", "prev env"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next env"),
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
			key.WithHelp("t", "toggle visibility"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

type Model struct {
	service         *service.SecretService
	currentView     view
	environments    []models.Environment
	currentEnvIndex int
	secrets         []secretItem
	selectedSecret  int
	showValues      bool
	width           int
	height          int
	keys            keyMap
	err             error

	// Input fields
	keyInput   textinput.Model
	valueInput textinput.Model
	inputMode  string // "key" or "value"
}

type secretItem struct {
	Key       string
	Value     string
	Masked    string
	CreatedAt string
}

func NewModel(service *service.SecretService) Model {
	keyInput := textinput.New()
	keyInput.Placeholder = "SECRET_KEY"
	keyInput.CharLimit = 128

	valueInput := textinput.New()
	valueInput.Placeholder = "secret value"
	valueInput.CharLimit = 1024
	valueInput.EchoMode = textinput.EchoPassword
	valueInput.EchoCharacter = '•'

	return Model{
		service:     service,
		currentView: viewEnvironments,
		keys:        newKeyMap(),
		keyInput:    keyInput,
		valueInput:  valueInput,
		showValues:  false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadEnvironments,
		textinput.Blink,
	)
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
		masked := strings.Repeat("•", len(value))
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case envsLoadedMsg:
		m.environments = msg.envs
		if len(m.environments) > 0 && m.currentView == viewEnvironments {
			return m, m.loadSecrets
		}
		return m, nil

	case secretsLoadedMsg:
		m.secrets = msg.secrets
		m.selectedSecret = 0
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		// Handle input mode separately
		if m.currentView == viewAddSecret || m.currentView == viewEditSecret {
			return m.handleInputView(msg)
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Back):
			if m.currentView == viewSecrets {
				m.currentView = viewEnvironments
				return m, nil
			}

		case key.Matches(msg, m.keys.Enter):
			if m.currentView == viewEnvironments && len(m.environments) > 0 {
				m.currentView = viewSecrets
				return m, m.loadSecrets
			}

		case key.Matches(msg, m.keys.Left):
			if m.currentEnvIndex > 0 {
				m.currentEnvIndex--
				return m, m.loadSecrets
			}

		case key.Matches(msg, m.keys.Right):
			if m.currentEnvIndex < len(m.environments)-1 {
				m.currentEnvIndex++
				return m, m.loadSecrets
			}

		case key.Matches(msg, m.keys.Up):
			if m.selectedSecret > 0 {
				m.selectedSecret--
			}

		case key.Matches(msg, m.keys.Down):
			if m.selectedSecret < len(m.secrets)-1 {
				m.selectedSecret++
			}

		case key.Matches(msg, m.keys.Toggle):
			if m.currentView == viewSecrets {
				m.showValues = !m.showValues
			}

		case key.Matches(msg, m.keys.Add):
			if m.currentView == viewSecrets {
				m.currentView = viewAddSecret
				m.inputMode = "key"
				m.keyInput.Focus()
				m.keyInput.SetValue("")
				m.valueInput.SetValue("")
			}
		}
	}

	return m, nil
}

func (m Model) handleInputView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.currentView = viewSecrets
		m.keyInput.Blur()
		m.valueInput.Blur()
		return m, nil

	case "tab":
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
		if m.keyInput.Value() != "" && m.valueInput.Value() != "" {
			env := m.environments[m.currentEnvIndex]
			err := m.service.SetSecret(env.Name, m.keyInput.Value(), m.valueInput.Value())
			if err != nil {
				m.err = err
			} else {
				m.currentView = viewSecrets
				m.keyInput.Blur()
				m.valueInput.Blur()
				return m, m.loadSecrets
			}
		}
		return m, nil
	}

	if m.inputMode == "key" {
		m.keyInput, cmd = m.keyInput.Update(msg)
	} else {
		m.valueInput, cmd = m.valueInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.currentView {
	case viewEnvironments:
		return m.viewEnvironmentsList()
	case viewSecrets:
		return m.viewSecretsList()
	case viewAddSecret:
		return m.viewAddSecretForm()
	default:
		return m.viewEnvironmentsList()
	}
}

func (m Model) viewEnvironmentsList() string {
	title := titleStyle.Render("🔐 ENVOKE")
	subtitle := subtitleStyle.Render("Encrypted Secret Manager")

	header := lipgloss.JoinVertical(lipgloss.Left, title, subtitle)

	if len(m.environments) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(colorTextDim).
			Padding(2, 0).
			Render("No environments found. Run 'envoke env create <name>' to get started.")

		content := lipgloss.JoinVertical(lipgloss.Left, header, "", empty)
		return borderStyle.Width(m.width - 4).Render(content)
	}

	// Environment list
	var envList strings.Builder
	envList.WriteString("\n")
	envList.WriteString(lipgloss.NewStyle().Foreground(colorSecondary).Bold(true).Render("ENVIRONMENTS"))
	envList.WriteString("\n\n")

	for i, env := range m.environments {
		prefix := "  "
		style := itemStyle

		if i == m.currentEnvIndex {
			prefix = "▶ "
			style = selectedItemStyle
		}

		status := ""
		if env.IsActive {
			status = envActiveStyle.Render(" [ACTIVE]")
		}

		line := fmt.Sprintf("%s%s%s", prefix, env.Name, status)
		envList.WriteString(style.Render(line))
		envList.WriteString("\n")
	}

	help := m.renderHelp([]string{
		"↑/↓ navigate",
		"enter select",
		"q quit",
	})

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		envList.String(),
		"",
		help,
	)

	return borderStyle.Width(m.width - 4).Render(content)
}

func (m Model) viewSecretsList() string {
	if len(m.environments) == 0 {
		return "No environments"
	}

	env := m.environments[m.currentEnvIndex]

	title := titleStyle.Render(fmt.Sprintf("🔐 %s", env.Name))
	subtitle := subtitleStyle.Render(fmt.Sprintf("%d secrets", len(m.secrets)))

	header := lipgloss.JoinVertical(lipgloss.Left, title, subtitle)

	if len(m.secrets) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(colorTextDim).
			Padding(2, 0).
			Render("No secrets in this environment. Press 'a' to add one.")

		help := m.renderHelp([]string{
			"a add",
			"esc back",
			"q quit",
		})

		content := lipgloss.JoinVertical(lipgloss.Left, header, "", empty, "", help)
		return activeBorderStyle.Width(m.width - 4).Render(content)
	}

	// Secrets list
	var secretList strings.Builder
	secretList.WriteString("\n")

	for i, secret := range m.secrets {
		prefix := "  "
		keyStyleLocal := itemStyle

		if i == m.selectedSecret {
			prefix = "▶ "
			keyStyleLocal = selectedItemStyle
		}

		displayValue := secret.Masked
		if m.showValues {
			displayValue = secret.Value
		}

		line := fmt.Sprintf("%s%s = %s",
			prefix,
			keyStyleLocal.Render(secret.Key),
			maskedStyle.Render(displayValue),
		)

		secretList.WriteString(line)
		secretList.WriteString("\n")
	}

	visibilityStatus := "hidden"
	if m.showValues {
		visibilityStatus = "visible"
	}

	statusBar := statusBarStyle.Render(fmt.Sprintf(" Values: %s ", visibilityStatus))

	help := m.renderHelp([]string{
		"↑/↓ navigate",
		"t toggle",
		"a add",
		"esc back",
		"q quit",
	})

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		secretList.String(),
		"",
		statusBar,
		"",
		help,
	)

	return activeBorderStyle.Width(m.width - 4).Render(content)
}

func (m Model) viewAddSecretForm() string {
	env := m.environments[m.currentEnvIndex]

	title := titleStyle.Render(fmt.Sprintf("➕ Add Secret to %s", env.Name))

	var form strings.Builder
	form.WriteString("\n")

	keyLabel := lipgloss.NewStyle().Foreground(colorSecondary).Bold(true).Render("Key:")
	form.WriteString(keyLabel)
	form.WriteString("\n")
	form.WriteString(m.keyInput.View())
	form.WriteString("\n\n")

	valueLabel := lipgloss.NewStyle().Foreground(colorSecondary).Bold(true).Render("Value:")
	form.WriteString(valueLabel)
	form.WriteString("\n")
	form.WriteString(m.valueInput.View())
	form.WriteString("\n\n")

	help := m.renderHelp([]string{
		"tab switch field",
		"enter save",
		"esc cancel",
	})

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		form.String(),
		help,
	)

	return activeBorderStyle.Width(m.width - 4).Render(content)
}

func (m Model) renderHelp(items []string) string {
	var parts []string
	for _, item := range items {
		parts = append(parts, helpStyle.Render(item))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

// Run starts the TUI
func Run(service *service.SecretService) error {
	p := tea.NewProgram(NewModel(service), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
