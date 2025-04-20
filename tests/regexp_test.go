package test

import (
	"log"
	"regexp"
	te "testing"
)

func TestRegExp(t *te.T) {
	r, err := regexp.Compile(`(?i)инжест`)
	if err != nil {
		t.Errorf("Regex could not be compiled!\n%s", err)
	}
	resp := r.FindString("одни пидорасы в ИнЖесТЕе")
	log.Printf("String matching: %v", resp)
	if resp == "" {
		t.Errorf("Expression not found")
	}
}
