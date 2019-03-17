package filters

import (
	"fmt"
	"github.com/tmortimer/urlfilter/config"
	"github.com/tmortimer/urlfilter/connectors"
)

func FilterFactory(config *config.Config) (Filter, error) {
	list := config.Filters
	var filter Filter = nil

	for i := (len(list) - 1); i >= 0; i-- {
		current, err := CreateFilter(list[i], config)

		if err != nil {
			return nil, err
		}

		current.AddSecondaryFilter(filter)
		filter = current
	}

	return filter, nil
}

func CreateFilter(name string, config *config.Config) (Filter, error) {
	switch name {
	case "fake":
		return NewFake(), nil
	case "redis":
		return NewDB(connectors.NewRedis(config.Redis)), nil
	}

	return nil, fmt.Errorf("Unknown filter %s", name)
}
