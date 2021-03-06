/*
Package textbackend simple usage example:

    textContent := TextContent{}
    textContent.Writeln("hello world")
    textContent.PushLevel()
    textContent.Writeln("one more line")
    textContent.PopLevel()
    theText := textContent.GetStringArray("\t")
*/
package textbackend

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// TextContent store a ContentLine list
type TextContent struct {
	Rows         []TextRow
	CurrentLevel uint8 // the current level
}

// NewTextContent initialize and return a TextContent object
func NewTextContent() *TextContent {
	c := new(TextContent)
	return c
}

// ReadFile read a file and return a TextContent object
func ReadFile(filename, whitespace, linebreak string) (*TextContent, error) {
	c := NewTextContent()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	lines := strings.Split(string(data), linebreak)
	for _, line := range lines {
		r := NewTextRowFromString(line, whitespace)
		c.AppendRow(*r)
	}
	return c, nil
}

// AppendRow add a TextRow object to the Rows array
func (t *TextContent) AppendRow(r TextRow) {
	t.Rows = append(t.Rows, r)
}

// Reset the content data and the current level value
func (t *TextContent) Reset() {
	t.Rows = []TextRow{}
	t.CurrentLevel = 0
}

// Writeln append a line to the content array
func (t *TextContent) Writeln(text string) {
	t.Rows = append(t.Rows, TextRow{t.CurrentLevel, text})
}

// Write writes text to the latest row
func (t *TextContent) Write(text string) {
	if len(t.Rows) != 0 {
		t.Rows[len(t.Rows)-1].Text += text
	} else {
		t.Writeln(text)
	}
}

// WriteTo writes text to a specific line.
func (t *TextContent) WriteTo(row int, text string) {
	if row >= 0 && row <= len(t.Rows) {
		t.Rows[row-1].Text += text
	}
}

// PushLevel pushes into the next level.
func (t *TextContent) PushLevel() {
	t.CurrentLevel++
}

// PopLevel pop out from the level.
func (t *TextContent) PopLevel() {
	if t.CurrentLevel > 0 {
		t.CurrentLevel--
	}
}

// GetLevel return the current level depth as int.
func (t *TextContent) GetLevel() uint8 {
	return t.CurrentLevel
}

// SetLevel set the level to a specific depth
func (t *TextContent) SetLevel(depth uint8) {
	t.CurrentLevel = depth
}

// GetTotalLines returns the number of lines.
func (t *TextContent) GetTotalLines() int {
	return len(t.Rows)
}

// Get renders out the content as string and you can set the type of linebreak and level char
func (t *TextContent) Get(linebreak, whitespace string) []byte {
	buf := []byte{}
	for i := 0; i < len(t.Rows); i++ {
		buf = append(buf, t.Rows[i].GetString(whitespace)...)
		buf = append(buf, linebreak...)
	}
	return buf
}

// GetStringArray return the content as array.
func (t *TextContent) GetStringArray(whitespace string) []string {
	tmp := []string{}
	for i := 0; i < len(t.Rows); i++ {
		tmp = append(tmp, t.Rows[i].GetString(whitespace))
	}
	return tmp
}

// GetJSON return the content as array with row objects.
func (t *TextContent) GetJSON() ([]byte, error) {
	return json.Marshal(t.Rows)
}

// WriteFile write the text content to a file.
func (t *TextContent) WriteFile(filename string, perm os.FileMode, linebreak, whitespace string) error {
	return ioutil.WriteFile(filename, t.Get(linebreak, whitespace), perm)
}
