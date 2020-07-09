package kucoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	. "github.com/georgexdz/ccxt/go/base"
)

type Kucoin struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Kucoin, err error) {
	ex = new(Kucoin)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		return
	}

	return
}

func (self *Kucoin) InitDescribe() (err error) {
	err = json.Unmarshal(self.Child.Describe(), &self.DescribeMap)
	if err != nil {
		return
	}

	err = self.DefineRestApi()
	if err != nil {
		return
	}

	publicUrl, err := NestedMapLookup(self.DescribeMap, "urls", "api", "public")
	if err != nil {
		return
	}
	privateUrl, err := NestedMapLookup(self.DescribeMap, "urls", "api", "private")
	if err != nil {
		return
	}
	self.ApiUrls = map[string]string{
		"private": privateUrl.(string),
		"public":  publicUrl.(string),
	}

	self.Options = self.DescribeMap["options"].(map[string]interface{})
	self.Urls = self.DescribeMap["urls"].(map[string]interface{})
	return
}

func (self *Kucoin) Describe() []byte {
	return []byte(`{
	"version": "v2",
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/57369448-3cc3aa80-7196-11e9-883e-5ebeb35e4f57.jpg",
        "referral": "https://www.kucoin.com/?rcode=E5wkqe",
        "api": {
            "public": "https://openapi-v2.kucoin.com",
            "private": "https://openapi-v2.kucoin.com"
        },
        "test": {
            "public": "https://openapi-sandbox.kucoin.com",
            "private": "https://openapi-sandbox.kucoin.com"
        },
        "www": "https://www.kucoin.com",
        "doc": [
            "https://docs.kucoin.com"
        ]
    },
    "api": {
        "public": {
            "get": [
                "timestamp",
                "status",
                "symbols",
                "markets",
                "market/allTickers",
                "market/orderbook/level{level}",
                "market/orderbook/level2",
                "market/orderbook/level2_20",
                "market/orderbook/level2_100",
                "market/orderbook/level3",
                "market/histories",
                "market/candles",
                "market/stats",
                "currencies",
                "currencies/{currency}",
                "prices",
                "mark-price/{symbol}/current",
                "margin/config"
            ],
            "post": [
                "bullet-public"
            ]
        },
        "private": {
            "get": [
                "accounts",
                "accounts/{accountId}",
                "accounts/{accountId}/ledgers",
                "accounts/{accountId}/holds",
                "accounts/transferable",
                "sub/user",
                "sub-accounts",
                "sub-accounts/{subUserId}",
                "deposit-addresses",
                "deposits",
                "hist-deposits",
                "hist-orders",
                "hist-withdrawals",
                "withdrawals",
                "withdrawals/quotas",
                "orders",
                "orders/{orderId}",
                "limit/orders",
                "fills",
                "limit/fills",
                "margin/account",
                "margin/borrow",
                "margin/borrow/outstanding",
                "margin/borrow/borrow/repaid",
                "margin/lend/active",
                "margin/lend/done",
                "margin/lend/trade/unsettled",
                "margin/lend/trade/settled",
                "margin/lend/assets",
                "margin/market",
                "margin/margin/trade/last"
            ],
            "post": [
                "accounts",
                "accounts/inner-transfer",
                "accounts/sub-transfer",
                "deposit-addresses",
                "withdrawals",
                "orders",
                "orders/multi",
                "margin/borrow",
                "margin/repay/all",
                "margin/repay/single",
                "margin/lend",
                "margin/toggle-auto-lend",
                "bullet-private"
            ],
            "delete": [
                "withdrawals/{withdrawalId}",
                "orders",
                "orders/{orderId}",
                "margin/lend/{orderId}"
            ]
        }
    },
	"options": {
		"version": "v1",
		"fetchBalance": {
			"type": "trade"
		}
	}
}
	`)
}

func (self *Kucoin) ApiFuncDecode(function string) (path string, api string, method string, err error) {
	// fmt.Println(self.ApiDecodeInfo)
	if info, ok := self.ApiDecodeInfo[function]; ok {
		return info.Path, info.Api, info.Method, nil
	}
	return "", "", "", errors.New("undefined function!")
}

