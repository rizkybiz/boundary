package schema

import (
	"context"
)

const nilVersion = -1

type migrationState struct {
	// devMigration is true if the database schema that would be applied by
	// InitStore would be from files in the /dev directory which indicates it would
	// not be safe to run in a non dev environment.
	devMigration bool

	// binarySchemaVersion provides the database schema version supported by
	// this binary.
	binarySchemaVersion int

	upMigrations   map[int][]byte
	downMigrations map[int][]byte
}

var migrationStates = make(map[string]migrationState)

// State contains information regarding the current state of a boundary database's schema.
type State struct {
	InitializationStarted bool
	Dirty                 bool
	CurrentSchemaVersion  int
	BinarySchemaVersion   int
}

// State provides the state of the boundary schema contained in the backing database.
func (b *Manager) State(ctx context.Context) (*State, error) {
	dbS := State{
		BinarySchemaVersion: BinarySchemaVersion(b.dialect),
	}
	v, dirty, err := b.driver.Version(ctx)
	if err != nil {
		return nil, err
	}
	if v == nilVersion {
		return &dbS, nil
	}
	dbS.InitializationStarted = true
	dbS.CurrentSchemaVersion = v
	dbS.Dirty = dirty
	return &dbS, nil
}

func DevMigration(dialect string) bool {
	return migrationStates[dialect].devMigration
}

func BinarySchemaVersion(dialect string) int {
	return migrationStates[dialect].binarySchemaVersion
}
