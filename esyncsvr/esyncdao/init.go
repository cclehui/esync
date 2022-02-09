package esyncdao

import daoongorm "github.com/cclehui/dao-on-gorm"

func InitDao() {
	daoongorm.RegisterModel(&EsyncEventDefaultDao{})
}
