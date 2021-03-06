package coct

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

type Dashboard struct {
	DayZero       time.Time   `json:"dayzero"`
	StatsAsAtWeek time.Time   `json:"statsAsAtWeek"`
	City          City        `json:"city"`
	Dams          Dams        `json:"dams"`
	CapeTonians   CapeTonians `json:"capetonians"`
	Other         Other       `json:"other"`
	Disclaimer    string      `json:"disclaimer"`
	Cached        bool        `json:"cached"`
}

type Other struct {
	Description string    `json:"description"`
	Projects    []Project `json:"projects"`
}

// Project represents an effort that the city is undergoing to increase water supply.
type Project struct {
	Area       string  `json:"area"`
	Type       string  `json:"type"`
	Percentage float64 `json:"percentage"`
	// Status will be -1, 0 or 1 - representing behind schedule, unknown and ahead of schedule.
	Status int `json:"status"`
}

type Trend struct {
	Amount float64 `json:"amount"`
	// Direction will be -1, 0 or 1 - representing negative, unknown or positive.
	Direction int `json:"direction"`
}

type Dams struct {
	Description    string  `json:"description"`
	DescriptionURL string  `json:"description_url"`
	Level          float64 `json:"level"`
	Trend          Trend   `json:"trend"`
}

type CapeTonians struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Trend       Trend   `json:"trend"`
}

type City struct {
	Description string    `json:"description"`
	Progress    float64   `json:"progress"`
	Projects    []Project `json:"projects"`
}

func GetCached(ctx context.Context, r *http.Request) (io.Reader, bool, error) {
	key := "api/dashboard"

	item, err := memcache.Get(ctx, key)
	if err == nil {
		reader := bytes.NewReader(item.Value)
		return reader, true, nil
	}

	fresh, err := Get(ctx)
	if err != nil {
		return bytes.NewReader([]byte("")), false, err
	}

	value, err := ioutil.ReadAll(fresh)
	if err != nil {
		return bytes.NewReader([]byte("")), false, err
	}

	newItem := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: time.Duration(time.Hour * 1),
	}
	memcache.Add(ctx, newItem)

	return bytes.NewReader(value), false, nil

}

func Get(ctx context.Context) (io.Reader, error) {
	client := urlfetch.Client(ctx)

	resp, err := client.Get("http://coct.co/water-dashboard/")
	if err != nil {
		return bytes.NewReader([]byte("")), err
	}

	if resp.StatusCode != http.StatusOK {
		return bytes.NewReader([]byte("")), errors.New("bad status: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes.NewReader([]byte("")), err
	}

	return bytes.NewReader(body), nil
}

func clean(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	return strings.Join(strings.Fields(s), " ")
}

func Parse(r io.Reader) (Dashboard, error) {
	var d Dashboard
	d.Disclaimer = "Data provided by the City of Cape Town (http://coct.co/water-dashboard/)"

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return d, err
	}

	d.City.Description = clean(doc.Find(".header").Eq(0).Find("p").Text())
	d.Dams.Description = clean(doc.Find(".header").Eq(2).Find("p").Text())
	d.CapeTonians.Description = clean(doc.Find(".header").Eq(3).Find("p").Text())
	d.Other.Description = clean(doc.Find(".header").Eq(4).Find("p").Text())
	href, exists := doc.Find(".header").Eq(2).Find("a").Attr("href")
	if exists {
		d.Dams.DescriptionURL = href
	}

	dayZero, err := getDayZero(doc)
	if err != nil {
		return d, err
	}
	d.DayZero = dayZero

	level, err := getDamLevel(doc)
	if err != nil {
		return d, err
	}
	d.Dams.Level = level

	damTrendAmount, err := getDamTrendAmount(doc)
	if err != nil {
		return d, err
	}
	d.Dams.Trend.Amount = damTrendAmount

	d.Dams.Trend.Direction = getDamTrendDirection(doc)

	amount, err := getCapeTonianAmount(doc)
	if err != nil {
		return d, err
	}
	d.CapeTonians.Amount = amount

	capeTonianTrendAmount, err := getCapeTonianTrendAmount(doc)
	if err != nil {
		return d, err
	}
	d.CapeTonians.Trend.Amount = capeTonianTrendAmount

	d.CapeTonians.Trend.Direction = getCapeTonianTrendDirection(doc)

	progress, err := getCityProgress(doc)
	if err != nil {
		return d, err
	}
	d.City.Progress = progress

	otherProjects, err := getOtherProjects(doc)
	if err != nil {
		return d, err
	}
	d.Other.Projects = otherProjects

	cityProjects, err := getCityProjects(doc)
	if err != nil {
		return d, err
	}
	d.City.Projects = cityProjects

	statsAsAtWeek, err := getStatsAsAtWeek(doc)
	if err != nil {
		return d, err
	}
	d.StatsAsAtWeek = statsAsAtWeek

	return d, nil
}

func areaAndType(s string) (string, string) {
	parts := strings.Split(s, "(")
	return strings.TrimSpace(parts[0]), strings.Replace(parts[1], ")", "", 1)
}

