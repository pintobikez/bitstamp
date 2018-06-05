package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// currency_pair btcusd, btceur, eurusd, xrpusd, xrpeur, xrpbtc, ltcusd, ltceur, ltcbtc, ethusd, etheur, ethbtc, bchusd, bcheur, bchbtc

var _url string = "https://www.bitstamp.net/api/v2"

type API struct {
	cfg *Config
}

type Config struct {
	User   string `yaml:"user,omitempty"`
	Key    string `yaml:"key,omitempty"`
	Secret string `yaml:"secret,omitempty"`
}

func New(c *Config) *API {
	return &API{cfg: c}
}

type AccountBalanceResult struct {
	UsdBalance   float64 `json:"usd_balance,string"`
	BtcBalance   float64 `json:"btc_balance,string"`
	BchBalance   float64 `json:"bch_balance,string"`
	EurBalance   float64 `json:"eur_balance,string"`
	XrpBalance   float64 `json:"xrp_balance,string"`
	LtcBalance   float64 `json:"ltc_balance,string"`
	EthBalance   float64 `json:"eth_balance,string"`
	UsdReserved  float64 `json:"usd_reserved,string"`
	BtcReserved  float64 `json:"btc_reserved,string"`
	BchReserved  float64 `json:"bch_reserved,string"`
	EurReserved  float64 `json:"eur_reserved,string"`
	XrpReserved  float64 `json:"xrp_reserved,string"`
	LtcReserved  float64 `json:"ltc_reserved,string"`
	EthReserved  float64 `json:"eth_reserved,string"`
	UsdAvailable float64 `json:"usd_available,string"`
	BtcAvailable float64 `json:"btc_available,string"`
	BchAvailable float64 `json:"bch_available,string"`
	EurAvailable float64 `json:"eur_available,string"`
	XrpAvailable float64 `json:"xrp_available,string"`
	LtcAvailable float64 `json:"ltc_available,string"`
	EthAvailable float64 `json:"eth_available,string"`
	BtcUsdFee    float64 `json:"btcusd_fee,string"`
	BtcEurFee    float64 `json:"btceur_fee,string"`
	EurUsdFee    float64 `json:"eurusd_fee,string"`
	XrpUsdFee    float64 `json:"xrpusd_fee,string"`
	XrpEurFee    float64 `json:"xrpeur_fee,string"`
	XrpBtcFee    float64 `json:"xrpbtc_fee,string"`
	LtcUsdFee    float64 `json:"ltcusd_fee,string"`
	LtcEurFee    float64 `json:"ltceur_fee,string"`
	LtcBtcFee    float64 `json:"ltcbtc_fee,string"`
	EthUsdFee    float64 `json:"ethusd_fee,string"`
	EthEurFee    float64 `json:"etheur_fee,string"`
	EthBtcFee    float64 `json:"ethbtc_fee,string"`
}

