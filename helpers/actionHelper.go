package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nsallis/gobot/action"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

var (
	Log *logrus.Logger
	DB  *redis.Client
)

func BootActions(pathOverride ...string) error {
	if DB == nil {
		Log.Error("DB must be initialized before booting actions")
		return errors.New("DB must be initialized before booting actions")
	}
	path := "../config/actions.json"
	if len(pathOverride) > 0 {
		path = pathOverride[0]
	}
	fmt.Sprintln("path: %v", path)
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		Log.Error("Could not read actions.json config file")
		return errors.New("Could not read actions.json config file")
	}
	actions := make([]action.Action, 0)
	err = json.Unmarshal(dat, &actions)
	if err != nil {
		Log.Error("Could not marshal actions.json into usable format")
		return errors.New("Could not marshal actions.json into usable format")
	}
	for _, actionFromConfig := range actions {
		actionFromConfig.Save()
	}
	return nil
}
