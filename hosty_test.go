package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestParseEntriesReturnsEntries(t *testing.T) {
	fileContent := "#hosty-example\n127.0.0.1 example.com"
	entries := parseEntries(fileContent)

	entriesExpected := map[string]string{
		"example": "127.0.0.1 example.com",
	}

	if !reflect.DeepEqual(entriesExpected, entries) {
		t.Errorf("Unexpected entries: %v", entries)
	}
}

func TestListPrintsEntriesEnabled(t *testing.T) {
	entries := map[string]string{
		"example": "127.0.0.1 example.com",
	}

	list(entries, func(a ...interface{}) (n int, err error) {
		if len(a) > 1 {
			t.Errorf("Unexpected arguments length: %d", len(a))
		}

		if "hosty entries:\n✔ example\t127.0.0.1 example.com\n" != a[0] {
			t.Errorf("Unexpected arguments: %v", a[0])
		}

		return 0, nil
	})
}

func TestListPrintsEntriesDisabled(t *testing.T) {
	entries := map[string]string{
		"example": "#127.0.0.1 example.com",
	}

	list(entries, func(a ...interface{}) (n int, err error) {
		if len(a) > 1 {
			t.Errorf("Unexpected arguments length: %d", len(a))
		}

		if "hosty entries:\n✖ example\t#127.0.0.1 example.com\n" != a[0] {
			t.Errorf("Unexpected arguments: %v", a[0])
		}

		return 0, nil
	})
}

func TestListPrintsNoEntries(t *testing.T) {
	entries := make(map[string]string)
	called := false

	flag.Usage = func() {
		called = true
	}

	list(entries, func(a ...interface{}) (n int, err error) {
		if len(a) > 1 {
			t.Errorf("Unexpected arguments length: %d", len(a))
		}

		if "hosty has no entries!\n" != a[0] {
			t.Errorf("Unexpected arguments: %v", a[0])
		}

		return 0, nil
	})

	if !called {
		t.Error("flag.Usage() not called ")
	}
}
