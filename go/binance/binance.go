
    
    package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/georgexdz/ccxt/go/base"
	"reflect"
	"strings"
)

type Binance struct {
	Exchange
}

func New(config *ExchangeConfig) (ex *Binance, err error) {
	ex = new(Binance)
	err = ex.Init(config)
	ex.Child = ex

	err = ex.InitDescribe()
	if err != nil {
		return
	}

	return
}

func (self *Binance) InitDescribe() (err error) {
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

    
    func (self *Binance) Describe() []byte {
	return []byte(`{
    "id": "binance",
    "name": "Binance",
    "countries": [
        "JP",
        "MT"
    ],
    "rateLimit": 500,
    "certified": true,
    "pro": true,
    "has": {
        "fetchDepositAddress": true,
        "CORS": false,
        "fetchBidsAsks": true,
        "fetchTickers": true,
        "fetchTime": true,
        "fetchOHLCV": true,
        "fetchMyTrades": true,
        "fetchOrder": true,
        "fetchOrders": true,
        "fetchOpenOrders": true,
        "fetchClosedOrders": "emulated",
        "withdraw": true,
        "fetchFundingFees": true,
        "fetchDeposits": true,
        "fetchWithdrawals": true,
        "fetchTransactions": false,
        "fetchTradingFee": true,
        "fetchTradingFees": true,
        "cancelAllOrders": true
    },
    "timeframes": {
        "1m": "1m",
        "3m": "3m",
        "5m": "5m",
        "15m": "15m",
        "30m": "30m",
        "1h": "1h",
        "2h": "2h",
        "4h": "4h",
        "6h": "6h",
        "8h": "8h",
        "12h": "12h",
        "1d": "1d",
        "3d": "3d",
        "1w": "1w",
        "1M": "1M"
    },
    "urls": {
        "logo": "https://user-images.githubusercontent.com/1294454/29604020-d5483cdc-87ee-11e7-94c7-d1a8d9169293.jpg",
        "test": {
            "fapiPublic": "https://testnet.binancefuture.com/fapi/v1",
            "fapiPrivate": "https://testnet.binancefuture.com/fapi/v1",
            "public": "https://testnet.binance.vision/api/v3",
            "private": "https://testnet.binance.vision/api/v3",
            "v3": "https://testnet.binance.vision/api/v3",
            "v1": "https://testnet.binance.vision/api/v1"
        },
        "api": {
            "wapi": "https://api.binance.com/wapi/v3",
            "sapi": "https://api.binance.com/sapi/v1",
            "fapiPublic": "https://fapi.binance.com/fapi/v1",
            "fapiPrivate": "https://fapi.binance.com/fapi/v1",
            "public": "https://api.binance.com/api/v3",
            "private": "https://api.binance.com/api/v3",
            "v3": "https://api.binance.com/api/v3",
            "v1": "https://api.binance.com/api/v1"
        },
        "www": "https://www.binance.com",
        "referral": "https://www.binance.com/?ref=10205187",
        "doc": [
            "https://binance-docs.github.io/apidocs/spot/en"
        ],
        "api_management": "https://www.binance.com/en/usercenter/settings/api-management",
        "fees": "https://www.binance.com/en/fee/schedule"
    },
    "api": {
        "sapi": {
            "get": [
                "accountSnapshot",
                "margin/asset",
                "margin/pair",
                "margin/allAssets",
                "margin/allPairs",
                "margin/priceIndex",
                "asset/assetDividend",
                "margin/loan",
                "margin/repay",
                "margin/account",
                "margin/transfer",
                "margin/interestHistory",
                "margin/forceLiquidationRec",
                "margin/order",
                "margin/openOrders",
                "margin/allOrders",
                "margin/myTrades",
                "margin/maxBorrowable",
                "margin/maxTransferable",
                "futures/transfer",
                "capital/config/getall",
                "capital/deposit/address",
                "capital/deposit/hisrec",
                "capital/deposit/subAddress",
                "capital/deposit/subHisrec",
                "capital/withdraw/history",
                "sub-account/futures/account",
                "sub-account/futures/accountSummary",
                "sub-account/futures/positionRisk",
                "sub-account/margin/account",
                "sub-account/margin/accountSummary",
                "sub-account/status",
                "sub-account/transfer/subUserHistory",
                "lending/daily/product/list",
                "lending/daily/userLeftQuota",
                "lending/daily/userRedemptionQuota",
                "lending/daily/token/position",
                "lending/union/account",
                "lending/union/purchaseRecord",
                "lending/union/redemptionRecord",
                "lending/union/interestHistory",
                "lending/project/list",
                "lending/project/position/list",
                "mining/pub/algoList",
                "mining/pub/coinList",
                "mining/worker/detail",
                "mining/worker/list",
                "mining/payment/list",
                "mining/statistics/user/status",
                "mining/statistics/user/list"
            ],
            "post": [
                "asset/dust",
                "account/disableFastWithdrawSwitch",
                "account/enableFastWithdrawSwitch",
                "capital/withdraw/apply",
                "margin/transfer",
                "margin/loan",
                "margin/repay",
                "margin/order",
                "sub-account/margin/enable",
                "sub-account/margin/enable",
                "sub-account/futures/enable",
                "userDataStream",
                "futures/transfer",
                "lending/customizedFixed/purchase",
                "lending/daily/purchase",
                "lending/daily/redeem"
            ],
            "put": [
                "userDataStream"
            ],
            "delete": [
                "margin/order",
                "userDataStream"
            ]
        },
        "wapi": {
            "post": [
                "withdraw",
                "sub-account/transfer"
            ],
            "get": [
                "depositHistory",
                "withdrawHistory",
                "depositAddress",
                "accountStatus",
                "systemStatus",
                "apiTradingStatus",
                "userAssetDribbletLog",
                "tradeFee",
                "assetDetail",
                "sub-account/list",
                "sub-account/transfer/history",
                "sub-account/assets"
            ]
        },
        "fapiPublic": {
            "get": [
                "ping",
                "time",
                "exchangeInfo",
                "depth",
                "trades",
                "historicalTrades",
                "aggTrades",
                "klines",
                "fundingRate",
                "premiumIndex",
                "ticker/24hr",
                "ticker/price",
                "ticker/bookTicker",
                "allForceOrders",
                "openInterest",
                "leverageBracket"
            ]
        },
        "fapiPrivate": {
            "get": [
                "allForceOrders",
                "allOrders",
                "openOrder",
                "openOrders",
                "order",
                "account",
                "balance",
                "positionMargin/history",
                "positionRisk",
                "positionSide/dual",
                "userTrades",
                "income"
            ],
            "post": [
                "batchOrders",
                "positionSide/dual",
                "positionMargin",
                "marginType",
                "order",
                "leverage",
                "listenKey",
                "countdownCancelAll"
            ],
            "put": [
                "listenKey"
            ],
            "delete": [
                "batchOrders",
                "order",
                "allOpenOrders",
                "listenKey"
            ]
        },
        "v3": {
            "get": [
                "ticker/price",
                "ticker/bookTicker"
            ]
        },
        "public": {
            "get": [
                "ping",
                "time",
                "depth",
                "trades",
                "aggTrades",
                "historicalTrades",
                "klines",
                "ticker/24hr",
                "ticker/price",
                "ticker/bookTicker",
                "exchangeInfo"
            ],
            "put": [
                "userDataStream"
            ],
            "post": [
                "userDataStream"
            ],
            "delete": [
                "userDataStream"
            ]
        },
        "private": {
            "get": [
                "allOrderList",
                "openOrderList",
                "orderList",
                "order",
                "openOrders",
                "allOrders",
                "account",
                "myTrades"
            ],
            "post": [
                "order/oco",
                "order",
                "order/test"
            ],
            "delete": [
                "openOrders",
                "orderList",
                "order"
            ]
        }
    },
    "fees": {
        "trading": {
            "tierBased": false,
            "percentage": true,
            "taker": 0.001,
            "maker": 0.001
        }
    },
    "commonCurrencies": {
        "BCC": "BCC",
        "YOYO": "YOYOW"
    },
    "options": {
        "fetchTradesMethod": "publicGetAggTrades",
        "fetchTickersMethod": "publicGetTicker24hr",
        "defaultTimeInForce": "GTC",
        "defaultType": "spot",
        "hasAlreadyAuthenticatedSuccessfully": false,
        "warnOnFetchOpenOrdersWithoutSymbol": true,
        "recvWindow": 5000,
        "timeDifference": 0,
        "adjustForTimeDifference": false,
        "parseOrderToPrecision": false,
        "newOrderRespType": {
            "market": "FULL",
            "limit": "RESULT"
        },
        "quoteOrderQty": true
    },
    "exceptions": {
        "API key does not exist": "AuthenticationError",
        "Order would trigger immediately.": "InvalidOrder",
        "Account has insufficient balance for requested action.": "InsufficientFunds",
        "Rest API trading is not enabled.": "ExchangeNotAvailable",
        "You don't have permission.": "PermissionDenied",
        "Market is closed.": "ExchangeNotAvailable",
        "-1000": "ExchangeNotAvailable",
        "-1003": "RateLimitExceeded",
        "-1013": "InvalidOrder",
        "-1021": "InvalidNonce",
        "-1022": "AuthenticationError",
        "-1100": "InvalidOrder",
        "-1104": "ExchangeError",
        "-1128": "ExchangeError",
        "-2010": "ExchangeError",
        "-2011": "OrderNotFound",
        "-2013": "OrderNotFound",
        "-2014": "AuthenticationError",
        "-2015": "AuthenticationError",
        "-3008": "InsufficientFunds",
        "-3010": "ExchangeError"
    }
}`)
	}
	
    func (self *Binance) FetchMarkets (params map[string]interface{}) (ret interface{}, err error) {
    
defaultType := self.SafeString2(self.Options,"fetchMarkets","defaultType","spot")

typ := self.SafeString(params,"type",defaultType)

query := self.Omit(params,"type")

if self.ToBool(typ != "spot" && typ != "future" && typ != "margin") {
err = errors.New(self.Id + " does not support "" + typ + "" type, set exchange.options["defaultType"] to "spot", "margin" or "future"");
return
}

method := self.IfThenElse(self.ToBool(typ == "future"), "fapiPublicGetExchangeInfo", "publicGetExchangeInfo")

response := self.Method(query)

if self.ToBool(self.Member(self.Options, "adjustForTimeDifference")) {
self.LoadTimeDifference()
}

markets := self.SafeValue(response,"symbols", nil)

result := []interface{}{}

for i := 0; i < self.Length(markets); i++ {
market := self.Member(markets, i)
future := self.InMap("maintMarginPercent", market)
spot := !self.ToBool(future)
marketType := self.IfThenElse(self.ToBool(spot), "spot", "future")
id := self.SafeString(market,"symbol", "")
lowercaseId := self.SafeStringLower(market,"symbol")
baseId := self.SafeString(market,"baseAsset", "")
quoteId := self.SafeString(market,"quoteAsset", "")
base := self.SafeCurrencyCode(baseId)
quote := self.SafeCurrencyCode(quoteId)
symbol := base + "/" + quote
filters := self.SafeValue(market,"filters",[]interface{}{})
filtersByType := self.IndexBy(filters,"filterType")
precision := map[string]interface{}{
"base": self.SafeInteger(market,"baseAssetPrecision", 0),
"quote": self.SafeInteger(market,"quotePrecision", 0),
"amount": self.SafeInteger(market,"baseAssetPrecision", 0),
"price": self.SafeInteger(market,"quotePrecision", 0),
}
status := self.SafeString(market,"status", "")
active := status == "TRADING"
margin := self.SafeValue(market,"isMarginTradingAllowed",future)
entry := map[string]interface{}{
"id": id,
"lowercaseId": lowercaseId,
"symbol": symbol,
"base": base,
"quote": quote,
"baseId": baseId,
"quoteId": quoteId,
"info": market,
"type": marketType,
"spot": spot,
"margin": margin,
"future": future,
"active": active,
"precision": precision,
"limits": map[string]interface{}{
"amount": map[string]interface{}{
"min": self.Member(Math, pow),
"max": nil,
},
"price": map[string]interface{}{
"min": nil,
"max": nil,
},
"cost": map[string]interface{}{
"min": nil,
"max": nil,
},
},
}
if self.ToBool(self.InMap("PRICE_FILTER", filtersByType)) {
filter := self.SafeValue(filtersByType,"PRICE_FILTER",map[string]interface{}{
})
self.SetValue(self.Member(entry, "limits"), "price", map[string]interface{}{
"min": self.SafeFloat(filter,"minPrice", 0),
"max": nil,
})
maxPrice := self.SafeFloat(filter,"maxPrice", 0)
if self.ToBool(!self.TestNil(maxPrice) && maxPrice > 0) {
self.SetValue(self.Member(self.Member(entry, "limits"), "price"), "max", maxPrice)
}
self.SetValue(self.Member(entry, "precision"), "price", self.PrecisionFromString(self.Member(filter, "tickSize")))
}
if self.ToBool(self.InMap("LOT_SIZE", filtersByType)) {
filter := self.SafeValue(filtersByType,"LOT_SIZE",map[string]interface{}{
})
stepSize := self.SafeString(filter,"stepSize", "")
self.SetValue(self.Member(entry, "precision"), "amount", self.PrecisionFromString(stepSize))
self.SetValue(self.Member(entry, "limits"), "amount", map[string]interface{}{
"min": self.SafeFloat(filter,"minQty", 0),
"max": self.SafeFloat(filter,"maxQty", 0),
})
}
if self.ToBool(self.InMap("MARKET_LOT_SIZE", filtersByType)) {
filter := self.SafeValue(filtersByType,"MARKET_LOT_SIZE",map[string]interface{}{
})
self.SetValue(self.Member(entry, "limits"), "market", map[string]interface{}{
"min": self.SafeFloat(filter,"minQty", 0),
"max": self.SafeFloat(filter,"maxQty", 0),
})
}
if self.ToBool(self.InMap("MIN_NOTIONAL", filtersByType)) {
filter := self.SafeValue(filtersByType,"MIN_NOTIONAL",map[string]interface{}{
})
self.SetValue(self.Member(self.Member(entry, "limits"), "cost"), "min", self.SafeFloat(filter,"minNotional", 0))
}
result = append(result, entry)
}

return result, nil

}
    
func (self *Binance) FetchBalance (params map[string]interface{}) (balanceResult *Account, err error) {
    
_, err = self.LoadMarkets()
 if err != nil {
 return 
}

defaultType := self.SafeString2(self.Options,"fetchBalance","defaultType","spot")

typ := self.SafeString(params,"type",defaultType)

method := "privateGetAccount"

if self.ToBool(typ == "future") {
method = "fapiPrivateGetAccount"
} else if self.ToBool(typ == "margin") {
method = "sapiGetMarginAccount"
}

query := self.Omit(params,"type")

response := self.Method(query)

result := map[string]interface{}{
"info": response,
}

if self.ToBool(typ == "spot" || typ == "margin") {
balances := self.SafeValue2(response,"balances","userAssets",[]interface{}{})
for i := 0; i < self.Length(balances); i++ {
balance := self.Member(balances, i)
currencyId := self.SafeString(balance,"asset", "")
code := self.SafeCurrencyCode(currencyId)
account := self.Account()
self.SetValue(account, "free", self.SafeFloat(balance,"free", 0))
self.SetValue(account, "used", self.SafeFloat(balance,"locked", 0))
self.SetValue(result, code, account)
}
} else {
balances := self.SafeValue(response,"assets",[]interface{}{})
for i := 0; i < self.Length(balances); i++ {
balance := self.Member(balances, i)
currencyId := self.SafeString(balance,"asset", "")
code := self.SafeCurrencyCode(currencyId)
account := self.Account()
self.SetValue(account, "used", self.SafeFloat(balance,"initialMargin", 0))
self.SetValue(account, "total", self.SafeFloat(balance,"marginBalance", 0))
self.SetValue(result, code, account)
}
}

return self.ParseBalance(result), nil

}
    
func (self *Binance) FetchOrderBook (symbol string, limit int64, params map[string]interface{}) (orderBook interface{}, err error) {
    
_, err = self.LoadMarkets()
 if err != nil {
 return 
}

market := self.Market(symbol)

request := map[string]interface{}{
"symbol": self.Member(market, "id"),
}

if self.ToBool(!self.TestNil(limit)) {
self.SetValue(request, "limit", limit)
}

method := self.IfThenElse(self.ToBool(self.Member(market, "spot")), "publicGetDepth", "fapiPublicGetDepth")

response := self.Method(self.Extend(request,params))

orderbook, err := self.ParseOrderBook(response)
if err != nil {
return nil, err
}

self.SetValue(orderbook, "nonce", self.SafeInteger(response,"lastUpdateId", 0))

return orderbook, nil

}
    
func (self *Binance) ParseOrder (order interface{}, market interface{}) (result interface{}) {
    
status := self.ParseOrderStatus(self.SafeString(order,"status", ""))

var symbol interface{}

marketId := self.SafeString(order,"symbol", "")

if self.ToBool(self.InMap(marketId, self.MarketsById)) {
market = self.Member(self.MarketsById, marketId)
}

if self.ToBool(!self.TestNil(market)) {
symbol = self.Member(market, "symbol")
}

var timestamp interface{}

if self.ToBool(self.InMap("time", order)) {
timestamp = self.SafeInteger(order,"time", 0)
} else if self.ToBool(self.InMap("transactTime", order)) {
timestamp = self.SafeInteger(order,"transactTime", 0)
}

price := self.SafeFloat(order,"price", 0)

amount := self.SafeFloat(order,"origQty", 0)

filled := self.SafeFloat(order,"executedQty", 0)

var remaining interface{}

cost := self.SafeFloat2(order,"cummulativeQuoteQty","cumQuote")

if self.ToBool(!self.TestNil(filled)) {
if self.ToBool(!self.TestNil(amount)) {
remaining = amount - filled
if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
remaining = parseFloat(self.AmountToPrecision(symbol,remaining))
}
remaining = self.Member(Math, max)
}
if self.ToBool(!self.TestNil(price)) {
if self.ToBool(self.TestNil(cost)) {
cost = price * filled
}
}
}

id := self.SafeString(order,"orderId", "")

typ := self.SafeStringLower(order,"type")

if self.ToBool(typ == "market") {
if self.ToBool(price == 0.0) {
if self.ToBool(!self.TestNil(cost) && !self.TestNil(filled)) {
if self.ToBool(cost > 0 && filled > 0) {
price = cost / filled
if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
price = parseFloat(self.PriceToPrecision(symbol,price))
}
}
}
}
} else if self.ToBool(typ == "limit_maker") {
typ = "limit"
}

