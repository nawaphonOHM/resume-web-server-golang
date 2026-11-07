// Package mongodb contains helpers for MongoDB BSON tests.
package mongodb

import (
	"testing"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

func marshalSPADocument(t *testing.T, doc domain.SPADocument) []byte {
	t.Helper()
	data, err := mongodb.Marshal(doc)
	if err != nil {
		t.Fatalf("failed to marshal SPADocument: %v", err)
	}
	return data
}

func unmarshalSPADocument(t *testing.T, data []byte, target any) {
	t.Helper()
	if err := mongodb.Unmarshal(data, target); err != nil {
		t.Fatalf("failed to unmarshal SPADocument: %v", err)
	}
}

func assertSPADocumentRoundTrip(t *testing.T, expected domain.SPADocument, data []byte) {
	t.Helper()
	var decoded domain.SPADocument
	unmarshalSPADocument(t, data, &decoded)
	assertEqual(t, expected.Location, decoded.Location)
	assertEqual(t, expected.FileName, decoded.FileName)
	assertEqual(t, expected.ID, decoded.ID)
}

func assertSPADocumentBSONFields(t *testing.T, doc domain.SPADocument, data []byte) {
	t.Helper()
	var fields map[string]any
	unmarshalSPADocument(t, data, &fields)
	assertStringBSONField(t, fields, "location", doc.Location)
	assertStringBSONField(t, fields, "fileName", doc.FileName)
}

func assertStringBSONField(t *testing.T, fields map[string]any, key, expected string) {
	t.Helper()
	value, ok := fields[key].(string)
	if !ok {
		t.Fatalf("expected %s BSON field to be a string, got %T", key, fields[key])
	}
	assertEqual(t, expected, value)
}

func assertZeroIDSPADocumentOmitted(t *testing.T, doc domain.SPADocument) {
	t.Helper()
	data := marshalSPADocument(t, domain.SPADocument{
		Location: doc.Location,
		FileName: doc.FileName,
	})
	var fields map[string]any
	unmarshalSPADocument(t, data, &fields)
	if _, ok := fields["_id"]; ok {
		t.Errorf("expected zero ID BSON field to be omitted, got %v", fields["_id"])
	}
}
