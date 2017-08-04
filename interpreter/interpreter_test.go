package interpreter

import (
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
	"time"
	. "github.com/ThoughtWorksStudios/bobcat/common"
)

func AssertShouldHaveField(t *testing.T, entity *generator.Generator, field dsl.Node) {
	AssertNotNil(t, entity.GetField(field.Name), "Expected entity to have field %s, but it did not", field.Name)
}

func AssertFieldShouldBeOverriden(t *testing.T, entity *generator.Generator, field dsl.Node) {
	AssertEqual(t, field.Value.(dsl.Node).Value, entity.GetField(field.Name).GenerateValue())
}

var validFields = dsl.NodeSet{
	FieldNode("name", BuiltinNode("string"), IntArgs(10)...),
	FieldNode("age", BuiltinNode("integer"), IntArgs(1, 10)...),
	FieldNode("weight", BuiltinNode("decimal"), FloatArgs(1.0, 200.0)...),
	FieldNode("dob", BuiltinNode("date"), DateArgs("2015-01-01", "2017-01-01")...),
	FieldNode("last_name", BuiltinNode("dict"), StringArgs("last_name")...),
	FieldNode("catch_phrase", StringNode("Grass.... Tastes bad")),
}

var nestedFields = dsl.NodeSet{
	FieldNode("name", BuiltinNode("string"), IntArgs(10)...),
	FieldNode("pet", IdNode("Goat"), IntArgs(2)...),
	FieldNode("friend", EntityNode("Horse", validFields), IntArgs(1)...),
}

var overridenFields = dsl.NodeSet{
	FieldNode("catch_phrase", StringNode("Grass.... Tastes good")),
}

func interp() *Interpreter {
	return New()
}

func TestScopingResolvesOtherEntities(t *testing.T) {
	scope := NewRootScope()
	i := interp()
	node := RootNode(EntityNode("person", dsl.NodeSet{
		FieldNode("pet", EntityNode("kitteh", overridenFields)),
		FieldNode("pets_can_have_pets_too", EntityNode("lolcat", dsl.NodeSet{
			FieldNode("cheezburgrz", StringNode("can has")),
			FieldNode("protoype", IdNode("kitteh")),
		})),
	}))
	err := i.Visit(node, scope)
	AssertNil(t, err, "`lolcat` should be able to resolve `kitteh` because it lives within the scope hierarchy. error was %v", err)

	// using same root scope to simulate previously defined symbols
	err = i.Visit(RootNode(GenerationNode(IdNode("person"), 2)), scope)
	AssertNil(t, err, "Should be able to resolve `person` because it is defined in the root scope. error was %v", err)

	// using same root scope to simulate previously defined symbols; here, `kitteh` was defined in a child scope of `person`,
	// but not at the root scope, so we should not be able to resolve it.
	ExpectsError(t, "Cannot resolve symbol \"kitteh\"", i.Visit(RootNode(GenerationNode(IdNode("kitteh"), 1)), scope))
}

func TestValidVisit(t *testing.T) {
	node := RootNode(EntityNode("person", validFields), GenerationNode(IdNode("person"), 2))
	i := interp()
	scope := NewRootScope()
	err := i.Visit(node, scope)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	for _, entry := range scope.symbols {
		entity := entry.Value.(*generator.Generator)
		for _, field := range validFields {
			AssertShouldHaveField(t, entity, field)
		}
	}
}

func TestValidVisitWithNesting(t *testing.T) {
	node := RootNode(EntityNode("Goat", validFields), EntityNode("person", nestedFields),
		GenerationNode(IdNode("person"), 2))
	i := interp()

	scope := NewRootScope()
	err := i.Visit(node, scope)
	if err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	person, _ := i.ResolveEntity(IdNode("person"), scope)
	for _, field := range nestedFields {
		AssertShouldHaveField(t, person, field)
	}
}

func TestValidVisitWithOverrides(t *testing.T) {
	node := RootNode(
		EntityNode("person", validFields),
		GenerationNode(
			EntityExtensionNode("lazyPerson", "person", overridenFields),
			2,
		),
	)

	i := interp()
	scope := NewRootScope()

	if err := i.Visit(node, scope); err != nil {
		t.Errorf("There was a problem generating entities: %v", err)
	}

	AssertEqual(t, 2, len(scope.symbols), "Should have 2 entities defined")

	for _, key := range []string{"person", "lazyPerson"} {
		_, isPresent := scope.symbols[key]
		// don't try to use AssertNotNil here; it won't work because it is unable to detect
		// whether a nil pointer passed as an interface{} param to AssertNotEqual is nil.
		// see this crazy shit: https://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
		Assert(t, isPresent, "`%v` should be defined in scope", key)

		if isPresent {
			entity, isGeneratorType := scope.symbols[key].Value.(*generator.Generator)
			Assert(t, isGeneratorType, "`key` should be defined")

			if entity.Name != "person" {
				for _, field := range overridenFields {
					AssertFieldShouldBeOverriden(t, entity, field)
				}
			}
		}
	}
}

func TestValidateFieldAmountWithNoArguments(t *testing.T) {
	i := interp()
	amount := dsl.NodeSet{}

	actual, _ := i.validateFieldAmount(amount)
	expected := Amount{1, 1}

	AssertEqual(t, expected, actual)
}

func TestValidateFieldAmountWithOneValidArgument(t *testing.T) {
	i := interp()
	amount := IntArgs(3)

	actual, _ := i.validateFieldAmount(amount)
	expected := Amount{3, 3}

	AssertEqual(t, expected, actual)
}

