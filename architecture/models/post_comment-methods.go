package models

import (
	"fmt"
	"strings"
)

func (c *PostComment) ValidateContent() error {
	if lng := len(c.Content); lng < 1 {
		return fmt.Errorf("content: invalid lenght (%d)", lng)
	}
	return nil
}

func (c *PostComment) PrepareContent() {
	c.Content = strings.Trim(c.Content, " ")
}

func (c *PostComment) Prepare() {
	c.PrepareContent()
}
