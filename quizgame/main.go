package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var file []byte
	var err error
	customFile := os.Args[1]
	defaultFile := "questions.csv"
	if customFile != "" {
		file, err = ioutil.ReadFile(customFile)
	} else {
		file, err = ioutil.ReadFile(defaultFile)
	}
	reader := bufio.NewScanner(os.Stdin)

	if err != nil {
		fmt.Println("No se pudo cargar el archivo de preguntas")
		return
	}
	csvFile := csv.NewReader(strings.NewReader(string(file)))

	questions, err := csvFile.ReadAll()

	success := 0

	for _, question := range questions {
		fmt.Print(string(question[0]), ": Resultado? ")
		reader.Scan()
		answer := reader.Text()
		if string(question[1]) == string(answer) {
			fmt.Println("Respuesta correcta")
			success++
		} else {
			fmt.Println("Respuesta incorrecta")
		}
	}
	fmt.Println("\nTotal preguntas contestadas:", len(questions))
	fmt.Println("Preguntas acertadas:", success)
}
