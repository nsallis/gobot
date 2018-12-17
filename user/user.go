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

// let's define the user object
type User struct {
	SessionID string
	UserHash  map[string]string
}

type FieldNames struct {
	UserHash string
}

func (u *User) Save() error {
	fmt.Sprintln("Logger: %v", Log)
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
	fieldNames := u.GetFieldNames()
	userHashString, err := json.Marshal(u.UserHash)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"UserHash": u.UserHash,
		}).Error("failed to marshal UserHash")
		return errors.New("failed to marshal userHash to string")
	}
	err = DB.MSet(
		fieldNames.UserHash,
		userHashString).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (u *User) GetFieldNames() FieldNames {
	fieldNames := FieldNames{
		UserHash: fmt.Sprintf("user:%s:userHash", u.SessionID),
	}
	return fieldNames

}

func New() *User {
	u := &User{
		SessionID: uuid.New().String(),
		UserHash:  map[string]string{},
	}
	return u
}

func Load(sessionID string) *User {
	u := &User{SessionID: sessionID}
	result := DB.Get(u.GetFieldNames().UserHash)
	if result.Err() != nil {
		Log.WithFields(logrus.Fields{ // TODO we should check if not found or real error
			"sessionID": sessionID,
		}).Error("No user with such sessionID")
	}
	err := json.Unmarshal([]byte(result.Val()), &u.UserHash)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"UserHash": u.UserHash,
			"error":    err,
		}).Error("could not unmarshal UserHash")
	}
	return u
}
