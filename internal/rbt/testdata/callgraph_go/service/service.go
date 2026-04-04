// internal/rbt/testdata/callgraph_go/service/service.go
package service

import "testapp/repo"

// Process calls repo.Save via the UserRepo interface.
// V2 cannot trace this; V3 (RTA) resolves to MySQLRepo.Save.
func Process(r repo.UserRepo) {
	r.Save("alice")
}
