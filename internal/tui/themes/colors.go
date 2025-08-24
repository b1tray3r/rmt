package themes

import "github.com/charmbracelet/lipgloss"

// Theme defines a color theme for the TUI package, mapping semantic names to lipgloss.Color values.
type Theme struct {
	Primary       lipgloss.Color // Theme.Primary is the main accent color for primary actions or highlights.
	Secondary     lipgloss.Color // Theme.Secondary is used for secondary actions or less prominent highlights.
	Success       lipgloss.Color // Theme.Success is used to indicate successful operations or positive statuses.
	Warning       lipgloss.Color // Theme.Warning is used to indicate warnings or cautionary statuses.
	Error         lipgloss.Color // Theme.Error is used to indicate errors or critical statuses.
	Info          lipgloss.Color // Theme.Info is used for informational messages.
	Background    lipgloss.Color // Theme.Background is the main background color.
	BackgroundAlt lipgloss.Color // Theme.BackgroundAlt is an alternative background color for contrast.
	Foreground    lipgloss.Color // Theme.Foreground is the main text color.
	Muted         lipgloss.Color // Theme.Muted is used for less prominent text, such as comments.
	Border        lipgloss.Color // Theme.Border is used for borders and separators.
	Highlight     lipgloss.Color // Theme.Highlight is used for highlighting selections or active elements.
	Link          lipgloss.Color // Theme.Link is used for hyperlinks or interactive text.
}

// TokyoNight provides a color theme inspired by the Tokyo Night color scheme.
var TokyoNight = Theme{
	Primary:       ColorTerminalBlack,
	Secondary:     ColorSkyBlue,
	Success:       ColorGreen,
	Warning:       ColorYellow,
	Error:         ColorRed,
	Info:          ColorBlueCyan,
	Background:    ColorBgNight,
	BackgroundAlt: ColorBgStorm,
	Foreground:    ColorWhite,
	Muted:         ColorComment,
	Border:        ColorTerminalBlack,
	Highlight:     ColorMagenta,
	Link:          ColorLightGreen,
}

// Color constants define the color palette used in the TUI package as lipgloss.Color values.
var (
	ColorRed           = lipgloss.Color("#f7768e") // ColorRed is used for keywords, HTML elements, regex group symbols, CSS units, and terminal red.
	ColorOrange        = lipgloss.Color("#ff9e64") // ColorOrange is used for number and boolean constants, and language support constants.
	ColorYellow        = lipgloss.Color("#e0af68") // ColorYellow is used for function parameters, regex character sets, and terminal yellow.
	ColorLightGray     = lipgloss.Color("#cfc9c2") // ColorLightGray is used for parameters inside functions (semantic highlighting only).
	ColorGreen         = lipgloss.Color("#9ece6a") // ColorGreen is used for strings and CSS class names.
	ColorLightGreen    = lipgloss.Color("#73daca") // ColorLightGreen is used for object literal keys, markdown links, and terminal green.
	ColorCyan          = lipgloss.Color("#b4f9f8") // ColorCyan is used for regex literal strings.
	ColorBlueCyan      = lipgloss.Color("#2ac3de") // ColorBlueCyan is used for language support functions and CSS HTML elements.
	ColorSkyBlue       = lipgloss.Color("#7dcfff") // ColorSkyBlue is used for object properties, regex quantifiers and flags, markdown headings, terminal cyan, markdown code, and import/export keywords.
	ColorBlue          = lipgloss.Color("#7aa2f7") // ColorBlue is used for function names, CSS property names, and terminal blue.
	ColorMagenta       = lipgloss.Color("#bb9af7") // ColorMagenta is used for control keywords, storage types, regex symbols and operators, HTML attributes, and terminal magenta.
	ColorWhite         = lipgloss.Color("#c0caf5") // ColorWhite is used for variables, class names, and terminal white.
	ColorEditorFg      = lipgloss.Color("#a9b1d6") // ColorEditorFg is used for the editor foreground.
	ColorMarkdownText  = lipgloss.Color("#9aa5ce") // ColorMarkdownText is used for markdown text and HTML text.
	ColorComment       = lipgloss.Color("#565f89") // ColorComment is used for comments.
	ColorTerminalBlack = lipgloss.Color("#414868") // ColorTerminalBlack is used for terminal black.
	ColorBgStorm       = lipgloss.Color("#24283b") // ColorBgStorm is used for the editor background (storm).
	ColorBgNight       = lipgloss.Color("#1a1b26") // ColorBgNight is used for the editor background (night).
)
