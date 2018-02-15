package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	filePtr := flag.String("file", "./problems.csv", "path to questions csv file")
	timePtr := flag.Int("time", 30, "time limit for quix in ms")

	timer := time.NewTimer(time.Duration(*timePtr) * time.Second)

	f, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		fmt.Printf("Error reading file: %s", err)
	}

	r := csv.NewReader(strings.NewReader(string(f)))
	records, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Error parsing csv: %s", err)
	}

	reader := bufio.NewReader(os.Stdin)
	var correct int

	start := time.Now()

	for _, q := range records {
		fmt.Printf("%v=\n", q[0])
		text, _ := reader.ReadString('\n')
		a, err := strconv.Atoi(strings.Replace(text, "\n", "", -1))

		if err != nil {
			log.Fatal("Error: Invalid user input %s", err)
		}
		b, _ := strconv.Atoi(q[1])

		if a == b {
			correct++
		}

		<-timer.C
		fmt.Print("expired")

		fmt.Println(a == b)
	}

	timer.Stop()

	t := time.Now()
	fmt.Printf("You took %v\n", t.Sub(start))
	fmt.Printf("You got %v/%v correct\n", correct, len(records))

}
