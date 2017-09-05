package templates

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/zpatrick/go-bytesize"
)

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 03:04")
}

func formatBinary(size int64) string {
	b := bytesize.Bytesize(size)
	return fmt.Sprintf("%.1f MB", b.Megabytes())
}

func formatLog(logs string) []string {
	out := strings.Split(logs, "\\n")
	return out
}

func previewLog(logs []string) []string {
	if len(logs) > 3 {
		return logs[0:3]
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
