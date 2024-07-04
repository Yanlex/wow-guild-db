package main

import (
	"kvd/deployments/db/structure"
)

// Запускаем создание структуры БД и заполнение стартовой информацией о гильдии с АПИ raid.io через defer firstfilldb.FirstFillDB()
func init() {
	structure.Init()
}

func main() {

}
