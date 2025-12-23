package output

import (
	"strings"

	"github.com/badimirzai/robostack-cli/internal/validate"
)

func RenderReport(r validate.Report) string {
	var b strings.Builder
	b.WriteString("robostack validate\n")
	b.WriteString("--------------\n")
	for _, f := range r.Findings {
		b.WriteString(string(f.Severity))
		b.WriteString(" ")
		b.WriteString(f.Code)
		b.WriteString(": ")
		b.WriteString(f.Message)
		b.WriteString("\n")
	}
	return b.String()
}
