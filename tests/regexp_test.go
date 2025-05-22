package test

import (
	"fmt"
	"log"
	"regexp"
	te "testing"
)

func TestRegExp(t *te.T) {
	r, err := regexp.Compile(`.*(?i)прокс.*|.*(?i)инжест.*|.*(?i)proxy.*`)
	if err != nil {
		t.Errorf("Regex could not be compiled!\n%s", err)
	}
	resp := r.FindString("инжест пожалуйста")
	log.Printf("String matching: %v", resp)
	if resp == "" {
		t.Errorf("Expression not found")
	}
}

func TestPrinting(t *te.T) {
	username := "name"
	text := "textextext"
	print := fmt.Sprintf("@%s\n\n%s", username, text)
	log.Print(print)
}
