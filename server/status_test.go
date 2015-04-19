package server

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	testStatsExpectedJSONResult = `{"startTime":"2006-01-02T13:24:56Z","reqCount":0,` +
		`"reqBytes":0,"connNumAvail":1234,"roomStats":{"engineering1":{"reqBytes":202,` +
		`"reqCounts":101},"engineering2":{"reqBytes":204,"reqCounts":103}}}`
)

func TestStatusNew(t *testing.T) {
	s := StatusNew()
	tp := reflect.TypeOf(s)

	if tp.Kind() != reflect.Ptr {
		t.Fatalf("Status not created as a pointer.")
	}

	tp = tp.Elem()
	if tp.Kind() != reflect.Struct {
		t.Fatalf("Status not created as a struct.")
	}
	if tp.Name() != "Status" {
		t.Fatalf("Status struct is not named correctly.")
	}
	if !(tp.NumField() > 0) {
		t.Fatalf("Status struct is empty.")
	}
}

func TestStatusIncrReqStats(t *testing.T) {
	t.Parallel()
	s := StatusNew()
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

func TestStatusIncrRoomStats(t *testing.T) {
	t.Parallel()
	s := StatusNew()
	s.IncrRoomStats("Engineering", -1)

	rs, ok := s.RoomStats["Engineering"]
	if !ok {
		t.Errorf(`RoomStats["Engineering"] entry not created correctly.`)
	}
	i := len(rs)
	if i != 1 {
		t.Errorf(`RoomStats["Engineering"] entry invalid size. Expected: 1 Actual: %d`, i)
	}

	j, ok := rs["reqCount"]
	if !ok {
		t.Errorf(`RoomStats["Engineering"]["reqCount"] entry not created correctly.`)
	}
	if j != 1 {
		t.Errorf(`RoomStats["Engineering"]["reqCount"] should have been incremented. Expected: 1 Actual: %d`, j)
	}

	_, ok = rs["reqBytes"]
	if ok {
		t.Errorf(`RoomStats["Engineering"]["reqBytes"] entry should not have been created.`)
	}

	s = StatusNew()
	s.IncrRoomStats("Engineering", -1)
	s.IncrRoomStats("Engineering", 201)
	s.IncrRoomStats("Engineering", 98)
	j = s.RoomStats["Engineering"]["reqCount"]
	if j != 3 {
		t.Errorf(`RoomStats["Engineering"]["reqCount"] not incremented correctly. Expected: 3 Actual: %d`, j)
	}
	j, ok = s.RoomStats["Engineering"]["reqBytes"]
	if !ok {
		t.Errorf(`RoomStats["Engineering"]["reqBytes"] entry should have been created.`)
	}
	if j != 299 {
		t.Errorf(`RoomStats["Engineering"]["reqBytes"] not incremented correctly. Expected: 299 Actual: %d`, j)
	}
}

func TestStatusString(t *testing.T) {
	t.Parallel()
	mockTime, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 13:24:56 -0000")
	s := StatusNew(func(sts *Status) {
		sts.Start = mockTime
		sts.ConnNumAvail = 1234
		sts.RoomStats = map[string]map[string]int64{
			"engineering1": map[string]int64{"reqCounts": 101, "reqBytes": 202},
			"engineering2": map[string]int64{"reqCounts": 103, "reqBytes": 204},
		}
	})
	actual := fmt.Sprint(s)
	if actual != testStatsExpectedJSONResult {
		t.Errorf("Status not converted to json string.\n\nExpected: %s\n\nActual: %s\n",
			testStatsExpectedJSONResult, actual)
	}
}
