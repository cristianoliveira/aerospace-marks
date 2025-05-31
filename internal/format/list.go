package format

import (
	"fmt"
	"strings"
)

// This package contains the formatters for the list
// objects returned by Aerospace.

// FormatTableList formats a list of Mark objects into a string
// Receive a | separated list of strings and make sure that
// the columns are aligned with the same width
//
// Example:
//
//	Input
//	1 | app1 | title1
//	2 | app2 super long | title2
//	3 | app3 | title3
//
//	Output
//	1 | app1             | title1
//	2 | app2 super long  | title2
//	3 | app3             | title3
func FormatTableList(list []string) string {
	var rows [][]string
	var colWidths []int

	// Parse and calculate max width for each column
	for _, line := range list {
		fields := strings.Split(line, "|")
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
			if len(colWidths) <= i {
				colWidths = append(colWidths, len(fields[i]))
			} else if len(fields[i]) > colWidths[i] {
				colWidths[i] = len(fields[i])
			}
		}
		rows = append(rows, fields)
	}

	var b strings.Builder
	for i, row := range rows {
		for j, field := range row {
			b.WriteString(fmt.Sprintf("%-*s", colWidths[j], field))
			if j < len(row)-1 {
				b.WriteString(" | ")
			}
		}
		if i < len(rows)-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}