type TickerResult struct {
	Last      float64 `json:"last,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Vwap      float64 `json:"vwap,string"`
	Volume    float64 `json:"volume,string"`
	Bid       float64 `json:"bid,string"`
	Ask       float64 `json:"ask,string"`
	Timestamp string  `json:"timestamp"`
	Open      float64 `json:"open,string"`
}

type BuyOrderResult struct {
	Id       int64   `json:"id,string"`
	DateTime string  `json:"datetime"`
	Type     int     `json:"type,string"`
	Price    float64 `json:"price,string"`
	Amount   float64 `json:"amount,string"`
}

type SellOrderResult struct {
	Id       int64   `json:"id,string"`
	DateTime string  `json:"datetime"`
	Type     int     `json:"type,string"`
	Price    float64 `json:"price,string"`
	Amount   float64 `json:"amount,string"`
}

type OpenOrder struct {
	Id           int64   `json:"id,string"`
	DateTime     string  `json:"datetime"`
	Type         int     `json:"type,string"`
	Price        float64 `json:"price,string"`
	Amount       float64 `json:"amount,string"`
	CurrencyPair string  `json:"currency_pair"`
}

func (a *API) SetAuth(clientId string, key string, secret string) {
	a.cfg = &Config{User: clientId, Key: key, Secret: secret}
}

// privateQuery submits an http.Request with key, sig & nonce
func (a *API) privateQuery(path string, values url.Values, v interface{}) error {

	// if no result interface, return
	if v == nil {
		return nil
	}

	// parse the bitstamp URL
	endpoint, err := url.Parse(_url)
	if err != nil {
		return err
	}

	// set the endpoint for this request
	endpoint.Path += path

	// add required key, signature & nonce to values
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	mac := hmac.New(sha256.New, []byte(a.cfg.Secret))
	mac.Write([]byte(nonce + a.cfg.User + a.cfg.Key))
	values.Set("key", a.cfg.Key)
	values.Set("signature", strings.ToUpper(hex.EncodeToString(mac.Sum(nil))))
	values.Set("nonce", nonce)

	// encode the url.Values in the body
	reqBody := strings.NewReader(values.Encode())

	// create the request
	//log.Println(endpoint.String(), values)
	req, err := http.NewRequest("POST", endpoint.String(), reqBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// submit the http request
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// read the body of the http message into a byte array
	return json.NewDecoder(r.Body).Decode(v)
}

func (a *API) AccountBalance() (*AccountBalanceResult, error) {
	balance := &AccountBalanceResult{}
	err := a.privateQuery("/balance/", url.Values{}, balance)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (a *API) Ticker(pair string) (*TickerResult, error) {
	ticker := &TickerResult{}
	err := a.privateQuery("/ticker/"+pair+"/", url.Values{}, ticker)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

func (a *API) BuyLimitOrder(pair string, amount float64, price float64, amountPrecision, pricePrecision int) (*BuyOrderResult, error) {
	// set params
	var v = url.Values{}
	v.Add("amount", strconv.FormatFloat(amount, 'f', amountPrecision, 64))
	v.Add("price", strconv.FormatFloat(price, 'f', pricePrecision, 64))

	// make request
	result := &BuyOrderResult{}
	err := a.privateQuery("/buy/"+pair+"/", v, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *API) BuyMarketOrder(pair string, amount float64) (*BuyOrderResult, error) {
	// set params
	var v = url.Values{}
	v.Add("amount", strconv.FormatFloat(amount, 'f', 8, 64))

	// make request
	result := &BuyOrderResult{}
	err := a.privateQuery("/buy/market/"+pair+"/", v, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *API) SellLimitOrder(pair string, amount float64, price float64, amountPrecision, pricePrecision int) (*SellOrderResult, error) {
	// set params
	var v = url.Values{}
	v.Add("amount", strconv.FormatFloat(amount, 'f', amountPrecision, 64))
	v.Add("price", strconv.FormatFloat(price, 'f', pricePrecision, 64))

	// make request
	result := &SellOrderResult{}
	err := a.privateQuery("/sell/"+pair+"/", v, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *API) SellMarketOrder(pair string, amount float64) (*SellOrderResult, error) {
	// set params
	var v = url.Values{}
	v.Add("amount", strconv.FormatFloat(amount, 'f', 8, 64))

	// make request
	result := &SellOrderResult{}
	err := a.privateQuery("/sell/market/"+pair+"/", v, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *API) CancelOrder(orderId int64) {
	// set params
	var v = url.Values{}
	v.Add("id", strconv.FormatInt(orderId, 10))

	// make request
	a.privateQuery("/cancel_order/", v, nil)
}

func (a *API) OpenOrders() (*[]OpenOrder, error) {
	// make request
	result := &[]OpenOrder{}
	err := a.privateQuery("/open_orders/all/", url.Values{}, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// key = "fdFhXhybSwglNYMExnJcwU6LfivSrlet"
// secret = "tl6IupZ6l9W5H0ZL5mfaiIKqptvQxtE7"