func (self *Kucoin) Nonce() int64 {
	return self.Milliseconds()
}

func (self *Kucoin) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}, err error) {

	versions := self.SafeValue(self.Options, "versions", map[string]interface{}{})

	apiVersions := self.SafeValue(versions, api, nil)

	methodVersions := self.SafeValue(apiVersions, method, map[string]interface{}{})

	defaultVersion := self.SafeString(methodVersions, path, self.Member(self.Options, "version"))

	version := self.SafeString(params, "version", defaultVersion)

	params = self.Omit(params, "version")

	endpoint := "/api/" + version + "/" + self.ImplodeParams(path, params)

	query := self.Omit(params, self.ExtractParams(path))

	endpart := ""

	headers = self.IfThenElse(self.ToBool(!self.TestNil(headers)), headers, map[string]interface{}{})

	if self.ToBool(self.Length(reflect.ValueOf(query).MapKeys())) {
		if self.ToBool(method != "GET") {
			endpart = self.Json(query)
			self.SetValue(headers, "Content-Type", "application/json")
		} else {
			endpoint += "?" + self.Urlencode(query)
		}
	}

	url := fmt.Sprintf("%v", self.Member(self.Member(self.Urls, "api"), api)) + endpoint

	if self.ToBool(api == "private") {
		self.CheckRequiredCredentials()
		timestamp := fmt.Sprintf("%v", self.Nonce())
		headers = self.Extend(map[string]interface{}{
			"KC-API-KEY":        self.ApiKey,
			"KC-API-TIMESTAMP":  timestamp,
			"KC-API-PASSPHRASE": self.Password,
		}, headers)
		payload := timestamp + method + endpoint + endpart
		signature, err := self.Hmac(self.Encode(payload), self.Encode(self.Secret), "sha256", "base64")
		if err != nil {
			return nil, err
		}
		self.SetValue(headers, "KC-API-SIGN", self.Decode(signature))
	}

	return map[string]interface{}{
		"url":     url,
		"method":  method,
		"body":    endpart,
		"headers": headers,
	}, nil

}

func (self *Kucoin) FetchMarkets(params map[string]interface{}) ([]*Market, error) {
	respJson, err := self.ApiFunc("publicGetSymbols", params, nil, nil)
	if err != nil {
		return nil, err
	}

	var result []*Market

	if respJson["code"] == "200000" {
		if symbolList, ok := respJson["data"].([]interface{}); ok {
			for _, oneSymbol := range symbolList {
				if oneSymbolInfo, ok := oneSymbol.(map[string]interface{}); ok {
					li := strings.Split(oneSymbolInfo["symbol"].(string), "-")
					oneMarket := &Market{
						Id:     oneSymbolInfo["symbol"].(string),
						Symbol: li[0] + "/" + li[1],
					}
					result = append(result, oneMarket)
				}
			}
		}
	}

	return result, nil
}

func (self *Kucoin) FetchOrderBook(symbol string, limit int, params map[string]interface{}) (orderBook *OrderBook, err error) {
	level := self.SafeInteger(params, "level", 2)

	levelLimit := fmt.Sprintf("%v", level)

	if self.ToBool(self.TestNil(levelLimit)) {
		if self.ToBool(self.TestNil(limit)) {
			if limit <= 20 {
				levelLimit += "_" + fmt.Sprintf("%v", 20)
			} else if limit <= 100 {
				levelLimit += "_" + fmt.Sprintf("%v", 100)
			}
		}
	}

	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	marketId, err := self.MarketId(symbol)
	if err != nil {
		return
	}

	request := map[string]interface{}{
		"symbol": marketId,
		"level":  levelLimit,
	}

	response, err := self.ApiFunc("publicGetMarketOrderbookLevelLevel", self.Extend(request, params), nil, nil)
	if err != nil {
		return
	}

	data := self.SafeValue(response, "data", map[string]interface{}{})

	timestamp := self.SafeInteger(data, "time", 0)

	orderbook, err := self.ParseOrderBook(data, timestamp, "bids", "asks", level-2, level-1)
	if err != nil {
		return
	}

	self.SetValue(orderbook, "nonce", self.SafeInteger(data, "sequence", 0))

	return orderbook, nil

}

