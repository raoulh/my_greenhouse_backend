package models

import (
	"sync"
	"time"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	logger "git.raoulh.pw/raoulh/my_greenhouse_backend/log"
	"git.raoulh.pw/raoulh/my_greenhouse_backend/models/orm"

	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/alitto/pond"
	"github.com/sirupsen/logrus"
)

var (
	db      *gorm.DB
	logging *logrus.Entry

	usersLock = make(UsersLock, 0)

	//a pool to start all tasks for refreshing data
	workerPool   *pond.WorkerPool
	quitRefresh  chan interface{}
	wgDone       sync.WaitGroup
	runningTasks sync.Map
)

type UsersLock map[uint]*UserLock

type UserLock struct {
	mu sync.Mutex
}

func (ul UsersLock) Lock(u *User) {
	if lock, ok := ul[u.ID]; ok {
		lock.mu.Lock()
	} else {
		ul[u.ID] = &UserLock{}
		ul[u.ID].mu.Lock()
	}
}

func (ul UsersLock) Unlock(u *User) {
	if _, ok := ul[u.ID]; !ok {
		return //nothing to unlock
	}
	ul[u.ID].mu.Unlock()
}

func init() {
	logging = logger.NewLogger("database")
}

//Init models
func Init() (err error) {
	logging.Infof("Opening database %s", config.Config.String("database.dsn"))
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: config.Config.String("database.dsn"),
	}), &gorm.Config{
		//Logger: logger.NewGorm(), //TODO: our logrus impl does not work
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return
	}

	migrateDb()

	quitRefresh = make(chan interface{})
	workerPool = pond.New(10, 10000, pond.MinWorkers(10), pond.Strategy(pond.Eager()))

	go refreshAllUsers()
	wgDone.Add(1)

	return
}

//Shutdown models
func Shutdown() {
	close(quitRefresh)
	workerPool.StopAndWait()
	wgDone.Wait()
}

func migrateDb() {
	//Migrate all tables
	db.AutoMigrate(
		&User{},
		&UnitMeasurements{},
	)

	logging.Infof("Migration did run successfully")
}

const (
	refreshTime = 5 * time.Minute
)

func refreshAllUsers() {
	defer wgDone.Done()

	for {
		select {
		case <-quitRefresh:
			logging.Debugln("Exit refreshAllUsers goroutine")
			return
		case <-time.After(refreshTime):
			var users []*User
			if err := orm.FindAll(db, &users); err != nil {
				logging.Errorf("failed to FindAll users: %v", err)
			}

			for _, u := range users {
				//Only run the task if it is not already in the task queue
				if _, ok := runningTasks.Load(u.ID); ok {
					continue //this user is already in the pool
				}

				runningTasks.Store(u.ID, true)
				workerPool.Submit(func() {
					refreshUserData(u.ID)

					//remove user from pool
					runningTasks.Delete(u.ID)
				})
			}
		}
	}

}
