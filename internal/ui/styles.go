package ui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary = "#1E90FF"
	colorSuccess = "#28A745"
	colorWarning = "#FFC107"
	colorDanger  = "#DC3545"
	colorInfo    = "#17A2B8"
	colorLight   = "#F8F9FA"
	colorDark    = "#343A40"

	colorGeneric     = "#874BFD"
	colorGenericDark = "#602bc9ff"
	colorText        = "#FAFAFA"
	colorTextDim     = "#d1d1d1ff"

	colorLogo       = "#FDF500"
	colorTraefik    = "#24A1C1"
	colorRedis      = "#D82C20"
	colorPostgres   = "#336791"
	colorMySQL      = "#4479A1"
	colorMongoDB    = "#47A248"
	colorNginx      = "#009639"
	colorMinio      = "#FFBD2E"
	colorMinecraft  = "#55AA55"
	colorGrafana    = "#F46800"
	colorPrometheus = "#E6522C"

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			MarginRight(1).
			MarginBottom(0)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorText)).
			Padding(0, 1).
			MarginBottom(0)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextDim)).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorText))

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorInfo)).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorDanger)).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorDanger)).
			Padding(1, 2)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorDark)).
			Italic(true)
)
