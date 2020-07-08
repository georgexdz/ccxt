package base

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/satori/go.uuid"
)

type JSONTime int64

type SignInfo struct {
	Url     string
	Method  string
	Body    interface{} // []byte or string
	Headers map[string]interface{}
}

// Market struct
type Market struct {
	Id             string      `json:"id"`     // exchange specific
	Symbol         string      `json:"symbol"` // ccxt unified
	Base           string      `json:"base"`
	BaseNumericId  string      `json:"baseNumericId"`
	Quote          string      `json:"quote"`
	QuoteNumericId string      `json:"quoteNumericId"`
	BaseId         string      `json:"baseId"`     // from bitmex
	QuoteId        string      `json:"quoteId"`    // from bitmex
	Active         bool        `json:"active"`     // from bitmex
	Taker          float64     `json:"taker"`      // from bitmex
	Maker          float64     `json:"maker"`      // from bitmex
	Type           string      `json:"type"`       // from bitmex
	Spot           bool        `json:"spot"`       // from bitmex
	Swap           bool        `json:"swap"`       // from bitmex
	Future         bool        `json:"future"`     // from bitmex
	Prediction     bool        `json:"prediction"` // from bitmex
	Precision      Precision   `json:"precision"`
	Limits         Limits      `json:"limits"`
	Lot            float64     `json:"lot"`
	Info           interface{} `json:"info"`
}

// Precision struct
type Precision struct {
	Amount int `json:"amount"`
	Base   int `json:"base"`
	Price  int `json:"price"`
	Cost   int `json:"cost"`
}

// Limits struct
type Limits struct {
	Amount MinMax `json:"amount"`
	Price  MinMax `json:"price"`
	Cost   MinMax `json:"cost"`
}

// MinMax struct
type MinMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// ExchangeConfig for main configuration
// Timeout takes json value in milliseconds
type ExchangeConfig struct {
	ApiKey          string        `json:"apiKey"`
	Password        string
	Secret          string        `json:"secret"`
	Timeout         time.Duration `json:"timeout"`
	EnableRateLimit bool          `json:"enableRateLimit"`
	Test            bool          `json:"test"`
}

// ExchangeInfo for the exchange
type ExchangeInfo struct {
	Id                               string            `json:"id"`
	Name                             string            `json:"name"`
	Countries                        StringSlice       `json:"countries"`
	Version                          string            `json:"version"`
	EnableRateLimit                  bool              `json:"enableRateLimit"`
	RateLimit                        int               `json:"rateLimit"`
	Has                              HasDescription    `json:"has"`
	Urls                             map[string]interface{}
	Api                              Apis              `json:"api"`
	Timeframes                       map[string]string `json:"timeframes"`
	Fees                             Fees              `json:"fees"`
	UserAgents                       map[string]string `json:"userAgents"`
	Header                           http.Header       `json:"header"`
	Proxy                            string            `json:"proxy"`
	Origin                           string            `json:"origin"`
	Verbose                          bool              `json:"verbose"`
	Limits                           Limits            `json:"limits"`
	Precision                        Precision         `json:"precision"`
	Exceptions                       Exceptions        `json:"exceptions"`
	DontGetUsedBalanceFromStaleCache bool              `json:"dontGetUsedBalanceFromStaleCache"`
}

// Apis public and private methods
type Apis struct {
	Public  ApiMethods `json:"public"`
	Private ApiMethods `json:"private"`
}

// ApiMethods get/post/put/delete urls
type ApiMethods struct {
	Get    map[string]string `json:"get"`
	Post   map[string]string `json:"post"`
	Put    map[string]string `json:"put"`
	Delete map[string]string `json:"delete"`
}

// Urls for exchange
type Urls struct {
	WWW  string      `json:"www"`
	Test string      `json:"test"`
	Api  string      `json:"api"`
	Logo StringSlice `json:"logo"`
	Doc  StringSlice `json:"doc"`
	Fees StringSlice `json:"fees"`
}

// Exceptions takes exact/broad errors and classifies
// them to general errors
type Exceptions struct {
	Exact Exception `json:"exact"`
	Broad Exception `json:"broad"`
}

// Exception takes the string and applies the error method
type Exception map[string]error

// UnmarshalJSON accepts strings and links to the appropaite error method:
func (e Exception) UnmarshalJSON(b []byte) error {
	var errTypes map[string]string
	err := json.Unmarshal(b, &errTypes)
	if err != nil {
		return err
	}
	for msg, errType := range errTypes {
		if e == nil {
			e = make(map[string]error)
		}
		e[msg] = TypedError(errType, msg)
	}
	return nil
}

// HasDescription for exchange functionality
type HasDescription struct {
	CancelAllOrders      bool `json:"cancelAllOrders"`
	CancelOrder          bool `json:"cancelOrder"`
	CancelOrders         bool `json:"cancelOrders"`
	CORS                 bool `json:"CORS"`
	CreateDepositAddress bool `json:"createDepositAddress"`
	CreateLimitOrder     bool `json:"createLimitOrder"`
	CreateMarketOrder    bool `json:"createMarketOrder"`
	CreateOrder          bool `json:"createOrder"`
	Deposit              bool `json:"deposit"`
	EditOrder            bool `json:"editOrder"`
	FetchBalance         bool `json:"fetchBalance"`
	FetchBidsAsks        bool `json:"fetchBidsAsks"`
	FetchClosedOrders    bool `json:"fetchClosedOrders"`
	FetchCurrencies      bool `json:"fetchCurrencies"`
	FetchDepositAddress  bool `json:"fetchDepositAddress"`
	FetchDeposits        bool `json:"fetchDeposits"`
	FetchFundingFees     bool `json:"fetchFundingFees"`
	FetchL2OrderBook     bool `json:"fetchL2OrderBook"`
	FetchLedger          bool `json:"fetchLedger"`
	FetchMarkets         bool `json:"fetchMarkets"`
	FetchMyTrades        bool `json:"fetchMyTrades"`
	FetchOHLCV           bool `json:"fetchOHLCV"`
	FetchOpenOrders      bool `json:"fetchOpenOrders"`
	FetchOrder           bool `json:"fetchOrder"`
	FetchOrderBook       bool `json:"fetchOrderBook"`
	FetchOrderBooks      bool `json:"fetchOrderBooks"`
	FetchOrders          bool `json:"fetchOrders"`
	FetchTicker          bool `json:"fetchTicker"`
	FetchTickers         bool `json:"fetchTickers"`
	FetchTrades          bool `json:"fetchTrades"`
	FetchTradingFee      bool `json:"fetchTradingFee"`
	FetchTradingFees     bool `json:"fetchTradingFees"`
	FetchTradingLimits   bool `json:"fetchTradingLimits"`
	FetchTransactions    bool `json:"fetchTransactions"`
	FetchWithdrawals     bool `json:"fetchWithdrawals"`
	PrivateApi           bool `json:"privateApi"`
	PublicApi            bool `json:"publicApi"`
	Withdraw             bool `json:"withdraw"`
}

