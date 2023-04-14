package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
	"xorm.io/core"
)

type MysqlOptions struct {
	HostName           string
	Port               string
	User               string
	Password           string
	DbName             string
	TablePrefix        string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    int
}

type MySQLConnector struct {
	options *MysqlOptions // the database config pointer
	tables  []interface{} //the table set for struct
	Db      *xorm.Engine  //the xorm pointer
}

func NewMySQLConnector(options *MysqlOptions, tables []interface{}) MySQLConnector {
	var connector MySQLConnector
	connector.options = options
	connector.tables = tables
	//setting the database connect url
	url := ""
	if options.HostName == "" || options.HostName == "127.0.0.1" {
		url = fmt.Sprintf(
			"%s:%s@/%s?charset=utf8&parseTime=True",
			options.User, options.Password, options.DbName)
	} else {
		url = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
			options.User, options.Password, options.HostName, options.Port, options.DbName)
	}
	//instance the mysql
	db, err := xorm.NewEngine("mysql", url)
	if err != nil {
		panic(fmt.Errorf("database init failed"))
	}
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, options.TablePrefix)
	db.SetTableMapper(tbMapper)
	db.DB().SetConnMaxLifetime(time.Duration(options.ConnMaxLifetime) * time.Second)
	db.DB().SetMaxOpenConns(options.MaxOpenConnections)
	db.DB().SetMaxIdleConns(options.MaxIdleConnections)
	//enable the mysql log for the terminal
	db.ShowSQL(true)
	if err = db.Ping(); err != nil {
		panic(fmt.Errorf("database conn failed %s", err.Error()))
	}
	connector.Db = db
	//create the tables:the policy is the table is not exist,then create the table
	if err := connector.createTables(); err != nil {
		panic(fmt.Errorf("create table failed %s", err.Error()))
	}
	return connector
}

func (m *MySQLConnector) createTables() error {
	if len(m.tables) == 0 {
		//not data tables need to create
		return nil
	}
	if err := m.Db.CreateTables(m.tables...); err != nil {
		return fmt.Errorf("create mysql table error: %s", err.Error())
	}
	//sync
	if err := m.Db.Sync2(m.tables...); err != nil {
		return fmt.Errorf("sync mysql table error: %s", err.Error())
	}
	return nil
}
