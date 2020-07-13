
package kucoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

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
	ex, err := New(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(ex.Exceptions)
	ex.Verbose = true

	get_test_config(ex)

	orderbook, err := ex.FetchOrderBook("BTC/USDT", 20, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("orderbook:", orderbook)

	markets := ex.LoadMarkets()
	fmt.Println("markets:", markets)

	ex.FetchBalance(nil)

	order, err := ex.CreateOrder("ETH/BTC", "limit", "buy", 0.0001, 0.024, nil)
	if err != nil {
		return
	}

	fmt.Println(ex.FetchOrder(order.Id, "ETH/BTC", nil))

	openOrders, err := ex.FetchOpenOrders("ETH/BTC", 0, 20, nil)
	if err == nil {
		fmt.Println("openorders", openOrders)
	}

	if err == nil {
		res, err := ex.CancelOrder(order.Id, "ETH/BTC", nil)
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