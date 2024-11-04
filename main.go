package main

import (
	`fmt`
)

func main() {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}

	w := NewWPool(3)

	go func() {
		for _, task := range list {
			w.jobs <- task
		}
		close(w.jobs)
	}()

	w.Add()
	w.Add()

	go func() {
		for output := range w.results {
			fmt.Println(output)
		}
	}()

	w.Remove()

	w.Close()
}
