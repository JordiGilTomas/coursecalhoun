package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

var questionsFile string
var countDown int
var shuffle bool

func init() {
	flag.StringVar(&questionsFile, "file", "questions.csv", "Archivo de preguntas")
	flag.IntVar(&countDown, "countdown", 30, "Tiempo cuenta atrás")
	flag.BoolVar(&shuffle, "shuffle", false, "Mezcla las preguntas")

}

func main() {
	flag.Parse()
	var file []byte
	var err error
	file, err = ioutil.ReadFile(questionsFile)
	reader := bufio.NewScanner(os.Stdin)

	if err != nil {
		fmt.Println("No se pudo cargar el archivo de preguntas")
		return
	}

	success := 0
	totalQuestions := 0
	var shuffleData [][]string
	if shuffle {
		csvFile := csv.NewReader(strings.NewReader(string(file)))
		file, err := os.Create(questionsFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		shuffleData, err = csvFile.ReadAll()
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(shuffleData), func(i, j int) {
			shuffleData[i], shuffleData[j] = shuffleData[j], shuffleData[i]
		})
		data := csv.NewWriter(file)
		data.WriteAll(shuffleData)
	}
	csvFile := csv.NewReader(strings.NewReader(string(file)))
	for {
		data, err := csvFile.Read()
		if err != nil {
			break
		}
		if totalQuestions == 0 {
			reader.Bytes()
		}

		question := string(data[0])
		correctAnswer := string(data[1])

		if totalQuestions == 0 {
			fmt.Println("Presiona Enter para empezar")
		} else {
			fmt.Println("Presiona Enter para la siguiente pregunta")
		}
		bufio.NewReader(os.Stdin).ReadRune()

		totalQuestions++
		fmt.Print(question, " Respuesta: ")
		timer := time.NewTimer(time.Duration(countDown) * time.Second)
		go func() {
			<-timer.C
			fmt.Println("\nSe acabó del tiempo!!!\nGame Over")
			gameOver(totalQuestions, success)
			os.Exit(3)
		}()

		reader.Scan()
		userAnswer := reader.Text()
		if correctAnswer == strings.TrimSpace(strings.ToLower(userAnswer)) {
			fmt.Println("Respuesta correcta")
			success++
		} else {
			fmt.Println("Respuesta incorrecta")
		}
		timer.Stop()
	}
	fmt.Println("\nJuego Completado! Enhorabuena!!!")
	gameOver(totalQuestions, success)
}

func gameOver(totalQuestions, success int) {
	fmt.Println("\nPreguntas realizadas:", totalQuestions)
	fmt.Println("Preguntas acertadas:", success)
}
