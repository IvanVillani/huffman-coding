package algo

import (
	"projects/golang/huffman/manager"
)

func GenerateHuffmanCodes(valueCodeLenghts *[]manager.ValueCodeLenght, destinationMap *map[byte]uint64) {
	huffmanCodeCounter := uint64(0)
	codeLenghtCounter := 1

	*destinationMap = make(map[byte]uint64)

	for i := 0; i < len(*valueCodeLenghts); i++ {
		for {
			if (*valueCodeLenghts)[i].ValueCodeLenght == codeLenghtCounter {
				break
			}
			huffmanCodeCounter = huffmanCodeCounter << 1
			codeLenghtCounter++
		}
		(*destinationMap)[(*valueCodeLenghts)[i].Value] = huffmanCodeCounter
		huffmanCodeCounter++
		// huffmanCodeCounter = huffmanCodeCounter << 1
	}
}
