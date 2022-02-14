package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AllocateFree(t *testing.T) {
	cm, err := New(5*time.Minute, 100)
	if err != nil {
		t.Fatal(err)
	}

	key := "some-key"

	// allocate the key; must result in single entry and reference
	entry1 := cm.AllocateEntry(key)
	assert.EqualValues(t, entry1.manager, cm, "entry points to its own CacheManager")
	assert.Equal(t, entry1.key, key, "entry key is saved")

	assert.Equal(t, len(cm.entries), 1, "must be a single entry allocated")
	assert.Equal(t, len(cm.references), 1, "must be a single reference count allocated")

	assert.Equal(t, cm.entries[key], entry1, "allocated entry must be for this key")
	assert.Equal(t, cm.references[key], 1, "reference count must be for this key")

	// allocate same key again; must result in same entry, but two references
	entry2 := cm.AllocateEntry(key)
	assert.EqualValues(t, entry1, entry2, "entries for same key must be the same")

	assert.Equal(t, len(cm.entries), 1, "must still be a single entry allocated")
	assert.Equal(t, len(cm.references), 1, "must still be single reference allocated")

	assert.Equal(t, cm.entries[key], entry1, "key must still have a entry")
	assert.Equal(t, cm.references[key], 2, "key must now have two references")

	// free the original key; reference count must be decremented, but entry must not be freed
	entry1.Free()
	assert.Equal(t, len(cm.entries), 1, "must still be a single entry allocated")
	assert.Equal(t, len(cm.references), 1, "must still be a single reference allocated")
	assert.Equal(t, cm.references[key], 1, "key must now have a single reference")

	key2 := "another-key"

	// allocate a different key; must result in two entries and two references
	entry3 := cm.AllocateEntry(key2)
	assert.EqualValues(t, entry3.manager, cm, "entry points to its own CacheManager")

	assert.Equal(t, len(cm.entries), 2, "must be two entries allocated")
	assert.Equal(t, len(cm.references), 2, "must be two references allocated")

	assert.Equal(t, cm.entries[key2], entry3, "key2 must have a entry")
	assert.Equal(t, cm.references[key2], 1, "key2 must have single reference")

	// free the original key again; its reference count and entry must go away, but key2 is still there
	entry2.Free()
	assert.Equal(t, len(cm.entries), 1, "must be only one entry allocated now")
	assert.Equal(t, len(cm.references), 1, "must be only one reference allocated now")

	assert.EqualValues(t, cm.entries[key2], entry3, "key2 entry must still exist")
	assert.Equal(t, cm.references[key2], 1, "key2 reference count must still be one")

	// free the second key; reference count drops to zero, so entry and reference counter are freed
	entry3.Free()
	assert.Equal(t, len(cm.entries), 0, "must be no entries allocated")
	assert.Equal(t, len(cm.references), 0, "must be no references allocated")
}
