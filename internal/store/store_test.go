package store

import "testing"

func TestSet(t *testing.T) {
	store := NewStore()

	t.Run("Successful Set", func(t *testing.T) {
		err := store.Set("key", "value")
		if err != nil {
			t.Errorf("expected no error %v", err)
		}
	})

	t.Run("Empty Key Error", func(t *testing.T) {
		err := store.Set("", "value")
		if err == nil {
			t.Errorf("expected error for empty key, got nil")
		}
	})
}

func TestGet(t *testing.T) {
	store := NewStore()
	store.Set("key", "value")
	store.Set("empty", "")

	t.Run("Existing Key", func(t *testing.T) {
		value, ok := store.Get("key")
		if !ok {
			t.Error("expected key to exist, but it didn't")
		}
		if value != "value" {
			t.Errorf("expected 'value', got '%s'", value)
		}
	})

	t.Run("Existing Key with Empty String", func(t *testing.T) {
		value, ok := store.Get("empty")
		if !ok {
			t.Error("expected key to exist, but it didn't")
		}
		if value != "" {
			t.Errorf("expected '', got '%s'", value)
		}
	})

	t.Run("Non-Existent Key", func(t *testing.T) {
		value, ok := store.Get("xyz")
		if ok {
			t.Errorf("expected key to not exist, but got value '%s'", value)
		}
		if value != "" {
			t.Errorf("expected empty string, got '%s'", value)
		}
	})

	t.Run("Empty String Key", func(t *testing.T) {
		value, ok := store.Get("")
		if ok {
			t.Errorf("expected empty string key to not exist, but got value '%s'", value)
		}
		if value != "" {
			t.Errorf("expected empty string, got '%s'", value)
		}
	})
}

func TestDelete(t *testing.T) {
	store := NewStore()
	store.Set("key1", "value1")
	store.Set("key2", "value2")
	store.Set("key3", "value3")

	t.Run("Delete Existing Key", func(t *testing.T) {
		err := store.Delete("key2")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		_, ok := store.Get("key2")
		if ok {
			t.Error("expected key to be deleted, but it still exists")
		}
	})

	t.Run("Delete Non Existing Key", func(t *testing.T) {
		err := store.Delete("key4")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
