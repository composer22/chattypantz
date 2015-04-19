package server

import (
	"regexp"
	"testing"
)

const (
	testV4UUIDRegExpFmt = "^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$"
)

func TestUtilsCreateV4UUID(t *testing.T) {
	t.Parallel()
	r, _ := regexp.Compile(testV4UUIDRegExpFmt)
	for i := 0; i < 10; i++ {
		uuid := createV4UUID()
		if !r.MatchString(uuid) {
			t.Errorf("UUID not V4 standard. Result: %s", uuid)
			break
		}
	}

	uuid1 := createV4UUID()
	uuid2 := createV4UUID()
	if uuid1 == uuid2 {
		t.Errorf("UUID not being created uniquely.\n\nuuid1: %s\n\nuuid2: %s", uuid1, uuid2)
	}
}
