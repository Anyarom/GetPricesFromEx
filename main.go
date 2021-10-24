package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type RecentTrade struct {
	Id        string    // ID транзакции
	Pair      string    // Торговая пара (из списка выше)
	Price     float64   // Цена транзакции
	Amount    float64   // Объём транзакции
	Side      string    // Как биржа засчитала эту сделку (как buy или как sell)
	Timestamp time.Time // Время транзакции
}

func main() {
	c, _, err := websocket.DefaultDialer.Dial("wss://api2.poloniex.com", nil)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = c.WriteMessage(websocket.TextMessage, []byte("{ \"command\":\"subscribe\", \"channel\": \"USDT_BTC,USDT_TRX,USDT_ETH\" }"))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer c.Close()

	// заведем map для pairs
	mapPairs := make(map[float64]string)
	for {
		//	var message string
		_, message, errRead := c.ReadMessage()
		if errRead != nil {
			log.Fatal().Err(errRead).Send()
		}

		var resp []interface{}
		errUnmarshal := json.Unmarshal(message, &resp)
		if errUnmarshal != nil {
			log.Error().Err(errUnmarshal).Send()
			continue
		}

		if len(resp) != 3 {
			continue
		}

		for _, event := range resp[2].([]interface{}) {
			currencyPair := resp[0].(float64)

			eventData := event.([]interface{})
			if eventData[0].(string) == "i" {
				mapPairs[currencyPair] = eventData[1].(map[string]interface{})["currencyPair"].(string)
			}

			if eventData[0].(string) == "t" {

				price, errParse := strconv.ParseFloat(eventData[3].(string), 64)
				if errParse != nil {
					log.Fatal().Err(errParse).Send()
				}

				amount, errParseAmount := strconv.ParseFloat(eventData[4].(string), 64)
				if errParseAmount != nil {
					log.Fatal().Err(errParseAmount).Send()
				}

				var sideNew string
				switch eventData[2].(float64) {
				case 0:
					sideNew = "sell"
				case 1:
					sideNew = "buy"
				}

				recentTrade := RecentTrade{
					Id:        eventData[1].(string),
					Pair:      mapPairs[resp[0].(float64)],
					Price:     price,
					Amount:    amount,
					Side:      sideNew,
					Timestamp: time.Unix(int64(eventData[5].(float64)), 0),
				}

				log.Debug().Interface("RecentTrade", recentTrade).Send()
			}
		}
	}
}
