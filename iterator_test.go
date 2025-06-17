package yit

import (
	"go.yaml.in/yaml/v3"
	"strings"
	"testing"
)

func TestFromNode(t *testing.T) {
	node := &yaml.Node{}
	next := FromNode(node)

	item, ok := next()
	if item != node {
		t.Errorf("expected item to be node, got %#v", item)
	}
	if !ok {
		t.Errorf("expected ok to be true")
	}

	item, ok = next()
	if item != nil {
		t.Errorf("expected item to be nil, got %#v", item)
	}
	if ok {
		t.Errorf("expected ok to be false")
	}
}

func TestRecurseNodes(t *testing.T) {
	tests := []struct {
		name   string
		yaml   string
		values []*yaml.Node
	}{
		{
			name:   "scalar",
			yaml:   "a",
			values: []*yaml.Node{docNode, scalarNode("a")},
		},
		{
			name:   "sequence",
			yaml:   "[a, b, c]",
			values: []*yaml.Node{docNode, seqNode, scalarNode("a"), scalarNode("b"), scalarNode("c")},
		},
		{
			name:   "map",
			yaml:   "{a: b}",
			values: []*yaml.Node{docNode, mapNode, scalarNode("a"), scalarNode("b")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := toYAMLTest(tt.yaml, t)
			next := FromNode(doc).RecurseNodes()
			for i, value := range tt.values {
				node, ok := next()
				if !ok {
					t.Errorf("expected ok to be true at index %d", i)
				}
				if node.Kind != value.Kind || node.Value != value.Value {
					t.Errorf("expected node.Kind=%d, Value=%q; got Kind=%d, Value=%q", value.Kind, value.Value, node.Kind, node.Value)
				}
			}
			_, ok := next()
			if ok {
				t.Errorf("expected ok to be false at end")
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name   string
		yaml   string
		values []*yaml.Node
	}{
		{
			name:   "scalar",
			yaml:   "",
			values: nil,
		},
		{
			name:   "sequence",
			yaml:   "[a, b, c]",
			values: []*yaml.Node{scalarNode("a"), scalarNode("b"), scalarNode("c")},
		},
		{
			name:   "map",
			yaml:   "a: b\nc: d",
			values: []*yaml.Node{scalarNode("a"), scalarNode("b"), scalarNode("c"), scalarNode("d")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := toYAMLTest(tt.yaml, t)
			next := FromNode(doc).Values().Values()
			for i, value := range tt.values {
				node, ok := next()
				if !ok {
					t.Errorf("expected ok to be true at index %d", i)
				}
				if node.Kind != value.Kind || node.Value != value.Value {
					t.Errorf("expected node.Kind=%d, Value=%q; got Kind=%d, Value=%q", value.Kind, value.Value, node.Kind, node.Value)
				}
			}
			_, ok := next()
			if ok {
				t.Errorf("expected ok to be false at end")
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Run("passes items through satisfying the predicate", func(t *testing.T) {
		next := FromNode(docNode).Filter(All)
		node, ok := next()
		if !ok {
			t.Errorf("expected ok to be true")
		}
		if node != docNode {
			t.Errorf("expected node to be docNode")
		}
	})

	t.Run("does not pass items that do not satisfy the predicate", func(t *testing.T) {
		next := FromNode(docNode).Filter(None)
		_, ok := next()
		if ok {
			t.Errorf("expected ok to be false")
		}
	})

	t.Run("predicate is not invoked when there are no items", func(t *testing.T) {
		empty := Iterator(func() (*yaml.Node, bool) {
			return nil, false
		})

		called := false
		next := empty.Filter(func(*yaml.Node) bool {
			called = true
			return true
		})

		_, ok := next()
		if ok {
			t.Errorf("expected ok to be false")
		}
		if called {
			t.Errorf("predicate should not be called")
		}
	})
}

func TestMapKeys(t *testing.T) {
	t.Run("returns the keys of a map", func(t *testing.T) {
		next := FromNode(toYAMLTest("a: b\nc: d\ne: f", t)).
			RecurseNodes().
			Filter(WithKind(yaml.MappingNode)).
			MapKeys()

		for _, value := range []string{"a", "c", "e"} {
			node, ok := next()
			if !ok {
				t.Errorf("expected ok to be true for key %q", value)
			}
			if node.Value != value {
				t.Errorf("expected node.Value=%q, got %q", value, node.Value)
			}
		}
		_, ok := next()
		if ok {
			t.Errorf("expected ok to be false at end")
		}
	})

	t.Run("returns nothing for sequences", func(t *testing.T) {
		next := FromNode(toYAMLTest("[a, b, c, d]", t)).
			RecurseNodes().
			Filter(WithKind(yaml.SequenceNode)).
			MapKeys()
		_, ok := next()
		if ok {
			t.Errorf("expected ok to be false for sequence")
		}
	})
}

func TestMapValues(t *testing.T) {
	t.Run("returns the values of a map", func(t *testing.T) {
		next := FromNode(toYAMLTest("a: b\nc: d\ne: f", t)).
			RecurseNodes().
			Filter(WithKind(yaml.MappingNode)).
			MapValues()

		for _, value := range []string{"b", "d", "f"} {
			node, ok := next()
			if !ok {
				t.Errorf("expected ok to be true for value %q", value)
			}
			if node.Value != value {
				t.Errorf("expected node.Value=%q, got %q", value, node.Value)
			}
		}
		_, ok := next()
		if ok {
			t.Errorf("expected ok to be false at end")
		}
	})

	t.Run("returns nothing for sequences", func(t *testing.T) {
		next := FromNode(toYAMLTest("[a, b, c, d]", t)).
			RecurseNodes().
			Filter(WithKind(yaml.SequenceNode)).
			MapValues()
		_, ok := next()
		if ok {
			t.Errorf("expected ok to be false for sequence")
		}
	})
}

func TestIterate(t *testing.T) {
	repeater := func(next Iterator) Iterator {
		return func() (node *yaml.Node, ok bool) {
			node, ok = next()
			if ok {
				node = scalarNode(strings.Repeat(node.Value, 2))
			}
			return
		}
	}

	next := FromNodes(scalarNode("a")).
		Iterate(repeater).
		Iterate(repeater)

	node, ok := next()
	if !ok {
		t.Errorf("expected ok to be true")
	}
	if node.Value != "aaaa" {
		t.Errorf("expected node.Value to be 'aaaa', got %q", node.Value)
	}
}

func TestValuesForMap(t *testing.T) {
	t.Run("returns the values of a map matching the key/value predicates", func(t *testing.T) {
		next := FromNode(toYAMLTest("a: b\nc: d\ne: f", t)).
			RecurseNodes().
			Filter(WithKind(yaml.MappingNode)).
			ValuesForMap(All, func(node *yaml.Node) bool {
				return node.Value == "d"
			})

		node, ok := next()
		if !ok {
			t.Errorf("expected ok to be true")
		}
		if node.Value != "d" {
			t.Errorf("expected node.Value to be 'd', got %q", node.Value)
		}

		_, ok = next()
		if ok {
			t.Errorf("expected ok to be false at end")
		}
	})

	t.Run("returns nothing for sequences", func(t *testing.T) {
		next := FromNode(toYAMLTest("[a, b, c, d]", t)).
			RecurseNodes().
			Filter(WithKind(yaml.SequenceNode)).
			ValuesForMap(All, All)
		_, ok := next()
		if ok {
			t.Errorf("expected ok to be false for sequence")
		}
	})
}

func TestFromIterators(t *testing.T) {
	next := FromIterators(
		FromNode(&yaml.Node{Value: "a"}),
		FromNode(&yaml.Node{Value: "b"}),
		FromNode(&yaml.Node{Value: "c"}),
	)

	for _, value := range []string{"a", "b", "c"} {
		node, ok := next()
		if !ok {
			t.Errorf("expected ok to be true for value %q", value)
		}
		if node.Value != value {
			t.Errorf("expected node.Value=%q, got %q", value, node.Value)
		}
	}
	_, ok := next()
	if ok {
		t.Errorf("expected ok to be false at end")
	}
}

// --- helpers ---

var mapNode = &yaml.Node{
	Kind: yaml.MappingNode,
}

var seqNode = &yaml.Node{
	Kind: yaml.SequenceNode,
}

var docNode = &yaml.Node{
	Kind: yaml.DocumentNode,
}

func scalarNode(val string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: val,
	}
}

func toYAMLTest(s string, t *testing.T) *yaml.Node {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(s), &node)
	if err != nil {
		t.Fatalf("failed to unmarshal yaml: %v", err)
	}
	return &node
}