side := self.SafeStringLower(order,"side")

var fee interface{}

var trades interface{}

fills := self.SafeValue(order,"fills", nil)

if self.ToBool(!self.TestNil(fills)) {
trades = self.ParseTrades(fills,market)
numTrades := self.Length(trades)
if self.ToBool(numTrades > 0) {
cost = self.Member(self.Member(trades, 0), "cost")
fee = map[string]interface{}{
"cost": self.Member(self.Member(self.Member(trades, 0), "fee"), "cost"),
"currency": self.Member(self.Member(self.Member(trades, 0), "fee"), "currency"),
}
for i := 1; i < self.Length(trades); i++ {
cost = self.Sum(cost,self.Member(self.Member(trades, i), "cost"))
self.SetValue(fee, "cost", self.Sum(self.Member(fee, "cost"),self.Member(self.Member(self.Member(trades, i), "fee"), "cost")))
}
}
}

var average interface{}

if self.ToBool(!self.TestNil(cost)) {
if self.ToBool(filled) {
average = cost / filled
if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
average = parseFloat(self.PriceToPrecision(symbol,average))
}
}
if self.ToBool(self.Member(self.Options, "parseOrderToPrecision")) {
cost = parseFloat(self.CostToPrecision(symbol,cost))
}
}

clientOrderId := self.SafeString(order,"clientOrderId", "")

