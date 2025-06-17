package yit

import (
	"go.yaml.in/yaml/v3"
	"testing"
)

func TestAnyMatch(t *testing.T) {
	t.Run("returns true if any element matches the predicate", func(t *testing.T) {
		result := FromNode(docNode).AnyMatch(All)
		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("returns false if no element matches the predicate", func(t *testing.T) {
		result := FromNode(docNode).AnyMatch(None)
		if result {
			t.Errorf("expected false, got true")
		}
	})
}

func TestAllMatch(t *testing.T) {
	t.Run("returns true if all elements match the predicate", func(t *testing.T) {
		result := FromNodes(
			scalarNode("a"),
			scalarNode("a"),
		).AllMatch(WithValue("a"))
		if !result {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("returns false if any element does not match the predicate", func(t *testing.T) {
		result := FromNodes(
			scalarNode("a"),
			scalarNode("b"),
		).AllMatch(WithValue("a"))
		if result {
			t.Errorf("expected false, got true")
		}
	})
}

func TestToArray(t *testing.T) {
	t.Run("adds all the iterated elements to an array", func(t *testing.T) {
		nodes := []*yaml.Node{
			{Value: "a"},
			{Value: "b"},
			{Value: "c"},
		}
		result := FromNodes(nodes...).ToArray()
		if len(result) != len(nodes) {
			t.Errorf("expected length %d, got %d", len(nodes), len(result))
		}
		for i := range nodes {
			if result[i].Value != nodes[i].Value {
				t.Errorf("expected value %q at index %d, got %q", nodes[i].Value, i, result[i].Value)
			}
		}
	})
}
