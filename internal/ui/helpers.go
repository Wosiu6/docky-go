package ui

import "strings"

type ContainerDetail interface {
	DetailFields() map[string]string
}

func TruncateString(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func StatusInfo(status string) (color string, icon string, text string) {
	switch strings.ToLower(status) {
	case "running":
		return colorSuccess, "\u25cf", "RUNNING"
	case "paused":
		return colorWarning, "\u275a\u275a", "PAUSED"
	case "restarting":
		return colorInfo, "\u21bb", "RESTARTING"
	case "exited":
		return colorDanger, "\u25a0", "EXITED"
	case "created":
		return colorLight, "\u25cb", "CREATED"
	case "dead":
		return colorDark, "\u2717", "DEAD"
	default:
		return colorDark, "\u2b58", strings.ToUpper(status)
	}
}
