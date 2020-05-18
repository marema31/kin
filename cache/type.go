package cache

import "github.com/hashicorp/go-memdb"

//ContainerInfo list of information to be printed for a container.
type ContainerInfo struct {
	Group string
	Name  string
	URL   string
}

//Cache give access to a in-memory database to allow asynchronous refresh of containers list.
type Cache struct {
	db *memdb.MemDB
}

//New return a new cache instance.
func New() (*Cache, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"container": {
				Name: "container",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return &Cache{db: db}, nil
}
