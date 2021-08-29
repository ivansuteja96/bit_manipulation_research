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

type TableWithoutBit struct {
	ID     int64  `db:"id"`
	TypeID int    `db:"type_id"`
	Value  string `db:"value"`
	Source int    `db:"source"`
	Status bool   `db:"status"`
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

			queryInput := TableWithoutBit{
				TypeID: int(typeID),
				Source: int(source),
				Value:  value,
				Status: true,
			}
			var queryOutput TableWithoutBit

			startTime := time.Now()
			err = db.Get(&queryOutput, `
				SELECT
					id,
					type_id,
					value,
					source,
					status
				FROM
					table_without_bit
				WHERE
					type_id = $1 AND value = $2 AND source = $3 AND status = $4
			`, queryInput.TypeID, queryInput.Value, queryInput.Source, queryInput.Status)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return
			}

			// fmt.Println(queryOutput, err)

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
