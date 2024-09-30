package structs

type Username struct {
	Username string `json:"name"`
}

type UserProfile struct {
	Token             int64   `json:"token"`
	Username          string  `json:"username"`
	Photos            []Photo `json:"photos"`
	NumberOfPhotos    int64   `json:"numberOfPhotos"`
	NumberOfFollowers int64   `json:"numberOfFollowers"`
	NumberOfFollowing int64   `json:"numberOfFollowing"`
	IsFollowed        bool    `json:"isFollowed"`
	IsOwner           bool    `json:"isOwner"`
	IsBanned          bool    `json:"isBanned"`
}

type Token struct {
	Identifier int64 `json:"identifier"`
}

type Message struct {
	Message string `json:"message"`
}

type CreatedCommentMessage struct {
	CommentId int64 `json:"comment_id"`
}

type Comment struct {
	Comment string `json:"comment"`
}

type FullDataComment struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Owner     string `json:"owner"`
	Photo     int64  `json:"photo"`
}

type Photo struct {
	Id               int64  `json:"id"`
	Owner            int64  `json:"owner"`
	OwnerUsername    string `json:"ownerUsername"`
	CreatedAt        string `json:"createdAt"`
	NumberOfLikes    int64  `json:"numberOfLikes"`
	NumberOfComments int64  `json:"numberOfComments"`
	IsLiked          bool   `json:"isLiked"`
}

