package storage

import (
	"io"
	"net/http"

	. "github.com/gophergala/nedomi/types"
)

// A unit of Storage
type Storage interface {
	// Returns a io.ReadCloser that will read from the `start`
	// of an object with ObjectId `id` to the `end`.
	Get(id ObjectID, start, end uint64) (io.ReadCloser, error)

	// Returns a io.ReadCloser that will read the whole file
	GetFullFile(id ObjectID) (io.ReadCloser, error)

	// Returns all headers for this object
	Headers(id ObjectID) (http.Header, error)

	// Discard an object from the storage
	Discard(id ObjectID) error

	// Discard an index of an Object from the storage
	DiscardIndex(index ObjectIndex) error
}
