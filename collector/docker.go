package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/marema31/kin/cache"
	"github.com/sirupsen/logrus"
)

func listContainers(log *logrus.Entry, cli *client.Client, db *cache.Cache) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Errorf("cannot retrieve list of container: %v", err)
		return err
	}

	log.Debugf("There is currently %d containers running", len(containers))

	ci := make([]cache.ContainerInfo, 0)

	for _, container := range containers {
		if _, ok := container.Labels["kin_name"]; ok {
			log.Debugf("Found %v", container.Names)

			ci = append(ci, cache.ContainerInfo{
				Name:  container.Labels["kin_name"],
				URL:   container.Labels["kin_url"],
				Type:  container.Labels["kin_type"],
				Group: container.Labels["kin_group"],
			})
		}
	}

	err = db.RefreshData(log, ci)
	if err != nil {
		return fmt.Errorf("cannot push test data in cache: %w", err)
	}

	return nil
}

//Run poll regurlarly the docker daemon and fill cache with updated list of container labelled for kin.
func Run(ctx context.Context, log *logrus.Entry, db *cache.Cache) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Errorf("cannot connect to docker daemon: %v", err)
		return err
	}

	ticker := time.NewTicker(10 * time.Second) //nolint: gomnd
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := listContainers(log, cli, db)
			if err != nil {
				return err
			}
		}
	}
}
