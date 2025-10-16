package inmemory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/takumi/personal-website/internal/model"
	"github.com/takumi/personal-website/internal/repository"
)

type blacklistRepository struct {
	mu      sync.RWMutex
	entries []model.BlacklistEntry
	nextID  int64
}

func NewBlacklistRepository() repository.BlacklistRepository {
	entries := make([]model.BlacklistEntry, len(defaultBlacklist))
	copy(entries, defaultBlacklist)
	var maxID int64
	for _, entry := range entries {
		if entry.ID > maxID {
			maxID = entry.ID
		}
	}
	return &blacklistRepository{
		entries: entries,
		nextID:  maxID + 1,
	}
}

func (r *blacklistRepository) ListBlacklistEntries(ctx context.Context) ([]model.BlacklistEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make([]model.BlacklistEntry, len(r.entries))
	copy(entries, r.entries)
	return entries, nil
}

func (r *blacklistRepository) AddBlacklistEntry(ctx context.Context, entry *model.BlacklistEntry) (*model.BlacklistEntry, error) {
	if entry == nil {
		return nil, repository.ErrInvalidInput
	}

	email := strings.ToLower(strings.TrimSpace(entry.Email))
	if email == "" {
		return nil, repository.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.entries {
		if strings.EqualFold(existing.Email, email) {
			return nil, repository.ErrDuplicate
		}
	}

	entry.ID = r.nextID
	r.nextID++
	entry.Email = email
	entry.CreatedAt = time.Now().UTC()
	r.entries = append(r.entries, *entry)

	added := *entry
	return &added, nil
}

func (r *blacklistRepository) RemoveBlacklistEntry(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, entry := range r.entries {
		if entry.ID == id {
			r.entries = append(r.entries[:idx], r.entries[idx+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func (r *blacklistRepository) FindBlacklistEntryByEmail(ctx context.Context, email string) (*model.BlacklistEntry, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, entry := range r.entries {
		if entry.Email == email {
			copyEntry := entry
			return &copyEntry, nil
		}
	}
	return nil, repository.ErrNotFound
}

var _ repository.BlacklistRepository = (*blacklistRepository)(nil)
