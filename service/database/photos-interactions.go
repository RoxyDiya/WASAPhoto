package database

// Helper function to execute a query that doesn't return rows
func (db *appdbimpl) execQuery(query string, args ...interface{}) error {
	_, err := db.c.Exec(query, args...)
	return err
}

// Helper function to execute a query that checks existence (count-based)
func (db *appdbimpl) checkExistence(query string, args ...interface{}) (bool, error) {
	var count int64
	err := db.c.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

// Posting and Deleting Photos
func (db *appdbimpl) PostPhoto(image []byte, token int64) error {
	return db.execQuery("INSERT INTO photo (owner, img) VALUES (?, ?)", token, image)
}

func (db *appdbimpl) DeletePhoto(token int64, photoId int64) error {
	return db.execQuery("DELETE FROM photo WHERE owner=? AND id=?", token, photoId)
}

// Retrieving Photo Data
func (db *appdbimpl) GetImage(photoId int64) ([]byte, error) {
	var image []byte
	err := db.c.QueryRow("SELECT img FROM photo WHERE id=?", photoId).Scan(&image)
	return image, err
}

func (db *appdbimpl) GetPhotoOwner(photoId int64) (int64, error) {
	var owner int64
	err := db.c.QueryRow("SELECT owner FROM photo WHERE id=?", photoId).Scan(&owner)
	return owner, err
}

func (db *appdbimpl) CheckPhotoOwner(token int64, photoId int64) (bool, error) {
	return db.checkExistence("SELECT count(*) FROM photo WHERE id=? AND owner=?", photoId, token)
}

// Liking and Unliking Photos
func (db *appdbimpl) LikePhoto(token int64, photoId int64) error {
	return db.execQuery("INSERT INTO like (owner, photo) VALUES (?, ?)", token, photoId)
}

func (db *appdbimpl) UnlikePhoto(token int64, photoId int64) error {
	return db.execQuery("DELETE FROM like WHERE owner=? AND photo=?", token, photoId)
}

func (db *appdbimpl) CheckLike(token int64, photoId int64) (bool, error) {
	return db.checkExistence("SELECT count(*) FROM like WHERE owner=? AND photo=?", token, photoId)
}

// Commenting on Photos
func (db *appdbimpl) CommentPhoto(token int64, photoId int64, content string) (int64, error) {
	res, err := db.c.Exec("INSERT INTO comment (owner, content, photo) VALUES (?, ?, ?)", token, content, photoId)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (db *appdbimpl) GetPhotoComments(photoId int64) ([]FullDataComment, error) {
	rows, err := db.c.Query("SELECT id, content, created_at, u.username, photo FROM comment JOIN user u ON u.token = comment.owner WHERE photo=? ORDER BY created_at DESC", photoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []FullDataComment
	for rows.Next() {
		var comment FullDataComment
		if err := rows.Scan(&comment.Id, &comment.Content, &comment.CreatedAt, &comment.Owner, &comment.Photo); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return comments, rows.Err()
	}

	return comments, nil
}

func (db *appdbimpl) GetCommentOwner(commentId int64) (int64, error) {
	var owner int64
	err := db.c.QueryRow("SELECT owner FROM comment WHERE id=?", commentId).Scan(&owner)
	return owner, err
}

func (db *appdbimpl) GetComment(commentId int64) (FullDataComment, error) {
	var comment FullDataComment
	err := db.c.QueryRow("SELECT * FROM comment WHERE id=?", commentId).Scan(&comment.Id, &comment.Content, &comment.CreatedAt, &comment.Owner, &comment.Photo)
	return comment, err
}

func (db *appdbimpl) DeleteComment(commentId int64) error {
	return db.execQuery("DELETE FROM comment WHERE id=?", commentId)
}

// Stream (Fetching Photos for the User's Stream)
func (db *appdbimpl) GetMyStream(token int64) ([]Photo, error) {
	rows, err := db.c.Query("SELECT id, owner, u.username, created_at FROM photo JOIN user u ON u.token = photo.owner WHERE owner NOT IN (SELECT banning FROM ban WHERE banned=?) AND owner != ? AND owner IN (SELECT followed FROM follow WHERE following=?) ORDER BY created_at DESC", token, token, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []Photo
	for rows.Next() {
		var photo Photo
		err := rows.Scan(&photo.Id, &photo.Owner, &photo.OwnerUsername, &photo.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Add additional photo data
		photo.NumberOfLikes, err = db.GetNumberOfLikes(photo.Id)
		if err != nil {
			return nil, err
		}
		photo.NumberOfComments, err = db.GetNumberOfComments(photo.Id)
		if err != nil {
			return nil, err
		}
		photo.IsLiked, err = db.CheckLike(token, photo.Id)
		if err != nil {
			return nil, err
		}

		photos = append(photos, photo)
	}

	if rows.Err() != nil {
		return photos, rows.Err()
	}

	return photos, nil
}

// Checking Photo Existence
func (db *appdbimpl) CheckPhotoExistence(photoId int64) (bool, error) {
	return db.checkExistence("SELECT count(*) FROM photo WHERE id=?", photoId)
}

// Additional Helper Functions for Likes and Comments
func (db *appdbimpl) GetNumberOfLikes(photoId int64) (int64, error) {
	var count int64
	err := db.c.QueryRow("SELECT count(*) FROM like WHERE photo=?", photoId).Scan(&count)
	return count, err
}

func (db *appdbimpl) GetNumberOfComments(photoId int64) (int64, error) {
	var count int64
	err := db.c.QueryRow("SELECT count(*) FROM comment WHERE photo=?", photoId).Scan(&count)
	return count, err
}
