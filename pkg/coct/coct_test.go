package coct

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"
)

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

func TestCityProjects(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if len(d.City.Projects) != 7 {
		t.Fatalf("expected `%v`, got `%v`", 7, len(d.City.Projects))
	}

	if d.City.Projects[0].Name != "Cape Town Harbour (Desalination)" {
		t.Fatalf("expected `%s`, got `%s`", "Cape Town Harbour (Desalination)", d.City.Projects[0].Name)
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

	if len(d.Other) != 12 {
		t.Fatalf("expected `%v`, got `%v`", 12, len(d.Other))
	}

	if d.Other[0].Name != "Hout Bay (Desalination)" {
		t.Fatalf("expected `%s`, got `%s`", "Hout Bay (Desalination)", d.Other[0].Name)
	}

	if d.Other[0].Percentage != 45 {
		t.Fatalf("expected `%v`, got `%v`", 45, d.Other[0].Percentage)
	}
}
