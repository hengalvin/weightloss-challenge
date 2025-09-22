package database

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id                int     `json:"id"`
	Email             string  `json:"email"`
	Name              string  `json:"name"`
	Password          string  `json:"-"`
	InitialWeightGr   int     `json:"initial_weight_gr"`
	CurrentWeightGr   int     `json:"current_weight_gr"`
	WeightLost        int     `json:"weight_lost"`
	WeightLostPercent float32 `json:"weight_lost_per"`
	CreatedAt         string  `json:"joined"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users (email, name, password_hash, initial_weight_gr, current_weight_gr, weight_lost_gr, weight_lost_per) VALUES (?1,?2,?3,?4,?5,?6,?7) RETURNING id, created_at"

	return m.DB.QueryRowContext(ctx, query, user.Email, user.Name, user.Password, user.InitialWeightGr, user.CurrentWeightGr, 0, 0).Scan(&user.Id, &user.CreatedAt)

}

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.InitialWeightGr, &user.CurrentWeightGr, &user.WeightLost, &user.WeightLostPercent, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	query := `SELECT * FROM users WHERE id = ?1`
	return m.getUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = ?1`
	return m.getUser(query, strings.ToLower(email))
}

func (m *UserModel) UpdateUserWeights(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user.WeightLost = calculateWeightLost(user)
	user.WeightLostPercent = calculateWeightLostPer(user)

	query := "UPDATE users SET weight_lost_gr = ?1, weight_lost_per = ?2, current_weight_gr = ?3 WHERE id = ?4"

	_, err := m.DB.ExecContext(ctx, query, user.WeightLost, user.WeightLostPercent, user.CurrentWeightGr, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) GetAllUserInfo() ([]userInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT name, initial_weight_gr, current_weight_gr, weight_lost_gr, weight_lost_per FROM users"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userInfoList []userInfo
	for rows.Next() {
		var userInfo userInfo
		err := rows.Scan(&userInfo.Name, &userInfo.InitialWeight, &userInfo.CurrentWeightGr, &userInfo.WeightLost, &userInfo.WeightLostPercent)
		if err != nil {
			return nil, err
		}

		userInfoList = append(userInfoList, userInfo)
	}

	return userInfoList, nil
}

func calculateWeightLost(user *User) int {
	return user.InitialWeightGr - user.CurrentWeightGr
}

func calculateWeightLostPer(user *User) float32 {
	return float32(user.WeightLost) / float32(user.InitialWeightGr) * 100
}
