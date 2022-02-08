package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	//TODO utk convert ke functions dengan trigger pubsub, masukkan ke method
	//func ExecPubSub(ctx context.Context, m PubSubMessage) error{}

	//dealing with Go env variable is pain in the butt... so just declare it in the code itself,
	//no one will be able to read the functions code anyway
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal(errEnv)
	}

	//declare this in functions
	functionsUrl := os.Getenv("ROMCAL_API_FUNCTIONS_URL")
	lineApiUrl, lineAccessToken := os.Getenv("LINE_API_URL"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	response, err := http.Get(functionsUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	data, _ := ioutil.ReadAll(response.Body)

	var liturgicalDays []LiturgicalDay
	json.Unmarshal(data, &liturgicalDays)

	var messageItems []MessageItem
	messageItems = append(
		messageItems, MessageItem{
			Type: "text",
			Text: "Hello! Today is: ",
		},
	)

	for i, day := range liturgicalDays {
		if i == 3 {
			break
		}
		//[day, date] //if memorial/feast/solemnity [rank] [name] in [seasonName] season using [color] color
		rankName, isHolyDayOfObligation, name, seasonNames, colorNames := day.RankName,
			day.IsHolyDayOfObligation, day.Name, day.SeasonNames, day.ColorName
		var text = "- "
		if isHolyDayOfObligation {
			text += fmt.Sprintf("A Holy Day of Obligation of ")
		}
		if rankName == "memorial" || rankName == "feast" || rankName == "solemnity" {
			text += fmt.Sprintf("%s of ", strings.Title(rankName))
			text += fmt.Sprintf("%s ", name)

			if len(seasonNames) > 0 {
				text += fmt.Sprintf("in %s season ", seasonNames[0])
			}
		} else {
			text += fmt.Sprintf("%s ", name)
		}
		if len(colorNames) > 0 {
			text += fmt.Sprintf("using %s color ", colorNames[0])
		}

		messageItem := MessageItem{
			Type: "text",
			Text: text,
		}
		messageItems = append(messageItems, messageItem)
	}

	messages := Messages{
		Messages: messageItems,
	}
	jsonValue, _ := json.Marshal(messages)
	req, err2 := http.NewRequest("POST", lineApiUrl, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lineAccessToken))
	if err2 != nil {
		log.Fatal(err2)
		return
	}
	resp, err3 := http.DefaultClient.Do(req)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	log.Printf("Finished executing with status %d", resp.StatusCode)

}
