package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitorTimes = 5
const delay = 5 * time.Second

func main() {
	for {
		input := menu()

		switch input {
		case 1:
			monitor()
		case 2:
			printLogs()
		case 0:
			fmt.Println("LEAVING")
			os.Exit(0)
		default:
			fmt.Println("INVALID COMMAND")
			os.Exit(-1)
		}
	}

}

func menu() int {
	fmt.Println("\n1 - Start monitoring")
	fmt.Println("2 - Show logs")
	fmt.Println("0 - Leave")

	var input int
	fmt.Scan(&input)

	return input
}

func monitor() {
	fmt.Println("\nSTART MONITORING ...")

	websites := readFile("websites.txt")

	for times := 0; times < monitorTimes; times++ {
		for i, website := range websites {
			fmt.Println("TESTING", i)
			testWebsite(website)
		}
		time.Sleep(delay)
		fmt.Println()
	}
}

func testWebsite(website string) {
	res, err := http.Get(website)

	if err != nil {
		fmt.Println("ERROR ON REQUEST", err)
		os.Exit(-1)
	}

	if res.StatusCode == 200 {
		fmt.Println("Website: ", website, "200 OK")
		log(website, true)
	} else {
		fmt.Println("Website: ", website, res.StatusCode, "ERROR")
		log(website, false)
	}
}

func readFile(fileName string) []string {
	var array []string

	file, openErr := os.Open(fileName)

	if openErr != nil {
		fmt.Println("ERROR READING FILE!", openErr)
		os.Exit(-1)
	}

	reader := bufio.NewReader(file)

	for {
		line, readErr := reader.ReadString('\n')
		array = append(array, strings.TrimSpace(line))

		if readErr == io.EOF {
			break
		}
	}

	file.Close()

	return array
}

func log(website string, status bool) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("FILE NOT FOUND", err)
		os.Exit(-1)
	}

	file.WriteString("[" + time.Now().Format("02/01/2006 15:04:05") + "] " + website + " -online: " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func printLogs() {
	fmt.Println("LOGGING ... ")

	logs, err := ioutil.ReadFile("logs.txt")

	if err != nil {
		fmt.Println("ERROR READING LOGS")
		os.Exit(-1)
	}

	fmt.Println(string(logs))
}
