package filters

import (
	"github.com/tmortimer/urlfilter/config"
	"testing"
)

//TOM test the TypeOf these filters

func TestCreateFilterSuccess(t *testing.T) {
	config := config.NewConfig()
	_, err := CreateFilter("fake", config)
	if err != nil {
		t.Errorf("Creating a Fake filter generated an error: %s", err)
	}

	_, err = CreateFilter("redis", config)
	if err != nil {
		t.Errorf("Creating a Redis filter generated an error: %s", err)
	}
}

func TestCreateFilterFailure(t *testing.T) {
	config := config.NewConfig()
	_, err := CreateFilter("wzzl", config)
	if err == nil {
		t.Errorf("Trying to create a filter type that does not exist failed to generate an error.")
	}
}

func TestFilterFactorySuccess(t *testing.T) {
	config := config.NewConfig()
	config.Filters = []string{"redis", "fake"}
	_, err := FilterFactory(config)
	if err != nil {
		t.Errorf("Creating a Fake filter generated an error: %s", err)
	}

	//TOM need to do some reflection shenanigans here too
	/*if filter.next == nil {
		t.Errorf("The filters were not chained together properly.")
	}

	if filter.next.next != nil {
		t.Errorf("There are extra filters in the chain.")
	}*/
}

func TestFilterFactoryFailure(t *testing.T) {
	config := config.NewConfig()
	config.Filters = []string{"redis", "wzzl"}
	_, err := FilterFactory(config)
	if err == nil {
		t.Errorf("Trying to create a filter chain with a filter type that does not exist failed to generate an error.")
	}
}
