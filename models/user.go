// models/user.go
package models

import (
	"database/sql"
	"errors"
	"log"

	"example.com/chat/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       int    
    Username string `form:"username" binding:"required"`
    Email    string `form:"email" binding:"required"`
    Password string `form:"password" binding:"required"`
}

func GetAllUsers() ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Save() error {
	stmt, err := db.DB.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(u.Username, u.Email, string(hashedPassword))
	if err != nil {
		return err
	}
	return nil
}

func Authenticate(email, password string) (*User, error) {
	user := User{}
	row := db.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return &user, nil
}

func GetUserByID(userID int) (User, error) {
    var user User
    query := "SELECT id, username FROM users WHERE id = ?"
    err := db.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username)
	if err == sql.ErrNoRows {
        log.Printf("No user found with ID %d", userID)
        return user, errors.New("user not found")
    }
    if err != nil {
        log.Printf("Error querying user with ID %d: %v", userID, err)
        return user, err
    }
    return user, err
}


func GetUserChannels(userID int) ([]string, error) {
    rows, err := db.DB.Query(`
        SELECT channels.name 
        FROM channels 
        JOIN user_channels ON channels.id = user_channels.channel_id 
        WHERE user_channels.user_id = ?`, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var channels []string
    for rows.Next() {
        var channelName string
        if err := rows.Scan(&channelName); err != nil {
            return nil, err
        }
        channels = append(channels, channelName)
    }
    return channels, nil
}

func GetUserIDByUsername(username string) (int, error) {
	var userID int
	query := "SELECT id FROM users WHERE username = ?"
	err := db.DB.QueryRow(query, username).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, errors.New("user not found")
	}
	return userID, err
}
