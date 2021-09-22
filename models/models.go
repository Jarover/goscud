//
// Package models - описание структуры данных внешних таблиц
//
package models

type FbEnv struct {
	ID  uint
	DT  string
	USR int
	DVS string
}

// Config структура для хранения последнего ID из внешней базы данных
type Config struct {
	RemoteID uint
}