func getCityProjects(doc *goquery.Document) ([]Project, error) {
	ps := []Project{}
	doc.Find(".box .areas").Find(".area").Each(func(index int, el *goquery.Selection) {
		var p Project
		area, typeOfProject := areaAndType(el.Find("p").Text())
		p.Area = area
		p.Type = typeOfProject
		percentS := el.Find(".pval").Text()
		percentS = strings.Replace(percentS, "%", "", 1)
		percent, err := strconv.ParseFloat(percentS, 64)
		if err != nil {
			return
		}
		p.Percentage = percent
		if el.Is(".behind_schedule") {
			p.Status = -1
		} else if el.Is(".on_schedule") {
			p.Status = 1
		}
		ps = append(ps, p)
	})
	return ps, nil
}

func getOtherProjects(doc *goquery.Document) ([]Project, error) {
	ps := []Project{}
	doc.Find(".other_projects").Eq(1).Find(".area").Each(func(index int, el *goquery.Selection) {
		var p Project
		area, typeOfProject := areaAndType(el.Find("h4").Text())
		p.Area = area
		p.Type = typeOfProject
		percentS := el.Find(".pval").Text()
		percentS = strings.Replace(percentS, "%", "", 1)
		percent, err := strconv.ParseFloat(percentS, 64)
		if err != nil {
			return
		}
		p.Percentage = percent
		ps = append(ps, p)
	})
	return ps, nil
}

func getDamTrendAmount(doc *goquery.Document) (float64, error) {
	amount := doc.Find(".box").Eq(1).Find(".footer span").Text()
	amount = strings.Replace(amount, "%", "", -1)
	return strconv.ParseFloat(amount, 64)
}

func getDamTrendDirection(doc *goquery.Document) int {
	span := doc.Find(".box").Eq(1).Find(".footer span")
	if span.HasClass("down") {
		return -1
	} else if span.HasClass("up") {
		return 1
	}
	return 0
}

func getCapeTonianTrendAmount(doc *goquery.Document) (float64, error) {
	amount := doc.Find(".box").Eq(2).Find(".footer span").Text()
	amount = strings.Replace(amount, "%", "", -1)
	return strconv.ParseFloat(amount, 64)
}

func getCapeTonianTrendDirection(doc *goquery.Document) int {
	span := doc.Find(".box").Eq(2).Find(".footer span")
	if span.HasClass("down") {
		return -1
	} else if span.HasClass("up") {
		return 1
	}
	return 0
}

func getCapeTonianAmount(doc *goquery.Document) (float64, error) {
	amountS := doc.Find(".percentage_label").Eq(2).Text()
	amountS = strings.Replace(amountS, "%", "", -1)

	if len(amountS) == 0 {
		return -1, nil
	}

	return strconv.ParseFloat(amountS, 64)
}

func getDamLevel(doc *goquery.Document) (float64, error) {
	levelS := doc.Find(".percentage_label").Eq(1).Text()
	levelS = strings.Replace(levelS, "%", "", 1)

	return strconv.ParseFloat(levelS, 64)
}

func getCityProgress(doc *goquery.Document) (float64, error) {
	levelS := doc.Find(".percentage_label").Eq(0).Text()
	levelS = strings.Replace(levelS, "%", "", 1)

	return strconv.ParseFloat(levelS, 64)
}

func getDayZero(doc *goquery.Document) (time.Time, error) {
	h3 := clean(doc.Find("h3").First().Text())
	h3 = strings.Replace(h3, " ", "", -1)
	h3 = strings.Replace(h3, "\n", "", -1)

	if len(h3) < 8 {
		return time.Now(), errors.New("invalid string length for <h3>")
	}

	dayS := h3[0:2]
	monthS := h3[2:4]
	yearS := h3[4:8]

	day, err := strconv.Atoi(dayS)
	if err != nil {
		return time.Now(), err
	}

	month, err := strconv.Atoi(monthS)
	if err != nil {
		return time.Now(), err
	}

	year, err := strconv.Atoi(yearS)
	if err != nil {
		return time.Now(), err
	}

	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		return time.Now(), err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc), nil
}

func monthFromString(s string) (time.Month, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	now := time.Now()
	for i := 1; i <= 12; i++ {
		t := time.Date(2000, time.Month(i), 1, 1, 1, 1, 1, now.Location())
		if strings.ToUpper(t.Format("January")) == s {
			return time.Month(i), nil
		}
	}
	return time.January, nil
}

func getStatsAsAtWeek(doc *goquery.Document) (time.Time, error) {
	status := doc.Find(".status p").Text()
	status = strings.Replace(status, "STATS AS AT WEEK ", "", 1)
	fmt.Println(status)

	parts := strings.Split(status, " ")

	if len(parts) != 3 {
		return time.Now(), errors.New("too few parts")
	}

	dayS := parts[0]
	yearS := parts[2]
	month, err := monthFromString(parts[1])
	if err != nil {
		return time.Now(), err
	}

	day, err := strconv.Atoi(dayS)
	if err != nil {
		return time.Now(), err
	}

	year, err := strconv.Atoi(yearS)
	if err != nil {
		return time.Now(), err
	}

	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		return time.Now(), err
	}

	return time.Date(year, month, day, 0, 0, 0, 0, loc), nil
}
