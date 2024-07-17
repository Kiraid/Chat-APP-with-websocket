package models

import (
	"errors"

	"example.com/chat/db"
)
type Channel struct{
	ID int
	Name string
}

func (c *Channel) Save() error {
	stmt, err := db.DB.Prepare("INSERT INTO channels(name) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(c.Name)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return nil
}

func GetChannelByName(name string) (*Channel, error) {
    var channel Channel
    err := db.DB.QueryRow("SELECT id, name FROM channels WHERE name = ?", name).Scan(&channel.ID, &channel.Name)
    if err != nil {
            return nil, errors.New("channel not found")
    }
    return &channel, nil
}


func GetChannelIDByName(name string) (int, error) {
	var id int
	err := db.DB.QueryRow("SELECT id FROM channels WHERE name = ?", name).Scan(&id)
	return id, err

}

func AddUsertoChannel(userID, channelID int) error {
	stmt, err := db.DB.Prepare("INSERT INTO user_channels (user_id, channel_id) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, channelID)
	if err != nil {
		return err
	}
	return nil

}