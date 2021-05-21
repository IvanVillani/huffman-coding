package algo

import (
	"bytes"
	"fmt"
	"math"
	"projects/golang/huffman/manager"
	"projects/golang/huffman/visual"
	"strconv"
)

var countOfDistinct int64
var codeLenghtSize int64

var decodedCodeLenghts []manager.ValueCodeLenght
var resultCodesMap map[byte]uint64
var decodingMap map[uint64]byte

var dataStartIndex int

var decodedResult string

func Decode(writer *manager.Writer, toDecode string) error {
	var b bytes.Buffer

	if err := splitAndExtractPrefix(toDecode, &b); err != nil {
		return err
	}

	b.WriteString("\n\n\n")

	GenerateHuffmanCodes(&decodedCodeLenghts, &resultCodesMap)

	b.WriteString("///////////STEP 2///////////\n\n")
	b.WriteString("\tThe next step is to create huffman codes from\n")
	b.WriteString("\tthe given code lenghts, using binary incrementation and shifting.\n")
	b.WriteString("\tThe Huffman codes should be with unique prefixes.\n\n")

	err := (*writer).Write(b.String())

	if err != nil {
		return err
	}

	b.Reset()

	str := visual.VisualizeResultTable(&decodedCodeLenghts, &resultCodesMap)

	err = (*writer).Write(str)

	if err != nil {
		return err
	}

	mirrorCodesMap()

	b.WriteString("///////////STEP 3///////////\n\n")
	b.WriteString("\tNow that we have everything that we need we start\n")
	b.WriteString("\titerating over the binary sequence and decode the string.\n")
	b.WriteString("\tA match is when the largest binary sequence persist in the table.\n\n")

	if err := extractData(toDecode); err != nil {
		return err
	}

	b.WriteString("  And the final decoded result is:\n")
	b.WriteString(decodedResult)
	b.WriteString("\n")

	err = (*writer).Write(b.String())

	if err != nil {
		return err
	}

	return nil
}

func splitAndExtractPrefix(toDecode string, b *bytes.Buffer) error {
	var err error

	b.WriteString("///////////STEP 2///////////\n\n")
	b.WriteString("\tThe process of decoding begins with extracting prefix\n")
	b.WriteString("\tfrom the encoded string. Doing this will allow us\n")
	b.WriteString("\tto recreate the tables of code lenghts and then the Huffman codes.\n\n")

	countOfDistinctBinary := toDecode[0:16]

	countOfDistinct, err = strconv.ParseInt(countOfDistinctBinary, 2, 17)

	if err != nil {
		panic(err)
	}

	b.WriteString("  --- Taking the first 16 bits, which are showing the number of distinct characters:\n\t")
	b.WriteString(fmt.Sprintf("binary -> [%s]   count of distinct -> %d\n", countOfDistinctBinary, countOfDistinct))

	codeLenghtSizeBinary := toDecode[16:24]

	codeLenghtSize, err = strconv.ParseInt(codeLenghtSizeBinary, 2, 9)

	if err != nil {
		panic(err)
	}

	b.WriteString("  --- Taking the next 8 bits, which are showing the maximum lenght of the encoded representation of code lenghts:\n\t")
	b.WriteString(fmt.Sprintf("binary -> [%s]   code lenght size -> %d\n", codeLenghtSizeBinary, codeLenghtSize))

	dataStartIndex = 24

	decodedCodeLenghts = make([]manager.ValueCodeLenght, 0)

	b.WriteString("  --- Taking the sequence of the characters and their associated code lenght, for example the first one:\n\t")

	for i := 0; i < int(countOfDistinct); i++ {
		charBinary := toDecode[dataStartIndex : dataStartIndex+8]

		char, err := strconv.ParseInt(charBinary, 2, 9)

		if err != nil {
			panic(err)
		}

		codeLenghtBinary := toDecode[dataStartIndex+8 : dataStartIndex+8+int(codeLenghtSize)]

		bitSize := getMinBitSize(codeLenghtSize)

		codeLenght, err := strconv.ParseInt(codeLenghtBinary, 2, bitSize)

		if err != nil {
			panic(err)
		}

		decodedCodeLenghts = append(decodedCodeLenghts, manager.ValueCodeLenght{
			Value:           byte(char),
			ValueCodeLenght: int(codeLenght),
		})

		if i == 0 {
			b.WriteString(fmt.Sprintf("character in binary -> [%s] = [%s]\n\tcode lenght in binary -> [%s] = %d\n", charBinary, string(byte(char)), codeLenghtBinary, int(codeLenght)))
		}

		dataStartIndex = dataStartIndex + 8 + int(codeLenghtSize)
	}

	return nil
}

func mirrorCodesMap() {
	decodingMap = make(map[uint64]byte)

	for k, v := range resultCodesMap {
		decodingMap[v] = k
	}
}

func extractData(toDecode string) error {
	var helper bytes.Buffer
	var result bytes.Buffer

	for i := dataStartIndex; i < len(toDecode); i++ {
		helper.WriteString(string(toDecode[i]))

		portion := helper.String()

		bitSize := getMinBitSize(int64(len(portion)))

		key, err := strconv.ParseInt(portion, 2, bitSize)

		if err != nil {
			panic(err)
		}

		if val, ok := decodingMap[uint64(key)]; ok {
			result.WriteByte(val)
			helper.Reset()
		}
	}

	decodedResult = result.String()

	return nil
}

func getMinBitSize(desired int64) int {
	var bitSize int

	for i := 1.0; i < 8.0; i += 2.0 {
		if math.Pow(2.0, i) > float64(desired) {
			bitSize = int(math.Pow(2.0, i))
			break
		}
	}

	return bitSize
}
