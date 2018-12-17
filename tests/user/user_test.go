package user

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/icrowley/fake"
	"github.com/nsallis/gobot/user"
	"github.com/sirupsen/logrus"
	"os"
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
	user.Log = log
	user.DB = db
	user.RootNodeID = "99999999999999999999999999999999"
}

func TestUserSetup(t *testing.T) {
	Setup()
	u := user.New()
	email := "nsallis@example.com"
	u.UserHash["email"] = email
	fmt.Println("User sessionID: " + u.SessionID)
	u.Save()
	fmt.Println("User saved")
	user := user.Load(u.SessionID)
	if user.UserHash == nil {
		t.Errorf(
			"user's userHash was nil after unmarshal. UserHash: %v",
			user.UserHash,
		)
	}
	if user.UserHash["email"] != email {
		t.Errorf(
			"user's email expected to be %s, got %s instead",
			email,
			user.UserHash["email"],
		)
	}
	fmt.Printf("Loaded user: %v", user)
}

// benchmarking /////////////////////////////////////////

func createFakeUserHash(fieldCount int) map[string]string {
	fields := map[string]string{}
	for i := 0; i < fieldCount; i++ {

		fields[fake.Word()] = fake.Word()
	}
	return fields
}

func BenchmarkNewUser(benchmark *testing.B) {
	Setup()
	for n := 0; n < benchmark.N; n++ {
		user.New()
	}
}

func BenchmarkSaveUserMinimal(benchmark *testing.B) {
	Setup()
	u := user.New()
	u.CurrentAction = "d41d8cd98f00b204e9800998ecf8427e"
	u.UserHash["email"] = "nsallis@example.com"
	for n := 0; n < benchmark.N; n++ {
		benchmark.StartTimer()
		u.Save()
		benchmark.StopTimer()
	}
}

func BenchmarkSaveUserLarge(benchmark *testing.B) {
	Setup()
	u := user.New()

	for n := 0; n < benchmark.N; n++ {
		fakeUserHash := createFakeUserHash(10)
		u.UserHash = fakeUserHash
		benchmark.StartTimer()
		u.Save()
		benchmark.StopTimer()
	}
}

func BenchmarkLoadUserSmall(benchmark *testing.B) {
	Setup()
	u := user.New()
	u.Save()
	sessionID := u.SessionID
	for n := 0; n < benchmark.N; n++ {
		benchmark.StartTimer()
		user.Load(sessionID)
		benchmark.StopTimer()
	}
}
