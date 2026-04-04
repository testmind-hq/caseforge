// internal/rbt/testdata/callgraph_go/repo/repo.go
package repo

// UserRepo is an interface — V2 name-matching breaks at this boundary.
type UserRepo interface {
	Save(name string)
}

// MySQLRepo is the concrete implementation.
type MySQLRepo struct{}

func (m *MySQLRepo) Save(name string) {}