func TestValidateFieldAmountWithOneInvalidArgument(t *testing.T) {
	i := interp()
	amount := StringArgs("nope")

	_, err := i.validateFieldAmount(amount)

	ExpectsError(t, "Expected nope to be an integer, but was string.", err)
}

func TestValidateFieldAmountWithTwoValidArguments(t *testing.T) {
	i := interp()
	amount := IntArgs(1, 3)

	actual, _ := i.validateFieldAmount(amount)
	expected := Amount{1, 3}

	AssertEqual(t, expected, actual)
}

func TestValidateFieldAmountWithTwoInvalidArguments(t *testing.T) {
	i := interp()
	amount := StringArgs("nope", "yup")

	_, err := i.validateFieldAmount(amount)

	ExpectsError(t, "Expected nope to be an integer, but was string.", err)
}

func TestValidateFieldAmountWithMaxLargerThanMin(t *testing.T) {
	i := interp()
	amount := IntArgs(3, 1)

	_, err := i.validateFieldAmount(amount)

	ExpectsError(t, "Max '1' cannot be less than min '3'", err)
}

func TestValidateFieldAmountWithTooManyValidArguments(t *testing.T) {
	i := interp()
	amount := IntArgs(1,2,3)

	_, err := i.validateFieldAmount(amount)

	ExpectsError(t, "Field amount must be one or two values only", err)
}

func TestInvalidGenerationNodeBadArgType(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	i.EntityFromNode(EntityNode("burp", validFields), scope)
	node := dsl.Node{Kind: "generation", Value: IdNode("burp"), Args: StringArgs("blah")}
	ExpectsError(t, `generate "burp" takes an integer count`, i.GenerateFromNode(node, scope))
}

func TestInvalidGenerationNodeBadCountArg(t *testing.T) {
	i := interp()
	scope := NewRootScope()
	i.EntityFromNode(EntityNode("person", validFields), scope)
	node := GenerationNode(IdNode("person"), 0)
	ExpectsError(t, "Must generate at least 1 `person` entity", i.GenerateFromNode(node, scope))
}

func TestEntityWithUndefinedParent(t *testing.T) {
	ent := EntityNode("person", validFields)
	unresolvable := IdNode("nope")
	ent.Related = &unresolvable
	_, err := interp().EntityFromNode(ent, NewRootScope())
	ExpectsError(t, `Cannot resolve parent entity "nope" for entity "person"`, err)
}

func TestGenerateEntitiesCannotResolveEntity(t *testing.T) {
	node := GenerationNode(IdNode("tree"), 2)
	ExpectsError(t, `Cannot resolve symbol "tree"`, interp().GenerateFromNode(node, NewRootScope()))
}

func TestDefaultArguments(t *testing.T) {
	i := interp()
	defaults := map[string]interface{}{
		"string":  5,
		"integer": [2]int{1, 10},
		"decimal": [2]float64{1, 10},
		"date":    [2]time.Time{UNIX_EPOCH, NOW},
	}

	for kind, expected_value := range defaults {
		actual, _ := i.defaultArgumentFor(kind)
		if actual != expected_value {
			t.Errorf("default value for argument type '%s' was expected to be %v but was %v", kind, expected_value, actual)
		}
	}
}

func TestDisallowNondeclaredEntityAsFieldIdentifier(t *testing.T) {
	i := interp()
	_, e := i.EntityFromNode(EntityNode("hiccup", nestedFields), NewRootScope())
	ExpectsError(t, `Cannot resolve symbol "Goat"`, e)

}

func TestDefaultArgumentsReturnsErrorOnUnsupportedFieldType(t *testing.T) {
	i := interp()
	arg, err := i.defaultArgumentFor("dict")
	if err == nil || err.Error() != "Field of type `dict` requires arguments" {
		t.Errorf("expected an error when getting a default Argument for an unsupported field Type")
	}
	AssertNil(t, arg, "defaultArgumentFor(\"dict\") Should not have returned anything")
}

func TestConfiguringFieldDiesWhenFieldWithoutArgsHasNoDefaults(t *testing.T) {
	i := interp()

	badNode := FieldNode("name", BuiltinNode("dict"))
	entity := generator.NewGenerator("cat", GetLogger(t))
	ExpectsError(t, "Field of type `dict` requires arguments", i.withDynamicField(entity, badNode, NewRootScope()))
}

func TestConfiguringFieldWithoutArguments(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	fieldNoArgs := FieldNode("last_name", BuiltinNode("string"))
	i.withDynamicField(testEntity, fieldNoArgs, NewRootScope())
	AssertShouldHaveField(t, testEntity, fieldNoArgs)
}

func TestConfiguringFieldsForEntityErrors(t *testing.T) {
	i := interp()
	testEntity := generator.NewGenerator("person", GetLogger(t))
	badNode := FieldNode("last_name", BuiltinNode("dict"), IntArgs(1, 10)...)
	ExpectsError(t, "Field type `dict` expected 1 args, but 2 found.", i.withDynamicField(testEntity, badNode, NewRootScope()))
}

func TestValInt(t *testing.T) {
	expected := 666
	actual := valInt(IntArgs(666)[0])
	AssertEqual(t, expected, actual)
}

func TestValStr(t *testing.T) {
	expected := "blah"
	actual := valStr(StringArgs("blah")[0])
	AssertEqual(t, expected, actual)
}

func TestValFloat(t *testing.T) {
	expected := 4.2
	actual := valFloat(FloatArgs(4.2)[0])
	AssertEqual(t, expected, actual)
}

func TestValTime(t *testing.T) {
	expected, _ := time.Parse("2006-01-02", "1945-01-01")
	actual := valTime(DateArgs("1945-01-01")[0])
	AssertEqual(t, expected, actual)
}
