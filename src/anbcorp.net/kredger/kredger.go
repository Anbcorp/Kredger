package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

type Trade struct {
    tx1 *Transaction
    tx2 *Transaction
}

func (tr Trade) fees() (fees float64, asset string) {

    if tr.tx1.fee > 0 {
        fees += tr.tx1.fee
        asset = tr.tx1.asset
    }

    if tr.tx2.fee > 0 {
        if asset != "" && tr.tx2.asset != asset {
            log.Printf("%f %s : %f %s", tr.tx1.fee, tr.tx1.asset, tr.tx2.fee, tr.tx2.asset)
            log.Fatal(tr.tx1.refid, ": Fees on both side of trade")
        }
        fees += tr.tx2.fee
        asset = tr.tx2.asset
    }

    return
}

func buildTx(record []string) (tx Transaction, err error) {
	var values [3]float64
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

	for i := 0; i < 3; i += 1 {
		values[i], err = strconv.ParseFloat(record[i+6], 64)
		if err != nil {
			return
		}
	}

	tx = Transaction{
		txid:    record[0],
		refid:   record[1],
		time:    record[2],
		txtype:  ttype,
		asset:   record[5],
		amount:  values[0],
		fee:     values[1],
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

    var tradeList []Trade
    var unmatchedQueue []Transaction
    var backlog []Transaction

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

        tx, err := buildTx(record)
        if err != nil {
            log.Println(record)
            log.Println(err)
            continue
        }

        if tx.txtype == Txtrade {
            if len(backlog)<=0 {
                backlog = append(backlog, tx)
                continue
            }
            backtx := backlog[len(backlog)-1]
            if backtx.refid == tx.refid {
                tradeList = append(tradeList, Trade{&tx, &backtx})
                // Remove the last element
                backlog = backlog[:len(backlog)-1]
                log.Printf("Matched %s with %s", tx.refid, backtx.refid)
            } else {
                unmatchedQueue = append(unmatchedQueue, tx)
            }
        }
	}
    fmt.Println(tradeList)
    fmt.Println(unmatchedQueue)

    for _, tr := range tradeList {
        fee, asset := tr.fees()
        fmt.Printf("Trade (%s) : %s(%f) vs %s(%f) for %s%f\n", tr.tx1.refid, tr.tx1.asset, tr.tx1.amount, tr.tx2.asset, tr.tx2.amount, asset, fee)
    }
}
