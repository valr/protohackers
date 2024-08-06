package problem2

type Asset struct {
	prices map[int32]int32
}

func NewAsset() Asset {
	return Asset{prices: make(map[int32]int32)}
}

func (asset Asset) AddPrice(timestamp, price int32) {
	asset.prices[timestamp] = price
}

func (asset Asset) MeanPrice(mintime, maxtime int32) (meanPrice int32) {
	var total, count int64
	for timestamp, price := range asset.prices {
		if mintime <= maxtime && mintime <= timestamp && timestamp <= maxtime {
			total += int64(price)
			count++
		}
	}
	if count > 0 {
		meanPrice = int32(total / count)
	}
	return
}
