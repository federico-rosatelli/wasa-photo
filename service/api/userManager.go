// 		---- Manager User ----

package api

import (
	"net/http"
	"regexp"
	"time"
	"wasa-photo/service/api/customErrors"

	"github.com/google/uuid"
)

// const key = "f57816d787b74374881e252127055088"

// struct from https://api.ipgeolocation.io
type jsonUserAgent struct {
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

type Credentials struct {
	Username string `json:"username"`
}

type User struct {
	Username string
	Id       string
	Data     userData
}

// User map used to store the user's
// informations and logs.
// The key is the Id of the user and the
// content is the struct User
var users = map[string]User{}

type Path struct {
	PathName  string
	TimeFirst time.Time
}

type userData struct {
	LastShoppingTime int64
	UserInfo         struct {
		Ip_client       string
		UserAgentString string
		UserAgent       jsonUserAgent
		Time_client     int64
		Times_visited   int
		Paths           []Path
	}
}

// Return the Id of the user identified by username.
// If the user exists it returns the id,
// if the user doesn't exist the function creates
// the user structure and pushes it into the database
func (cred Credentials) returnID(rt _router) (string, error) {
	for key, value := range sessions {
		if value.Username == cred.Username {
			return key, nil
		}
	}
	newUser := User{
		Username: cred.Username,
		Id:       uuid.NewString(),
	}

	users[newUser.Id] = newUser
	newSessionToken := uuid.NewString()
	newSession := Session{
		Username:  cred.Username,
		Id:        newUser.Id,
		IdSession: newSessionToken,
	}

	sessions[newSessionToken] = newSession

	err := NewProfile(cred.Username, newUser.Id, rt)
	if err != nil {
		return "", err
	}
	if rt.db != nil {
		err = newUser.newUserDB(rt)
		if err != nil {
			return "", err
		}
		err = newSession.newSessionDB(newSessionToken, rt)
		if err != nil {
			return "", err
		}
	}
	return newSessionToken, nil
}

func (u User) newUserDB(rt _router) error {
	return rt.db.InsertOneUsers(u.converUser())
}

// Insert the new session in the database
func (s Session) newSessionDB(newSessionToken string, rt _router) error {
	return rt.db.InsertOneSession(s.converSession(newSessionToken))
}

// Return a string and a boolean. It'll return the id of the user and true if exist a username
// that match the string searched by the user
func findIdByUsername(username string) (string, bool) {
	for key, value := range users {
		if value.Username == username {
			return key, true
		}
	}
	return "", false
}

// Basic regex match.
// The function try to match the given string searched
// with a username. It'll return a boolean based on that
// match
func matchRegex(query string, username string) bool {
	match, _ := regexp.MatchString("^"+query, username)
	return match
}

// Function for finding a user in the search bar.
// It'll return a list of UltraBasicProfile struct
// from the functionalities package.
func searchUsername(query string) []UltraBasicProfile {
	var matching []UltraBasicProfile
	for _, username := range users {
		if matchRegex(query, username.Username) {
			if userId, ok := findIdByUsername(username.Username); ok {
				data := GetUltraBasicProfile(userId)
				matching = append(matching, data)
			}
		}
	}
	return matching
}

// Check if user exists in users map
// It'll return a boolean based on the
// existence of the id passed in the map
// users
func userExists(id string) bool {
	_, ok := users[id]
	return ok
}

// Update the username of the user
func (u *User) updateUsername(newUsername string, rt _router) error {
	var sessionToken string
	for key, value := range sessions {
		if value.Username == u.Username {
			sessionToken = key
		}
	}
	for _, value := range users {
		if value.Username == newUsername {
			return customErrors.NewErrStatus("User already exists")
		}
	}
	u.Username = newUsername
	session := sessions[sessionToken]
	session.updateUsernameSession(newUsername, rt)
	profile := GetProfile(u.Id)
	users[u.Id] = *u
	err := profile.SetMyUsername(newUsername, rt)
	return err
}

// Extract the user-agent string information
// from https://api.ipgeolocation.io
// It'll return a jsonAgent struct or an error
// func user_agent_extractor(user_agent string) (jsonUserAgent, error) {
// 	log.Printf("USER AGENT: %s", user_agent)
// 	client := http.Client{}
// 	url := `https://api.ipgeolocation.io/user-agent?apiKey=` + key
// 	request, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return jsonUserAgent{}, customError.NewErrStatus("Not able to request on -api.ipgeolocation.io-")
// 	}

// 	request.Header.Set("User-Agent", user_agent)

// 	res, err := client.Do(request)
// 	if err != nil {
// 		return jsonUserAgent{}, customError.NewErrStatus("User-Agent Header Set went wrong")
// 	}
// 	returnJson := json.NewDecoder(res.Body)
// 	var dataStruct jsonUserAgent
// 	errDecoding := returnJson.Decode(&dataStruct)
// 	if errDecoding != nil {
// 		return jsonUserAgent{}, customError.NewErrStatus("InternalServerError on json.Decode")
// 	}
// 	return dataStruct, nil
// }

// Update the info of the client such as:
// ip-address, user-agent, user-agent-information
// and the times visited of the user
func (user *userData) updateInfo(r *http.Request, url string) {
	user.UserInfo.Ip_client = r.RemoteAddr
	user_agent := r.UserAgent()
	if user.UserInfo.UserAgentString != user_agent {
		user.UserInfo.UserAgentString = user_agent
		// log.Printf("NUOVO USER-AGENT")
		// json_user_extract, ok := user_agent_extractor(user_agent)
		// if ok != nil {
		// 	log.Printf(" ERROR %s\n", ok)
		// } else {
		// 	user.UserInfo.UserAgent = json_user_extract
		// }
	}
	user.UserInfo.Times_visited += 1
	path := Path{
		PathName:  url,
		TimeFirst: time.Now(),
	}
	user.UserInfo.Paths = append(user.UserInfo.Paths, path)
}

// func (user *User) updateID(newID string) {
// 	user.Id = newID
// }
