package middleware

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	prettyDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	prettyBold    = lipgloss.NewStyle().Bold(true)
	prettyLabel   = lipgloss.NewStyle().Bold(true).Width(9)
	prettySection = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	prettyDivider = prettyDim.Render(strings.Repeat("-", 64))
)

type requestLog struct {
	RequestID, Method, Path, Route, IP, Username, UserID, Error string
	Status                                                      int
	Duration                                                    time.Duration
	Query                                                       map[string]any
	Body                                                        any
}

func (entry *requestLog) attributes() []any {
	attributes := []any{
		"request_id", entry.RequestID,
		"method", entry.Method,
		"path", entry.Path,
		"route", entry.Route,
		"status", entry.Status,
		"duration_ms", entry.Duration.Milliseconds(),
		"ip", entry.IP,
		"user", entry.Username,
		"user_id", entry.UserID,
	}

	if len(entry.Query) > 0 {
		attributes = append(attributes, "query_params", entry.Query)
	}

	if entry.Body != nil {
		attributes = append(attributes, "body", entry.Body)
	}

	if len(entry.Error) > 0 {
		attributes = append(attributes, "errors", entry.Error)
	}

	return attributes
}

func prettyBadge(color string) lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Padding(0, 1).Foreground(lipgloss.Color("0")).Background(lipgloss.Color(color))
}

func statusColor(status int) string {
	switch {
	case status >= 500:
		return "9" // Red
	case status >= 400:
		return "11" // Yellow
	case status >= 300:
		return "14" // Blue
	default:
		return "10" // Green
	}
}

func methodColor(method string) string {
	switch method {
	case "GET":
		return "12" // Blue
	case "POST":
		return "10" // Green
	case "PUT":
		return "11" // Yellow
	case "DELETE":
		return "9" // Red
	default:
		return "8" // Gray
	}
}

func printPretty(r requestLog) {
	var b strings.Builder

	row := func(key, value string) {
		fmt.Fprintf(&b, " %s %s\n", prettyLabel.Render(key), value)
	}

	fmt.Fprintln(&b, prettyDivider)
	fmt.Fprintf(&b, "%s %s  %s  %s\n",
		prettyBadge(methodColor(r.Method)).Render(r.Method),
		prettyBold.Render(r.Path),
		prettyBadge(statusColor(r.Status)).Render(fmt.Sprintf("%d", r.Status)),
		prettyDim.Render(fmt.Sprintf("(%d ms)", r.Duration.Milliseconds())))
	fmt.Fprintln(&b, prettyDivider)

	row("Route", r.Route)
	row("User / User ID", fmt.Sprintf("%s / %s", r.Username, r.UserID))
	row("IP", r.IP)
	row("time", time.Now().Format("2006-01-01 15:04:05"))
	row("ReqID", prettyDim.Render(r.RequestID))

	if len(r.Query) > 0 {
		fmt.Fprintf(&b, "\n %s\n", prettySection.Render("🔎 QUERY"))
		for _, key := range slices.Sorted(maps.Keys(r.Query)) {
			row(key, fmt.Sprint(r.Query[key]))
		}
	}

	if r.Body != nil {
		fmt.Fprintf(&b, "\n %s\n", prettySection.Render("📦 BODY"))
		pretty, err := json.MarshalIndent(r.Body, " ", " ")

		if err != nil {
			pretty = []byte(fmt.Sprint(r.Body))
		}
		fmt.Fprintf(&b, " %s\n", pretty)
	}

	if r.Error != "" {
		errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(statusColor(r.Status)))

		fmt.Fprintf(&b, "\n %s %s\n", errStyle.Render("❌ Error"), r.Error)
	}

	fmt.Fprint(os.Stdout, b.String())

}