// StringSlice a custom type for handling variable JSON
type StringSlice []string

// UnmarshalJSON accepts both forms for StringSlice:
//   - ["s1", "s2"...]
//   - "s"
// For the latter, ss will hold a slice of one element "s"
// todo: unify to array form ?
func (ss *StringSlice) UnmarshalJSON(b []byte) (err error) {
	// try slice unmarshal
	var slice []string
	err = json.Unmarshal(b, &slice)
	if err == nil {
		*ss = slice
		return nil
	}
	// try string unmarshal
	var s string
	err2 := json.Unmarshal(b, &s)
	if err2 == nil {
		*ss = append(*ss, s)
		return nil
	}
	// return original error
	return err
}

// ApiUrls for different types of urls
type ApiUrls struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}

// UnmarshalJSON accepts both forms for ApiUrls:
//   - {"public": "urlpub", "private": "urlpriv"} or
//   - "url"
// For the latter, "url" is assigned to both a.Private and a.Public
// todo: unify to struct form ?
func (a *ApiUrls) UnmarshalJSON(b []byte) (err error) {
	// default struct unmarshal to compatible type
	type t ApiUrls
	err = json.Unmarshal(b, (*t)(a))
	if err == nil {
		return nil
	}
	// string unmarshal
	var s string
	err2 := json.Unmarshal(b, &s)
	if err2 == nil {
		a.Private = s
		a.Public = s
		return nil
	}
	// return original error
	return err
}

// Fees for using the exchange
type Fees struct {
	Trading TradingFees `json:"trading"`
	Funding FundingFees `json:"funding"`
}

// TradingFees on the exchange
type TradingFees struct {
	TierBased  bool             `json:"tierBased"`
	Percentage bool             `json:"percentage"`
	Maker      float64          `json:"maker"`
	Taker      float64          `json:"taker"`
	Tiers      TradingFeesTiers `json:"tiers"`
}

// TradingFeesTiers on the exchange
type TradingFeesTiers struct {
	Taker [][2]float64 `json:"taker"`
	Maker [][2]float64 `json:"maker"`
}

// FundingFees for using the exchange
type FundingFees struct {
	TierBased  bool               `json:"tierBased"`
	Percentage bool               `json:"percentage"`
	Deposit    map[string]float64 `json:"deposit"`
	Withdraw   map[string]float64 `json:"withdraw"`
}

// Account details
type Account struct {
	Free  map[string]float64 `json:"free"`
	Used  map[string]float64 `json:"used"`
	Total map[string]float64 `json:"total"`
}

