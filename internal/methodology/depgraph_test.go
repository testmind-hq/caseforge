// internal/methodology/depgraph_test.go
package methodology

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/spec"
)

func crudOps() []*spec.Operation {
	return []*spec.Operation{
		{
			OperationID: "createUser",
			Method:      "POST",
			Path:        "/users",
			Responses: map[string]*spec.Response{
				"201": {
					Headers: map[string]string{},
					Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"id":    {Type: "integer"},
								"email": {Type: "string"},
							},
						}},
					},
				},
			},
		},
		{
			OperationID: "getUser",
			Method:      "GET",
			Path:        "/users/{userId}",
			Parameters:  []*spec.Parameter{{Name: "userId", In: "path", Required: true}},
			Responses:   map[string]*spec.Response{"200": {}},
		},
		{
			OperationID: "updateUser",
			Method:      "PUT",
			Path:        "/users/{userId}",
			Parameters:  []*spec.Parameter{{Name: "userId", In: "path", Required: true}},
			Responses:   map[string]*spec.Response{"200": {}},
		},
		{
			OperationID: "deleteUser",
			Method:      "DELETE",
			Path:        "/users/{userId}",
			Parameters:  []*spec.Parameter{{Name: "userId", In: "path", Required: true}},
			Responses:   map[string]*spec.Response{"204": {}},
		},
	}
}

func TestBuildDepGraph_FindsCreatorConsumerEdges(t *testing.T) {
	g := BuildDepGraph(crudOps())
	require.NotNil(t, g)
	assert.GreaterOrEqual(t, len(g.Edges), 3, "expect GET, PUT, DELETE all linked to POST /users")

	var consumers []string
	for _, e := range g.Edges {
		assert.Equal(t, "POST", e.Creator.Method)
		assert.Equal(t, "/users", e.Creator.Path)
		consumers = append(consumers, e.Consumer.Method)
		assert.Equal(t, "userId", e.PathParam)
		assert.Equal(t, "jsonpath $.id", e.CaptureFrom)
	}
	assert.ElementsMatch(t, []string{"GET", "PUT", "DELETE"}, consumers)
}

func TestBuildDepGraph_LocationHeaderCapture(t *testing.T) {
	ops := []*spec.Operation{
		{
			OperationID: "createItem",
			Method:      "POST",
			Path:        "/items",
			Responses: map[string]*spec.Response{
				"201": {Headers: map[string]string{"Location": "string"}},
			},
		},
		{
			OperationID: "getItem",
			Method:      "GET",
			Path:        "/items/{itemId}",
			Responses:   map[string]*spec.Response{"200": {}},
		},
	}
	g := BuildDepGraph(ops)
	require.Len(t, g.Edges, 1)
	assert.Equal(t, "header Location", g.Edges[0].CaptureFrom)
}

func TestBuildDepGraph_NoEdgesWithoutPost(t *testing.T) {
	ops := []*spec.Operation{
		{Method: "GET", Path: "/users", Responses: map[string]*spec.Response{"200": {}}},
		{Method: "GET", Path: "/users/{id}", Responses: map[string]*spec.Response{"200": {}}},
	}
	g := BuildDepGraph(ops)
	assert.Empty(t, g.Edges)
}

func TestBuildDepGraph_NestedIDPath(t *testing.T) {
	ops := []*spec.Operation{
		{
			OperationID: "createOrder",
			Method:      "POST",
			Path:        "/orders",
			Responses: map[string]*spec.Response{
				"201": {
					Headers: map[string]string{},
					Content: map[string]*spec.MediaType{
						"application/json": {Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"data": {
									Type: "object",
									Properties: map[string]*spec.Schema{
										"id":     {Type: "string"},
										"status": {Type: "string"},
									},
								},
							},
						}},
					},
				},
			},
		},
		{
			OperationID: "getOrder",
			Method:      "GET",
			Path:        "/orders/{orderId}",
			Responses:   map[string]*spec.Response{"200": {}},
		},
	}
	g := BuildDepGraph(ops)
	require.Len(t, g.Edges, 1)
	assert.Equal(t, "jsonpath $.data.id", g.Edges[0].CaptureFrom)
	assert.Equal(t, "data.id", g.Edges[0].IDField, "IDField must match CaptureFrom nested path")
}

func TestBuildDepGraph_Deterministic(t *testing.T) {
	ops := crudOps()
	g1 := BuildDepGraph(ops)
	g2 := BuildDepGraph(ops)
	require.Equal(t, len(g1.Edges), len(g2.Edges))
	for i := range g1.Edges {
		assert.Equal(t, g1.Edges[i].Consumer.OperationID, g2.Edges[i].Consumer.OperationID)
	}
}
