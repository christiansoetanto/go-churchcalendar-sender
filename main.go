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
	"time"
)

func getText(liturgicalDays []LiturgicalDay, locale string) string {
	var text string
	for _, day := range liturgicalDays {

		text += "• "
		//[day, date] //if memorial/feast/solemnity [rank] [name] in [seasonName] season using [color] color
		rank, rankName, isHolyDayOfObligation, name, seasonNames, colorNames := strings.ToLower(day.Rank), day.RankName,
			day.IsHolyDayOfObligation, day.Name, day.SeasonNames, day.ColorName
		if rank == "memorial" || rank == "feast" || rank == "solemnity" {
			if locale == "en" {
				text += fmt.Sprintf("%s of %s", strings.Title(rankName), name)
			} else {
				text += fmt.Sprintf("%s %s", strings.Title(rankName), name)
			}
			if len(seasonNames) > 0 {
				if locale == "en" {
					text += fmt.Sprintf(" in the %s", seasonNames[0])
				} else {
					text += fmt.Sprintf(" %s", seasonNames[0])
				}
			}
		} else {
			text += fmt.Sprintf("%s", name)
		}
		if len(colorNames) > 0 {
			text += fmt.Sprintf(". %s", strings.Title(colorNames[0]))
		}

		if isHolyDayOfObligation && locale == "en" {
			text += fmt.Sprintf(". A Holy Day of Obligation.")
		}

		text += "\n"
	}
	return text

}
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

	var allLiturgicalDays AllLiturgicalDays
	errUnmarshal := json.Unmarshal(data, &allLiturgicalDays)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
		return
	}

	currentTime := time.Now()
	greetingText := fmt.Sprintf(
		"Hello! Today is %s, %d %s\n\nThe Roman Catholic Church is celebrating: \n", currentTime.Weekday(),
		currentTime.Day(), currentTime.Month(),
	)
	textEn, textLa := getText(allLiturgicalDays.LiturgicalDaysEn, "en"), getText(
		allLiturgicalDays.LiturgicalDaysLa,
		"la",
	)

	messages := Messages{
		Messages: []MessageItem{
			{
				Type: "text", Text: fmt.Sprintf(
					"%s\n%s\n%s", greetingText, textEn, textLa,
				),
			},
		},
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
