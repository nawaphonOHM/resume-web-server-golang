// Package mongodb tests MongoDB file repository logic and query construction.
package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	testIndexFileName = "index.html"
)

var (
	errDatabaseTimeout = errors.New("database connection timeout")
)

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestNewMongoFileRepository(t *testing.T) {
	repo := NewMongoFileRepository(nil)
	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
}

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
	expected := mongodb.D{{Key: "fileName", Value: "assets/main.js"}}

	assertEqual(t, len(expected), len(filter))
	assertEqual(t, expected[0].Key, filter[0].Key)
	assertEqual(t, expected[0].Value, filter[0].Value)
}

func TestNilCollection(t *testing.T) {
	repo := NewMongoFileRepository(nil)
	ctx := context.Background()

	doc, err := repo.FindFile(ctx, testIndexFileName)
	if doc != nil {
		t.Errorf("expected nil document, got %v", doc)
	}
	if !errors.Is(err, ErrNilCollection) {
		t.Errorf("expected ErrNilCollection, got %v", err)
	}
}

func TestMapFindErrorNoDocuments(t *testing.T) {
	mapped := mapFindError(mongodb.ErrNoDocuments, testIndexFileName)
	if !errors.Is(mapped, domain.ErrFileNotFound) {
		t.Errorf("expected ErrFileNotFound, got %v", mapped)
	}
}

func TestMapFindErrorOther(t *testing.T) {
	mapped := mapFindError(errDatabaseTimeout, testIndexFileName)
	if !errors.Is(mapped, errDatabaseTimeout) {
		t.Errorf("expected wrapped errDatabaseTimeout, got %v", mapped)
	}
	if errors.Is(mapped, domain.ErrFileNotFound) {
		t.Error("did not expect ErrFileNotFound for timeout error")
	}
}

func TestSPADocumentBSON(t *testing.T) {
	doc := domain.SPADocument{
		HashValueFile: "abc123hash",
		Location:      "gs://my-bucket/index.html",
		FileName:      testIndexFileName,
	}

	data, err := mongodb.Marshal(doc)
	if err != nil {
		t.Fatalf("failed to marshal SPADocument: %v", err)
	}

	var decoded domain.SPADocument
	if unmarshalErr := mongodb.Unmarshal(data, &decoded); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal SPADocument: %v", unmarshalErr)
	}

	assertEqual(t, doc.HashValueFile, decoded.HashValueFile)
	assertEqual(t, doc.Location, decoded.Location)
	assertEqual(t, doc.FileName, decoded.FileName)
}
