package templates

import (
	"github.com/dustin/go-humanize"
	"html/template"
	"strings"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 03:04:05")
}

func formatBinary(size int64) string {
	return humanize.Bytes(uint64(size))
}

func formatLog(logs string) []string {
	out := strings.Split(logs, "\\n")
	return out
}

func previewLog(logs []string) []string {
	if len(logs) > 5 {
		return logs[0:5]
	}
	return logs
}

var TplFuncMap = template.FuncMap{
	"formatTime":   formatTime,
	"formatBinary": formatBinary,
	"safeURL":      func(u string) template.URL { return template.URL(u) },
	"formatLog":    formatLog,
	"previewLog":   previewLog,
}
