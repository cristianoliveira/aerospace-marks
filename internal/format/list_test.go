package format

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableOutFormatSameLength(t *testing.T) {
	list := []string{
		"1 | app1 | title1",
		"2 | app2 super long | title2",
		"3 | another | title3",
		"212 | app4 | title4",
	}
	result := FormatTableList(list)
	lines := strings.Split(result, "\n")

	assert.Equal(t, 4, len(lines))
	assert.Equal(t, lines[0], "1   | app1            | title1")
	assert.Equal(t, lines[1], "2   | app2 super long | title2")
	assert.Equal(t, lines[2], "3   | another         | title3")
	assert.Equal(t, lines[3], "212 | app4            | title4")
}

func TestTableOutFormatWithSixColumns(t *testing.T) {
	list := []string{
		"mark1|1|app1|title1|workspace1|bundle1\r\n",
		"mark2|2|app2 super long|title2|_|bundle2\r\n",
		"mark3|3|app3|_|workspace3|_\r\n",
	}
	result := FormatTableList(list)
	lines := strings.Split(result, "\n")

	assert.Equal(t, 3, len(lines))
	assert.Equal(t, lines[0], "mark1 | 1 | app1            | title1 | workspace1 | bundle1")
	assert.Equal(t, lines[1], "mark2 | 2 | app2 super long | title2 | _          | bundle2")
	assert.Equal(t, lines[2], "mark3 | 3 | app3            | _      | workspace3 | _      ")
}
