package main

import "testing"

func TestContains(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		result := contains([]string{}, "needle")
		if result != false {
			t.Error("Expected false")
		}
	})

	t.Run("one-element slice", func(t *testing.T) {
		slice := []string{"needle"}
		result := contains(slice, "needle")
		if result != true {
			t.Error("Expected true for needle")
		}

		result = contains(slice, "fake needle")
		if result != false {
			t.Error("Expected false for fake needle")
		}
	})

	t.Run("many-element slice", func(t *testing.T) {
		slice := []string{"needle1", "needle2"}
		result := contains(slice, "needle1")
		if result != true {
			t.Error("Expected true for needle1")
		}

		result = contains(slice, "needle2")
		if result != true {
			t.Error("Expected true for needle1")
		}

		result = contains(slice, "fake needle")
		if result != false {
			t.Error("Expected false for fake needle")
		}
	})
}
