package inputany

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// InputAny[T] is a generic form input field.
type InputAny[T any] struct {
	value    *T
	strValue *string
	key      string

	// customization
	title       string
	description string
	inline      bool

	// marshaling T<->string
	marshal   func(T) string
	unmarshal func(string) (T, error)

	// error handling
	validate func(T) error
	err      error

	// model
	textinput textinput.Model

	// state
	focused bool

	// options
	width      int
	height     int
	accessible bool
	theme      *huh.Theme
	keymap     huh.InputKeyMap
}

// NewInputAny returns a new input field.
func NewInputAny[T any](marshal func(T) string, unmarshal func(string) (T, error)) *InputAny[T] {
	input := textinput.New()

	i := &InputAny[T]{
		value:     new(T),
		strValue:  new(string),
		textinput: input,
		validate:  func(T) error { return nil },
		marshal:   marshal,
		unmarshal: unmarshal,
	}

	return i
}

// Value sets the value of the input field.
func (i *InputAny[T]) Value(value *T) *InputAny[T] {
	i.value = value
	strValue := i.marshal(*value)
	i.strValue = &strValue
	i.textinput.SetValue(*i.strValue)
	return i
}

// Key sets the key of the input field.
func (i *InputAny[T]) Key(key string) *InputAny[T] {
	i.key = key
	return i
}

// Title sets the title of the input field.
func (i *InputAny[T]) Title(title string) *InputAny[T] {
	i.title = title
	return i
}

// Description sets the description of the input field.
func (i *InputAny[T]) Description(description string) *InputAny[T] {
	i.description = description
	return i
}

// Prompt sets the prompt of the input field.
func (i *InputAny[T]) Prompt(prompt string) *InputAny[T] {
	i.textinput.Prompt = prompt
	return i
}

// CharLimit sets the character limit of the input field.
func (i *InputAny[T]) CharLimit(charlimit int) *InputAny[T] {
	i.textinput.CharLimit = charlimit
	return i
}

// Suggestions sets the suggestions to display for autocomplete in the input
// field.
func (i *InputAny[T]) Suggestions(suggestions []string) *InputAny[T] {
	i.textinput.ShowSuggestions = len(suggestions) > 0
	i.textinput.KeyMap.AcceptSuggestion.SetEnabled(len(suggestions) > 0)
	i.textinput.SetSuggestions(suggestions)
	return i
}

// EchoMode sets the input behavior of the text Input field.
type EchoMode textinput.EchoMode

const (
	// EchoNormal displays text as is.
	// This is the default behavior.
	EchoModeNormal EchoMode = EchoMode(textinput.EchoNormal)

	// EchoPassword displays the EchoCharacter mask instead of actual characters.
	// This is commonly used for password fields.
	EchoModePassword EchoMode = EchoMode(textinput.EchoPassword)

	// EchoNone displays nothing as characters are entered.
	// This is commonly seen for password fields on the command line.
	EchoModeNone EchoMode = EchoMode(textinput.EchoNone)
)

// EchoMode sets the echo mode of the input.
func (i *InputAny[T]) EchoMode(mode EchoMode) *InputAny[T] {
	i.textinput.EchoMode = textinput.EchoMode(mode)
	return i
}

// Password sets whether or not to hide the input while the user is typing.
//
// Deprecated: use EchoMode(EchoPassword) instead.
func (i *InputAny[T]) Password(password bool) *InputAny[T] {
	if password {
		i.textinput.EchoMode = textinput.EchoPassword
	} else {
		i.textinput.EchoMode = textinput.EchoNormal
	}
	return i
}

// Placeholder sets the placeholder of the text input.
func (i *InputAny[T]) Placeholder(str string) *InputAny[T] {
	i.textinput.Placeholder = str
	return i
}

// Inline sets whether the title and input should be on the same line.
func (i *InputAny[T]) Inline(inline bool) *InputAny[T] {
	i.inline = inline
	return i
}

// Validate sets the validation function of the input field.
func (i *InputAny[T]) Validate(validate func(T) error) *InputAny[T] {
	i.validate = validate
	return i
}

// Error returns the error of the input field.
func (i *InputAny[T]) Error() error {
	return i.err
}

// Skip returns whether the input should be skipped or should be blocking.
func (*InputAny[T]) Skip() bool {
	return false
}

// Zoom returns whether the input should be zoomed.
func (*InputAny[T]) Zoom() bool {
	return false
}

// Focus focuses the input field.
func (i *InputAny[T]) Focus() tea.Cmd {
	i.focused = true
	return i.textinput.Focus()
}

// Blur blurs the input field.
func (i *InputAny[T]) Blur() tea.Cmd {
	i.focused = false
	*i.strValue = i.textinput.Value()
	i.textinput.Blur()
	*i.value, i.err = i.unmarshal(*i.strValue)
	if i.err == nil {
		i.err = i.validate(*i.value)
	}
	return nil
}

// KeyBinds returns the help message for the input field.
func (i *InputAny[T]) KeyBinds() []key.Binding {
	if i.textinput.ShowSuggestions {
		return []key.Binding{i.keymap.AcceptSuggestion, i.keymap.Prev, i.keymap.Submit, i.keymap.Next}
	}
	return []key.Binding{i.keymap.Prev, i.keymap.Submit, i.keymap.Next}
}

