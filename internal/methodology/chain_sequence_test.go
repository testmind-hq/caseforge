package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func TestScoreFieldSimilarity_SameToken(t *testing.T) {
	assert.Greater(t, scoreFieldSimilarity("userId", "id"), 0.0)
}

func TestScoreFieldSimilarity_NoOverlap(t *testing.T) {
	assert.Equal(t, 0.0, scoreFieldSimilarity("price", "name"))
}

func TestTokenizeFieldName_CamelCase(t *testing.T) {
	tokens := tokenizeFieldName("orderId")
	assert.Equal(t, []string{"order", "id"}, tokens)
}

func TestTokenizeFieldName_SnakeCase(t *testing.T) {
	tokens := tokenizeFieldName("order_item_id")
	assert.Equal(t, []string{"order", "item", "id"}, tokens)
}

func makeNonCRUDSpec() *spec.ParsedSpec {
	return &spec.ParsedSpec{
		Operations: []*spec.Operation{
			{
				Method: "POST", Path: "/carts",
				Responses: map[string]*spec.Response{
					"201": {Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"cartId": {Type: "string"},
								"total":  {Type: "number"},
							},
						}},
					}},
				},
			},
			{
				Method: "POST", Path: "/carts/{cartId}/items",
				Parameters: []*spec.Parameter{{Name: "cartId", In: "path", Required: true}},
				Responses:  map[string]*spec.Response{"201": {Description: "item added"}},
			},
		},
	}
}

func TestChainSequenceTechnique_DetectsNonCRUDChain(t *testing.T) {
	s := makeNonCRUDSpec()
	tech := NewChainSequenceTechnique()
	cases, err := tech.Generate(s)
	require.NoError(t, err)
	assert.Greater(t, len(cases), 0, "expected at least one chain_sequence case")
	if len(cases) > 0 {
		assert.Equal(t, "chain", cases[0].Kind)
		assert.Len(t, cases[0].Steps, 2)
	}
}

func TestChainSequenceTechnique_Name(t *testing.T) {
	assert.Equal(t, "chain_sequence", NewChainSequenceTechnique().Name())
}
