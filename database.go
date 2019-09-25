package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database struct
type Database struct {
	Connection *sql.DB
	DBHost     string
	DBPort     int
	DBUser     string
	DBPass     string
	DBName     string
}

// User struct
type User struct {
	ID             int
	Upload         uint64
	Download       uint64
	Port           int
	Method         string
	Password       string
	Enable         int
	TransferEnable uint64
}

// Open database connection
func (database *Database) Open() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=10s", database.DBUser, database.DBPass, database.DBHost, database.DBPort, database.DBName))
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(4)

	database.Connection = db
	return nil
}

// Close database connection
func (database *Database) Close() error {
	if database.Connection != nil {
		return database.Connection.Close()
	}

	return nil
}

// GetUser RT.
func (database *Database) GetUser() ([]User, error) {
	results, err := database.Connection.Query("SELECT pid, u, d, port, method, passwd, enable, transfer_enable FROM user WHERE enable=1")
	if err != nil {
		return nil, err
	}

	users := make([]User, 65535)
	count := 0
	for results.Next() {
		var user User

		err = results.Scan(&user.ID, &user.Upload, &user.Download, &user.Port, &user.Method, &user.Password, &user.Enable, &user.TransferEnable)
		if err != nil {
			return nil, err
		}

		users[count] = user
		count++
	}

	return users[:count], nil
}

// UpdateBandwidth R.T.
func (database *Database) UpdateBandwidth(list map[int]*Instance) ([]int, error) {
	when1 := ""
	when2 := ""
	in := ""

	users := make([]int, 65535)
	count := 0
	for _, instance := range list {
		if (instance.Bandwidth.Upload != 0 || instance.Bandwidth.Download != 0) && time.Now().Unix()-instance.Bandwidth.Last > 10 {
			users[count] = instance.Port
			count++

			when1 += fmt.Sprintf(" WHEN %d THEN u+%d", instance.Port, uint64(float64(instance.Bandwidth.Upload)*flags.NodeRate))
			when2 += fmt.Sprintf(" WHEN %d THEN d+%d", instance.Port, uint64(float64(instance.Bandwidth.Download)*flags.NodeRate))

			if in == "" {
				in = fmt.Sprintf("%d", instance.Port)
			} else {
				in += fmt.Sprintf(", %d", instance.Port)
			}
		}
	}

	if when1 == "" {
		return nil, nil
	}

	_, err := database.Connection.Query(fmt.Sprintf("UPDATE user SET u = CASE port %s END, d = CASE port %s END, t = %d WHERE port IN (%s)", when1, when2, time.Now().Unix(), in))
	return users[:count], err
}

// UpdateUserBandwidth RT.
func (database *Database) UpdateUserBandwidth(instance *Instance) error {
	log.Printf("Reporting %d uploaded %d downloaded %d to database", instance.UserID, instance.Bandwidth.Upload, instance.Bandwidth.Download)

	upload := uint64(float64(instance.Bandwidth.Upload) * flags.NodeRate)
	download := uint64(float64(instance.Bandwidth.Download) * flags.NodeRate)

	// _, err := database.Connection.Query(fmt.Sprintf("UPDATE chart SET upload=concat(upload, \",%d\"), download=concat(download, \",%d\"), date=%d WHERE pid=%d", upload, download, time.Now().Unix(), instance.UserID))
	// if err != nil {
	// 	return err
	// }

	_, err := database.Connection.Query(fmt.Sprintf("UPDATE user SET u=u+%d, d=d+%d, t=%d WHERE pid=%d", upload, download, time.Now().Unix(), instance.UserID))
	return err
}

func newDatabase(host string, port int, user, pass, name string) *Database {
	database := Database{}
	database.DBHost = host
	database.DBPort = port
	database.DBUser = user
	database.DBPass = pass
	database.DBName = name

	return &database
}
