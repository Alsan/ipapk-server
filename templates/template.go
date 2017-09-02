package templates

import (
	"github.com/dustin/go-humanize"
	"html/template"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format("2016-01-02 03:04:05")
}

func formatBinary(size int64) string {
	return humanize.Bytes(uint64(size))
}

var TplFuncMap = template.FuncMap{
	"formatTime":   formatTime,
	"formatBinary": formatBinary,
}
