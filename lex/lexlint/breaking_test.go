package lexlint

import (
	"testing"

	"github.com/bluesky-social/indigo/atproto/lexicon"

	"github.com/stretchr/testify/assert"
)

func TestBreakingDefs_InnerTypeChange(t *testing.T) {
	assert := assert.New(t)

	// When a field changes type (e.g. string -> object), breakingDefs should
	// report a type-change error rather than panicking on a type assertion.
	local := lexicon.SchemaDef{Inner: lexicon.SchemaString{}}
	remote := lexicon.SchemaDef{Inner: lexicon.SchemaObject{}}

	issues := breakingDefs("com.example.test", "testField", local, remote)

	assert.Len(issues, 1)
	assert.Equal("type-change", issues[0].LintName)
	assert.Equal("error", issues[0].LintLevel)
}

func TestBreakingDefs_Space(t *testing.T) {
	assert := assert.New(t)

	before := lexicon.SchemaDef{Inner: lexicon.SchemaSpace{Type: "space", Key: "any", Collections: []string{"com.example.a", "com.example.b"}}}

	// adding a collection is not breaking
	added := lexicon.SchemaDef{Inner: lexicon.SchemaSpace{Type: "space", Key: "any", Collections: []string{"com.example.b", "com.example.a", "com.example.c"}}}
	assert.Empty(breakingDefs("com.example.group", "main", before, added))

	changed := lexicon.SchemaDef{Inner: lexicon.SchemaSpace{Type: "space", Key: "tid", Collections: []string{"com.example.b"}}}
	issues := breakingDefs("com.example.group", "main", before, changed)
	names := []string{}
	for _, iss := range issues {
		names = append(names, iss.LintName)
	}
	assert.ElementsMatch([]string{"space-key-type", "space-collection-removed"}, names)
}
