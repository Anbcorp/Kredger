package main

import (
	"encoding/csv"
    "strconv"
    "io"
    "log"
    "fmt"
    "os"
    "bufio"
)

type Txtype int

const (
	Txdeposit = iota
	Txwithdraw
	Txtrade
)

type Transaction struct {
	txid    string
	refid   string
	time    string
	txtype  Txtype
	asset   string
	amount  float64
	fee     float64
	balance float64
}

func main() {
    csvFile, err := os.Open("history.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer csvFile.Close()

    r := csv.NewReader(bufio.NewReader(csvFile))

upper:
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

        tx := Transaction{ }
        tx.txid = record[0]
        tx.refid = record[1]
        tx.time = record[2]
        switch record[3] {
            case "withdrawal":
                tx.txtype = Txwithdraw
            case "deposit":
                tx.txtype = Txdeposit
            case "trade":
                tx.txtype = Txtrade
            default:
                tx.txtype = -1
        }
        tx.asset = record[5]
        var val [3]float64
        for i:= 0; i<3; i+=1 {
            val[i], err = strconv.ParseFloat(record[i+6], 64)
            if err != nil {
                log.Println(record)
                log.Println(err)
                continue upper
            }
        }
        tx.amount = val[0]
        tx.fee = val[1]
        tx.balance = val[2]
		fmt.Println(tx)
	}
}
