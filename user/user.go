package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var DB *redis.Client
var RootNodeID string

// let's define the user object
type User struct {
	SessionID     string
	UserHash      map[string]string
	CurrentAction string // action id
}

func New() *User {
	u := &User{
		SessionID:     uuid.New().String(),
		UserHash:      map[string]string{},
		CurrentAction: RootNodeID,
	}
	return u
}

func (u *User) Save() error {
	if DB == nil {
		Log.Error("DB is not set of user package")
		return errors.New("DB is not set")
	}
	if u.SessionID == "" {
		Log.WithFields(logrus.Fields{
			"user": u,
		}).Error("User not valid")
		return errors.New("user does not have a sessionID")
	}
	userHashString, err := json.Marshal(u.UserHash)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"UserHash": u.UserHash,
		}).Error("failed to marshal UserHash")
		return errors.New("failed to marshal userHash to string")
	}
	err = DB.HMSet(u.GetFieldName(), map[string]interface{}{
		"userHash":      userHashString,
		"currentAction": u.CurrentAction,
	}).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func Load(sessionID string) *User {
	u := &User{SessionID: sessionID}
	result := DB.HGetAll(u.GetFieldName())
	if result.Err() != nil {
		Log.WithFields(logrus.Fields{ // TODO we should check if not found or real error
			"sessionID": sessionID,
		}).Error("No user with such sessionID")
	}
	err := json.Unmarshal([]byte(result.Val()["userHash"]), &u.UserHash)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"UserHash": u.UserHash,
			"error":    err,
		}).Error("could not unmarshal UserHash")
	}
	u.CurrentAction = result.Val()["currentAction"]
	if u.CurrentAction == "" {
		Log.Error("currentAction is empty upon loading user. Assigning root action")
		u.CurrentAction = RootNodeID
	}
	return u
}

func (u *User) GetFieldName() string {
	return fmt.Sprintf("user:%s", u.SessionID)
}