return map[string]interface{}{
"info": order,
"id": id,
"clientOrderId": clientOrderId,
"timestamp": timestamp,
"datetime": self.Iso8601(timestamp),
"lastTradeTimestamp": nil,
"symbol": symbol,
"type": typ,
"side": side,
"price": price,
"amount": amount,
"cost": cost,
"average": average,
"filled": filled,
"remaining": remaining,
"status": status,
"fee": fee,
"trades": trades,
}

}
    
func (self *Binance) CreateOrder (symbol string, typ string, side string, amount float64, price float64, params map[string]interface{}) (order map[string]interface{}, err error) {
    
_, err = self.LoadMarkets()
 if err != nil {
 return 
}

market := self.Market(symbol)

defaultType := self.SafeString2(self.Options,"createOrder","defaultType",self.Member(market, "type"))

orderType := self.SafeString(params,"type",defaultType)

clientOrderId := self.SafeString2(params,"newClientOrderId","clientOrderId")

params = self.Omit(params,[]interface{}{"type","newClientOrderId","clientOrderId"})

method := "privatePostOrder"

if self.ToBool(orderType == "future") {
method = "fapiPrivatePostOrder"
} else if self.ToBool(orderType == "margin") {
method = "sapiPostMarginOrder"
}

if self.ToBool(self.Member(market, "spot")) {
test := self.SafeValue(params,"test",false)
if self.ToBool(test) {
method += "Test"
}
params = self.Omit(params,"test")
}

uppercaseType := self.Member(typ, toUpperCase)

validOrderTypes := self.SafeValue(self.Member(market, "info"),"orderTypes", nil)

if self.ToBool(!self.ToBool(self.InArray(uppercaseType,validOrderTypes))) {
err = errors.New(self.Id + " " + typ + " is not a valid order type in " + self.Member(market, "type") + " market " + symbol);
return
}

request := map[string]interface{}{
"symbol": self.Member(market, "id"),
"type": uppercaseType,
"side": self.Member(side, toUpperCase),
}

if self.ToBool(!self.TestNil(clientOrderId)) {
self.SetValue(request, "newClientOrderId", clientOrderId)
}

if self.ToBool(self.Member(market, "spot")) {
self.SetValue(request, "newOrderRespType", self.SafeValue(self.Member(self.Options, "newOrderRespType"),typ,"RESULT"))
}

timeInForceIsRequired := false

priceIsRequired := false

stopPriceIsRequired := false

quantityIsRequired := false

if self.ToBool(uppercaseType == "MARKET") {
quoteOrderQty := self.SafeValue(self.Options,"quoteOrderQty",false)
if self.ToBool(quoteOrderQty) {
quoteOrderQty := self.SafeFloat(params,"quoteOrderQty", 0)
precision := self.Member(self.Member(market, "precision"), "price")
if self.ToBool(!self.TestNil(quoteOrderQty)) {
self.SetValue(request, "quoteOrderQty", self.DecimalToPrecision(quoteOrderQty,TRUNCATE,precision,self.PrecisionMode))
params = self.Omit(params,"quoteOrderQty")
} else if self.ToBool(!self.TestNil(price)) {
self.SetValue(request, "quoteOrderQty", self.DecimalToPrecision(amount * price,TRUNCATE,precision,self.PrecisionMode))
} else {
quantityIsRequired = true
}
} else {
quantityIsRequired = true
}
} else if self.ToBool(uppercaseType == "LIMIT") {
priceIsRequired = true
timeInForceIsRequired = true
quantityIsRequired = true
} else if self.ToBool(uppercaseType == "STOP_LOSS" || uppercaseType == "TAKE_PROFIT") {
stopPriceIsRequired = true
quantityIsRequired = true
if self.ToBool(self.Member(market, "future")) {
priceIsRequired = true
}
} else if self.ToBool(uppercaseType == "STOP_LOSS_LIMIT" || uppercaseType == "TAKE_PROFIT_LIMIT") {
quantityIsRequired = true
stopPriceIsRequired = true
priceIsRequired = true
timeInForceIsRequired = true
} else if self.ToBool(uppercaseType == "LIMIT_MAKER") {
priceIsRequired = true
quantityIsRequired = true
} else if self.ToBool(uppercaseType == "STOP") {
quantityIsRequired = true
stopPriceIsRequired = true
priceIsRequired = true
} else if self.ToBool(uppercaseType == "STOP_MARKET" || uppercaseType == "TAKE_PROFIT_MARKET") {
closePosition := self.SafeValue(params,"closePosition", nil)
if self.ToBool(self.TestNil(closePosition)) {
quantityIsRequired = true
}
stopPriceIsRequired = true
}

if self.ToBool(quantityIsRequired) {
self.SetValue(request, "quantity", self.AmountToPrecision(symbol,amount))
}

if self.ToBool(priceIsRequired) {
if self.ToBool(self.TestNil(price)) {
err = errors.New(self.Id + " createOrder method requires a price argument for a " + typ + " order");
return
}
self.SetValue(request, "price", self.PriceToPrecision(symbol,price))
}

if self.ToBool(timeInForceIsRequired) {
self.SetValue(request, "timeInForce", self.Member(self.Options, "defaultTimeInForce"))
}

if self.ToBool(stopPriceIsRequired) {
stopPrice := self.SafeFloat(params,"stopPrice", 0)
if self.ToBool(self.TestNil(stopPrice)) {
err = errors.New(self.Id + " createOrder method requires a stopPrice extra param for a " + typ + " order");
return
} else {
params = self.Omit(params,"stopPrice")
self.SetValue(request, "stopPrice", self.PriceToPrecision(symbol,stopPrice))
}
}

response := self.Method(self.Extend(request,params))

return self.ParseOrder(response,market), nil

}
    
