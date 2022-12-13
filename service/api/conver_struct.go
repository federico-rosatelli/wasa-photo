package api

import "wasa-photo/service/database"

func (p Profile) converProfile() database.Profile {
	return database.Profile{
		Username:       p.Username,
		Id:             p.Id,
		Name:           p.Name,
		Surname:        p.Surname,
		ProfilePicture: database.ProfilePicture{},
		Followers:      []database.UserFollow{},
		Followings:     []database.UserFollow{},
		Bans:           []database.UserFollow{},
		Images:         []database.Image{},
		AlreadySeen:    map[string]int64{},
	}
}

func (u User) converUser() database.User {
	return database.User{
		Username: u.Username,
		Id:       u.Id,
		Data:     database.UserData{},
	}
}

func (s Session) converSession(idSession string) database.Session {
	return database.Session{
		Username:  s.Username,
		Id:        s.Id,
		IdSession: idSession,
	}
}
