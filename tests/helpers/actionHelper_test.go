package helpers

import (
	"github.com/go-redis/redis"
	"github.com/nsallis/gobot/action"
	"github.com/nsallis/gobot/helpers"
	"github.com/sirupsen/logrus"
	"os"
	// "reflect"
	"testing"
)

// Set up our app for integration tests
var log *logrus.Logger

// initialize packages
func Setup() {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	log = logrus.New()
	log.Out = os.Stdout
	action.DB = db
	action.Log = log
	helpers.DB = db
	helpers.Log = log
}

func TestActionBoot(t *testing.T) {
	Setup()
	helpers.BootActions(t, "./testActions.json")
	bootedAction, err := action.Load("00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Errorf(
			"Failed to find an action in the test actions config %v", err,
		)
	}
	if bootedAction.Key != "%r*" {
		t.Errorf(
			"Test action has an unexpected Key value, %v", bootedAction.Key,
		)
	}
	if !existsInSlice(bootedAction.Children, "asd") {
		t.Errorf(
			"Test action has an unexpected Child value, %v", bootedAction,
		)
	}
	if !existsInSlice(bootedAction.Children, "sdf") {
		t.Errorf(
			"Test action has an unexpected Child value, %v", bootedAction,
		)
	}
}

func existsInSlice(slice []string, key string) bool {
	for _, v := range slice {
		if v == key {
			return true
		}
	}
	return false
}