func (self *Binance) FetchOrder (id string, symbol string, params map[string]interface{}) (order interface{}, err error) {
    
if self.ToBool(self.TestNil(symbol)) {
err = errors.New(self.Id + " fetchOrder requires a symbol argument");
return
}

_, err = self.LoadMarkets()
 if err != nil {
 return 
}

market := self.Market(symbol)

defaultType := self.SafeString2(self.Options,"fetchOrder","defaultType",self.Member(market, "type"))

typ := self.SafeString(params,"type",defaultType)

method := "privateGetOrder"

if self.ToBool(typ == "future") {
method = "fapiPrivateGetOrder"
} else if self.ToBool(typ == "margin") {
method = "sapiGetMarginOrder"
}

request := map[string]interface{}{
"symbol": self.Member(market, "id"),
}

clientOrderId := self.SafeValue2(params,"origClientOrderId","clientOrderId")

if self.ToBool(!self.TestNil(clientOrderId)) {
self.SetValue(request, "origClientOrderId", clientOrderId)
} else {
self.SetValue(request, "orderId", parseInt(id))
}

query := self.Omit(params,[]interface{}{"type","clientOrderId","origClientOrderId"})

response := self.Method(self.Extend(request,query))

return self.ParseOrder(response,market), nil

}
    
func (self *Binance) FetchOpenOrders (symbol string, since int64, limit int64, params map[string]interface{}) (orders interface{}, err error) {
    
_, err = self.LoadMarkets()
 if err != nil {
 return 
}

var market interface{}

var query interface{}

var typ interface{}

request := map[string]interface{}{
}

if self.ToBool(!self.TestNil(symbol)) {
market = self.Market(symbol)
self.SetValue(request, "symbol", self.Member(market, "id"))
defaultType := self.SafeString2(self.Options,"fetchOpenOrders","defaultType",self.Member(market, "type"))
typ = self.SafeString(params,"type",defaultType)
query = self.Omit(params,"type")
} else if self.ToBool(self.Member(self.Options, "warnOnFetchOpenOrdersWithoutSymbol")) {
symbols := self.Symbols
numSymbols := self.Length(symbols)
fetchOpenOrdersRateLimit := parseInt(numSymbols / 2)
err = errors.New(self.Id + " fetchOpenOrders WARNING: fetching open orders without specifying a symbol is rate-limited to one call per " + fmt.Sprintf("%v", fetchOpenOrdersRateLimit) + " seconds. Do not call this method frequently to avoid ban. Set " + self.Id + ".options["warnOnFetchOpenOrdersWithoutSymbol"] = false to suppress this warning message.");
return
} else {
defaultType := self.SafeString2(self.Options,"fetchOpenOrders","defaultType","spot")
typ = self.SafeString(params,"type",defaultType)
query = self.Omit(params,"type")
}

method := "privateGetOpenOrders"

if self.ToBool(typ == "future") {
method = "fapiPrivateGetOpenOrders"
} else if self.ToBool(typ == "margin") {
method = "sapiGetMarginOpenOrders"
}

response := self.Method(self.Extend(request,query))

return self.ParseOrders(response,market,since,limit), nil

}
    
