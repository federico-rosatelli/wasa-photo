package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"wasa-photo/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

const (
	StatusUnauthorized        = "StatusUnauthorized"
	StatusInternalServerError = "StatusInternalServerError"
	StatusBadRequest          = "Empty Token"
)

func (rt *_router) SearchProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	query := r.URL.Query().Get("query")
	precise := r.URL.Query().Get("precise")
	ua := r.Header.Get("Token")
	_, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	var userQuery []UltraBasicProfile
	if precise == "1" {
		dataId, v := findIdByUsername(query)
		if !v {
			if errJson := json.NewEncoder(w).Encode(userQuery); errJson != nil {
				http.Error(w, errJson.Error(), http.StatusBadRequest)
				return
			}
		}
		dataUserQuery := GetUltraBasicProfile(dataId)
		userQuery = append(userQuery, dataUserQuery)
	} else {
		userQuery = searchUsername(query)
	}
	if errJson := json.NewEncoder(w).Encode(userQuery); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetProfileImageInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	id := ps.ByName("id")
	imageId := ps.ByName("imageid")
	if !userExists(id) {
		http.Error(w, "Invalid Params", http.StatusNotFound)
		return
	}
	ua := r.Header.Get("Token")
	_, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	profile := GetProfile(id)
	image, err := profile.GetImageInfo(imageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errJson := json.NewEncoder(w).Encode(image); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetProfileFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	id := ps.ByName("id")
	if !userExists(id) {
		http.Error(w, "Invalid Params", http.StatusNotFound)
		return
	}
	ua := r.Header.Get("Token")
	_, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	profile := GetProfile(id)
	followersdata := profile.GetBasicUserFollowers()
	if errJson := json.NewEncoder(w).Encode(followersdata); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetProfileFollowings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	id := ps.ByName("id")
	if !userExists(id) {
		http.Error(w, "Invalid Params", http.StatusNotFound)
		return
	}
	ua := r.Header.Get("Token")
	_, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	profile := GetProfile(id)
	followings := profile.GetBasicUserFollowings()
	if errJson := json.NewEncoder(w).Encode(followings); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetBasicProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	id := ps.ByName("id")
	if !userExists(id) {
		http.Error(w, "Invalid Params", http.StatusNotFound)
		return
	}
	ua := r.Header.Get("Token")
	_, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	data := GetProfileBasicInfo(id)
	if errJson := json.NewEncoder(w).Encode(data); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetUltraBasicProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	id := ps.ByName("id")
	if !userExists(id) {
		http.Error(w, "Invalid Params", http.StatusNotFound)
		return
	}
	ua := r.Header.Get("Token")
	_, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	data := GetUltraBasicProfile(id)
	if errJson := json.NewEncoder(w).Encode(data); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UpdateProfileInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var prof ProfileUpdate
	err := json.NewDecoder(r.Body).Decode(&prof)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	if prof.NewUsername != "" {
		user := users[idSession]
		err := user.updateUsername(prof.NewUsername, *rt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		err = profile.UpdateProfileInfo(prof, *rt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddPhotoProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	id := session.Id
	prof := PhotoAdd{
		Text:                   r.PostFormValue("text"),
		ProfilePictureLocation: r.PostFormValue("profilePicture"),
	}
	profile := GetProfile(id)
	if prof.ProfilePictureLocation != "" {

		profile.UpdateProfilePicture(prof.ProfilePictureLocation)
		if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
			http.Error(w, errJson.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	idImage, err := profile.AddPhoto(prof.Text, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastImageId := idImage
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, _, err := r.FormFile("myFile")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := os.OpenFile("./public/images/"+lastImageId+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error: ", err.Error())
		f.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeletePhotoProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	imageId := ps.ByName("imageid")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	id := session.Id
	profile := GetProfile(id)
	err = profile.DeletePhoto(imageId, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddCommentProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof CommentAdd
	id := ps.ByName("id")
	imageId := ps.ByName("imageId")
	err := json.NewDecoder(r.Body).Decode(&prof)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(id)
	log.Println("ID SESSION", idSession, "IMAGE ID", imageId, prof.Comment)
	err = profile.AddPhotoComment(idSession, imageId, prof.Comment, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeleteCommentProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	imageId := ps.ByName("imageId")
	prof := r.URL.Query().Get("index")
	index, err := strconv.Atoi(prof)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	sessionId := session.Id
	profile := GetProfile(id)
	log.Println("session:", sessionId, "imageId", imageId, "index", index)
	err = profile.DeletePhotoComment(sessionId, imageId, index, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeleteLikeProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	imageId := ps.ByName("imageId")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	sessionId := session.Id
	profile := GetProfile(id)
	err = profile.DeletePhotoLike(sessionId, imageId, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddLikeProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	imageId := ps.ByName("imageid")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(id)
	err = profile.AddPhotoLike(idSession, imageId, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	ua := r.Header.Get("Token")
	log.Println(ua)
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	err = profile.AddFollowings(id, *rt)
	if err != nil {
		log.Println("QUA CON TUTTO", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	basic := GetProfileBasicInfo(id)

	if errJson := json.NewEncoder(w).Encode(basic); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UnFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	err = profile.UnFollowers(id, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	basic := GetProfileBasicInfo(id)

	if errJson := json.NewEncoder(w).Encode(basic); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetBanFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	response := profile.IsBan(id)

	if errJson := json.NewEncoder(w).Encode(response); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) BanFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	log.Println(profile.Bans)
	err = profile.AddBans(id, *rt)
	log.Println(profile.Bans)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UnBanFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	err = profile.UnBans(id, *rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) ProfileInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	id := session.Id
	basic := GetProfileBasicInfo(id)

	if errJson := json.NewEncoder(w).Encode(basic); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := creds.returnID(*rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(userID)
	if errJson := json.NewEncoder(w).Encode(userID); errJson != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddSeen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var prof PhotoAdd
	err := json.NewDecoder(r.Body).Decode(&prof)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	id := session.Id
	profile := GetProfile(id)
	profile.AddAlreadySeen(prof.IdImage)
	if errJson := json.NewEncoder(w).Encode("OK"); errJson != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) Welcome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ua := r.Header.Get("Token")
	log.Println(ua)
	session, err := returnSessionFromId(ua)
	if err != nil {
		if err.Error() == StatusUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if err.Error() == StatusInternalServerError {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err.Error() == StatusBadRequest {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}
	id := session.Id
	if rt.db != nil {
		entry := users[id]
		userData := users[id].Data
		userData.updateInfo(r, "/")
		entry.Data = userData
		users[id] = entry
	}
	profile := GetProfile(id)
	streamList := profile.GetNewStream()
	if len(streamList) == 0 {
		return
	}
	if errJson := json.NewEncoder(w).Encode(streamList); errJson != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// Serve image for frontend
func (s *_router) ServeImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	path, _ := os.Getwd()
	buf, err := os.ReadFile(path + "/public/images/" + id)
	if err != nil {

		http.Error(w, "File Not Found "+path+": "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(buf)
	if err != nil {
		path, _ := os.Getwd()
		http.Error(w, "File Not Found "+path, http.StatusInternalServerError)
		return
	}
}
