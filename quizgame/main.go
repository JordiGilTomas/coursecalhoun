package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
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
	reader := bufio.NewScanner(os.Stdin)
	success := 0
	totalQuestions := 0
	csvFile := loadQuestionsFile()

	if shuffle {
		shuffleQuestions(csvFile)
		csvFile = loadQuestionsFile()
	}

	for {
		if totalQuestions == 0 {
			fmt.Println("Presiona Enter para empezar")
		} else {
			fmt.Println("Presiona Enter para la siguiente pregunta")
		}
		bufio.NewReader(os.Stdin).ReadRune()

		question, correctAnswer, err := loadNextQuestion(csvFile)
		if err != nil {
			break
		}
		totalQuestions++
		fmt.Print(question, " Respuesta: ")

		timer := time.NewTimer(time.Duration(countDown) * time.Second)
		go func() {
			<-timer.C
			fmt.Println("\nSe acabó el tiempo!!!\nGame Over")
			gameOver(totalQuestions, success)
			os.Exit(3)
		}()

		checkAnswer(reader, correctAnswer, &success)
		timer.Stop()
	}
	fmt.Println("\nJuego Completado! Enhorabuena!!!")
	gameOver(totalQuestions, success)
}

func checkAnswer(reader *bufio.Scanner, correctAnswer string, success *int) {
	reader.Scan()
	userAnswer := reader.Text()
	if correctAnswer == strings.TrimSpace(strings.ToLower(userAnswer)) {
		fmt.Println("Respuesta correcta")
		*success++
	} else {
		fmt.Println("Respuesta incorrecta")
	}
}

func loadNextQuestion(csvFile *csv.Reader) (question, correctAnswer string, err error) {
	data, err := csvFile.Read()
	if err != nil {
		return
	}

	question = string(data[0])
	correctAnswer = string(data[1])

	return

}

func loadQuestionsFile() *csv.Reader {
	file, err := os.Open(questionsFile)

	if err != nil {
		fmt.Println("No se pudo cargar el archivo de preguntas")
		os.Exit(1)
	}
	return csv.NewReader(file)
}

func shuffleQuestions(csvFile *csv.Reader) {
	file, err := os.Create(questionsFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	shuffleData, err := csvFile.ReadAll()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffleData), func(i, j int) {
		shuffleData[i], shuffleData[j] = shuffleData[j], shuffleData[i]
	})
	data := csv.NewWriter(file)
	data.WriteAll(shuffleData)
}

func gameOver(totalQuestions, success int) {
	fmt.Println("\nPreguntas realizadas:", totalQuestions)
	fmt.Println("Preguntas acertadas:", success)
}
