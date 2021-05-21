package visual

import (
	"bytes"
	"fmt"
	"projects/golang/huffman/manager"
	"strconv"
)

type Visualizer struct {
	Writer manager.Writer
}

func VisualizeFrequencyTable(sortedValues *manager.Values, frequencyMap *map[byte]int, maxFrequency int) string {
	var b bytes.Buffer

	title := "Frequency"

	b.WriteString("///////////STEP 1///////////\n\n")
	b.WriteString("\tThe process of encoding begins with separation of the ASCII symbols\n")
	b.WriteString("\tand then counting the frequency of every symbol.\n")
	b.WriteString("\tBelow are the sorted symbols of your input text, sorted\n")
	b.WriteString("\tin according to the ASCII table and their frequencies.\n\n")

	width := 17 + maxFrequency

	printFrame(width, &b)

	printTitle(width, &b, title)

	printEmptyTableLine(width, &b)

	for _, v := range *sortedValues {
		b.WriteString(fmt.Sprintf("|    [%s] -> %d", string(v), (*frequencyMap)[v]))

		diff := 5 - len(fmt.Sprint((*frequencyMap)[v])) + (maxFrequency - 1)

		for i := 0; i < diff; i++ {
			b.WriteString(" ")
		}
		b.WriteString("|\n")
	}

	printEmptyTableLine(width, &b)

	printFrame(width, &b)

	b.WriteString("\n\n\n")

	return b.String()
}

func VisualizeResultTable(finalValueCodeLenghts *[]manager.ValueCodeLenght, valueCodes *map[byte]uint64) string {
	var b bytes.Buffer

	title := "Huffman Codes"

	maxCodeLenght := len(fmt.Sprint((*finalValueCodeLenghts)[len(*finalValueCodeLenghts)-1].ValueCodeLenght))

	width := 26 + maxCodeLenght + len(*finalValueCodeLenghts)

	printFrame(width, &b)

	printTitle(width, &b, title)

	printEmptyTableLine(width, &b)

	for _, v := range *finalValueCodeLenghts {
		b.WriteString(fmt.Sprintf("|    [%s] -> %d", string(v.Value), v.ValueCodeLenght))

		diff1 := 5 - len(fmt.Sprint(v.ValueCodeLenght))

		for i := 0; i < diff1; i++ {
			b.WriteString(" ")
		}

		b.WriteString(fmt.Sprintf("===> %s", strconv.FormatUint((*valueCodes)[v.Value], 2)))

		diff2 := len(*finalValueCodeLenghts) + 4 - len(strconv.FormatUint((*valueCodes)[v.Value], 2))

		for i := 0; i < diff2; i++ {
			b.WriteString(" ")
		}

		b.WriteString("|\n")
	}

	printEmptyTableLine(width, &b)

	printFrame(width, &b)

	b.WriteString("\n\n\n")

	return b.String()
}

func printFrame(n int, b *bytes.Buffer) {
	(*b).WriteString(" ")

	for i := 0; i < n-2; i++ {
		(*b).WriteString("-")
	}

	(*b).WriteString(" \n")
}

func printEmptyTableLine(n int, b *bytes.Buffer) {
	(*b).WriteString("|")

	for i := 0; i < n-2; i++ {
		(*b).WriteString(" ")
	}

	(*b).WriteString("|\n")
}

func printTitle(n int, b *bytes.Buffer, title string) {
	(*b).WriteString("|")

	padding := int((n - len(title)) / 2)

	for i := 0; i < padding; i++ {
		(*b).WriteString(" ")
	}

	(*b).WriteString(title)

	for i := padding + len(title); i < n-2; i++ {
		(*b).WriteString(" ")
	}

	(*b).WriteString("|\n")
}
