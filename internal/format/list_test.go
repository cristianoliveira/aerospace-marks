package format_test

import (
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/format"
	"github.com/stretchr/testify/assert"
)

func TestTableOutFormatSameLength(t *testing.T) {
	list := []string{
		"1 | app1 | title1",
		"2 | app2 super long | title2",
		"3 | another | title3",
		"212 | app4 | title4",
	}
	result := format.FormatTableList(list)
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 4)
	assert.Equal(t, "1   | app1            | title1", lines[0])
	assert.Equal(t, "2   | app2 super long | title2", lines[1])
	assert.Equal(t, "3   | another         | title3", lines[2])
	assert.Equal(t, "212 | app4            | title4", lines[3])
}

func TestTableOutFormatWithSixColumns(t *testing.T) {
	list := []string{
		"mark1|1|app1|title1|workspace1|bundle1\r\n",
		"mark2|2|app2 super long|title2|_|bundle2\r\n",
		"mark3|3|app3|_|workspace3|_\r\n",
	}
	result := format.FormatTableList(list)
	lines := strings.Split(result, "\n")

	assert.Len(t, lines, 3)
	assert.Equal(t, "mark1 | 1 | app1            | title1 | workspace1 | bundle1", lines[0])
	assert.Equal(t, "mark2 | 2 | app2 super long | title2 | _          | bundle2", lines[1])
	assert.Equal(t, "mark3 | 3 | app3            | _      | workspace3 | _      ", lines[2])
}
