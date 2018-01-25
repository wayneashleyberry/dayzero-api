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

func TestCapeTonianAmount(t *testing.T) {
	b, _ := ioutil.ReadFile("./test.html")
	d, _ := Parse(bytes.NewReader(b))

	if d.CapeTonians.Amount != 41.0 {
		t.Fatalf("expected `%v`, got `%v`", 41.0, d.CapeTonians.Amount)
	}
}
