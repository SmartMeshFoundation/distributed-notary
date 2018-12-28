package models

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //for gorm
	_ "github.com/jinzhu/gorm/dialects/postgres" //for gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite"   //for gorm
)

// DB :
type DB struct {
	*gorm.DB
}

//SetUpDB init db
func SetUpDB(dbtype, path string) (mdb *DB) {
	var err error
	db, err := gorm.Open(dbtype, path)
	if err != nil {
		panic("failed to connect database")
	}
	if false {
		db = db.Debug()
		db.LogMode(true)
	}
	//db.SetLogger(gorm.Logger{revel.TRACE})
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	db.AutoMigrate(&PrivateKeyInfoModel{})
	db.AutoMigrate(&SignMessgeModel{})
	db.AutoMigrate(&NotaryInfo{})
	//db.Model(&ChannelParticipantInfo{}).AddForeignKey("channel_id", "channels(channel_id)", "CASCADE", "CASCADE") // Foreign key need to define manually
	//db.AutoMigrate(&SettledChannel{})
	//db.AutoMigrate(&latestBlockNumber{})
	//db.AutoMigrate(&tokenNetwork{})
	//db.AutoMigrate(&AccountFee{}, &AccountTokenFee{}, &TokenFee{}, &NodeStatus{})
	//
	//db.FirstOrCreate(lb)
	return &DB{db}
}

//CloseDB release connection
func CloseDB(db *DB) {
	err := db.Close()
	if err != nil {
		log.Printf(fmt.Sprintf("closedb err %s", err))
	}
}

//SetupTestDB for test only
func SetupTestDB() (db *DB) {
	return SetupTestDB2("test.db")
}

// SetupTestDB2 :
func SetupTestDB2(name string) (db *DB) {
	dbPath := filepath.Join(os.TempDir(), name)
	err := os.Remove(dbPath)
	if err != nil {
		//ignore
	}
	return SetUpDB("sqlite3", dbPath)
}
