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

func buildTx(record []string) (tx Transaction, err error) {
        var values [3] float64
        var ttype Txtype

        switch record[3] {
            case "withdrawal":
                ttype = Txwithdraw
            case "deposit":
                ttype = Txdeposit
            case "trade":
                ttype = Txtrade
            default:
                ttype = -1
        }

        for i:= 0; i<3; i+=1 {
            values[i], err = strconv.ParseFloat(record[i+6], 64)
            if err != nil {
                return
            }
        }

        tx = Transaction{
            txid: record[0],
            refid: record[1],
            time: record[2],
            txtype: ttype,
            asset: record[5],
            amount: values[0],
            fee: values[1],
            balance: values[2],
        }

        return tx, nil

}

func main() {
    csvFile, err := os.Open("history.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer csvFile.Close()

    r := csv.NewReader(bufio.NewReader(csvFile))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buildTx(record))
	}
}