// Balance details
type Balance struct {
	Free  float64 `json:"free"`
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

// Order structure
type Order struct {
	Id        string      `json:"id"`
	Timestamp JSONTime    `json:"timestamp"`
	Datetime  string      `json:"datetime"`
	Symbol    string      `json:"symbol"`
	Status    string      `json:"status"`
	Type      string      `json:"type"`
	Side      string      `json:"side"`
	Price     float64     `json:"price"`
	Cost      float64     `json:"cost"`
	Amount    float64     `json:"amount"`
	Filled    float64     `json:"filled"`
	Remaining float64     `json:"remaining"`
	Fee       float64     `json:"fee"`
	Info      interface{} `json:"info"`
}

func (o Order) String() string {
	return fmt.Sprintf("%s %f %s @%f (filled: %f)",
		o.Side, o.Amount, o.Symbol, o.Price, o.Filled)
}

// OrderBook struct
type OrderBook struct {
	Asks      [][]float64
	Bids      [][]float64
	Timestamp int64
	Datetime  string
	Nonce     string
}

// BookEntry struct
type BookEntry struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

// UnmarshalJSON accepts both forms for BookEntry:
//   - []float64 of size 2 or
//   - default BookEntry struct
func (o *BookEntry) UnmarshalJSON(b []byte) (err error) {
	// []float64 unmarshal
	var f []float64
	err = json.Unmarshal(b, &f)
	if err == nil {
		if len(f) != 2 {
			return fmt.Errorf("UnmarshalJSON: expected 2 floats for BookEntry, got %d", len(f))
		}
		o.Price, o.Amount = f[0], f[1]
		return nil
	}
	// default struct unmarshal to compatible type
	type t BookEntry
	err2 := json.Unmarshal(b, (*t)(o))
	if err2 == nil {
		return nil
	}
	return err
}

// Trade struct
type Trade struct {
	Id        string      `json:"id"`
	Symbol    string      `json:"symbol"`
	Amount    float64     `json:"amount"`
	Price     float64     `json:"price"`
	Timestamp JSONTime    `json:"timestamp"`
	Datetime  string      `json:"datetime"`
	Order     string      `json:"order"`
	Type      string      `json:"type"`
	Side      string      `json:"side"`
	Info      interface{} `json:"info"`
}

// Ticker struct
type Ticker struct {
	Symbol      string      `json:"symbol"`
	Ask         float64     `json:"ask"`
	Bid         float64     `json:"bid"`
	High        float64     `json:"high"`
	Low         float64     `json:"low"`
	Average     float64     `json:"average"`
	BaseVolume  float64     `json:"baseVolume"`
	QuoteVolume float64     `json:"quoteVolume"`
	Change      float64     `json:"change"`
	Open        float64     `json:"open"`
	Close       float64     `json:"close"`
	First       float64     `json:"first"`
	Last        float64     `json:"last"`
	Percentage  float64     `json:"percentage"`
	VWAP        float64     `json:"vwap"`
	Timestamp   JSONTime    `json:"timestamp"`
	Datetime    string      `json:"datetime"`
	Info        interface{} `json:"info"`
}

// Currency struct
type Currency struct {
	Id        string `json:"id"`
	Code      string `json:"code"`
	NumericId string `json:"numericId"`
	Precision int    `json:"precision"`
}

// DepositAddress struct
type DepositAddress struct {
	Currency string      `json:"currency"`
	Address  string      `json:"address"`
	Status   string      `json:"status"`
	Info     interface{} `json:"info"`
}

type ApiDecode struct {
	Path string
	Api string
	Method string
}

// OHLCV open, high, low, close, volume
type OHLCV struct {
	Timestamp JSONTime `json:"timestamp"`
	O         float64  `json:"o"`
	H         float64  `json:"h"`
	L         float64  `json:"l"`
	C         float64  `json:"c"`
	V         float64  `json:"v"`
}

// UnmarshalJSON accepts both forms for OHLCV:
//   - default struct or
//   - []float64 of size 6
func (o *OHLCV) UnmarshalJSON(b []byte) (err error) {
	// default struct unmarshal to compatible type
	type t OHLCV
	err = json.Unmarshal(b, (*t)(o))
	if err == nil {
		return nil
	}
	// []float64 unmarshal
	var f []float64
	err2 := json.Unmarshal(b, &f)
	if err2 != nil {
		return err2 // return float64 unmarshal error as it's the expected form
	}
	if len(f) != 6 {
		return fmt.Errorf("UnmarshalJSON: expected 6 floats for OHLCV, got %d", len(f))
	}

	err2 = json.Unmarshal([]byte(fmt.Sprintf("%f", f[0])), &o.Timestamp)
	if err2 != nil {
		return fmt.Errorf("UnmarshalJSON: couldn't unmarshal timestamp: %s", err2)
	}
	o.O, o.H, o.L, o.C, o.V = f[1], f[2], f[3], f[4], f[5]
	return nil
}

// Exchange is a common interface of methods
type ExchangeInterface interface {
	// FetchTickers(symbols []string, params map[string]interface{}) (map[string]Ticker, error)
	// FetchTicker(symbol string, params map[string]interface{}) (Ticker, error)
	// FetchOHLCV(symbol, tf string, since *JSONTime, limit *int, params map[string]interface{}) ([]OHLCV, error)
	FetchOrderBook(symbol string, limit int, params map[string]interface{}) (*OrderBook, error)
	// FetchL2OrderBook(symbol string, limit *int, params map[string]interface{}) (OrderBook, error)
	// FetchTrades(symbol string, since *JSONTime, params map[string]interface{}) ([]Trade, error)
	// FetchOrder(id string, symbol *string, params map[string]interface{}) (Order, error)
	// FetchOrders(symbol *string, since *JSONTime, limit *int, params map[string]interface{}) ([]Order, error)
	// FetchOpenOrders(symbol *string, since *JSONTime, limit *int, params map[string]interface{}) ([]Order, error)
	// FetchClosedOrders(symbol *string, since *JSONTime, limit *int, params map[string]interface{}) ([]Order, error)
	// FetchMyTrades(symbol *string, since *JSONTime, limit *int, params map[string]interface{}) ([]Trade, error)
	// FetchBalance(params map[string]interface{}) (Balances, error)
	//FetchCurrencies() (map[string]*Currency, error)
	FetchMarkets(params map[string]interface{}) ([]*Market, error)

	// CreateOrder(symbol, t, side string, amount float64, price *float64, params map[string]interface{}) (Order, error)
	// CancelOrder(id string, symbol *string, params map[string]interface{}) error

	// Describe() []byte
	//GetMarkets() map[string]*Market
	SetMarkets([]*Market, map[string]*Currency) (map[string]*Market, error)
	//GetMarketsById() map[string]Market
	//SetMarketsById(map[string]Market)
	//GetCurrencies() map[string]Currency
	//SetCurrencies(map[string]Currency)
	//GetCurrenciesById() map[string]Currency
	//SetCurrenciesById(map[string]Currency)
	//SetSymbols([]string)
	//SetIds([]string)
	// GetOrders() []Order
	LoadMarkets() (map[string]*Market, error)
	// LoadMarkets(reload bool, params map[string]interface{}) (map[string]*Market, error)
	// GetMarket(symbol string) (Market, error)
	// CreateLimitBuyOrder(symbol string, amount float64, price *float64, params map[string]interface{}) (Order, error)
	// CreateLimitSellOrder(symbol string, amount float64, price *float64, params map[string]interface{}) (Order, error)
	// CreateMarketBuyOrder(symbol string, amount float64, params map[string]interface{}) (Order, error)
	// CreateMarketSellOrder(symbol string, amount float64, params map[string]interface{}) (Order, error)
}

type ExchangeInterfaceInternal interface {
	ExchangeInterface
	Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (interface{}, error)
	ApiFuncDecode(function string) (path string, api string, method string, err error)
	ApiFunc(function string, params interface{}, headers map[string]interface{}, body interface{}) (response map[string]interface{}, err error)
	Fetch(url string, method string, headers map[string]interface{}, body interface{}) (response []byte, err error)
	Request(path string, api string, method string, params map[string]interface{}, headers map[string]interface{}, body interface{}) (response []byte, err error)
	Describe() []byte
}

// Exchange struct
type Exchange struct {
	ExchangeInfo
	ExchangeConfig

	Client         *http.Client
	Markets        map[string]*Market
	MarketsById    map[string]*Market
	Ids            []string
	Symbols        []string
	Currencies     map[string]*Currency
	CurrenciesById map[string]*Currency

	Child ExchangeInterfaceInternal
	ApiDecodeInfo map[string]*ApiDecode
	ApiUrls map[string]string
	DescribeMap map[string]interface{}
	Options map[string]interface{}
}

func (self *Exchange) Init(config *ExchangeConfig) (err error) {
	self.Child = self

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	self.Client = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10, //超时时间
	}

	return
}

