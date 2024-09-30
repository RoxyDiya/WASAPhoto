package database

import (
	"database/sql"
	"errors"
)

// GetUserToken retrieves or creates a user token based on the provided username.
func (db *appdbimpl) GetUserToken(username string) (int64, error) {
	token, err := db.GetUserTokenOnly(username)
	if errors.Is(err, sql.ErrNoRows) {
		// No user found, create a new one
		token, err = db.addUser(username)
		if err != nil {
			return 0, err
		}
	}
	return token, err
}

// GetUserProfile retrieves the user profile along with photos, followers, and following info.
func (db *appdbimpl) GetUserProfile(username string, requestUser int64) (UserProfile, error) {
	var profile UserProfile

	// Fetch the user token
	var token int64
	err := db.c.QueryRow("SELECT token FROM user WHERE username=?", username).Scan(&token)
	if err != nil {
		return profile, err
	}

	// Fetch user data (token and username)
	profile.Token, profile.Username, err = db.getUserData(token)
	if err != nil {
		return profile, err
	}

	// Get user photos
	profile.Photos, err = db.getListOfPhotos(profile.Token, requestUser)
	if err != nil {
		return profile, err
	}
	profile.NumberOfPhotos = int64(len(profile.Photos))

	// Get followers and following counts
	profile.NumberOfFollowers, err = db.getNumberFollowers(profile.Token)
	if err != nil {
		return profile, err
	}

	profile.NumberOfFollowing, err = db.getNumberFollowing(profile.Token)
	if err != nil {
		return profile, err
	}

	// Check relationship between the requesting user and the profile owner
	if requestUser == profile.Token {
		profile.IsOwner = true
		profile.IsFollowed = false
		profile.IsBanned = false
	} else {
		profile.IsOwner = false
		profile.IsFollowed, err = db.CheckFollow(requestUser, profile.Token)
		if err != nil {
			return profile, err
		}
		profile.IsBanned, err = db.CheckBan(requestUser, profile.Token)
		if err != nil {
			return profile, err
		}
	}

	return profile, nil
}

// SetUserName updates the username for a given user token.
func (db *appdbimpl) SetUserName(token int64, username string) error {
	_, err := db.c.Exec("UPDATE user SET username=? WHERE token=?", username, token)
	return err
}

// GetUsersList retrieves a list of users whose usernames contain the provided substring.
func (db *appdbimpl) GetUsersList(substring string) ([]string, error) {
	var users []string
	rows, err := db.c.Query("SELECT username FROM user WHERE username LIKE ?", "%"+substring+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		users = append(users, username)
	}
	return users, rows.Err()
}

// CheckToken verifies if a given user token exists in the database.
func (db *appdbimpl) CheckToken(token int64) bool {
	var count int64
	err := db.c.QueryRow("SELECT COUNT(*) FROM user WHERE token=?", token).Scan(&count)
	return err == nil && count == 1
}

// Get user data (token and username).
func (db *appdbimpl) getUserData(token int64) (int64, string, error) {
	var username string
	err := db.c.QueryRow("SELECT token, username FROM user WHERE token=?", token).Scan(&token, &username)
	if err != nil {
		return 0, "", err
	}
	return token, username, nil
}

// GetNumberOfLikes returns the number of likes for a given photo.
func (db *appdbimpl) GetNumberOfLikes(photoId int64) (int64, error) {
	var count int64
	err := db.c.QueryRow("SELECT count(*) FROM like WHERE photo=?", photoId).Scan(&count)
	return count, err
}

// GetNumberOfComments returns the number of comments for a given photo.
func (db *appdbimpl) GetNumberOfComments(photoId int64) (int64, error) {
	var count int64
	err := db.c.QueryRow("SELECT count(*) FROM comment WHERE photo=?", photoId).Scan(&count)
	return count, err
}

// GetNumberFollowers returns the number of followers for a user.
func (db *appdbimpl) getNumberFollowers(token int64) (int64, error) {
	var count int64
	err := db.c.QueryRow("SELECT count(*) FROM follow WHERE followed=?", token).Scan(&count)
	return count, err
}

// GetNumberFollowing returns the number of users the given token is following.
func (db *appdbimpl) getNumberFollowing(token int64) (int64, error) {
	var count int64
	err := db.c.QueryRow("SELECT count(*) FROM follow WHERE following=?", token).Scan(&count)
	return count, err
}

// GetListOfPhotos retrieves the list of photos for a user, along with likes, comments, and whether the requesting user liked them.
func (db *appdbimpl) getListOfPhotos(userToken int64, requestUser int64) ([]Photo, error) {
	rows, err := db.c.Query("SELECT id, owner, u.username, created_at FROM photo JOIN user u ON u.token = photo.owner WHERE owner=?", userToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.Id, &photo.Owner, &photo.OwnerUsername, &photo.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Fetch additional photo details
		photo.NumberOfLikes, err = db.GetNumberOfLikes(photo.Id)
		if err != nil {
			return nil, err
		}
		photo.NumberOfComments, err = db.GetNumberOfComments(photo.Id)
		if err != nil {
			return nil, err
		}
		photo.IsLiked, err = db.CheckLike(requestUser, photo.Id)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, rows.Err()
}

// AddUser adds a new user to the database and returns the newly created user's token.
func (db *appdbimpl) addUser(username string) (int64, error) {
	res, err := db.c.Exec("INSERT INTO user (username) VALUES (?)", username)
	if err != nil {
		return 0, err
	}
	token, err := res.LastInsertId()
	return token, err
}

// GetUserTokenOnly retrieves the token for a given username.
func (db *appdbimpl) GetUserTokenOnly(username string) (int64, error) {
	var token int64
	err := db.c.QueryRow("SELECT token FROM user WHERE username=?", username).Scan(&token)
	return token, err
}

// CheckUsernameExistence checks if a username exists in the database.
func (db *appdbimpl) CheckUsernameExistence(username string) (int64, error) {
	var count int64
	err := db.c.QueryRow("SELECT count(*) FROM user WHERE username=?", username).Scan(&count)
	return count, err
}
