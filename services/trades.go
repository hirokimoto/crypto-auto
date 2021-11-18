package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hirokimoto/crypto-auto/utils"
)

func StableTokens(wg *sync.WaitGroup, pairs utils.Pairs) {
	for _, item := range pairs.Data.Pairs {
		defer wg.Done()
		c := make(chan string, 1)
		go utils.Post(c, "swaps", item.Id)
		fmt.Print(".")
		stableToken(c, item.Id)
	}
}

func TradableTokens(wg *sync.WaitGroup, pairs utils.Pairs, target chan string) {
	defer wg.Done()

	for _, item := range pairs.Data.Pairs {
		c := make(chan string, 1)
		go utils.Post(c, "swaps", item.Id)
		fmt.Print(".")
		tradableToken(c, item.Id, target)
	}
}

func StoreAndRemovePair(pair string) (err error) {
	if IsExist(pair) {
		err = RemoveOnePair(pair)
		if err == nil {
			// Alert("Removed!", pair)
		}
	} else {
		err = WriteOnePair(pair)
		if err == nil {
			// Alert("Saved!", pair)
		}
	}
	return err
}

func stableToken(pings <-chan string, id string) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		min, max, _, _, _, _ := minMax(swaps)
		last, _ := priceOfSwap(swaps.Data.Swaps[0])
		_, _, period := periodOfSwaps(swaps)
		howold := howMuchOld(swaps)

		if (max-min)/last < 0.1 && period > 24 && howold < 24 {
			fmt.Println(id)
		}
	}
}

func tradableToken(pings <-chan string, id string, target chan string) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		min, max, _, _, _, _ := minMax(swaps)
		last, _ := priceOfSwap(swaps.Data.Swaps[0])
		_, _, period := periodOfSwaps(swaps)
		howold := howMuchOld(swaps)

		if (max-min)/last > 0.1 && period < 6 && howold < 24 {
			target <- id
		}
	}
}
