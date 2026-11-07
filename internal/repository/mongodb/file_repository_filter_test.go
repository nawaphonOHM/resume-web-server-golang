// Package mongodb tests MongoDB file repository filters.
package mongodb

import (
	"testing"

	mongotest "github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

func TestCleanFileName(t *testing.T) {
	assertEqual(t, testIndexFileName, CleanFileName(testIndexFileName))
	assertEqual(t, testIndexFileName, CleanFileName("/"+testIndexFileName))
	assertEqual(t, "app.js", CleanFileName("///app.js"))
	assertEqual(t, "style.css", CleanFileName("  style.css  "))
	assertEqual(t, "", CleanFileName(""))
	assertEqual(t, "", CleanFileName("///"))
}

func TestBuildFilter(t *testing.T) {
	filter := BuildFilter("/assets/main.js")
	expected := mongotest.D{{Key: "fileName", Value: "assets/main.js"}}

	assertEqual(t, len(expected), len(filter))
	assertEqual(t, expected[0].Key, filter[0].Key)
	assertEqual(t, expected[0].Value, filter[0].Value)
}
