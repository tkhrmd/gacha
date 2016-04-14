package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

const (
	WEIGHT_UNIT = 1000
)

type Gacha struct {
	Classes []Class         `json:"classes"`
	Items   map[string]Item `json:"items"`
}

type Class struct {
	Weight int    `json:"weight"`
	Filter Filter `json:"filter"`
}

type Filter struct {
	Group       string         `json:"group"`
	ItemWeights map[string]int `json:"item_weights"`
}

type Item struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

type Stack map[string]int

func main() {

	// json読み込み
	setting, _ := ioutil.ReadFile("setting.json")

	// 構造体にマッピング
	var gacha Gacha
	json.Unmarshal(setting, &gacha)

	// ガチャ回す
	id := draw(gacha)
	fmt.Println(id)
}

// ガチャ
func draw(gacha Gacha) string {

	// 当選対象アイテムリスト
	bingoList, totalStack := generateBingoList(gacha)

	// 抽選
	var bingo string
	rand.Seed(time.Now().UnixNano())
	dice := rand.Intn(totalStack - 1)
	for id, weight := range bingoList {
		bingo = id
		if dice < weight {
			break
		}
		dice -= weight
	}

	return bingo
}

// 全対象アイテムの重み付け
func generateBingoList(gacha Gacha) (Stack, int) {

	list := Stack{}
	totalStack := 0

	for _, class := range gacha.Classes {

		// 当選アイテムを選定
		classList, classStack := filterItems(gacha.Items, class.Filter)

		for id, weight := range classList {
			itemStack := (weight * class.Weight * 1000 / classStack)
			list[id] = itemStack
			totalStack += itemStack
		}
	}

	return list, totalStack
}

// アイテムの選定とクラス内での重み付け
func filterItems(items map[string]Item, filter Filter) (Stack, int) {

	list := Stack{}
	totalStack := 0

	for id, item := range items {

		// 重み
		var weight int

		if itemWeight, ok := filter.ItemWeights[id]; ok {
			// アイテムの重みを直接指定
			weight = itemWeight
		} else {
			// グループフィルタ
			if item.Group != filter.Group {
				continue
			}

			// ここにフィルタを追加する

			weight = WEIGHT_UNIT
		}

		list[id] = weight
		totalStack += weight
	}

	return list, totalStack
}