func (self *Kucoin) FetchBalance(params map[string]interface{}) (balanceResult *Account, err error) {
	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	var typ interface{}

	request := map[string]interface{}{}

	if self.ToBool(self.InMap("type", params)) {
		typ = self.Member(params, "type")
		if self.ToBool(self.TestNil(typ)) {
			self.SetValue(request, "type", typ)
		}
		params = self.Omit(params, "type")
	} else {
		options := self.SafeValue(self.Options, "fetchBalance", map[string]interface{}{})
		typ = self.SafeString(options, "type", "trade")
	}

	response, err := self.ApiFunc("privateGetAccounts", self.Extend(request, params), nil, nil)
	if err != nil {
		return
	}

	data := self.SafeValue(response, "data", []interface{}{})

	result := map[string]interface{}{
		"info": response,
	}

	for i := 0; i < self.Length(data); i++ {
		balance := self.Member(data, i)
		balanceType := self.SafeString(balance, "type", "")
		if self.ToBool(self.TestNil(balanceType)) {
			currencyId := self.SafeString(balance, "currency", "")
			code := self.SafeCurrencyCode(currencyId)
			account := self.Account()
			self.SetValue(account, "total", self.SafeFloat(balance, "balance", 0))
			self.SetValue(account, "free", self.SafeFloat(balance, "available", 0))
			self.SetValue(account, "used", self.SafeFloat(balance, "holds", 0))
			self.SetValue(result, code, account)
		}
	}

	return self.ParseBalance(result), nil

}

