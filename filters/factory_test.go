package filters

import (
	"github.com/tmortimer/urlfilter/config"
	"github.com/tmortimer/urlfilter/connectors"
	"testing"
)

//TOM test the TypeOf these filters

func TestCreateFakeFilterSuccess(t *testing.T) {
	config := config.NewConfig()
	filter, err := CreateFilter("fake", config)

	if err != nil {
		t.Fatalf("Creating a Fake filter generated an error: %s", err)
	}

	_, ok := filter.(*Fake)
	if !ok {
		t.Fatalf("A filter other than Fake was created.")
	}
}

func TestCreateRedisFilterSuccess(t *testing.T) {
	config := config.NewConfig()
	filter, err := CreateFilter("redis", config)

	if err != nil {
		t.Fatalf("Creating a Redis filter generated an error: %s", err)
	}

	db, ok := filter.(*DB)
	if !ok {
		t.Fatalf("A filter other than DB was created.")
	}

	_, ok = db.conn.(*connectors.Redis)
	if !ok {
		t.Fatalf("A handler other than Redis was created.")
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
	filter, err := FilterFactory(config)

	if err != nil {
		t.Errorf("Creating a Fake filter generated an error: %s", err)
	}

	db, ok := filter.(*DB)
	if !ok {
		t.Fatalf("A filter other than DB was created.")
	}

	if db.next == nil {
		t.Errorf("The filters were not chained together properly.")
	}

	_, ok = db.next.(*Fake)
	if !ok {
		t.Fatalf("A filter other than Fake was created.")
	}

	// No Secondary filter in fake, so G2G.
}

func TestFilterFactoryFailure(t *testing.T) {
	config := config.NewConfig()
	config.Filters = []string{"redis", "wzzl"}
	_, err := FilterFactory(config)
	if err == nil {
		t.Errorf("Trying to create a filter chain with a filter type that does not exist failed to generate an error.")
	}
}
