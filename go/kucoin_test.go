package kucoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

//func testFetchMarkets(ex *ccxt.Kucoin) {
	//markets, err := ex.FetchMarkets(nil)
	//if err == nil {
		//fmt.Println(markets)
	//}
//}

func get_test_config(ex *Kucoin) {
	plan, err := ioutil.ReadFile("test_config.json")
	if err != nil {
		return
	}

	var data interface{}
	err = json.Unmarshal(plan, &data)
	if err != nil {
		return
	}
	
	fmt.Println(data)

	if json_config, ok := data.(map[string]interface{}); ok {
        ex.Urls = map[string]interface{}{
        	"api": map[string]interface{}{
        		"public": json_config["url"],
        		"private": json_config["url"],
			},
        }
		ex.ApiUrls["private"] = json_config["url"].(string)
		ex.ApiUrls["public"] = json_config["url"].(string)
		ex.ApiKey = json_config["key"].(string)
		ex.Secret = json_config["secret"].(string)
		ex.Password = json_config["password"].(string)
	}
}

func TestFetchOrderBook(t *testing.T) {
	ex, _ := New(nil)
	fmt.Println(ex.ApiDecodeInfo)
	ex.Verbose = true

	get_test_config(ex)

	markets, err := ex.LoadMarkets()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println("markets:", markets)

	orderbook, err := ex.FetchOrderBook("BTC/USDT", 20, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("orderbook:", orderbook)

	ex.FetchBalance(nil)

	order, err := ex.CreateOrder("ETH/BTC", "limit", "buy", 0.0001, 0.024, nil)
	if err != nil {
		return
	}

	fmt.Println(ex.FetchOrder(order["id"].(string), "ETH/BTC", nil))

	openOrders, err := ex.FetchOpenOrders("ETH/BTC", 0, 20, nil)
	if err == nil {
		fmt.Println("openorders", openOrders)
	}

	if err == nil {
		res, err := ex.CancelOrder(order["id"].(string), "ETH/BTC", nil)
		fmt.Println(res, err)
	}
}

//func main() {
	//ex := &ccxt.Kucoin{}
	//ex.Init()
	//// testFetchMarkets(ex)
	//fmt.Println("enter")
	//testFetchOrderBook(ex)
//}