func (self *Kucoin) CreateOrder(symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (result *Order, err error) {
	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	marketId, err := self.MarketId(symbol)
	if err != nil {
		return
	}

	clientOid := self.Uuid()

	request := map[string]interface{}{
		"clientOid": clientOid,
		"side":      side,
		"symbol":    marketId,
		"type":      typ,
	}

	if self.ToBool(typ != "market") {
		self.SetValue(request, "price", self.PriceToPrecision(symbol, price))
		self.SetValue(request, "size", self.AmountToPrecision(symbol, amount))
	} else {
		if self.ToBool(self.SafeValue(params, "quoteAmount", nil)) {
			self.SetValue(request, "funds", self.AmountToPrecision(symbol, amount))
		} else {
			self.SetValue(request, "size", self.AmountToPrecision(symbol, amount))
		}
	}

	response, err := self.ApiFunc("privatePostOrders", self.Extend(request, params), nil, nil)
	if err != nil {
		return
	}

	data := self.SafeValue(response, "data", map[string]interface{}{})

	timestamp := self.Milliseconds()

	order := map[string]interface{}{
		"id":            self.SafeString(data, "orderId", ""),
		"symbol":        symbol,
		"type":          typ,
		"side":          side,
		"price":         price,
		"cost":          nil,
		"filled":        nil,
		"remaining":     nil,
		"timestamp":     timestamp,
		"datetime":      self.Iso8601(timestamp),
		"fee":           nil,
		"status":        "open",
		"clientOrderId": clientOid,
		"info":          data,
	}

	if self.ToBool(!self.ToBool(self.SafeValue(params, "quoteAmount", nil))) {
		self.SetValue(order, "amount", amount)
	}

	return (&Order{}).InitFromMap(order)
}

func (self *Kucoin) FetchOrder(id string, symbol string, params map[string]interface{}) (order *Order, err error) {
	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	request := map[string]interface{}{
		"orderId": id,
	}

	var market interface{}

	if self.ToBool(self.TestNil(symbol)) {
		market = self.Market(symbol)
	}

	response, err := self.ApiFunc("privateGetOrdersOrderId", self.Extend(request, params), nil, nil)
	if err != nil {
		return
	}

	responseData := self.Member(response, "data")

	return (&Order{}).InitFromMap(self.ParseOrder(responseData, market))

}

func (self *Kucoin) ParseOrder(order interface{}, market interface{}) map[string]interface{} {
	var symbol interface{}

	marketId := self.SafeString(order, "symbol", "")

	if self.ToBool(!self.TestNil(marketId)) {
		if self.ToBool(self.InMap(marketId, self.MarketsById)) {
			market = self.Member(self.MarketsById, marketId)
			symbol = self.Member(market, "symbol")
		} else {
			baseId, quoteId := self.Unpack2(strings.Split(marketId, "-"))
			base := self.SafeCurrencyCode(baseId)
			quote := self.SafeCurrencyCode(quoteId)
			symbol = base + "/" + quote
		}
		market = self.SafeValue(self.MarketsById, marketId, nil)
	}

	if self.ToBool(self.TestNil(symbol)) {
		if self.ToBool(!self.TestNil(market)) {
			symbol = self.Member(market, "symbol")
		}
	}

	orderId := self.SafeString(order, "id", "")

	typ := self.SafeString(order, "type", "")

	timestamp := self.SafeInteger(order, "createdAt", 0)

	datetime := self.Iso8601(timestamp)

	price := self.SafeFloat(order, "price", 0)

	side := self.SafeString(order, "side", "")

	feeCurrencyId := self.SafeString(order, "feeCurrency", "")

	feeCurrency := self.SafeCurrencyCode(feeCurrencyId)

	feeCost := self.SafeFloat(order, "fee", 0)

	amount := self.SafeFloat(order, "size", 0)

	filled := self.SafeFloat(order, "dealSize", 0)

	cost := self.SafeFloat(order, "dealFunds", 0)

	remaining := amount - filled

	status := self.IfThenElse(self.ToBool(self.Member(order, "isActive")), "open", "closed")

	status = self.IfThenElse(self.ToBool(self.Member(order, "cancelExist")), "canceled", status)

	fee := map[string]interface{}{
		"currency": feeCurrency,
		"cost":     feeCost,
	}

	if self.ToBool(typ == "market") {
		if self.ToBool(price == 0.0) {
			if self.ToBool(!self.TestNil(cost) && !self.TestNil(filled)) {
				if self.ToBool(cost > 0 && filled > 0) {
					price = cost / filled
				}
			}
		}
	}

	clientOrderId := self.SafeString(order, "clientOid", "")

	return map[string]interface{}{
		"id":                 orderId,
		"clientOrderId":      clientOrderId,
		"symbol":             symbol,
		"type":               typ,
		"side":               side,
		"amount":             amount,
		"price":              price,
		"cost":               cost,
		"filled":             filled,
		"remaining":          remaining,
		"timestamp":          timestamp,
		"datetime":           datetime,
		"fee":                fee,
		"status":             status,
		"info":               order,
		"lastTradeTimestamp": nil,
		"average":            nil,
		"trades":             nil,
	}

}

func (self *Kucoin) CancelOrder(id string, symbol string, params map[string]interface{}) (response interface{}, err error) {

	request := map[string]interface{}{
		"orderId": id,
	}

	response, err = self.ApiFunc("privateDeleteOrdersOrderId", self.Extend(request, params), nil, nil)
	if err != nil {
		return
	}

	return response, nil

}

func (self *Kucoin) FetchOrdersByStatus(status string, symbol string, since int64, limit int64, params map[string]interface{}) (result interface{}, err error) {

	_, err = self.LoadMarkets()
	if err != nil {
		return
	}

	request := map[string]interface{}{
		"status": status,
	}

	var market interface{}

	if self.ToBool(self.TestNil(symbol)) {
		market = self.Market(symbol)
		self.SetValue(request, "symbol", self.Member(market, "id"))
	}

	if self.ToBool(self.TestNil(since)) {
		self.SetValue(request, "startAt", since)
	}

	if self.ToBool(self.TestNil(limit)) {
		self.SetValue(request, "pageSize", limit)
	}

	response, err := self.ApiFunc("privateGetOrders", self.Extend(request, params), nil, nil)
	if err != nil {
		return
	}

	responseData := self.SafeValue(response, "data", map[string]interface{}{})

	orders := self.SafeValue(responseData, "items", []interface{}{})

	return self.ParseOrders(orders, market, since, limit), nil
}

func (self *Kucoin) FetchOpenOrders(symbol string, since int64, limit int64, params map[string]interface{}) (orders []*Order, err error) {
	resp, err := self.FetchOrdersByStatus("active", symbol, since, limit, params)
	if err != nil {
		return
	}
	return self.ToOrders(resp)
}
