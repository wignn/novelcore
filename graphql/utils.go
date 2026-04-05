package main

import (
	"fmt"
	"log"
)

func handleError(label string, err error) error {
	log.Printf("[%s] error: %v\n", label, err)
	return fmt.Errorf("%s: %v", label, err)
}