// Init initializes the input field.
func (i *InputAny[T]) Init() tea.Cmd {
	i.textinput.Blur()
	return nil
}

// Update updates the input field.
func (i *InputAny[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	i.textinput, cmd = i.textinput.Update(msg)
	cmds = append(cmds, cmd)
	*i.strValue = i.textinput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		i.err = nil

		switch {
		case key.Matches(msg, i.keymap.Prev):
			i.saveValue()
			if i.err != nil {
				return i, nil
			}
			cmds = append(cmds, huh.PrevField)
		case key.Matches(msg, i.keymap.Next, i.keymap.Submit):
			i.saveValue()
			if i.err != nil {
				return i, nil
			}
			cmds = append(cmds, huh.NextField)
		}
	}

	return i, tea.Batch(cmds...)
}

func (i *InputAny[T]) saveValue() {
	strValue := i.textinput.Value()
	value, err := i.unmarshal(strValue)
	if err == nil {
		err = i.validate(value)
		if err != nil {
			*i.value = value
		}
	}

	i.err = err
}

func (i *InputAny[T]) activeStyles() *huh.FieldStyles {
	theme := i.theme
	if theme == nil {
		theme = huh.ThemeCharm()
	}
	if i.focused {
		return &theme.Focused
	}
	return &theme.Blurred
}

// View renders the input field.
func (i *InputAny[T]) View() string {
	styles := i.activeStyles()

	// NB: since the method is on a pointer receiver these are being mutated.
	// Because this runs on every render this shouldn't matter in practice,
	// however.
	i.textinput.PlaceholderStyle = styles.TextInput.Placeholder
	i.textinput.PromptStyle = styles.TextInput.Prompt
	i.textinput.Cursor.Style = styles.TextInput.Cursor
	i.textinput.TextStyle = styles.TextInput.Text

	var sb strings.Builder
	if i.title != "" {
		sb.WriteString(styles.Title.Render(i.title))
		if !i.inline {
			sb.WriteString("\n")
		}
	}
	if i.description != "" {
		sb.WriteString(styles.Description.Render(i.description))
		if !i.inline {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(i.textinput.View())

	return styles.Base.Render(sb.String())
}

// Run runs the input field in accessible mode.
func (i *InputAny[T]) Run() error {
	if i.accessible {
		return i.runAccessible()
	}
	return i.run()
}

// run runs the input field.
func (i *InputAny[T]) run() error {
	return huh.Run(i)
}

// runAccessible runs the input field in accessible mode.
func (i *InputAny[T]) runAccessible() error {
	styles := i.activeStyles()
	fmt.Println(styles.Title.Render(i.title))
	fmt.Println()
	*i.strValue = accessibility.PromptString("Input: ", func(input string) error {
		value, err := i.unmarshal(input)
		if err != nil {
			return err
		}
		err = i.validate(value)
		if err != nil {
			return err
		}

		return nil
	})
	fmt.Println(styles.SelectedOption.Render("Input: " + *i.strValue + "\n"))
	return nil
}

// WithKeyMap sets the keymap on an input field.
func (i *InputAny[T]) WithKeyMap(k *huh.KeyMap) huh.Field {
	i.keymap = k.Input
	i.textinput.KeyMap.AcceptSuggestion = i.keymap.AcceptSuggestion
	return i
}

// WithAccessible sets the accessible mode of the input field.
func (i *InputAny[T]) WithAccessible(accessible bool) huh.Field {
	i.accessible = accessible
	return i
}

// WithTheme sets the theme of the input field.
func (i *InputAny[T]) WithTheme(theme *huh.Theme) huh.Field {
	if i.theme != nil {
		return i
	}
	i.theme = theme
	return i
}

// WithWidth sets the width of the input field.
func (i *InputAny[T]) WithWidth(width int) huh.Field {
	styles := i.activeStyles()
	i.width = width
	frameSize := styles.Base.GetHorizontalFrameSize()
	promptWidth := lipgloss.Width(i.textinput.PromptStyle.Render(i.textinput.Prompt))
	titleWidth := lipgloss.Width(styles.Title.Render(i.title))
	descriptionWidth := lipgloss.Width(styles.Description.Render(i.description))
	i.textinput.Width = width - frameSize - promptWidth - 1
	if i.inline {
		i.textinput.Width -= titleWidth
		i.textinput.Width -= descriptionWidth
	}
	return i
}

// WithHeight sets the height of the input field.
func (i *InputAny[T]) WithHeight(height int) huh.Field {
	i.height = height
	return i
}

// WithPosition sets the position of the input field.
func (i *InputAny[T]) WithPosition(p huh.FieldPosition) huh.Field {
	i.keymap.Prev.SetEnabled(!p.IsFirst())
	i.keymap.Next.SetEnabled(!p.IsLast())
	i.keymap.Submit.SetEnabled(p.IsLast())
	return i
}

// GetKey returns the key of the field.
func (i *InputAny[T]) GetKey() string {
	return i.key
}

// GetValue returns the value of the field.
func (i *InputAny[T]) GetValue() any {
	return *i.value
}
