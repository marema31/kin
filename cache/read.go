package cache

import (
	"github.com/sirupsen/logrus"
)

//RetrieveData provides the list of all containers to be printed with and the corresponding metadata.
func (c *Cache) RetrieveData(log *logrus.Entry) ([]ContainerInfo, error) {
	ci := make([]ContainerInfo, 0)

	// Create read-only transaction
	txn := c.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("container", "id")
	if err != nil {
		log.Errorf("Cannot retrieve container list in cache: %v", err)
		return ci, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		c := obj.(ContainerInfo)
		ci = append(ci, c)
	}

	return ci, nil
}
