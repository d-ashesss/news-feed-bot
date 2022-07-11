package model

import (
	"testing"
)

func TestSubscriber_AddCategory(t *testing.T) {
	s := &Subscriber{}
	c := Category{}
	s.AddCategory(c)
	if len(s.Categories) != 1 {
		t.Errorf("AddCategory(): category was not added")
	}
}

func TestSubscriber_RemoveCategory(t *testing.T) {
	c1 := Category{ID: "test-cat-1"}
	c2 := Category{ID: "test-cat-2"}
	s := &Subscriber{Categories: []Category{c1, c2}}
	s.RemoveCategory(c1)
	if len(s.Categories) != 1 {
		t.Errorf("RemoveCategory(): category was not removed")
	}
}

func TestSubscriber_HasCategory(t *testing.T) {
	c1 := Category{ID: "test-cat-1"}
	c2 := Category{ID: "test-cat-2"}
	s := &Subscriber{Categories: []Category{c1, c2}}

	t.Run("does not have", func(t *testing.T) {
		cat := Category{ID: "test-cat-3"}
		if s.HasCategory(cat) {
			t.Errorf("want subscriber to not have category %v", cat)
		}
	})

	t.Run("has", func(t *testing.T) {
		cat := Category{ID: "test-cat-2"}
		if !s.HasCategory(cat) {
			t.Errorf("want subscriber to have category %v", cat)
		}
	})
}
