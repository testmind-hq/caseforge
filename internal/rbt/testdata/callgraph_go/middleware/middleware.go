// internal/rbt/testdata/callgraph_go/middleware/middleware.go
package middleware

import "testapp/deep"

// Process delegates to deep.DoWork.
// Called by handler.CreateUser, making deep.go two hops away from handler.
func Process() { deep.DoWork() }
