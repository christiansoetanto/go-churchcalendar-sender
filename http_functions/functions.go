package http_functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type AllLiturgicalDays struct {
	LiturgicalDaysEn []LiturgicalDay
	LiturgicalDaysLa []LiturgicalDay
}

type LiturgicalDay struct {
	Key                   string        `json:"key"`
	Date                  string        `json:"date"`
	Precedence            string        `json:"precedence"`
	Rank                  string        `json:"rank"`
	IsHolyDayOfObligation bool          `json:"isHolyDayOfObligation"`
	IsOptional            bool          `json:"isOptional"`
	Martyrology           []Martyrology `json:"martyrology"`
	Titles                []string      `json:"titles"`
	Calendar              Calendar      `json:"calendar"`
	Cycles                Cycles        `json:"cycles"`
	Name                  string        `json:"name"`
	RankName              string        `json:"rankName"`
	ColorName             []string      `json:"colorName"`
	SeasonNames           []string      `json:"seasonNames"`
}
type Calendar struct {
	WeekOfSeason          int    `json:"weekOfSeason,omitempty"`
	DayOfSeason           int    `json:"dayOfSeason,omitempty"`
	DayOfWeek             int    `json:"dayOfWeek,omitempty"`
	NthDayOfWeekInMonth   int    `json:"nthDayOfWeekInMonth,omitempty"`
	StartOfSeason         string `json:"startOfSeason,omitempty"`
	EndOfSeason           string `json:"endOfSeason,omitempty"`
	StartOfLiturgicalYear string `json:"startOfLiturgicalYear,omitempty"`
	EndOfLiturgicalYear   string `json:"endOfLiturgicalYear,omitempty"`
}
type Cycles struct {
	ProperCycle  string `json:"properCycle"`
	SundayCycle  string `json:"sundayCycle"`
	WeekdayCycle string `json:"weekdayCycle"`
	PsalterWeek  string `json:"psalterWeek"`
}
type Martyrology struct {
	Key               string   `json:"key"`
	CanonizationLevel string   `json:"canonizationLevel"`
	DateOfDeath       int      `json:"dateOfDeath"`
	Titles            []string `json:"titles,omitempty"`
}

type Messages struct {
	Messages []MessageItem `json:"messages"`
}
type MessageItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func Exec(_ http.ResponseWriter, _ *http.Request) {

	ROMCAL_API_FUNCTIONS_URL := "omitted"
	LINE_CHANNEL_ACCESS_TOKEN := "omitted"
	LINE_API_URL := "omitted"

	getText := func(liturgicalDays []LiturgicalDay, locale string) string {
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

	response, errFetch := http.Get(ROMCAL_API_FUNCTIONS_URL)
	if errFetch != nil {
		log.Fatal(errFetch)
	}

	data, errRead := ioutil.ReadAll(response.Body)

	if errRead != nil {
		log.Fatal(errRead)
	}
	var allLiturgicalDays AllLiturgicalDays
	err := json.Unmarshal(data, &allLiturgicalDays)
	if err != nil {
		log.Fatal(err)
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
	req, err2 := http.NewRequest("POST", LINE_API_URL, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", LINE_CHANNEL_ACCESS_TOKEN))
	if err2 != nil {
		log.Fatal(err2)
	}
	resp, err3 := http.DefaultClient.Do(req)
	if err3 != nil {
		log.Fatal(err3)
	}
	log.Printf("Finished executing with status %d", resp.StatusCode)
}
