package database

import (
	"context"
	"database/sql"
	"time"
)

type WeightModel struct {
	DB *sql.DB
}

type Weight struct {
	Id     int `json:"-"`
	UserId int `json:"user_id"`
	Weight int `json:"weight" binding:"required,gt=10000"`
}

type userWeights struct {
	LogTime string `json:"logged_at"`
	Weight  int    `json:"weight"`
}

type userInfo struct {
	Name              string  `json:"name`
	InitialWeight     int     `json:"initial_weight`
	CurrentWeightGr   int     `json:"current_weight_gr"`
	WeightLost        int     `json:"weight_lost"`
	WeightLostPercent float32 `json:"weight_lost_per"`
}

func (m *WeightModel) Insert(weight *Weight) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO weight (user_id, weight_gr, logged_at) VALUES ($1, $2, $3) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, weight.UserId, weight.Weight, time.Now().Unix()).Scan(&weight.Id)
}

func (m *WeightModel) GetWeightLogByUser(userId int) ([]userWeights, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT logged_at, weight_gr FROM weight WHERE user_id = $1"

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weights []userWeights
	for rows.Next() {
		var myWeight userWeights
		err := rows.Scan(&myWeight.LogTime, &myWeight.Weight)
		if err != nil {
			return nil, err
		}
		weights = append(weights, myWeight)
	}

	return weights, nil
}

func (m *WeightModel) GetCurrentWeight(userId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT weight_gr FROM weight WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1"

	var weight int
	err := m.DB.QueryRowContext(ctx, query, userId).Scan(&weight)
	if err != nil {
		return 0, err
	}

	return weight, nil
}


