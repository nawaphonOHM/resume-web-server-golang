// Package mongodb implements MongoDB data access for SPA file metadata.
package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/domain"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
)

// Sentinel errors for MongoDB repository.
var (
	ErrNilCollection = errors.New("mongodb collection is nil")
)

// MongoFileRepository implements domain.FileRepository using MongoDB driver.
type MongoFileRepository struct {
	collection *mongodb.Collection
}

// NewMongoFileRepository creates a new MongoFileRepository instance.
func NewMongoFileRepository(col *mongodb.Collection) *MongoFileRepository {
	return &MongoFileRepository{
		collection: col,
	}
}

// CleanFileName normalizes the input file name by trimming whitespace and leading slashes.
// For example, "/index.html" becomes "index.html" and "///app.js" becomes "app.js".
func CleanFileName(fileName string) string {
	return domain.CleanFileName(fileName)
}

// BuildFilter constructs the MongoDB filter for matching a file by name.
// It matches documents where fileName equals the normalized file name (without leading slashes).
func BuildFilter(fileName string) mongodb.D {
	cleaned := CleanFileName(fileName)
	return mongodb.D{{Key: "fileName", Value: cleaned}}
}

// FindFile queries the MongoDB collection for file metadata by file name.
func (r *MongoFileRepository) FindFile(ctx context.Context, fileName string) (*domain.SPADocument, error) {
	if r.collection == nil {
		return nil, ErrNilCollection
	}

	filter := BuildFilter(fileName)
	var doc domain.SPADocument
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, mapFindError(err, fileName)
	}

	return &doc, nil
}

func mapFindError(err error, fileName string) error {
	if errors.Is(err, mongodb.ErrNoDocuments) {
		return domain.ErrFileNotFound
	}
	return fmt.Errorf("find file %q in mongodb failed: %w", fileName, err)
}
