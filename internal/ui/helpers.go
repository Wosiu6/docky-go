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
		return "#00FF00", "\u25cf", "RUNNING"
	case "paused":
		return "#FFA500", "\u275a\u275a", "PAUSED"
	case "restarting":
		return "#FFFF00", "\u21bb", "RESTARTING"
	case "exited":
		return "#FF0000", "\u25a0", "EXITED"
	case "created":
		return "#00BFFF", "\u25cb", "CREATED"
	case "dead":
		return "#FF0000", "\u2717", "DEAD"
	default:
		return "#FF0000", "\u2b58", strings.ToUpper(status)
	}
}
