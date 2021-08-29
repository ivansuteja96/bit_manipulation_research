package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	var (
		split              int
		totalGeneratedData int
	)

	flag.IntVar(&totalGeneratedData, "total_data", 0, "Total Data")
	flag.IntVar(&split, "split", 1, "Split Into")
	flag.Parse()
	if split < 1 {
		log.Println("Split must >= 1")
		return
	}
	if totalGeneratedData < 1 {
		log.Println("total_data must >= 1")
		return
	}

	part := totalGeneratedData / split

	for i := 1; i <= split; i++ {
		startIndex := part * (i - 1)
		endIndex := part * i
		if endIndex > totalGeneratedData {
			endIndex = totalGeneratedData
		}
		subTotalGeneratedData := endIndex - startIndex

		subMain(startIndex, i, subTotalGeneratedData)
	}

}

func subMain(i1 int, split int, totalGeneratedData int) {
	starti1 := i1

	var i int

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	file1, _ := os.Create(fmt.Sprintf("result1-%+v.csv", split))
	defer file1.Close()
	file2, _ := os.Create(fmt.Sprintf("result2-%+v.csv", split))
	defer file2.Close()
	file3, _ := os.Create(fmt.Sprintf("result3-%+v.csv", split))
	defer file3.Close()

	writer1 := csv.NewWriter(file1)
	defer writer1.Flush()
	writer2 := csv.NewWriter(file2)
	defer writer2.Flush()
	writer3 := csv.NewWriter(file3)
	defer writer3.Flush()

	writer1.Write([]string{"id", "type_id", "value", "source", "status"})
	writer2.Write([]string{"id", "type_id", "value", "source"})
	writer3.Write([]string{"id", "type_id", "value", "source"})

	for {
		myUUID := uuid.New()
		value := myUUID.String()
		typeID := r.Intn(25) + 1
		sourceProbability := r.Intn(100) + 1

		var sourceCount int
		if sourceProbability <= 25 {
			sourceCount = 1
		} else if sourceProbability <= 75 {
			sourceCount = 2
		} else {
			sourceCount = 3
		}

		if i+sourceCount > totalGeneratedData {
			sourceCount = totalGeneratedData - i
		}

		var listSourceID []int

		mapSourceExist := make(map[int]bool)
		var sourceBinary int64
		for sourceCount != 0 {
			sourceID := r.Intn(25) + 1
			if mapSourceExist[sourceID] {
				continue
			}
			mapSourceExist[sourceID] = true
			listSourceID = append(listSourceID, sourceID)
			sourceCount--
			i++
			sourceBinary += int64(math.Pow(2, float64(sourceID-1)))
			writer1.Write([]string{fmt.Sprintf("%+v", i1), fmt.Sprintf("%+v", typeID), value, fmt.Sprintf("%+v", sourceID), "TRUE"})
			i1++
		}

		writer2.Write([]string{fmt.Sprintf("%+v", i1), fmt.Sprintf("%+v", typeID), value, fmt.Sprintf("%+v", sourceBinary)})
		sourceBit, _ := SourceToStringBit(listSourceID, 1000)
		writer3.Write([]string{fmt.Sprintf("%+v", i1), fmt.Sprintf("%+v", typeID), value, sourceBit})
		fmt.Printf("Processing on : %+v/%+v, split : %+v\n", i1-starti1, totalGeneratedData, split)

		if i == totalGeneratedData {
			i = 0
			break
		}
	}
}

//SourceToStringBit
func SourceToStringBit(source []int, totalBit int) (stringBit string, err error) {

	if len(source) == 0 {
		return stringBit, fmt.Errorf("source must be filled")
	}

	keys := make(map[int]bool)

	maxSource, minSource := math.MinInt32, math.MaxInt32
	for _, entry := range source {
		if !keys[entry] {
			keys[entry] = true
			if maxSource < entry {
				maxSource = entry
			}
			if minSource > entry {
				minSource = entry
			}
		}
	}

	if minSource <= 0 {
		return stringBit, fmt.Errorf("source <= 0")
	} else if maxSource > totalBit {
		return stringBit, fmt.Errorf("source > totalBit")
	}

	for i := 1; i <= maxSource; i++ {
		if keys[i] {
			stringBit = "1" + stringBit
			continue
		}
		stringBit = "0" + stringBit
	}

	stringBit = StringBitAppender(stringBit, totalBit)

	return
}

func StringBitAppender(stringBit string, totalBit int) (res string) {

	validation := fmt.Sprintf("%%0%ds", totalBit)

	res = fmt.Sprintf(validation, stringBit)

	return
}
