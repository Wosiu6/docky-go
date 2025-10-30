package model

import "strings"

type BaseContainerInfo struct {
	ID         string
	Names      []string
	Image      string
	CPUPercent float64
	Mem        uint64
	Status     string
}

func ParseEnv(env []string) map[string]string {
	out := make(map[string]string, len(env))
	for _, e := range env {
		if eq := strings.IndexByte(e, '='); eq > 0 {
			out[e[:eq]] = e[eq+1:]
		}
	}
	return out
}
