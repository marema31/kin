package cache

import "github.com/sirupsen/logrus"

//RefreshData replace the list of all containers in cache by the ones to be printed with and the corresponding metadata.
func (c *Cache) RefreshData(log *logrus.Entry, containers []ContainerInfo) error {
	// Create write transaction
	txn := c.db.Txn(true)

	if _, err := txn.DeleteAll("container", "id"); err != nil {
		log.Errorf("Cannot empty container list in cache: %v", err)
		return err
	}

	for _, container := range containers {
		if err := txn.Insert("container", container); err != nil {
			log.Errorf("Cannot push container list in cache: %v", err)
			return err
		}
	}

	txn.Commit()

	return nil
}
