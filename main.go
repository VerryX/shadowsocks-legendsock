package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/shadowsocks-server/shadowsocks-legendsock/core"
)

var (
	ErrLastUse = errors.New("Less than 10 seconds from the last use")
)

var flags struct {
	ListCipher   bool
	DBHost       string
	DBPort       int
	DBUser       string
	DBPass       string
	DBName       string
	NodeRate     float64
	SyncInterval int
}

func purge(instanceList map[int]*Instance, users []User) {
	for _, instance := range instanceList {
		contains := false

		for _, v := range users {
			if instance.Port == v.Port {
				contains = v.TransferEnable > v.Upload+v.Download
				break
			}
		}

		if !contains && instance.Started {
			instance.Stop()
		}
	}
}

func report(instance *Instance, database *Database) error {
	if instance.Bandwidth.Upload != 0 || instance.Bandwidth.Download != 0 {
		if time.Now().Unix()-instance.Bandwidth.Last > 10 {
			err := database.UpdateUserBandwidth(instance)
			if err == nil {
				instance.Bandwidth.Reset()

				return nil
			}

			return err
		}

		return ErrLastUse
	}

	return nil
}

func update(instance *Instance, method, password string) {
	if instance.Method != method || instance.Password != password {
		instance.Method = method
		instance.Password = password

		restart(instance)
	}
}

func restart(instance *Instance) {
	if instance.Started {
		instance.Stop()
	}

	instance.Start()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.BoolVar(&flags.ListCipher, "listcipher", false, "list cipher")
	flag.StringVar(&flags.DBHost, "dbhost", "localhost", "database host")
	flag.IntVar(&flags.DBPort, "dbport", 3306, "database port")
	flag.StringVar(&flags.DBUser, "dbuser", "legendsock", "database user")
	flag.StringVar(&flags.DBPass, "dbpass", "legendsock", "database pass")
	flag.StringVar(&flags.DBName, "dbname", "legendsock", "database name")
	flag.Float64Var(&flags.NodeRate, "noderate", 1, "node rate")
	flag.IntVar(&flags.SyncInterval, "syncinterval", 300, "sync interval")
	flag.Parse()

	if flags.ListCipher {
		for _, v := range core.ListCipher() {
			fmt.Println(v)
		}

		return
	}

	log.Println("Starting shadowsocks-legendsock")
	log.Println("Version: 1.0.2")

	instanceList := make(map[int]*Instance, 65535)
	database := newDatabase(flags.DBHost, flags.DBPort, flags.DBUser, flags.DBPass, flags.DBName)

	log.Println("Started")

	first := true
	for {
		if !first {
			log.Printf("Wait %d seconds for sync users", flags.SyncInterval)
			time.Sleep(time.Duration(flags.SyncInterval) * time.Second)
		} else {
			first = false
		}

		log.Println("Start syncing")

		log.Println("Opening database connection")
		if err := database.Open(); err != nil {
			log.Println(err)

			database.Close()
			continue
		}

		log.Println("Get database users")
		users, err := database.GetUser()
		if err != nil {
			log.Println(err)

			database.Close()
			continue
		}

		log.Println("Purge server users")
		purge(instanceList, users)

		for _, user := range users {
			if instance, ok := instanceList[user.Port]; ok {
				if user.TransferEnable > user.Upload+user.Download {
					update(instance, user.Method, user.Password)
				} else {
					if instance.Started {
						instance.Stop()
					}

					err = report(instance, database)
					if err != nil && err != ErrLastUse {
						log.Println(err)
						continue
					}

					delete(instanceList, user.Port)
				}
			} else if user.TransferEnable > user.Upload+user.Download {
				instance := newInstance(user.ID, user.Port, user.Method, user.Password)

				log.Printf("Starting new instance for %d", user.ID)
				instance.Start()

				instanceList[user.Port] = instance
			}
		}

		log.Println("Reporting user bandwidth")
		ports, err := database.UpdateBandwidth(instanceList)
		if err != nil {
			log.Println(err)
		} else {
			for _, p := range ports {
				instanceList[p].Bandwidth.Reset()
			}
		}

		log.Println("Sync done")
		database.Close()
	}
}
