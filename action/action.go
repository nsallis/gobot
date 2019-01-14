package action

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"regexp"
)

var Log *logrus.Logger
var DB *redis.Client

type Action struct {
	ID          string
	Key         string   // string used for matching agains a given user's text
	Action      string   // executed by JS engine
	Children    []string // slice of action ids
	IsDirective bool
}

func New() *Action {
	a := &Action{
		ID:          uuid.New().String(),
		Key:         "%r*",          // defaults to matching everything. % is our escape character. %r denotes regex
		Action:      "return null;", // don't do anything by default
		Children:    []string{},
		IsDirective: false,
	}
	return a
}

func (a *Action) Save() error {
	if DB == nil {
		Log.Error("DB is not set for Action package")
		return errors.New("DB is not set")
	}
	validationErrors := a.Validate()
	if len(validationErrors) > 0 {
		Log.WithFields(logrus.Fields{
			"action": a,
		}).Error("Action is invalid")
		return errors.New("Action is invalid")
	}
	childrenString, err := json.Marshal(a.Children)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Action": a.Action,
		}).Error("failed to marshal action to string")
		return errors.New("failed to marshal action to string")
	}
	err = DB.HMSet(a.GetFieldName(), map[string]interface{}{
		"key":      a.Key,
		"action":   a.Action,
		"children": childrenString,
	}).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func Load(id string) (*Action, error) {
	a := &Action{ID: id}
	result := DB.HGetAll(a.GetFieldName())
	if result.Err() != nil {
		Log.WithFields(logrus.Fields{
			"id": id,
		}).Error("No action with such id")
	}
	a.Key = result.Val()["key"]
	a.Action = result.Val()["action"]
	err := json.Unmarshal([]byte(result.Val()["children"]), &a.Children)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"children": result.Val()["children"],
		}).Error("could not unmarshal children for action")
		return nil, err
	}
	return a, nil
}

func (a *Action) GetFieldName() string {
	return "action:" + a.ID
}

func (a *Action) Validate() map[string]string {
	errs := map[string]string{}
	idMatch, _ := regexp.Compile("^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}")
	if !idMatch.MatchString(a.ID) {
		errs["ID"] = "Invalid"
	}
	if a.IsDirective && len(a.Children) > 1 {
		errs["children"] = "directives may only have one child"
	}
	return errs
}