func (self *Exchange) Describe() ([]byte) {
	return nil
}

func (self *Exchange) FetchMarkets(params map[string]interface{}) ([]*Market, error) {
	return nil, nil
}
func (self *Exchange) FetchOrderBook(symbol string, limit int, params map[string]interface{}) (*OrderBook, error) {
	return nil, errors.New("FetchOrderBook not supported yet")
}

func (self *Exchange) Sign(path string, api string, method string, params map[string]interface{}, headers interface{}, body interface{}) (interface{}, error) {
	return nil, nil
}

func (self *Exchange) MarketId(symbol string) (string, error) {
	return self.Markets[symbol].Id, nil
}

func (self *Exchange) SetMarkets(markets []*Market, currencies map[string]*Currency) (map[string]*Market, error) {
	symbols := make([]string, len(markets))
	Ids := make([]string, len(markets))
	marketsBySymbol := make(map[string]*Market, len(markets))
	marketsById := make(map[string]*Market, len(markets))
	baseCurrencies := make([]*Currency, 0)
	quoteCurrencies := make([]*Currency, 0)

	for i, market := range markets {
		marketsBySymbol[market.Symbol] = market
		marketsById[market.Id] = market
		symbols[i] = market.Symbol
		Ids[i] = market.Id
		// currency logic
		baseCurrency := new(Currency)
		if market.Base != "" {
			baseCurrency.Id = market.BaseId
			if baseCurrency.Id == "" {
				baseCurrency.Id = market.Base
			}
			baseCurrency.NumericId = market.BaseNumericId
			baseCurrency.Code = market.Base
			switch {
			case market.Precision.Base != 0:
				baseCurrency.Precision = market.Precision.Base
			case market.Precision.Amount != 0:
				baseCurrency.Precision = market.Precision.Amount
			default:
				baseCurrency.Precision = 8
			}
		}
		baseCurrencies = append(baseCurrencies, baseCurrency)
		quoteCurrency := new(Currency)
		if market.Quote != "" {
			quoteCurrency.Id = market.QuoteId
			if baseCurrency.Id == "" {
				quoteCurrency.Id = market.Quote
			}
			quoteCurrency.NumericId = market.QuoteNumericId
			quoteCurrency.Code = market.Quote
			switch {
			case market.Precision.Base != 0:
				quoteCurrency.Precision = market.Precision.Base
			case market.Precision.Amount != 0:
				quoteCurrency.Precision = market.Precision.Amount
			default:
				quoteCurrency.Precision = 8
			}
		}
		quoteCurrencies = append(quoteCurrencies, quoteCurrency)
	}
	allCurrencies := append(baseCurrencies, quoteCurrencies...)
	groupedCurrencies := make(map[string][]*Currency)
	for _, currency := range allCurrencies {
		groupedCurrencies[currency.Code] = append(groupedCurrencies[currency.Code], currency)
	}
	sortedCurrencies := make(map[string]*Currency)
	for code, currencies := range groupedCurrencies {
		for _, currency := range currencies {
			if sortedCurrencies[code] == nil {
				continue
			}
			if sortedCurrencies[code].Id == "" {
				sortedCurrencies[code] = currency
			}
			if sortedCurrencies[code].Precision < currency.Precision {
				sortedCurrencies[code] = currency
			}
		}
	}

	sort.Strings(symbols)
	sort.Strings(Ids)

	self.Symbols = symbols
	self.Ids = Ids
	self.MarketsById = marketsById
	self.Markets = marketsBySymbol

	if len(currencies) == 0 {
		xCurrencies := self.Currencies
		if xCurrencies == nil {
			xCurrencies = make(map[string]*Currency)
		}
		for code, currency := range sortedCurrencies {
			xCurrencies[code] = currency
		}
		self.Currencies = xCurrencies
	} else {
		self.Currencies = sortedCurrencies
	}
	currenciesById := self.CurrenciesById
	if len(currenciesById) == 0 {
		currenciesById = make(map[string]*Currency, len(currencies))
	}
	for _, currency := range sortedCurrencies {
		currenciesById[currency.Id] = currency
	}
	self.CurrenciesById = currenciesById
	return self.Markets, nil
}

// func (self *Exchange) LoadMarkets(reload bool, params map[string]interface{}) (map[string]*Market, error) {
func (self *Exchange) LoadMarkets() (map[string]*Market, error) {
	var currencies map[string]*Currency
	if self.Markets == nil {
		markets, err := self.Child.FetchMarkets(nil)
		if err != nil {
			return nil, err
		}
		return self.Child.SetMarkets(markets, currencies)
	}
	return self.Markets, nil
	/*
		if !reload && self.Markets != nil {
			if self.MarketsById == nil {
				var marketsSlice []Market
				for _, one := range self.Markets {
					marketsSlice = append(marketsSlice, one)
				}
				return self.SetMarkets(marketsSlice, nil)
			}
			return self.Markets, nil
		}
		var currencies map[string]Currency
		var err error
		if self.GetInfo().Has.FetchCurrencies {
			currencies, err = self.FetchCurrencies()
			if err != nil {
				return nil, err
			}
		}
		markets, err := self.FetchMarkets(params)
		if err != nil {
			return nil, err
		}

		return self.SetMarkets(markets, currencies)
	*/
}

