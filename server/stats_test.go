package server

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	testStatsExpectedJSONResult = `{"startTime":"2006-01-02T13:24:56Z","reqCount":0,` +
		`"reqBytes":0,"routeStats":{"route1":{"requesBytes":202,"requestCounts":101},` +
		`"route2":{"requesBytes":204,"requestCounts":103}},"chatterStats":[],"roomStats":[]}`
)

func TestStatsNew(t *testing.T) {
	s := StatsNew()
	tp := reflect.TypeOf(s)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Stats not created as a pointer.")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Stats not created as a struct.")
	}
	if tp.Name() != "Stats" {
		t.Fatalf("Stats struct is not named correctly.")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Stats struct is empty.")
	}
}

func TestStatsIncrReqStats(t *testing.T) {
	t.Parallel()
	s := StatsNew()
	s.IncrReqStats(-1)
	if s.ReqCount != 1 {
		t.Errorf("ReqCount not incremented correctly. Expected: 1 Actual: %d", s.ReqCount)
	}
	if s.ReqBytes != 0 {
		t.Errorf("ReqBytes should not have been incremented or decremented. Expected: 0 Actual: %d", s.ReqBytes)
	}

	s.IncrReqStats(101)
	s.IncrReqStats(99)
	if s.ReqCount != 3 {
		t.Errorf("ReqCount not incremented correctly. Expected: 3 Actual: %d", s.ReqCount)
	}
	if s.ReqBytes != 200 {
		t.Errorf("ReqBytes should not have been incremented or decremented. Expected: 200 Actual: %d", s.ReqBytes)
	}
}

func TestStatsIncrRouteStats(t *testing.T) {
	t.Parallel()
	s := StatsNew()
	s.IncrRouteStats("Route1", -1)

	rs, ok := s.RouteStats["Route1"]
	if !ok {
		t.Errorf(`Stats RouteStats["Route1"] entry not created correctly.`)
	}
	if len(rs) != 1 {
		t.Errorf(`Stats RouteStats["Route1"] entry invalid size.`)
	}

	rc, ok := rs["requestCount"]
	if !ok {
		t.Errorf(`Stats RouteStats["Route1"]["requestCount"] entry not created correctly.`)
	}
	if rc != 1 {
		t.Errorf(`Stats RouteStats["Route1"]["requestCount"] should have been incremented.`)
	}

	rc, ok = rs["requestBytes"]
	if ok {
		t.Errorf(`Stats RouteStats["Route1"]["requestBytes"] entry should not have been created.`)
	}

	s = StatsNew()
	s.IncrRouteStats("Route2", -1)
	s.IncrRouteStats("Route2", 201)
	s.IncrRouteStats("Route2", 98)
	if s.RouteStats["Route2"]["requestCount"] != 3 {
		t.Errorf(`Stats["Route2"]["requestCount"] not incremented correctly.`)
	}
	_, ok = s.RouteStats["Route2"]["requestBytes"]
	if !ok {
		t.Errorf(`Stats RouteStats["Route1"]["requestBytes"] entry should have been created.`)
	}
	if s.RouteStats["Route2"]["requestBytes"] != 299 {
		t.Errorf(`Stats["Route2"]["requestBytes"] not incremented correctly.`)
	}
}

func TestStatString(t *testing.T) {
	t.Parallel()
	mockTime, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 13:24:56 -0000")
	s := StatsNew(func(sts *Stats) {
		sts.Start = mockTime
		sts.RouteStats = map[string]map[string]int64{
			"route1": map[string]int64{"requestCounts": 101, "requesBytes": 202},
			"route2": map[string]int64{"requestCounts": 103, "requesBytes": 204},
		}
	})
	actual := fmt.Sprint(s)
	if actual != testStatsExpectedJSONResult {
		t.Errorf("Stats not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testStatsExpectedJSONResult, actual)
	}
}
