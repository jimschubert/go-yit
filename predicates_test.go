package yit

import (
	"go.yaml.in/yaml/v3"
	"testing"
)

func TestWithKind(t *testing.T) {
	predicate := WithKind(yaml.ScalarNode)

	if !predicate(&yaml.Node{Kind: yaml.ScalarNode}) {
		t.Errorf("should return true when nodes match the supplied Kind")
	}
	if predicate(&yaml.Node{Kind: yaml.MappingNode}) {
		t.Errorf("should return false when nodes don't match the supplied Kind")
	}
}

func TestUnionAndIntersect(t *testing.T) {
	type testCase struct {
		op       func(...Predicate) Predicate
		a, b     bool
		expected bool
	}

	cases := []testCase{
		{Union, true, true, true},
		{Union, true, false, true},
		{Union, false, true, true},
		{Union, false, false, false},
		{Intersect, true, true, true},
		{Intersect, true, false, false},
		{Intersect, false, true, false},
		{Intersect, false, false, false},
	}

	for _, c := range cases {
		actual := c.op(
			func(node *yaml.Node) bool { return c.a },
			func(node *yaml.Node) bool { return c.b },
		)
		if actual(nil) != c.expected {
			t.Errorf("expected %v, got %v for op", c.expected, actual(nil))
		}
	}
}

func TestWithShortTag(t *testing.T) {
	predicate := WithShortTag("booooo")

	if !predicate(&yaml.Node{Tag: "booooo"}) {
		t.Errorf("should return true when nodes match the tag")
	}
	if predicate(&yaml.Node{Tag: "not boooo"}) {
		t.Errorf("should return false when nodes do not match the tag")
	}
}

func TestWithMapKey(t *testing.T) {
	predicate := WithMapKey("a")

	nodeWithKey := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "a"},
		{Kind: yaml.ScalarNode, Value: "b"},
	}}
	if !predicate(nodeWithKey) {
		t.Errorf("should return true when the map has a specific key")
	}

	nodeWithoutKey := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "c"},
		{Kind: yaml.ScalarNode, Value: "d"},
	}}
	if predicate(nodeWithoutKey) {
		t.Errorf("should return false when the map doesn't have a specific key")
	}

	notAMap := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "a"},
	}}
	if predicate(notAMap) {
		t.Errorf("should return false when the node isn't a map")
	}
}

func TestWithMapValue(t *testing.T) {
	predicate := WithMapValue("b")

	nodeWithValue := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "a"},
		{Kind: yaml.ScalarNode, Value: "b"},
	}}
	if !predicate(nodeWithValue) {
		t.Errorf("should return true when the map has a specific value")
	}

	nodeWithoutValue := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "c"},
		{Kind: yaml.ScalarNode, Value: "d"},
	}}
	if predicate(nodeWithoutValue) {
		t.Errorf("should return false when the map doesn't have a specific value")
	}

	notAMap := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "b"},
	}}
	if predicate(notAMap) {
		t.Errorf("should return false when the node isn't a map")
	}
}

func TestWithMapKeyValue(t *testing.T) {
	predicate := WithMapKeyValue(
		WithStringValue("a"),
		WithStringValue("b"),
	)

	nodeWithPair := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "a"},
		{Kind: yaml.ScalarNode, Value: "b"},
	}}
	if !predicate(nodeWithPair) {
		t.Errorf("should return true when the map has a specific key value pair")
	}

	nodeWithoutPair := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "a"},
		{Kind: yaml.ScalarNode, Value: "c"},
	}}
	if predicate(nodeWithoutPair) {
		t.Errorf("should return false when the map doesn't have a specific key value pair")
	}

	notAMap := &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: "a"},
		{Kind: yaml.ScalarNode, Value: "b"},
	}}
	if predicate(notAMap) {
		t.Errorf("should return false when the node isn't a map")
	}
}

func TestWithPrefix(t *testing.T) {
	predicate := WithPrefix("pre")

	nodeWithPrefix := &yaml.Node{Value: "prefix"}
	if !predicate(nodeWithPrefix) {
		t.Errorf("should return true when the node's value has a prefix")
	}

	nodeWithoutPrefix := &yaml.Node{Value: "postfix"}
	if predicate(nodeWithoutPrefix) {
		t.Errorf("should return false when the node's value does not have a prefix")
	}
}

func TestWithSuffix(t *testing.T) {
	predicate := WithSuffix("fix")

	nodeWithSuffix := &yaml.Node{Value: "prefix"}
	if !predicate(nodeWithSuffix) {
		t.Errorf("should return true when the node's value has a suffix")
	}

	nodeWithoutSuffix := &yaml.Node{Value: "fixpost"}
	if predicate(nodeWithoutSuffix) {
		t.Errorf("should return false when the node's value does not have a suffix")
	}
}

func TestNegate(t *testing.T) {
	if Negate(All)(nil) {
		t.Errorf("Negate(All) should be false")
	}
	if !Negate(None)(nil) {
		t.Errorf("Negate(None) should be true")
	}
}
