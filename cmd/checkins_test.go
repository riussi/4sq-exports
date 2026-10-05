package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	kml "github.com/twpayne/go-kml/v3"
)

func TestAddCheckins(t *testing.T) {
	time.Local = time.UTC

	body, err := os.ReadFile(filepath.Join("testdata", "checkins.json"))
	if err != nil {
		t.Fatal(err)
	}

	kDoc := kml.Document()
	addCheckins(body, kDoc)

	var got bytes.Buffer
	if err := kml.KML(kDoc).WriteIndent(&got, "", "  "); err != nil {
		t.Fatal(err)
	}

	goldenPath := filepath.Join("testdata", "checkins.kml")
	if os.Getenv("UPDATE_GOLDEN") != "" {
		if err := os.WriteFile(goldenPath, got.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Errorf("KML output mismatch\ngot:\n%s\nwant:\n%s", got.Bytes(), want)
	}
}
