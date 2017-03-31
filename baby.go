package main

import (
	"fmt"
	"time"
)

// Baby holds all information of a baby
type Baby struct {
	ID        int
	Name      string
	Type      string
	Birthdate time.Time
	Size      int // Size is in centimeter
	Weight    int // Weight in gramm
}

func (b *Baby) String() string {
	return fmt.Sprintf("%s: %s (%s) - %d g - %d cm", b.Gender(), b.Name, b.Birthdate.Format("02.01.2006 - 15:04"), b.Weight, b.Size)
}

// Gender returns a human readble name
func (b *Baby) Gender() string {
	if b.Type == "male" {
		return "Mann"
	}
	if b.Type == "female" {
		return "Frau"
	}
	if b.Type == "twins" {
		return "Zwillinge"
	}
	return "transgender"
}
