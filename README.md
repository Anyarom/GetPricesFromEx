## Получение данных о сделках с биржи Poloniex

По Websocket получаю данные о состоявшихся сделках с биржи Poloniex (https://docs.poloniex.com/#price-aggregated-book)

На входе словарь в виде json-строки:
{"poloniex":["BTC_USDT", "TRX_USDT", "ETH_USDT"]}

Данные получаю по подписке

На выходе(лог в консоль):

```
type RecentTrade struct {
    Id        string    // ID транзакции
    Pair      string    // Торговая пара
    Price     float64   // Цена транзакции
    Amount    float64   // Объём транзакции
    Side      string    // Как биржа засчитала эту сделку (как buy или как sell)
    Timestamp time.Time // Время транзакции
}
```