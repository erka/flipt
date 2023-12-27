package object

import (
	"context"
	"sync"
	"time"

	"go.flipt.io/flipt/internal/config"
	"go.flipt.io/flipt/internal/containers"
	"go.flipt.io/flipt/internal/storage"
	storagefs "go.flipt.io/flipt/internal/storage/fs"

	"go.flipt.io/flipt/internal/storage/fs/object/blob"
	"go.uber.org/zap"
)

var _ storagefs.SnapshotStore = (*SnapshotStore)(nil)

// SnapshotStore represents an implementation of storage.SnapshotStore
// This implementation is backed by an Azure Blob Storage container
type SnapshotStore struct {
	*storagefs.Poller

	logger *zap.Logger

	mu   sync.RWMutex
	snap storage.ReadOnlyStore

	fsURL    string
	fsType   string
	bucket   string
	prefix   string
	pollOpts []containers.Option[storagefs.Poller]
}

// View accepts a function which takes a *StoreSnapshot.
// The SnapshotStore will supply a snapshot which is valid
// for the lifetime of the provided function call.
func (s *SnapshotStore) View(fn func(storage.ReadOnlyStore) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fn(s.snap)
}

// NewSnapshotStore constructs a Store
func NewSnapshotStore(ctx context.Context, logger *zap.Logger, opts ...containers.Option[SnapshotStore]) (*SnapshotStore, error) {
	s := &SnapshotStore{
		logger: logger,
		pollOpts: []containers.Option[storagefs.Poller]{
			storagefs.WithInterval(60 * time.Second),
		},
	}

	containers.ApplyAll(s, opts...)

	// fetch snapshot at-least once before returning store
	// to ensure we have some state to serve
	if _, err := s.update(ctx); err != nil {
		return nil, err
	}

	s.Poller = storagefs.NewPoller(s.logger, ctx, s.update, s.pollOpts...)

	go s.Poll()

	return s, nil
}

// WithFS configures the fs type and url
func WithFS(fsType config.ObjectSubStorageType, fsURL string) containers.Option[SnapshotStore] {
	return func(s *SnapshotStore) {
		s.fsURL = fsURL
		s.fsType = string(fsType)
	}
}

// WithBucket configures the bucket and prefix
func WithBucket(bucket, prefix string) containers.Option[SnapshotStore] {
	return func(s *SnapshotStore) {
		s.bucket = bucket
		s.prefix = prefix
	}
}

// WithPollOptions configures the poller options used when periodically updating snapshot state
func WithPollOptions(opts ...containers.Option[storagefs.Poller]) containers.Option[SnapshotStore] {
	return func(s *SnapshotStore) {
		s.pollOpts = append(s.pollOpts, opts...)
	}
}

// Update fetches a new snapshot and swaps it out for the current one.
func (s *SnapshotStore) update(ctx context.Context) (bool, error) {
	fs, err := blob.NewFS(ctx, s.fsURL, s.bucket, s.prefix)
	if err != nil {
		return false, err
	}

	snap, err := storagefs.SnapshotFromFS(s.logger, fs)
	if err != nil {
		return false, err
	}

	s.mu.Lock()
	s.snap = snap
	s.mu.Unlock()

	return true, nil
}

// String returns an identifier string for the store type.
func (s *SnapshotStore) String() string {
	return s.fsType
}
