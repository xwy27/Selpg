package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	flag "github.com/spf13/pflag"
)

/*============================== Global Var ==============================*/
const MAX_INT = int(^uint(0) >> 1)

var fileName string
var err error

/*============================== Option Args ==============================*/
var startPage = flag.Int("s", -1, "Start page of file")
var endPage = flag.Int("e", -1, "End page of file")
var linePage = flag.Int("l", 72, "lines in one page")
var flagPage = flag.Bool("f", false, "flag splits page")
var printDst = flag.String("d", "default", "name of printer")

/*============================== Error Handle ==============================*/
func errorHandler(err error) {
	if err != nil {
		fmt.Println("\n Error!\n", err)
		os.Exit(1)
	}
}

/*============================== MAIN FUNCTION ==============================*/
func main() {
	flag.Parse()

	// handle options
	err = handle_args()
	errorHandler(err)

	// read selected file
	str, err := readAndWrite()
	errorHandler(err)
	fmt.Println(str)

	fmt.Println("Print Completed!")
}

/*============================== Args Handle FUNCTION ==============================*/

func handle_args() error {
	// Option Args num
	if *startPage == -1 || *endPage == -1 {
		return errors.New("Not enough arguments!\n")
	}

	// Start Page Option
	if *startPage < 1 || *startPage > (MAX_INT-1) {
		return errors.New("Start page invalid\n")
	}

	// End Page Option
	if *endPage < 1 || *endPage > (MAX_INT-1) ||
		*endPage < *startPage {
		return errors.New("End page invalid\n")
	}

	// Page Split Option
	if *flagPage == false {
		// line split mode
		if *linePage < 1 || *linePage > (MAX_INT-1) ||
			*linePage < *startPage {
			return errors.New("Page length invalid\n")
		}
	}

	return nil // Options right
}

/*============================== Read & Write FUNCTION ==============================*/

func readAndWrite() (string, error) {
	result := ""
	pageCount := 1
	lineCount := 0
	reader := bufio.NewReader(os.Stdin)

	// set the input source
	if flag.NArg() > 0 {
		// accept one file each time
		fileName = flag.Arg(0)
		file, err := os.Open(fileName)

		if err != nil {
			return "", err
		}

		defer file.Close()

		reader = bufio.NewReader(file)
	}

	// process the input
	if *flagPage {
		pageCount = 1
		for {
			str, err := reader.ReadString('\f')
			if err == io.EOF {
				break
			} else if err != nil {
				return "", err
			}
			pageCount++
			if pageCount >= *startPage && pageCount <= *endPage {
				result = strings.Join([]string{result, str}, "")
			}
		}
	} else {
		pageCount = 1
		lineCount = 0

		for {
			str, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				return "", err
			}
			lineCount++
			if lineCount > *linePage {
				pageCount++
				lineCount = 1
			}
			if pageCount >= *startPage && pageCount <= *endPage {
				result = strings.Join([]string{result, str}, "")
			}
		}
	}

	// handle invalid input option
	if pageCount < *startPage {
		msg := fmt.Sprintf("start page: (%d) greater than total pages: (%d)",
			*startPage, pageCount)
		return "", errors.New(msg)
	} else if pageCount < *endPage {
		msg := fmt.Sprintf("end page: (%d) greater than total pages: (%d)",
			*endPage, pageCount)
		return "", errors.New(msg)
	}

	// set the output source
	if *printDst != "default" {
		cmd := exec.Command("lp", "-d"+*printDst)
		cmd.Stdin = strings.NewReader(result)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return "", errors.New(fmt.Sprint(err) + " : " + stderr.String())
		}
	}

	return result, nil
}
