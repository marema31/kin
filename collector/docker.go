package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/marema31/kin/cache"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func parseLabel(log *logrus.Entry, name string, labels map[string]string) (cache.ContainerInfo, bool) {
	if _, ok := labels["kin_name"]; ok {
		log.Debugf("Found %v", name)

		return cache.ContainerInfo{
			Name:  labels["kin_name"],
			URL:   labels["kin_url"],
			Type:  labels["kin_type"],
			Group: labels["kin_group"],
		}, true
	}

	return cache.ContainerInfo{}, false
}

func listServices(ctx context.Context, log *logrus.Entry, cli *client.Client) (*[]cache.ContainerInfo, error) {
	services, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		log.Errorf("cannot retrieve list of container: %v", err)
		return nil, err
	}

	log.Debugf("There is currently %d containers running", len(services))

	ci := make([]cache.ContainerInfo, 0)

	for _, service := range services {
		if infos, ok := parseLabel(log, service.Spec.Name, service.Spec.Annotations.Labels); ok {
			ci = append(ci, infos)
		}
	}

	return &ci, nil
}

func listContainers(ctx context.Context, log *logrus.Entry, cli *client.Client) (*[]cache.ContainerInfo, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Errorf("cannot retrieve list of container: %v", err)
		return nil, err
	}

	log.Debugf("There is currently %d containers running", len(containers))

	ci := make([]cache.ContainerInfo, 0)

	for _, container := range containers {
		if infos, ok := parseLabel(log, container.Names[0], container.Labels); ok {
			ci = append(ci, infos)
		}
	}

	return &ci, nil
}

func refreshList(ctx context.Context, log *logrus.Entry, cli *client.Client, db *cache.Cache, swarmMode bool) error {
	var (
		ci  *[]cache.ContainerInfo
		err error
	)

	switch {
	case swarmMode:
		ci, err = listServices(ctx, log, cli)
	default:
		ci, err = listContainers(ctx, log, cli)
	}

	if err != nil {
		return err
	}

	err = db.RefreshData(log, *ci)
	if err != nil {
		return fmt.Errorf("cannot push test data in cache: %w", err)
	}

	return nil
}

//Run poll regurlarly the docker daemon and fill cache with updated list of container labelled for kin.
func Run(ctx context.Context, log *logrus.Entry, db *cache.Cache, swarmMode bool) error {
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
			err := refreshList(ctx, log, cli, db, swarmMode)
			if err != nil {
				return err
			}
		}
	}
}
