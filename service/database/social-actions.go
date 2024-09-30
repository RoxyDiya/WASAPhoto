package database

// AddFollow adds a follow relationship between the following user and the user being followed.
func (db *appdbimpl) AddFollow(following int64, followUsername string) error {
	var followed int64

	// Get the token of the user to follow
	err := db.c.QueryRow("SELECT token FROM user WHERE username=?", followUsername).Scan(&followed)
	if err != nil {
		return err
	}

	// Insert the follow relationship into the database
	_, err = db.c.Exec("INSERT INTO follow (following, followed) VALUES (?, ?)", following, followed)
	return err
}

// RemoveFollow removes a follow relationship between the following user and the user being unfollowed.
func (db *appdbimpl) RemoveFollow(following int64, unfollowUsername string) error {
	var followed int64

	// Get the token of the user to unfollow
	err := db.c.QueryRow("SELECT token FROM user WHERE username=?", unfollowUsername).Scan(&followed)
	if err != nil {
		return err
	}

	// Remove the follow relationship from the database
	_, err = db.c.Exec("DELETE FROM follow WHERE following=? AND followed=?", following, followed)
	return err
}

// AddBan adds a ban relationship between the banning user and the user being banned.
func (db *appdbimpl) AddBan(banning int64, banUsername string) error {
	var banned int64

	// Get the token of the user to ban
	err := db.c.QueryRow("SELECT token FROM user WHERE username=?", banUsername).Scan(&banned)
	if err != nil {
		return err
	}

	// Insert the ban relationship into the database
	_, err = db.c.Exec("INSERT INTO ban (banning, banned) VALUES (?, ?)", banning, banned)
	return err
}

// RemoveBan removes a ban relationship between the banning user and the banned user.
func (db *appdbimpl) RemoveBan(banning int64, banUsername string) error {
	var banned int64

	// Get the token of the user to unban
	err := db.c.QueryRow("SELECT token FROM user WHERE username=?", banUsername).Scan(&banned)
	if err != nil {
		return err
	}

	// Remove the ban relationship from the database
	_, err = db.c.Exec("DELETE FROM ban WHERE banning=? AND banned=?", banning, banned)
	return err
}

// CheckFollow checks if a follow relationship exists between two users.
func (db *appdbimpl) CheckFollow(following int64, followed int64) (bool, error) {
	var count int

	// Check if the follow relationship exists
	err := db.c.QueryRow("SELECT count(*) FROM follow WHERE following=? AND followed=?", following, followed).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckBan checks if a ban relationship exists between two users.
func (db *appdbimpl) CheckBan(banning int64, banned int64) (bool, error) {
	var count int

	// Check if the ban relationship exists
	err := db.c.QueryRow("SELECT count(*) FROM ban WHERE banning=? AND banned=?", banning, banned).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