func (self *Exchange) Request(
	path string,
	api string,
	method string,
	params map[string]interface{},
	headers map[string]interface{},
	body interface{},
) (response []byte, err error) {
	signInfo, err := self.Child.Sign(path, api, method, params, headers, body)
	if self.Verbose {
		fmt.Println("signinfo", signInfo)
	}
	if err != nil {
		return
	}
	fmt.Println(signInfo)
	return self.Child.Fetch(
		self.Member(signInfo, "url").(string),
		self.Member(signInfo, "method").(string),
		self.Member(signInfo, "headers").(map[string]interface{}),
		self.Member(signInfo, "body"),
		)
}

func (self *Exchange) PrepareRequestHeaders(req* http.Request, headers map[string]interface{}) error {
	//req.Header.Set("Accept-Encoding", "gzip, deflate")

	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}

	return nil
}

func (self *Exchange) Fetch(url string, method string, headers map[string]interface{}, body interface{}) (response []byte, err error) {
	var rbody []byte
	if body != nil {
		switch body.(type) {
		case string:
			rbody = []byte(body.(string))
		case []byte:
			rbody = body.([]byte)
		default:
			err = fmt.Errorf("Invalid Argument body: %v", body)
			return
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(rbody))
	if err != nil {
		return
	}

	err = self.PrepareRequestHeaders(req, headers)
	if err != nil {
		return nil, err
	}

	if self.Verbose {
		log.Println("\nRequest:", method, url, headers, body)
	}

	resp, err := self.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)

	if self.Verbose {
		log.Println("\nResponse:", method, url, resp.StatusCode, resp.Header, string(response))
	}

	if resp.StatusCode != 200 {
		err = errors.New("not 200")
		return
	}

	return
}

func (self *Exchange) RegSplit(text string, delimeter string) (result []string) {
    reg := regexp.MustCompile(delimeter)
    indexes := reg.FindAllStringIndex(text, -1)
    laststart := 0
    result = make([]string, len(indexes) + 1)
    for i, element := range indexes {
            result[i] = text[laststart:element[0]]
            laststart = element[1]
    }
    result[len(indexes)] = text[laststart:len(text)]
    return result
}

func (self *Exchange) DefineRestApi() (err error) {
	self.ApiDecodeInfo = make(map[string]*ApiDecode)

	if jsonApiInfo, ok := self.DescribeMap["api"].(map[string]interface{}); ok {
		for strApi, apiInfo := range jsonApiInfo {
			if methodInfo, ok:= apiInfo.(map[string]interface{}); ok {
				for strMethod, methodInfo := range methodInfo {
					if pathList, ok := methodInfo.([]interface{}); ok {
						for _, path := range pathList {
							if strPath, ok := path.(string); ok {
								var strDealPath string
								splitParts := self.RegSplit(strPath, "[^a-zA-Z0-9]")
								for _, part := range splitParts {
									strDealPath += strings.Title(part)
								}
								self.ApiDecodeInfo[strApi+strings.Title(strMethod)+strDealPath] = &ApiDecode{Api: strApi, Method: strings.ToUpper(strMethod), Path: strPath}

								if self.Verbose {
									log.Println("\napiDecodeInfo:", strApi, strPath, strMethod, strDealPath)
								}
							}
						}
					}
				}
			}
		}
	}

	return
}

func (self *Exchange) ApiFuncDecode(function string) (path string, api string, method string, err error) {
	return
}

func (self *Exchange) ApiFunc(function string, params interface{}, headers map[string]interface{}, body interface{}) (result map[string]interface{}, err error) {
	path, api, method, err := self.Child.ApiFuncDecode(function)
	if err != nil {
		return
	}

	resp, err := self.Child.Request(path, api, method, params.(map[string]interface{}), headers, body)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &result)
	if self.Verbose {
		fmt.Println(result)
	}

	return
}

func (self *Exchange) Iso8601(milliseconds int64) string {
	var seconds int64
	seconds = milliseconds / 1000
	return time.Unix(seconds, 0).Format("2006-01-02T15:04:05-0700")
}

func (self *Exchange) Milliseconds() int64 {
	return time.Now().Unix() * 1000
}

/*
// Init Exchange
func Init(conf ccxt.ExchangeConfig) (*Exchange, error) {
	var info ccxt.ExchangeInfo
	configFile := "/Users/stefan/Github/ccxt/templates/internal/app/bitmex/bitmex.json"
	f, err := os.Open(configFile)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	json.NewDecoder(f).Decode(&info)
	timeout := 10 * time.Second
	if conf.Timeout > 0 {
		timeout = conf.Timeout * time.Millisecond
	}
	client := &http.Client{Timeout: timeout}
	exchange := Exchange{
		Config: conf,
		Client: client,
		Info:   info,
	}
	return &exchange, nil
}
*/

// Exchanges returns the available exchanges
func Exchanges() []string {
	available := []string{"bitmex"}
	return available
}

// Info on the exchange
func Info(x Exchange) {
	var info interface{}
	info = x
	msg, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Printf("Error JSONMarshal: %v", err)
	}
	fmt.Println(string(msg))
}

func MapValues(input interface{}) []interface{} {
	v := reflect.ValueOf(input)
	keys := v.MapKeys()
	output := []interface{}{}
	for i, l := 0, v.Len(); i < l; i++ {
		output = append(output, v.MapIndex(keys[i]))
	}
	return output
}

/*
// LoadMarkets from the exchange that are active and available
func (self *Exchange) LoadMarkets(reload bool, params map[string]interface{}) (map[string]Market, error) {
	if !reload && self.Markets != nil {
		if self.MarketsById == nil {
			var marketsSlice []Market
			for _, one := range self.Markets {
				marketsSlice = append(marketsSlice, one)
			}
			return self.SetMarkets(marketsSlice, nil)
		}
		return self.Markets, nil
	}
	var currencies map[string]Currency
	var err error
	if self.GetInfo().Has.FetchCurrencies {
		currencies, err = self.FetchCurrencies()
		if err != nil {
			return nil, err
		}
	}
	markets, err := self.FetchMarkets(params)
	if err != nil {
		return nil, err
	}
	return self.SetMarkets(markets, currencies)
}

*/

