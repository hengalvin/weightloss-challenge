package database

import "database/sql"

type Models struct {
	Users     UserModel
	Weight    WeightModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:     UserModel{DB: db},
		Weight:    WeightModel{DB: db},
	}
}
