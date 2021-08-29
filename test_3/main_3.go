package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TableWithBitString struct {
	ID     int64  `db:"id"`
	TypeID int    `db:"type_id"`
	Value  string `db:"value"`
	Source string `db:"source"`
}

func main() {

	var (
		mutex     sync.Mutex
		lineLock  sync.Mutex
		dbConfig  string
		totalLine int64
		qps       int
		csvFile   string
	)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	flag.StringVar(&dbConfig, "db_config", "", "DB Config")
	flag.StringVar(&csvFile, "csv_file", "", "CSV File")
	flag.Int64Var(&totalLine, "total_line", 0, "Total Line")
	flag.IntVar(&qps, "qps", 0, "QPS")

	flag.Parse()

	if dbConfig == "" {
		log.Println("Empty db_config")
		return
	} else if csvFile == "" {
		log.Println("Empty csv_file")
		return
	} else if totalLine == 0 {
		log.Println("Empty total_line")
		return
	} else if qps == 0 {
		log.Println("Empty QPS")
		return
	}

	db, err := sqlx.Connect("postgres", dbConfig)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(0)

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	br := bufio.NewReader(file)
	br.ReadString('\n')

	var (
		i             int
		countChecked  int64
		count         int64
		countDuration time.Duration
	)

	startTime := time.Now()
	for {
		go func() {
			var (
				typeID int64
				source int64
				value  string
			)

			lineLock.Lock()
			lines, err := br.ReadString('\n')
			if err != nil {
				log.Println(err)
				lineLock.Unlock()
				return
			}
			lineLock.Unlock()

			col := strings.Split(lines, ",")
			strTypeID, strSource := col[1], col[3]
			value = col[2]
			typeID, _ = strconv.ParseInt(strTypeID, 10, 64)
			source, _ = strconv.ParseInt(strSource, 10, 64)
			probability := r.Intn(100) + 1
			if probability <= 90 {
				myUUID := uuid.New()
				value = myUUID.String()
				typeID = int64(r.Intn(25) + 1)
				source = int64(r.Intn(25) + 1)
			}

			queryInput := TableWithBitString{
				TypeID: int(typeID),
				Value:  value,
			}
			var queryOutput TableWithBitString

			startTime := time.Now()
			err = db.Get(&queryOutput, `
				SELECT
					id,
					type_id,
					value,
					source
				FROM
				table_with_bit_string
				WHERE
					type_id = $1 AND value = $2
			`, queryInput.TypeID, queryInput.Value)

			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return
			} else if err != sql.ErrNoRows {
				strBit := queryOutput.Source
				replaceStrBit, err := StringBitReplacer(strBit, int(source), true)
				if err != nil {
					log.Println(err)
					return
				} else if strBit == replaceStrBit {
				}
			}

			mutex.Lock()
			count++
			countDuration += time.Since(startTime)
			mutex.Unlock()
		}()
		i++
		countChecked++

		if i == qps {
			go func() {
				mutex.Lock()
				if count != 0 {
					fmt.Printf("Avg processing time : %+v\n", countDuration/time.Duration(count))
				}
				mutex.Unlock()
			}()
			i = 0

			if delay := time.Second - time.Since(startTime); delay > 0 {
				time.Sleep(delay)
			}

			startTime = time.Now()
		}

		if countChecked == totalLine {
			fmt.Println("done")
			time.Sleep(10 * time.Second)
			break
		}
	}

}

func StringBitAppender(stringBit string, totalBit int) (res string) {
	validation := fmt.Sprintf("%%0%ds", totalBit)

	res = fmt.Sprintf(validation, stringBit)

	return
}

func StringBitReplacer(stringBit string, source int, blacklistStatus bool) (res string, err error) {

	if len(stringBit) < source {
		return res, fmt.Errorf("source > stringBit length")
	} else if source < 0 {
		return res, fmt.Errorf("source < 0")
	}

	byteData := []byte(stringBit)

	if blacklistStatus {
		byteData[len(stringBit)-source] = '1'
	} else {
		byteData[len(stringBit)-source] = '0'
	}

	res = string(byteData)

	return
}