func getCurrencyUsedOnOpenOrders(currency string) float64 {
	// TODO: implement getCurrencyUsedOnOpenOrders()
	return 0.0
}

/*
// setMarkets sets Markets, MarketsById, Symbols, SymbolsByIds, Ids,
// Currencies, CurrenciesByIds into the Exchange struct
func (self *Exchange) SetMarkets(markets []Market, currencies map[string]Currency) (map[string]Market, error) {
	symbols := make([]string, len(markets))
	Ids := make([]string, len(markets))
	marketsBySymbol := make(map[string]Market, len(markets))
	marketsById := make(map[string]Market, len(markets))
	baseCurrencies := make([]Currency, 0)
	quoteCurrencies := make([]Currency, 0)

	for i, market := range markets {
		marketsBySymbol[market.Symbol] = market
		marketsById[market.Id] = market
		symbols[i] = market.Symbol
		Ids[i] = market.Id
		// currency logic
		baseCurrency := new(Currency)
		if market.Base != "" {
			baseCurrency.Id = market.BaseId
			if baseCurrency.Id == "" {
				baseCurrency.Id = market.Base
			}
			baseCurrency.NumericId = market.BaseNumericId
			baseCurrency.Code = market.Base
			switch {
			case market.Precision.Base != 0:
				baseCurrency.Precision = market.Precision.Base
			case market.Precision.Amount != 0:
				baseCurrency.Precision = market.Precision.Amount
			default:
				baseCurrency.Precision = 8
			}
		}
		baseCurrencies = append(baseCurrencies, *baseCurrency)
		quoteCurrency := new(Currency)
		if market.Quote != "" {
			quoteCurrency.Id = market.QuoteId
			if baseCurrency.Id == "" {
				quoteCurrency.Id = market.Quote
			}
			quoteCurrency.NumericId = market.QuoteNumericId
			quoteCurrency.Code = market.Quote
			switch {
			case market.Precision.Base != 0:
				quoteCurrency.Precision = market.Precision.Base
			case market.Precision.Amount != 0:
				quoteCurrency.Precision = market.Precision.Amount
			default:
				quoteCurrency.Precision = 8
			}
		}
		quoteCurrencies = append(quoteCurrencies, *quoteCurrency)
	}
	allCurrencies := append(baseCurrencies, quoteCurrencies...)
	groupedCurrencies := make(map[string][]Currency)
	for _, currency := range allCurrencies {
		groupedCurrencies[currency.Code] = append(groupedCurrencies[currency.Code], currency)
	}
	sortedCurrencies := make(map[string]Currency)
	for code, currencies := range groupedCurrencies {
		for _, currency := range currencies {
			if sortedCurrencies[code].Id == "" {
				sortedCurrencies[code] = currency
			}
			if sortedCurrencies[code].Precision < currency.Precision {
				sortedCurrencies[code] = currency
			}
		}
	}

	sort.Strings(symbols)
	sort.Strings(Ids)

	self.Symbols = symbols
	self.Ids = Ids
	self.MarketsById = marketsById
	self.Markets = marketsBySymbol

	if len(currencies) == 0 {
		xCurrencies := self.Currencies
		if xCurrencies == nil {
			xCurrencies = make(map[string]Currency)
		}
		for code, currency := range sortedCurrencies {
			xCurrencies[code] = currency
		}
		self.Currencies = xCurrencies
	} else {
		self.Currencies = sortedCurrencies
	}
	currenciesById := self.CurrenciesById
	if len(currenciesById) == 0 {
		currenciesById = make(map[string]Currency, len(currencies))
	}
	for _, currency := range sortedCurrencies {
		currenciesById[currency.Id] = currency
	}
	self.CurrenciesById = currenciesById
	return self.Markets, nil
}
*/

func SortSliceByIndex(s [][]float64, idx int, descending bool) {
	if !descending {
		sort.Slice(s, func(i, j int) bool {
			// edge cases
			if len(s[i]) == 0 && len(s[j]) == 0 {
				return false // two empty slices - so one is not less than other i.e. false
			}
			if len(s[i]) == 0 || len(s[j]) == 0 {
				return len(s[i]) == 0 // empty slice listed "first" (change to != 0 to put them last)
			}

			// both slices len() > 0, so can test this now:
			return s[i][idx] < s[j][idx]
		})
	} else {
		sort.Slice(s, func(i, j int) bool {
			// edge cases
			if len(s[i]) == 0 && len(s[j]) == 0 {
				return false // two empty slices - so one is not less than other i.e. false
			}
			if len(s[i]) == 0 || len(s[j]) == 0 {
				return len(s[i]) == 0 // empty slice listed "first" (change to != 0 to put them last)
			}

			// both slices len() > 0, so can test this now:
			return s[i][idx] >= s[j][idx]
		})
	}
}

func ToFloat(x interface{}) (float64, error) {
	if str_x, ok := x.(string); ok {
		if y, err := strconv.ParseFloat(str_x, 64); err == nil {
			return y, err
		}
	}
	if y, ok := x.(float64); ok {
		return y, nil
	}
	return 0, nil
}

func (self *Exchange) ParseBidsAsks(bidsAsks []interface{}, priceKey int64, amountKey int64, out *[][]float64) error {
	if len(bidsAsks) == 0 {
		return nil
	}

	if _, ok := bidsAsks[0].([]interface{}); ok {
		for _, one := range bidsAsks {
			if bidAsk, ok := one.([]interface{}); ok {
				price := bidAsk[priceKey]
				amount := bidAsk[amountKey]
				if price != "" && amount != "" {
					priceF, err1 := ToFloat(price)
					amountF, err2 := ToFloat(amount)
					if err1 == nil && err2 == nil {
						*out = append(*out, []float64{priceF, amountF})
					}
				}
			}
		}
	} else {
		if _, ok := bidsAsks[0].(map[int]interface{}); ok {
		}
	}

	return nil
}


