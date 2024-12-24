package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

var notification_hours = 2

func main() {
	//ファイル読み込み
	file, err := os.Open("sample.json") //fileとエラーの2つが返る。
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer file.Close()

	//データを受け取る
	var data []map[string]interface{}

	//dataに格納,できたらerr=nil
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	//fmt.Println("Data:", data[0]["QR-data"]) データ取り出し例

	//ソートを行う
	sort.Slice(data, func(i, j int) bool {
		//文字列化
		date_stringI, okI := data[i]["date"].(string)
		date_stringJ, okJ := data[j]["date"].(string)
		if !okI || !okJ {
			fmt.Println("Error", err)
			return false
		}

		//時間形式に変換
		timeI, errI := time.Parse(time.RFC3339, date_stringI)
		timeJ, errJ := time.Parse(time.RFC3339, date_stringJ)
		if errI != nil || errJ != nil {
			fmt.Println("Error", err)
			return false
		}

		// ソートの条件指定
		return timeI.Before(timeJ)
	})

	//ソート後のデータ表示（確認用）
	/* for _, item := range data {
		fmt.Println(item)
		fmt.Println()
	} */

	//JSONファイル作る
	outputFile, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	//書き込み
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ") // インデント
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("output.jsonに書き込みました")

	fmt.Println("直近", notification_hours, "時間の予定は、")
	// 現在時刻
	now := time.Now()

	for _, item := range data {
		dateString, ok := item["date"].(string) //json->文字列変換
		if !ok {
			fmt.Println("Error")
			continue
		}

		eventTime, err := time.Parse(time.RFC3339, dateString) //時間形式に変換
		if err != nil {
			fmt.Println("Error", err)
			continue
		}

		//eventTime.After(now)　現在より先の予定か？
		//eventTime.Sub(now) 予定時刻-現在が
		if eventTime.After(now) && eventTime.Sub(now) <= time.Duration(notification_hours)*time.Hour {
			fmt.Printf("%sに予定があります\n", eventTime.Format(time.RFC3339))
			break // 最初の条件を満たすイベントのみ通知
		}
	}
}
