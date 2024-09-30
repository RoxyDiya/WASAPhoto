package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high-level interface for the DB.
type AppDatabase interface {
	Ping() error
	GetUserTokenOnly(username string) (int64, error)
	CheckUsernameExistence(username string) (int64, error)
	CheckPhotoOwner(token int64, photoId int64) (bool, error)
	GetPhotoOwner(photoId int64) (int64, error)
	CheckLike(token int64, photoId int64) (bool, error)
	GetNumberOfLikes(photoId int64) (int64, error)
	GetNumberOfComments(photoId int64) (int64, error)
	CheckPhotoExistence(photoId int64) (bool, error)

	GetUserToken(username string) (int64, error)
	SetUserName(token int64, username string) error
	CheckToken(token int64) bool
	GetUserProfile(username string, requestUser int64) (UserProfile, error)
	GetUsersList(username string) ([]string, error)

	AddFollow(following int64, follow string) error
	RemoveFollow(following int64, follow string) error
	AddBan(banning int64, ban string) error
	RemoveBan(banning int64, ban string) error
	CheckFollow(u1 int64, u2 int64) (bool, error)
	CheckBan(u1 int64, u2 int64) (bool, error)

	PostPhoto(image []byte, token int64) error
	DeletePhoto(token int64, photoId int64) error
	GetImage(photoId int64) ([]byte, error)
	LikePhoto(token int64, photoId int64) error
	UnlikePhoto(token int64, photoId int64) error
	CommentPhoto(token int64, photoId int64, content string) (int64, error)
	GetPhotoComments(photoId int64) ([]FullDataComment, error)
	GetCommentOwner(commentId int64) (int64, error)
	DeleteComment(commentId int64) error
	GetMyStream(token int64) ([]Photo, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New creates a new instance of AppDatabase based on the provided SQLite connection `db`.
// An error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// Check if the user table exists; if not, create the database structure.
	if err := checkAndCreateTables(db); err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

// Ping checks the database connection.
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// checkAndCreateTables verifies if the necessary tables exist, and creates them if they do not.
func checkAndCreateTables(db *sql.DB) error {
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `
			CREATE TABLE user (
				token    INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT NOT NULL UNIQUE
			);

			CREATE TABLE photo (
				id         INTEGER PRIMARY KEY AUTOINCREMENT,
				owner      INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
				img        BLOB NOT NULL,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);

			CREATE TABLE likes (
				owner INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
				photo INTEGER NOT NULL REFERENCES photo ON DELETE CASCADE,
				PRIMARY KEY (owner, photo)
			);

			CREATE TABLE follow (
				following INTEGER NOT NULL REFERENCES user,
				followed  INTEGER NOT NULL REFERENCES user,
				PRIMARY KEY (following, followed),
				CHECK (following != followed)
			);

			CREATE TABLE comment (
				id         INTEGER PRIMARY KEY AUTOINCREMENT,
				content    TEXT NOT NULL,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
				owner      INTEGER NOT NULL REFERENCES user,
				photo      INTEGER NOT NULL REFERENCES photo
			);

			CREATE TABLE ban (
				banning INTEGER NOT NULL REFERENCES user,
				banned  INTEGER NOT NULL REFERENCES user,
				PRIMARY KEY (banning, banned),
				CHECK (banning != banned)
			);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return nil
}