func (self *Exchange) Extend(maps ...interface{}) interface{} {
	size := len(maps)
	if size == 0 {
		return maps
	}
	if size == 1 {
		return maps[0]
	}
	output := make(map[string]interface{})
	for _, m := range maps {
		if oneMap, ok := m.(map[string]interface{}); ok {
			for k, v := range oneMap {
				output[k] = v
			}
		}
	}
	return output
}

func (self *Exchange) InMap(k interface{}, m interface{}) bool {
	if strMap, ok := m.(map[string]interface{}); ok {
		if _, ok := strMap[k.(string)]; ok {
			return true
		}
	}
	return false
}

func (self *Exchange) ToBool(v interface{}) bool {
	if v != nil {
		switch v.(type) {
		case bool:
			vv := v.(bool)
			return vv
		case string:
			v = v.(string)
			if v == "" {
				return false
			} else {
				return true
			}
		case int64:
			v = v.(int64)
			if v == 0 {
				return false
			} else {
				return true
			}
		case int:
			v = v.(int)
			if v == 0 {
				return false
			} else {
				return true
			}
		default:
			return true
		}
	} else {
		return false
	}
}

func (self *Exchange) SafeList(m map[string]interface{}, key string, defaultVal []interface{}) (val []interface{}) {
	if val, ok := m[key]; ok {
		if listVal, ok := val.([]interface{}); ok {
			return listVal
		}
	}
	return defaultVal
}

func (self *Exchange) SafeValue(m interface{}, key string, defaultVal interface{}) (val interface{}) {
	if mm, ok := m.(map[string]interface{}); ok {
		if val, ok := mm[key]; ok {
			return val
		}
	}
	return defaultVal
}

func NestedMapLookup(m map[string]interface{}, ks ...string) (rval interface{}, err error) {
	var ok bool

	if len(ks) == 0 { // degenerate input
		return nil, fmt.Errorf("NestedMapLookup needs at least one key")
	}
	if rval, ok = m[ks[0]]; !ok {
		return nil, fmt.Errorf("key not found; remaining keys: %v", ks)
	} else if len(ks) == 1 { // we've reached the final key
		return rval, nil
	} else if m, ok = rval.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("malformed structure at %#v", rval)
	} else { // 1+ more keys
		return NestedMapLookup(m, ks[1:]...)
	}
}

func (self *Exchange) ParseOrderBook(orderBook interface{}, timeStamp int64, bidsKey string, asksKey string, priceKey int64, amountKey int64) (*OrderBook, error) {
	var result OrderBook

	if orderBookMap, ok := orderBook.(map[string]interface{}); ok {
		if bids, ok := orderBookMap[bidsKey]; ok {
			if bidsList, ok := bids.([]interface{}); ok {
				err := self.ParseBidsAsks(bidsList, priceKey, amountKey, &result.Bids)
				if err != nil {
					SortSliceByIndex(result.Bids, 0, true)
				}
			}
		}
		if asks, ok := orderBookMap[asksKey]; ok {
			if asksList, ok := asks.([]interface{}); ok {
				err := self.ParseBidsAsks(asksList, priceKey, amountKey, &result.Asks)
				if err != nil {
					SortSliceByIndex(result.Asks, 0, false)
				}
			}
		}
		result.Timestamp = timeStamp
		// result.Datetime = time.Unix(timeStamp/1000, 0).Format("2006-01-02 15:04:05")

		return &result, nil
	}

	return nil, nil
}

func (self *Exchange) SafeInteger(d interface{}, key string, defaultVal int64) (ret int64) {
	if d, ok := d.(map[string]interface{}); ok {
		if val, ok := d[key]; ok {
			if intVal, ok := val.(int64); ok {
				return intVal
				}
			}
		}
	return defaultVal
}

func (self *Exchange) SafeString(d interface{}, key string, DefaultValue interface{}) (result string) {
	if d, ok := d.(map[string]interface{}); ok {
		if val, ok := d[key]; ok {
			if strVal, ok := val.(string); ok {
		        return strVal
		    }
		}
	}
	return DefaultValue.(string)
}

func (self *Exchange) SafeFloat(d interface{}, key string, DefaultValue float64) (result float64) {
	if d, ok := d.(map[string]interface{}); ok {
		if val, ok := d[key]; ok {
			if strVal, ok := val.(string); ok {
				fVal, err := strconv.ParseFloat(strVal, 64)
				if err == nil {
					return fVal
				}
			}
		}
	}
	return DefaultValue
}

func (self *Exchange) Omit(d map[string]interface{}, args interface{}) (result map[string]interface{}) {
	if argList, ok := args.([]string); ok {
		for _, arg := range argList {
			if _, ok := d[arg]; ok {
				delete(d, arg)
			}
		}
		return d
	}

	if arg, ok := args.(string); ok {
		delete(d, arg)
		return d
	}

	return d
}

func (self *Exchange) ExtractParams(s string) (result []string) {
	r := regexp.MustCompile(`{([^{}]*)}`)
	matches := r.FindAllStringSubmatch(s, -1)
	for _, v := range matches {
		result = append(result, v[1])
		fmt.Println(v[1])
	}
	return result
}

func (self *Exchange) ImplodeParams(s string, params interface{}) string {
	if paramsMap, ok := params.(map[string]interface{}); ok {
		for k, v := range paramsMap {
			s = strings.ReplaceAll(s, "{" + k + "}", fmt.Sprintf("%v", v))
		}
	}
	return s
}
var hashers = map[string]func() hash.Hash{
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
	"sha384": sha512.New384,
	"md5":    md5.New,
}

// Hash encodes the payload based on the available hashing algo
func Hash(payload, algo, encoding string) (string, error) {
	if hashers[algo] == nil {
		return "", fmt.Errorf("hash: unsupported algo \"%s\"", algo)
	}
	h := hashers[algo]()
	_, err := h.Write([]byte(payload))
	if err != nil {
		return "", fmt.Errorf("hash: %s", err)
	}
	return string(encode(h.Sum(nil), encoding)), nil
}

