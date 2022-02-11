package main

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
