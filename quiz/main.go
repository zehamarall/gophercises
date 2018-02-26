package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Problem struct to store question an answer.
type Problem struct {
	question string
	answer   string
}

var ch = make(chan bool, 1)

var csvFilename *string
var timeout *int

func init() {
	csvFilename = flag.String("csv", "problems.csv", "Usage -csv=problems.csv")
	timeout = flag.Int("timeout", 30, "timeout for all questions -timeout=30")
	flag.Parse()
}

func main() {

	csvFile, err := os.Open(*csvFilename)

	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	problems := readQuestions(lines)
	var countCorret int

	go makeQuiz(problems, &countCorret)
	defer close(ch)
	timer := time.NewTimer(time.Duration(*timeout) * time.Second)

	select {
	case <-ch:
		fmt.Println("Congratulation you winner")
	case <-timer.C:
		fmt.Println(" Timeout of game")
	}

	fmt.Printf("Number of question correct %v of %v \n", countCorret, len(problems))
}
func makeQuiz(problems []Problem, countCorret *int) {
	for i, q := range problems {
		fmt.Printf("Question number %v: %s:", i, q.question)
		var read string
		fmt.Scanf("%s", &read)
		if read == q.answer {
			*countCorret++
		}
	}
	ch <- true
}

func readQuestions(lines [][]string) []Problem {
	problems := make([]Problem, len(lines))

	for i, line := range lines {
		problems[i] = Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems

}
