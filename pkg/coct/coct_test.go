package coct

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"
)

func TestClean(t *testing.T) {
	expected := "foo bar"
	actual := clean(`
		foo    
		   bar

		`)

	if actual != expected {
		t.Fatalf("expected `%v`, got `%v`", expected, actual)
	}
}

func TestDayZero(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.DayZero.Year() != 2018 {
		t.Fatalf("expected `%d`, got `%d`", 2018, d.DayZero.Year())
	}

	if d.DayZero.Month() != time.Month(4) {
		t.Fatalf("expected `%d`, got `%d`", time.Month(4), d.DayZero.Month())
	}

	if d.DayZero.Day() != 12 {
		t.Fatalf("expected `%d`, got `%d`", 12, d.DayZero.Day())
	}
}

func TestDamLevel(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.Dams.Level != 27.2 {
		t.Fatalf("expected `%v`, got `%v`", 27.2, d.Dams.Level)
	}
}

func TestDamTrendAmount(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.Dams.Trend.Amount != 1.5 {
		t.Fatalf("expected `%v`, got `%v`", 1.5, d.Dams.Trend.Amount)
	}
}

func TestDamTrendDirection(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.Dams.Trend.Direction != -1 {
		t.Fatalf("expected `%v`, got `%v`", -1, d.Dams.Trend.Direction)
	}
}

func TestCapeTonianAmount(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.CapeTonians.Amount != 41.0 {
		t.Fatalf("expected `%v`, got `%v`", 41.0, d.CapeTonians.Amount)
	}
}

func TestCapeTonianAmountUnderReview(t *testing.T) {
	b, _ := ioutil.ReadFile("./test2.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.CapeTonians.Amount != -1 {
		t.Fatalf("expected `%v`, got `%v`", -1, d.CapeTonians.Amount)
	}
}

func TestCapeTonianTrendAmount(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.CapeTonians.Trend.Amount != 2 {
		t.Fatalf("expected `%v`, got `%v`", 2, d.CapeTonians.Trend.Amount)
	}
}

func TestCapeTonianTrendDirection(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.CapeTonians.Trend.Direction != 1 {
		t.Fatalf("expected `%v`, got `%v`", 1, d.CapeTonians.Trend.Direction)
	}
}

func TestCityProgress(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.City.Progress != 57.0 {
		t.Fatalf("expected `%v`, got `%v`", 57.0, d.City.Progress)
	}
}

func TestFindAreaAndType(t *testing.T) {
	a, ty := areaAndType("Cape Town Harbour (Desalination)")
	if a != "Cape Town Harbour" {
		t.Fatalf("expected `%v`, got `%v`", "Cape Town Harbour", a)
	}
	if ty != "Desalination" {
		t.Fatalf("expected `%v`, got `%v`", "Desalination", ty)
	}
}

func TestCityProjects(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if len(d.City.Projects) != 7 {
		t.Fatalf("expected `%v`, got `%v`", 7, len(d.City.Projects))
	}

	if d.City.Projects[0].Area != "Cape Town Harbour" {
		t.Fatalf("expected `%s`, got `%s`", "Cape Town Harbour", d.City.Projects[0].Area)
	}

	if d.City.Projects[0].Type != "Desalination" {
		t.Fatalf("expected `%s`, got `%s`", "Desalination", d.City.Projects[0].Type)
	}

	if d.City.Projects[0].Percentage != 50 {
		t.Fatalf("expected `%v`, got `%v`", 50, d.City.Projects[0].Percentage)
	}

	if d.City.Projects[0].Status != -1 {
		t.Fatalf("expected `%v`, got `%v`", -1, d.City.Projects[0].Status)
	}

	if d.City.Projects[3].Status != 1 {
		t.Fatalf("expected `%v`, got `%v`", 1, d.City.Projects[3].Status)
	}
}

func TestOtherProjects(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if len(d.Other.Projects) != 12 {
		t.Fatalf("expected `%v`, got `%v`", 12, len(d.Other.Projects))
	}

	if d.Other.Projects[0].Area != "Hout Bay" {
		t.Fatalf("expected `%s`, got `%s`", "Hout Bay", d.Other.Projects[0].Area)
	}

	if d.Other.Projects[0].Type != "Desalination" {
		t.Fatalf("expected `%s`, got `%s`", "Desalination", d.Other.Projects[0].Type)
	}

	if d.Other.Projects[0].Percentage != 45 {
		t.Fatalf("expected `%v`, got `%v`", 45, d.Other.Projects[0].Percentage)
	}
}

func TestCapeTonianDescription(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))
	expected := "Percentage of residents using 87 l or less per day."
	actual := d.CapeTonians.Description

	if actual != expected {
		t.Fatalf("expected `%v`, got `%v`", expected, actual)
	}
}

func TestCityDescription(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))
	expected := "The City's progress on securing alternative water sources."
	actual := d.City.Description

	if actual != expected {
		t.Fatalf("expected `%v`, got `%v`", expected, actual)
	}
}

func TestDamsDescription(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))
	expected := "Combined level of dams supplying the city. For more info click here."
	actual := d.Dams.Description

	if actual != expected {
		t.Fatalf("expected `%v`, got `%v`", expected, actual)
	}

	expectedURL := "http://www.capetown.gov.za/damlevels"
	actualURL := d.Dams.DescriptionURL

	if actualURL != expectedURL {
		t.Fatalf("expected `%v`, got `%v`", expectedURL, actualURL)
	}
}

func TestOtherDescription(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))
	expected := "Additional projects in advanced stage of planning that are ready to proceed if required."
	actual := d.Other.Description

	if actual != expected {
		t.Fatalf("expected `%v`, got `%v`", expected, actual)
	}
}

func TestMonthFromString(t *testing.T) {
	testCases := []struct {
		s    string
		want time.Month
	}{
		{"JANUARY", time.January},
		{"MARCH", time.March},
		{"DECEMBER", time.December},
	}
	for _, tc := range testCases {
		got, err := monthFromString(tc.s)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if got != tc.want {
			t.Errorf("wanted '%s' got '%s'", tc.want, got)
		}
	}

}

func TestStatsAsAtWeek(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))
	loc, _ := time.LoadLocation("Africa/Johannesburg")

	want := time.Date(2018, time.January, 22, 0, 0, 0, 0, loc)

	if !want.Equal(d.StatsAsAtWeek) {
		t.Fatalf("expected `%+v`, got `%+v`", want, d.StatsAsAtWeek)
	}
}
