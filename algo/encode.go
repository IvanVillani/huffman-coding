package algo

import (
	"bytes"
	"fmt"
	"math"
	"projects/golang/huffman/manager"
	"projects/golang/huffman/visual"
	"sort"
	"strconv"
)

var codeArrays []manager.ValuesCodeLengthsArray
var finalValueCodeLenghts []manager.ValueCodeLenght
var valueCodes map[byte]uint64

var maxFrequency int

var valuesFreq map[byte]int
var values manager.Values

func Encode(writer *manager.Writer, toEncode string, withPrefixForDecoding bool) error {
	calculateFrequency(toEncode)

	sortValuesFreq()

	str := visual.VisualizeFrequencyTable(&values, &valuesFreq, maxFrequency)

	err := (*writer).Write(str)

	if err != nil {
		return err
	}

	str = findValuesCodeLenghts()

	err = (*writer).Write(str)

	if err != nil {
		return err
	}

	GenerateHuffmanCodes(&finalValueCodeLenghts, &valueCodes)

	var b bytes.Buffer

	b.WriteString("///////////STEP 3///////////\n\n")
	b.WriteString("\tThe next step is to create huffman codes from\n")
	b.WriteString("\tthe given code lenghts, using binary incrementation and shifting.\n")
	b.WriteString("\tThe Huffman codes should be with unique prefixes to enable decoding.\n\n")

	err = (*writer).Write(b.String())

	if err != nil {
		return err
	}

	str = visual.VisualizeResultTable(&finalValueCodeLenghts, &valueCodes)

	err = (*writer).Write(str)

	if err != nil {
		return err
	}

	str = constructEncodedMessage(toEncode)

	err = (*writer).Write(str)

	if err != nil {
		return err
	}

	return nil
}

func calculateFrequency(str string) {
	valuesFreq = make(map[byte]int)
	values = make(manager.Values, 0)

	maxFrequency = 1

	byteArray := []byte(str)

	for _, char := range byteArray {
		if val, ok := valuesFreq[char]; ok {
			if len(fmt.Sprint(val+1)) > maxFrequency {
				maxFrequency = len(fmt.Sprint(val + 1))
			}
			valuesFreq[char] = val + 1
		} else {
			valuesFreq[char] = 1
			values = append(values, char)
		}
	}
}

func sortValuesFreq() {
	sort.Sort(values)
}

func findValuesCodeLenghts() string {
	initializeStartingArray()

	remaining := len(codeArrays)

	var lastArrayIndex int

	var b bytes.Buffer

	b.WriteString("///////////STEP 2///////////\n\n")
	b.WriteString("\tThe next step is to find the code lenghts of the values.\n")
	b.WriteString("\tIn order to do that, we perform merging of lists. We create\n")
	b.WriteString("\tthese lists by assigning zero code lenghts to the values.\n")
	b.WriteString("\tThe merging happens between the two lists with the smallest\n")
	b.WriteString("\tfrequencies and by that the frequency of the new list becomes\n")
	b.WriteString("\tthe sum of the merged lists frequencies.\n\n")

	iterator := 1

	for {
		b.WriteString(fmt.Sprintf("--- ITERATION %d:\n", iterator))

		for i, val := range codeArrays {
			if val.Frequency != -1 {
				b.WriteString(fmt.Sprintf("\t        List %d:\n\t\t       ", i))
				for _, v := range val.ValuesCodeLenghts {
					b.WriteString(fmt.Sprintf("['%s' %d] ", string(v.Value), v.ValueCodeLenght))
				}
				b.WriteString(fmt.Sprintf("with freq = %d\n", val.Frequency))
			}
		}

		if remaining == 1 {
			break
		}

		mergeCodeArrays(&lastArrayIndex, &b)

		remaining--

		iterator++
	}

	b.WriteString("\n\n\n")

	finalValueCodeLenghts = make([]manager.ValueCodeLenght, 0)

	for i := 0; i < len(codeArrays[lastArrayIndex].ValuesCodeLenghts); i++ {
		index := 0
		min := math.MaxInt64

		for i := len(codeArrays[lastArrayIndex].ValuesCodeLenghts) - 1; i >= 0; i-- {
			if codeArrays[lastArrayIndex].ValuesCodeLenghts[i].ValueCodeLenght == -1 {
				continue
			}

			if min > codeArrays[lastArrayIndex].ValuesCodeLenghts[i].ValueCodeLenght {
				min = codeArrays[lastArrayIndex].ValuesCodeLenghts[i].ValueCodeLenght
				index = i
			}
		}

		finalValueCodeLenghts = append(finalValueCodeLenghts, manager.ValueCodeLenght{
			Value:           codeArrays[lastArrayIndex].ValuesCodeLenghts[index].Value,
			ValueCodeLenght: codeArrays[lastArrayIndex].ValuesCodeLenghts[index].ValueCodeLenght,
		})

		codeArrays[lastArrayIndex].ValuesCodeLenghts[index].ValueCodeLenght = -1
	}

	return b.String()
}

