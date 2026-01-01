package output

import (
	"fmt"
	"strings"

	"github.com/badimirzai/robotics-verifier-cli/internal/ui"
	"github.com/badimirzai/robotics-verifier-cli/internal/validate"
)

func RenderReport(r validate.Report) string {
	var b strings.Builder
	b.WriteString(ui.Colorize("HEADER", "rv check"))
	b.WriteString("\n")
	b.WriteString(ui.Colorize("HEADER", "--------------"))
	b.WriteString("\n")
	for _, f := range r.Findings {
		severity := string(f.Severity)
		b.WriteString(ui.Colorize(severity, severity))
		b.WriteString(" ")
		b.WriteString(f.Code)
		b.WriteString(": ")
		if f.Location != nil {
			b.WriteString(f.Location.File)
			if f.Location.Line > 0 {
				b.WriteString(fmt.Sprintf(":%d", f.Location.Line))
			}
			b.WriteString(" ")
		}
		b.WriteString(f.Message)
		b.WriteString("\n")
	}
	return b.String()
}