// HMAC encodes the payload based on the available hashing algo
func (self *Exchange) Hmac(payload, key, algo, encoding string) (string, error) {
	if hashers[algo] == nil {
		return "", fmt.Errorf("HMAC: unsupported hashing algo \"%s\"", algo)
	}
	h := hmac.New(hashers[algo], []byte(key))
	_, err := h.Write([]byte(payload))
	if err != nil {
		return "", fmt.Errorf("hmac: %s", err)
	}
	return string(encode(h.Sum(nil), encoding)), nil
}

// JWT creates a new signed token
func JWT(claims map[string]interface{}, secret string) (string, error) {
	var signer jwt.SigningMethod = jwt.SigningMethodHS256
	token := jwt.New(signer)
	token.Claims = jwt.MapClaims(claims)
	return token.SignedString([]byte(secret))
}

func encode(payload []byte, encoding string) []byte {
	var result []byte
	switch encoding {
	case "hex":
		result = []byte(hex.EncodeToString(payload))
	case "base64":
		buf := make([]byte, base64.StdEncoding.EncodedLen(len(payload)))
		base64.StdEncoding.Encode(buf, payload)
		result = buf
	default:
		result = payload
	}
	return result
}

func (self *Exchange) ParseBalance(balances map[string]interface{}) (pAccount *Account) {
	var account Account
	account.Free = make(map[string]float64)
	account.Used = make(map[string]float64)
	account.Total = make(map[string]float64)

	for currency, balance := range self.Omit(balances, []string{"info", "free", "used", "total"}) {
		if balance, ok := balance.(Balance); ok {
			account.Free[currency] = balance.Free
			account.Used[currency] = balance.Used
			account.Total[currency] = balance.Total
		}
	}

	return &account
}

func (self *Exchange) Uuid() string {
	return fmt.Sprintf("%v", uuid.NewV4())
}

func (self *Exchange) PriceToPrecision(symbol string, price float64) string {
	ret, _ := DecimalToPrecision(price, Round, 8, DecimalPlaces, NoPadding)
	return ret
}

func (self *Exchange) AmountToPrecision(symbol string, amount float64) string {
	ret, _ := DecimalToPrecision(amount, Truncate, 8, DecimalPlaces, NoPadding)
	return ret
}

func (self *Exchange) Account() (map[string]interface{}) {
	return map[string]interface{} {
		"free": nil,
		"used": nil,
		"total": nil,
	}
}

func (self *Exchange) SafeCurrencyCode(x interface{}) string {
	if _, ok := x.(string); ok {
		return "BTC"
	}
	return ""
}

func (self *Exchange) Length(o interface{}) int {
	switch reflect.TypeOf(o).Kind() {
	case reflect.Slice:
		return reflect.ValueOf(o).Len()
	case reflect.Map:
		return reflect.ValueOf(o).Len()
	default:
		return 0
	}
}

func (self *Exchange) Member (o interface{}, idx interface{}) interface{} {
	switch reflect.TypeOf(o).Kind() {
    case reflect.Slice:
        return reflect.ValueOf(o).Index(idx.(int)).Interface()
    case reflect.Map:
        return reflect.ValueOf(o).MapIndex(reflect.ValueOf(idx)).Interface()
    case reflect.Struct:
        return reflect.ValueOf(o).FieldByName(idx.(string)).Interface()
    }

	return nil
}

func (self *Exchange) Market (symbol string) (*Market){
	if self.Markets == nil {
		return nil
	}

	v, ok := self.Markets[symbol]
	if ok {
		return v
	} else {
		return nil
	}
}

func (self *Exchange)Unpack2(l interface{}) (interface{}, interface{}) {
	switch l.(type) {
	case []string:
		if ll, ok := l.([]string); ok {
			return ll[0], ll[1]
		}
	case []int64:
		if ll, ok := l.([]int64); ok {
			return ll[0], ll[1]
		}
	case []int:
		if ll, ok := l.([]int); ok {
			return ll[0], ll[1]
		}
	default:
		return nil, nil
	}
	return  nil,nil
}

func (self *Exchange) IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
    if condition {
        return a
    }
    return b
}

func (self *Exchange) TestNil (x interface{}) bool {
	switch x.(type) {
	case string:
		return x.(string) == ""
	case int:
		return x.(int) == 0
	case int64:
		return x.(int64) == 0
	case float64:
		return x.(float64) == 0.
	case float32:
		return x.(float32) == 0.
	default:
		return true
	}
}

func (self *Exchange) SetValue (x interface{}, k string, v interface{}) {
	if m, ok := x.(map[string]interface{}); ok {
		m[k] = v
	}
}

func (self *Exchange) CheckRequiredCredentials () {

}

func (self *Exchange) Urlencode (i interface{}) string {
	if m, ok := i.(map[string]interface{}); ok {
		v := url.Values{}
		for k, val := range m {
			v.Add(k, fmt.Sprintf("%v", val))
		}
		return v.Encode()
	}
	return ""
}

func (self *Exchange) Json (i interface{}) string {
	ret, err := json.Marshal(i)
	if err == nil {
		return string(ret)
	}
	return ""
}

func (self *Exchange) Encode (s interface{}) string {
	return s.(string)
}
func (self *Exchange) Decode (s interface{}) interface{}  {
	return s
}

func (self* Exchange) AddTwoInterface (a interface{}, b interface{}) interface{} {
	if a == nil || b == nil {
		return nil
	}

	switch a.(type) {
	case string:
		return a.(string) + b.(string)
	case int:
		return a.(int) + b.(int)
	case int64:
		return a.(int64) + b.(int64)
	case float64:
		return a.(float64) + b.(float64)
	case float32:
		return a.(float32) + b.(float32)
	default:
		 return nil
	}
}
