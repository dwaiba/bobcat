package dsl

import (
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	"testing"
	"time"
)

var cnt = &current{
	pos:         position{line: 4, col: 3, offset: 42},
	text:        []byte("wubba lubba dub dub!!!!"),
	globalStore: map[string]interface{}{"filename": "whatever.spec"},
}

var location = NewLocation("whatever.spec", 4, 3, 42)

func TestRootNodeReturnsExpectedNode(t *testing.T) {
	node1 := Node{}
	node2 := Node{}
	kids := NodeSet{node1, node2}
	expected := Node{Kind: "root", Children: kids, Ref: location}
	var statements interface{} = []interface{}{node1, node2}
	actual, err := rootNode(cnt, statements)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestEntityNodeReturnsExpectedNode(t *testing.T) {
	beth := Node{Kind: "field", Name: "beth"}
	morty := Node{Kind: "argument", Name: "morty"}
	kids := NodeSet{beth, morty}
	expected := Node{Kind: "definition", Name: "Rick", Children: kids, Ref: location}
	actual, err := entityNode(cnt, Node{Value: "Rick"}, kids)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestGenNodeReturnsExpectedNodeWithArgs(t *testing.T) {
	summer := Node{Kind: "field", Name: "summer"}
	morty := Node{Kind: "argument", Name: "morty"}
	kids := NodeSet{summer, morty}
	expected := Node{Kind: "generation", Name: "Beth", Args: kids, Ref: location, Children: NodeSet{}}
	actual, err := genNode(cnt, Node{Value: "Beth"}, nil, kids)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestGenNodeReturnsExpectedNodeWithoutArgs(t *testing.T) {
	expected := Node{Kind: "generation", Name: "Beth", Args: NodeSet{}, Ref: location, Children: NodeSet{}}
	actual, err := genNode(cnt, Node{Value: "Beth"}, nil, nil)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestStaticFieldNode(t *testing.T) {
	morty := Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := Node{Kind: "field", Ref: location, Name: "Rick", Value: morty}
	actual, err := staticFieldNode(cnt, Node{Value: "Rick"}, morty)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithoutArgs(t *testing.T) {
	morty := Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	expected := Node{Kind: "field", Ref: location, Name: "Rick", Value: morty, Args: NodeSet{}}
	actual, err := dynamicFieldNode(cnt, Node{Value: "Rick"}, morty, nil)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDynamicNodeWithArgs(t *testing.T) {
	morty := Node{Kind: "builtin", Name: "grandson", Value: "morty"}
	args := NodeSet{Node{}}
	expected := Node{Kind: "field", Ref: location, Name: "Rick", Value: morty, Args: args}
	actual, err := dynamicFieldNode(cnt, Node{Value: "Rick"}, morty, args)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestIDNode(t *testing.T) {
	expected := Node{Kind: "identifier", Ref: location, Value: "wubba lubba dub dub!!!!"}
	actual, err := idNode(cnt)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBuiltinNode(t *testing.T) {
	expected := Node{Kind: "builtin", Ref: location, Value: "wubba lubba dub dub!!!!"}
	actual, err := builtinNode(cnt)

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDateLiteralNode(t *testing.T) {
	fullDate, _ := time.Parse("2006-01-02", "2017-07-19")
	expected := Node{Kind: "literal-date", Ref: location, Value: fullDate}
	date := "2017-07-19"
	actual, err := dateLiteralNode(cnt, date, []string{})

	AssertNil(t, err, "Got an error constructing root node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestDateLiteralNodeReturnsError(t *testing.T) {
	_, err := dateLiteralNode(cnt, "2017-07-19", []string{"13:00:00-0700"})

	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

func TestIntLiteralNode(t *testing.T) {
	expected := Node{Kind: "literal-int", Ref: location, Value: 5}
	actual, err := intLiteralNode(cnt, "5")

	AssertNil(t, err, "Got an error constructing int literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestIntLiteralNodeError(t *testing.T) {
	_, err := intLiteralNode(cnt, string(5))
	AssertNotNil(t, err, "Expected an error, but got none")
}

func TestFloatLiteralNode(t *testing.T) {
	expected := Node{Kind: "literal-float", Ref: location, Value: float64(5)}
	actual, err := floatLiteralNode(cnt, "5")

	AssertNil(t, err, "Got an error constructing float literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestNullLiteralNode(t *testing.T) {
	expected := Node{Kind: "literal-null", Ref: location}
	actual, err := nullLiteralNode(cnt)
	AssertNil(t, err, "Got an error constructing null literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBoolLiteralNode(t *testing.T) {
	expected := Node{Kind: "literal-bool", Ref: location, Value: true}
	c := *cnt
	c.text = []byte("true")
	actual, err := boolLiteralNode(&c)
	AssertNil(t, err, "Got an error constructing bool literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestBoolLiteralNodeReturnsError(t *testing.T) {
	c := *cnt
	c.text = []byte("eek")
	_, err := boolLiteralNode(&c)
	AssertNotNil(t, err, "Expected Error, but didn't get one")
}

func TestStrLiteralNode(t *testing.T) {
	expected := Node{Kind: "literal-string", Ref: location, Value: "v"}
	c := *cnt
	c.text = []byte("\"v\"")
	actual, err := strLiteralNode(&c)

	AssertNil(t, err, "Got an error constructing string literal node: %v", err)
	AssertEqual(t, expected.String(), actual.String())
}

func TestStrLiteralNodeReturnsError(t *testing.T) {
	_, err := strLiteralNode(cnt)

	AssertNotNil(t, err, "Expected an error but didn't get one")
}
