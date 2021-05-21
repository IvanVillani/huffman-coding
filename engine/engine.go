package engine

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"projects/golang/huffman/algo"
	"projects/golang/huffman/manager"
)

func Start() {
	fmt.Println()
	fmt.Println()
	fmt.Print("---Welcome to the exploration of Huffman Coding implementation!---\n\n")

	fmt.Print("Usage:\n\n")
	fmt.Print("\tBefore each encoding/decoding you will be prompted to enter/specify\n")
	fmt.Print("\ta few options in order to control what action and how it should be executed\n")
	fmt.Print("\twith your input information.\n\n")

	fmt.Print("The commands are:\n\n")

	fmt.Print("\tACTION_TYPE      DESCRIPTION\n")
	fmt.Print("\tencode           encodes input text with Huffman Coding algorithm\n")
	fmt.Print("\tdecode           decodes Huffman-encoded binary input to output text\n")
	fmt.Println()

	fmt.Print("\tINPUT_TYPE       DESCRIPTION\n")
	fmt.Print("\tfile             input data will be read from a file, located in the 'files' folder of this project directory; for multi-line text\n")
	fmt.Print("\tconsole          input data will be read directly from the console; works only for single line text\n")
	fmt.Print("\n")

	fmt.Print("\nIMPORTANT!:")
	fmt.Print("\n  - In order for algorithm to function correctly, the binary input,")
	fmt.Print("\n    specified during decode, must be encoded with this exact program!\n")
	fmt.Print("\n  - If /file/ INPUT_TYPE is selected, then the file must be provided")
	fmt.Print("\n    in the 'files' folder in the project directory\n")

	fmt.Print("\n")
	fmt.Print("Starting Huffman Coding program...\n\n")

	reader := bufio.NewReader(os.Stdin)

	var writer manager.Writer

	for {
		var encode bool
		var outputToFile bool
		var inputFromFile bool

		if err := ensureSpecifiedOption(&encode, "encode", "decode",
			`Please specify the action you want to perform (ex. "encode" or "decode"): `,
			"ERROR: You specified wrong action! Action '%s' is not supported! Please try again!",
			*reader); err != nil {
			fmt.Println(err)
			break
		}

		if err := ensureSpecifiedOption(&outputToFile, "file", "console",
			`Please specify the output destination type of your data (ex. "file" or "console"): `,
			"ERROR: You specified wrong output type! Output '%s' is not supported! Please try again!",
			*reader); err != nil {
			fmt.Println(err)
			break
		}

		if err := ensureSpecifiedOption(&inputFromFile, "file", "console",
			`Please specify the input location type of your data (ex. "file" or "console"): `,
			"ERROR: You specified wrong input type! Input '%s' is not supported! Please try again!",
			*reader); err != nil {
			fmt.Println(err)
			break
		}

		var inputData string

		if inputFromFile {
			if err := ensureSpecifiedFile(&inputData, *reader); err != nil {
				fmt.Println(err)
				break
			}
		} else {
			var err error

			fmt.Println("\nEnter text here:")

			inputData, err = reader.ReadString('\n')
			inputData = strings.Replace(inputData, "\n", "", -1)

			if err != nil {
				fmt.Println(err)
				break
			} else if inputData == "" {
				fmt.Println("Empty file/console input!")
				continue
			}
		}

		if outputToFile {
			var actionType string

			if encode {
				actionType = "encoding"
			} else {
				actionType = "decoding"
			}

			file, err := os.OpenFile(fmt.Sprintf("./files/huffman_%s_info", actionType),
				os.O_TRUNC|os.O_CREATE|os.O_WRONLY, fs.FileMode(0666))

			if err != nil {
				fmt.Println(err)
				return
			}

			fileWriter := manager.FileWriter{
				File: *file,
			}

			writer = &fileWriter
		} else {
			consoleWriter := manager.ConsoleWriter{}

			writer = &consoleWriter
		}

		if encode {
			err := algo.Encode(&writer, inputData, true)

			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			err := algo.Decode(&writer, inputData)

			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	fw, ok := writer.(*manager.FileWriter)

	if ok {
		err := fw.File.Close()

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func ensureSpecifiedOption(param *bool, opt1, opt2, msg, errMsg string, reader bufio.Reader) error {
	for {
		fmt.Print(msg)
		input, err := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)

		if err != nil {
			return err
		}

		switch input {
		case opt1:
			*param = true
			return nil
		case opt2:
			*param = false
			return nil
		default:
			fmt.Print("\n***")
			fmt.Printf(errMsg, input)
			fmt.Print("***\n\n")
			fmt.Print("Do you wish to continue?(y/n):")

			answer, err := reader.ReadString('\n')
			answer = strings.Replace(answer, "\n", "", -1)

			if err != nil {
				fmt.Println(err)
				return err
			}

			if answer == "y" || answer == "yes" {
				continue
			}

			return errors.New("EXIT: abort process")
		}
	}
}

func ensureSpecifiedFile(data *string, reader bufio.Reader) error {
	for {
		fmt.Print(`Please specify the name of the file, located in "files" folder (ex. "data.txt"): `)
		fileName, err := reader.ReadString('\n')
		fileName = strings.Replace(fileName, "\n", "", -1)

		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("./files/%s", fileName)

		if _, err := os.Stat(filePath); err == nil {
			res, err := ioutil.ReadFile(filePath)

			if err != nil {
				return err
			}

			*data = string(res)

			break
		} else if os.IsNotExist(err) {
			fmt.Print("\n***")
			fmt.Printf("ERROR: You specified wrong file name! File '%s' is not found in ./files! Please try again!", fileName)
			fmt.Print("***\n\n")
			fmt.Print("Do you wish to continue?(y/n):")

			answer, err := reader.ReadString('\n')
			answer = strings.Replace(answer, "\n", "", -1)

			if err != nil {
				fmt.Println(err)
				return err
			}

			if answer == "y" || answer == "yes" {
				continue
			}

			return errors.New("aborting process")
		} else {
			return err
		}
	}

	return nil
}
