package gou

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type aliasedStringSlice []string
type aliasedIntSlice []int
type aliasedFloatSlice []float64

func TestCoerce(t *testing.T) {

	data := map[string]interface{}{
		"int":                4,
		"float":              45.3,
		"string":             "22",
		"stringf":            "22.2",
		"boolt":              true,
		"emptystring":        "",
		"stringslice":        []string{"one", "two", "three"},
		"aliasedstringslice": aliasedStringSlice{"one", "two", "three"},
		"aliasedintslice":    aliasedIntSlice{1, 2, 3},
		"aliasedfloatslice":  aliasedFloatSlice{1.1, 2.2, 3.3},
	}
	assert.True(t, CoerceStringShort(data["int"]) == "4", "get int as string")
	assert.True(t, CoerceStringShort(data["float"]) == "45.3", "get float as string: %v", data["float"])
	assert.True(t, CoerceStringShort(data["string"]) == "22", "get string as string: %v", data["string"])
	assert.True(t, CoerceStringShort(data["stringf"]) == "22.2", "get stringf as string: %v", data["stringf"])
	assert.Equal(t, "true", CoerceStringShort(data["boolt"]))

	assert.True(t, CoerceIntShort(data["int"]) == 4, "get int as int: %v", data["int"])

	assert.EqualValues(t, []int{1, 2, 3}, CoerceInts(data["aliasedintslice"]))

	assert.True(t, CoerceIntShort(data["float"]) == 45, "get float as int: %v", data["float"])
	assert.True(t, CoerceIntShort(data["string"]) == 22, "get string as int: %v", data["string"])
	assert.True(t, CoerceIntShort(data["stringf"]) == 22, "get stringf as int: %v", data["stringf"])

	assert.Equal(t, 0, len(CoerceStrings(data["emptystring"])), "get emptystring as []string: %v", data["emptystring"])
	assert.Equal(t, []string{"22"}, CoerceStrings(data["string"]), "get string as []string: %v", data["string"])
	assert.Equal(t, []string{"4"}, CoerceStrings(data["int"]), "get int as []string: %v", data["int"])

	assert.EqualValues(t, []string{"one", "two", "three"}, CoerceStrings(data["stringslice"]))
	assert.EqualValues(t, []string{"one", "two", "three"}, CoerceStrings(data["aliasedstringslice"]))
	assert.Equal(t, []string{"4"}, CoerceStrings(data["int"]), "get int as []string: %v", data["int"])

	assert.Equal(t, []float64{float64(4)}, CoerceFloats(data["int"]), "get int as []float64: %v", data["int"])
	assert.Equal(t, []float64{float64(45.3)}, CoerceFloats(data["float"]), "get float as []float64: %v", data["float"])
	assert.EqualValues(t, []float64{1.1, 2.2, 3.3}, CoerceFloats(data["aliasedfloatslice"]))
}
