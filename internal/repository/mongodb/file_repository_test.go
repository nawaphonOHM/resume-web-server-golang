// Package mongodb tests MongoDB file repository logic and query construction.
package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	mongotypes "github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
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
		Location: "gs://my-bucket/index.html",
		FileName: testIndexFileName,
		ID:       mongotypes.NewObjectID(),
	}
	data := marshalSPADocument(t, doc)
	assertSPADocumentRoundTrip(t, doc, data)
	assertSPADocumentBSONFields(t, doc, data)
	assertZeroIDSPADocumentOmitted(t, doc)
}
