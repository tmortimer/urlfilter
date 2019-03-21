package filters

import (
	"fmt"
	"github.com/tmortimer/urlfilter/config"
	"github.com/tmortimer/urlfilter/connectors"
)

// Generage a chain of URL filter caches and then a final url
// store based on the provided config.
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

// Create each filter in the filter chain.
func CreateFilter(name string, config *config.Config) (Filter, error) {
	switch name {
	case "fake":
		return NewFake(), nil
	case "redis":
		return NewDB(connectors.NewRedis(config.Redis)), nil
	case "mysql":
		connector, err := connectors.NewMySQL(config.MySQL)
		if err != nil {
			return nil, err
		}
		return NewDB(connector), nil
	case "redismysqlbloom":
		loader, err := connectors.NewMySQL(config.RedisMySQLBloom.MySQL)
		if err != nil {
			return nil, err
		}
		return NewBloom(connectors.NewRedis(config.RedisMySQLBloom.Redis), loader), nil
	}

	return nil, fmt.Errorf("Unknown filter %s", name)
}
