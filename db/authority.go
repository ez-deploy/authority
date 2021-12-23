package db

// Authority db.
// ATTENTION: Need to add a joint unique index, the triples cannot be repeated.
type Authority struct {
	Identity string `json:"identity,omitempty" db:"identity,VARCHAR(128)"`
	Resource string `json:"resource,omitempty" db:"resource,VARCHAR(128)"`
	Action   string `json:"action,omitempty" db:"action,VARCHAR(128)"`
}

// AuthorityRawModel for sqlm orm framework.
func AuthorityRawModel() interface{} {
	return &Authority{}
}