func (self *Binance) CancelOrder (id string, symbol string, params map[string]interface{}) (response interface{}, err error) {
    
if self.ToBool(self.TestNil(symbol)) {
err = errors.New(self.Id + " cancelOrder requires a symbol argument");
return
}

_, err = self.LoadMarkets()
 if err != nil {
 return 
}

market := self.Market(symbol)

defaultType := self.SafeString2(self.Options,"fetchOpenOrders","defaultType",self.Member(market, "type"))

typ := self.SafeString(params,"type",defaultType)

origClientOrderId := self.SafeValue2(params,"origClientOrderId","clientOrderId")

request := map[string]interface{}{
"symbol": self.Member(market, "id"),
}

if self.ToBool(self.TestNil(origClientOrderId)) {
self.SetValue(request, "orderId", parseInt(id))
} else {
self.SetValue(request, "origClientOrderId", origClientOrderId)
}

method := "privateDeleteOrder"

if self.ToBool(typ == "future") {
method = "fapiPrivateDeleteOrder"
} else if self.ToBool(typ == "margin") {
method = "sapiDeleteMarginOrder"
}

query := self.Omit(params,[]interface{}{"type","origClientOrderId","clientOrderId"})

response = self.Method(self.Extend(request,query))

return self.ParseOrder(response), nil

}
    
