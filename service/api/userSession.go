// 		---- Session User ----

package api

import (
	"wasa-photo/service/api/errors"
)

type Session struct {
	Username string
	Id       string
}

// sessions represent a set of struct Session
// For easy usage it'll use Id as the key
var sessions = map[string]Session{}

// validateSession is used for validate a current session struct,
// matching the string of the Username in the session with the user
// provided
// func (s Session) validateSession(username string) bool {
// 	return s.Username == username
// }

// updateIDSession is used in a complex version of the server
// where the key of the map sessions is updated every (int) seconds
// Note used in this version
// func updateIDSession(oldID string) string {
// 	newID := uuid.NewString()
// 	newSession := Session{
// 		Username: sessions[oldID].Username,
// 		Id:       oldID,
// 	}
// 	sessions[newID] = newSession
// 	delete(sessions, oldID)
// 	return newID
// }

// updateUsernameSession is used for the update of the username
// of the client in the map sessions, identified by id
func (s *Session) updateUsernameSession(newUsername string) {
	s.Username = newUsername
}

// validateUserByUsernameID is used as a verification for a logged
// user. It takes the user's username and id and returns a error
// if the id doesn't exist or the user is different. It'll return
// nil if the username and id match the session parameters
// func validateUserByUsernameID(username string, id string) error {
// 	session, err := sessions[id]
// 	//log.Println(username, id)
// 	if !err {
// 		return errors.NewErrStatus("ID Not Found for user " + username)
// 	}
// 	if !session.validateSession(username) {
// 		return errors.NewErrStatus("Username Not Matching with ID " + id)
// 	}
// 	// if session.isExpired() {
// 	// 	return errors.NewErrStatus("SessionExpired")
// 	// }
// 	sessions[id] = session
// 	return nil
// }

func returnSessionFromId(id string) (Session, error) {
	session, err := sessions[id]
	if id == "" {
		return Session{}, errors.NewErrStatus("Token Error: " + id)
	}
	if !err {
		return Session{}, errors.NewErrStatus("StatusUnauthorized")
	}
	return session, nil
}
