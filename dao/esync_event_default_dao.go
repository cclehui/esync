package dao

import (
	"context"
	"time"

	daoongorm "github.com/cclehui/dao-on-gorm"
	"github.com/cclehui/esync/esyncsvr"
)

type EsyncEventDefaulDao struct {
	ID           int64     `gorm:"column:id;primaryKey" structs:"id" json:"id"`
	EventDate    int       `gorm:"column:event_date" structs:"event_date" json:"event_date"`
	EventType    string    `gorm:"column:event_type" structs:"event_type" json:"event_type"`
	UniqKey      string    `gorm:"column:uniq_key" structs:"uniq_key" json:"uniq_key"`
	UniqKeyCRC32 int64     `gorm:"column:uniq_key_crc32" structs:"uniq_key_crc32" json:"uniq_key_crc32"`
	EventOption  string    `gorm:"column:event_option" structs:"event_option" json:"event_option"`
	EventData    string    `gorm:"column:event_data" structs:"event_data" json:"event_data"`
	EStatus      int       `gorm:"column:e_status" structs:"e_status" json:"e_status"`
	HandlerInfo  string    `gorm:"column:handler_info" structs:"handler_info" json:"handler_info"`
	CreatedAt    time.Time `gorm:"column:created_at" structs:"created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" structs:"updated_at" json:"updated_at"`

	daoBase *daoongorm.DaoBase
}

func NewEsyncEventDefaulDao(ctx context.Context, myDao *EsyncEventDefaulDao, readOnly bool, options ...daoongorm.Option) (*EsyncEventDefaulDao, error) {
	daoBase, err := daoongorm.NewDaoBase(ctx, myDao, readOnly, options...)

	myDao.daoBase = daoBase

	return myDao, err
}

// 支持事务
func NewEsyncEventDefaulDaoWithTX(ctx context.Context,
	myDao *EsyncEventDefaulDao, tx *daoongorm.DBClient, options ...daoongorm.Option) (*EsyncEventDefaulDao, error) {

	daoBase, err := daoongorm.NewDaoBaseWithTX(ctx, myDao, tx, options...)

	myDao.daoBase = daoBase

	return myDao, err
}

func (myDao *EsyncEventDefaulDao) DBName() string {
	return esyncsvr.GetServer().GetMysqlClient().GetDBClientConfig().DSN.DBName
}

func (myDao *EsyncEventDefaulDao) TableName() string {
	return EventDefaultTableName
}

func (myDao *EsyncEventDefaulDao) DBClient() daoongorm.DBClientInterface {
	return esyncsvr.GetServer().GetMysqlClient()
}

func (myDao *EsyncEventDefaulDao) GetDaoBase() *daoongorm.DaoBase {
	return myDao.daoBase
}

func (myDao *EsyncEventDefaulDao) SetDaoBase(myDaoBase *daoongorm.DaoBase) {
	myDao.daoBase = myDaoBase
}