func (self *Binance) Sign (path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (ret interface{}, err error) {
    
if self.ToBool(!self.ToBool(self.InMap(api, self.Member(self.Urls, "api")))) {
err = errors.New(self.Id + " does not have a testnet/sandbox URL for " + api + " endpoints");
return
}

url := self.Member(self.Member(self.Urls, "api"), api)

url += "/" + path

if self.ToBool(api == "wapi") {
url += ".html"
}

userDataStream := path == "userDataStream" || path == "listenKey"

if self.ToBool(path == "historicalTrades") {
if self.ToBool(self.ApiKey) {
headers = map[string]interface{}{
"X-MBX-APIKEY": self.ApiKey,
}
} else {
err = errors.New(self.Id + " historicalTrades endpoint requires `apiKey` credential");
return
}
} else if self.ToBool(userDataStream) {
if self.ToBool(self.ApiKey) {
body = self.Urlencode(params)
headers = map[string]interface{}{
"X-MBX-APIKEY": self.ApiKey,
"Content-Type": "application/x-www-form-urlencoded",
}
} else {
err = errors.New(self.Id + " userDataStream endpoint requires `apiKey` credential");
return
}
}

if self.ToBool(api == "private" || api == "sapi" || api == "wapi" && path != "systemStatus" || api == "fapiPrivate") {
self.CheckRequiredCredentials()
var query interface{}
if self.ToBool(api == "sapi" && path == "asset/dust") {
query = self.UrlencodeWithArrayRepeat(self.Extend(map[string]interface{}{
"timestamp": self.Nonce(),
"recvWindow": self.Member(self.Options, "recvWindow"),
},params))
} else if self.ToBool(path == "batchOrders") {
query = self.Rawencode(self.Extend(map[string]interface{}{
"timestamp": self.Nonce(),
"recvWindow": self.Member(self.Options, "recvWindow"),
},params))
} else {
query = self.Urlencode(self.Extend(map[string]interface{}{
"timestamp": self.Nonce(),
"recvWindow": self.Member(self.Options, "recvWindow"),
},params))
}
signature, err := self.Hmac(self.Encode(query),self.Encode(self.Secret))
if err != nil {
return nil, err
}
query += "&" + "signature=" + signature
headers = map[string]interface{}{
"X-MBX-APIKEY": self.ApiKey,
}
if self.ToBool(method == "GET" || method == "DELETE" || api == "wapi") {
url += "?" + query
} else {
body = query
self.SetValue(headers, "Content-Type", "application/x-www-form-urlencoded")
}
} else {
if self.ToBool(!self.ToBool(userDataStream)) {
if self.ToBool(self.Length(reflect.ValueOf(query).MapKeys())) {
url += "?" + self.Urlencode(params)
}
}
}

return map[string]interface{}{
"url": url,
"method": method,
"body": body,
"headers": headers,
}, nil

}
    

    