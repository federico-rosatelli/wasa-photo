package database

import "time"

type Profile struct {
	Username       string
	Id             string
	Name           string
	Surname        string
	ProfilePicture ProfilePicture
	Followers      []UserFollow
	Followings     []UserFollow
	Bans           []UserFollow
	Images         []Image
	AlreadySeen    map[string]string
}

type UserFollow struct {
	IdUser string
	Time   string
}

type Like struct {
	UserIdLike string
	Time       string
}
type Comment struct {
	UserIdComment string
	Time          string
	Content       string
}

type Image struct {
	IdImage  string
	Location string
	Text     string
	Time     string
	Likes    []Like
	Comments []Comment
}

type ProfilePicture struct {
	Location string
	Time     string
}

type Session struct {
	Username  string
	Id        string
	IdSession string
}

type JsonUserAgent struct {
	UserAgentString string `json:"userAgentString"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Version         string `json:"version"`
	VersionMajor    string `json:"versionMajor"`
	Device          struct {
		Name  string      `json:"name"`
		Type  string      `json:"type"`
		Brand string      `json:"brand"`
		CPU   interface{} `json:"CPU"`
	} `json:"device"`
	Engine struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Version      string `json:"version"`
		VersionMajor string `json:"versionMajor"`
	} `json:"engine"`
	OperatingSystem struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		Version      string `json:"version"`
		VersionMajor string `json:"versionMajor"`
	} `json:"operatingSystem"`
}

type User struct {
	Username string
	Id       string
	Data     UserData
}

type Path struct {
	PathName  string
	TimeFirst time.Time
}

type UserData struct {
	LastShoppingTime time.Time
	UserInfo         struct {
		Ip_client       string
		UserAgentString string
		UserAgent       JsonUserAgent
		Time_client     time.Time
		Times_visited   int
		Paths           []Path
	}
}