func initializeStartingArray() {
	codeArrays = make([]manager.ValuesCodeLengthsArray, 0)

	for _, val := range values {
		newValuesCodeLenghtsArray := manager.ValuesCodeLengthsArray{
			ValuesCodeLenghts: make([]manager.ValueCodeLenght, 0),
			Frequency:         valuesFreq[val],
		}

		newValuesCodeLenghtsArray.ValuesCodeLenghts = append(newValuesCodeLenghtsArray.ValuesCodeLenghts, manager.ValueCodeLenght{
			Value:           val,
			ValueCodeLenght: 0,
		})

		codeArrays = append(codeArrays, newValuesCodeLenghtsArray)
	}
}

func mergeCodeArrays(lastArrayIndex *int, b *bytes.Buffer) {
	maxValue := math.MaxInt64

	firstArrayIndex := maxValue
	secondArrayIndex := maxValue

	for i := len(codeArrays) - 1; i >= 0; i-- {
		if codeArrays[i].Frequency == -1 && codeArrays[i].ValuesCodeLenghts == nil {
			continue
		}

		if codeArrays[i].Frequency < maxValue && firstArrayIndex == maxValue {
			firstArrayIndex = i
		} else if codeArrays[i].Frequency < maxValue && firstArrayIndex < maxValue {
			secondArrayIndex = i
			break
		}
	}

	for i := len(codeArrays) - 1; i >= 0; i-- {
		if codeArrays[i].Frequency == -1 && codeArrays[i].ValuesCodeLenghts == nil {
			continue
		}

		if codeArrays[i].Frequency < codeArrays[firstArrayIndex].Frequency {
			secondArrayIndex = firstArrayIndex

			firstArrayIndex = i
		} else if codeArrays[i].Frequency < codeArrays[secondArrayIndex].Frequency && i != firstArrayIndex {
			secondArrayIndex = i
		}
	}

	b.WriteString(fmt.Sprintf("  ***Based on the combined frequency of the lists we decide to merge list %d and list %d\n\n\n\n", firstArrayIndex, secondArrayIndex))

	codeArrays[firstArrayIndex].Frequency = codeArrays[firstArrayIndex].Frequency + codeArrays[secondArrayIndex].Frequency
	codeArrays[firstArrayIndex].ValuesCodeLenghts = append(codeArrays[firstArrayIndex].ValuesCodeLenghts,
		codeArrays[secondArrayIndex].ValuesCodeLenghts...)

	for i := range codeArrays[firstArrayIndex].ValuesCodeLenghts {
		codeArrays[firstArrayIndex].ValuesCodeLenghts[i].ValueCodeLenght++
	}

	codeArrays[secondArrayIndex].ValuesCodeLenghts = nil
	codeArrays[secondArrayIndex].Frequency = -1

	*lastArrayIndex = firstArrayIndex
}

func constructEncodedMessage(toEncode string) string {
	var res bytes.Buffer

	var b bytes.Buffer

	b.WriteString("///////////STEP 4///////////\n\n")
	b.WriteString("\tThe final step is to assemble the encoded binary representation of\n")
	b.WriteString("\tthe input text. Except the encoded text, a prefix must be appended to\n")
	b.WriteString("\tthe result, containing information for decoding it.\n\n")

	prefix := createPrefix(&b)

	res.WriteString(prefix)

	for i := 0; i < len(toEncode); i++ {
		char := toEncode[i]

		binary := strconv.FormatInt(int64(valueCodes[char]), 2)

		res.WriteString(binary)
	}

	b.WriteString("  And the final encoded result is:\n")
	b.WriteString(res.String())
	b.WriteString("\n")

	return b.String()
}

func createPrefix(b *bytes.Buffer) string {
	var res bytes.Buffer

	countOfDistinct := strconv.FormatInt(int64(len(finalValueCodeLenghts)), 2)

	b.WriteString("  --- The number of distinct characters in the encoded text (16-bit number):\n\t[")

	for i := 0; i < 16-len(countOfDistinct); i++ {
		b.WriteString("0")
		res.WriteString("0")
	}

	b.WriteString(countOfDistinct)
	b.WriteString("]\n")
	res.WriteString(countOfDistinct)

	maxCodeLenght := len(strconv.FormatInt(int64(finalValueCodeLenghts[len(finalValueCodeLenghts)-1].ValueCodeLenght), 2))

	maxCodeLenghtBinary := strconv.FormatInt(int64(maxCodeLenght), 2)

	b.WriteString("  --- The maximum lenght of the encoded representation of code lenghts:\n\t[")

	for i := 0; i < 8-len(maxCodeLenghtBinary); i++ {
		b.WriteString("0")
		res.WriteString("0")
	}

	b.WriteString(maxCodeLenghtBinary)
	b.WriteString("]\n")
	res.WriteString(maxCodeLenghtBinary)

	b.WriteString("  --- The sequence of the characters and their associated code lenght, for example the first one:\n\t")

	for i := range finalValueCodeLenghts {
		value := fmt.Sprintf("%08b", finalValueCodeLenghts[i].Value)

		res.WriteString(value)

		codeLenght := strconv.FormatInt(int64(finalValueCodeLenghts[i].ValueCodeLenght), 2)

		for i := 0; i < maxCodeLenght-len(codeLenght); i++ {
			res.WriteString("0")
		}

		res.WriteString(codeLenght)

		if i == 0 {
			b.WriteString(fmt.Sprintf("character -> [%s]   code lenght -> [%s]\n\n", value, codeLenght))
		}
	}

	return res.String()
}
