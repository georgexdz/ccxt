package main

import (
	"ccxt-master/go"
	"fmt"
)

func testFetchMarkets(ex *ccxt.Kucoin) {
	markets, err := ex.FetchMarkets(nil)
	if err == nil {
		fmt.Println(markets)
	}
}

func testFetchOrderBook(ex *ccxt.Kucoin) {
	orderbook, err := ex.FetchOrderBook("BTC/USDT", 20, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(orderbook.Bids)
	fmt.Println(orderbook.Asks)
	fmt.Println(orderbook.Nonce)
	fmt.Println(orderbook.Timestamp)
}

func main() {
	ex := &ccxt.Kucoin{}
	ex.Init()
	// testFetchMarkets(ex)
	fmt.Println("enter")
	testFetchOrderBook(ex)
}
