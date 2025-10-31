package firestoredb

import (
	"fmt"
	"strings"
)

// CollectionName composes a collection name using an optional prefix.
func CollectionName(prefix, base string) string {
	trimmed := strings.TrimSpace(prefix)
	if trimmed == "" {
		return base
	}
	return fmt.Sprintf("%s_%s", trimmed, base)
}
