package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	customError "wasa-photo/service/api/errors"
	"wasa-photo/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

const StatusUnauthorized = "StatusUnauthorized"
const StatusInternalServerError = "StatusInternalServerError"

func (rt *_router) SearchProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	query := r.URL.Query().Get("query")
	precise := r.URL.Query().Get("precise")
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
	//params := mux.Vars(r)
	id := ps.ByName("id")
	imageId := ps.ByName("imageid")
	if !userExists(id) {
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
	//params := mux.Vars(r)
	id := ps.ByName("id")
	log.Println(id)
	if !userExists(id) {
		return
	}
	profile := GetProfile(id)
	followers := profile.GetBasicUserFollowers()
	if errJson := json.NewEncoder(w).Encode(followers); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetProfileFollowings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	//params := mux.Vars(r)
	id := ps.ByName("id")
	if !userExists(id) {
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
	//params := mux.Vars(r)
	id := ps.ByName("id")
	if !userExists(id) {
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
	//params := mux.Vars(r)
	id := ps.ByName("id")
	if !userExists(id) {
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
	json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)
		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	if prof.NewUsername != "" {
		user := users[idSession]
		user.updateUsername(prof.NewUsername, *rt)
	} else {
		err = profile.UpdateProfileInfo(prof, *rt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// w.Header().Set("Content-type", "application/json")
	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddPhotoProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

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

	prof := r.URL.Query().Get("imageid")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

		}
		return
	}
	id := session.Id
	profile := GetProfile(id)
	err = profile.DeletePhoto(prof, *rt)
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
	json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

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
	imageId := ps.ByName("imageid")
	prof := r.URL.Query().Get("index")
	index, err := strconv.Atoi(prof)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

		}
		return
	}

	sessionId := session.Id
	profile := GetProfile(id)
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
	imageId := ps.ByName("imageid")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

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
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

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

func (rt *_router) AddFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	ua := r.Header.Get("Token")
	log.Println(ua)
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	err = profile.AddFollowings(id, *rt)
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

func (rt *_router) UnFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	id := ps.ByName("id")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	profile.UnFollowers(id)

	basic := GetProfileBasicInfo(id)

	if errJson := json.NewEncoder(w).Encode(basic); errJson != nil {
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
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

		}
		return
	}
	idSession := session.Id
	profile := GetProfile(idSession)
	err = profile.AddBans(id, *rt)
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
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

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
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)

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

func (rt *_router) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	htmlFile := "./views/login.html"
	http.ServeFile(w, r, htmlFile)
}

func (rt *_router) AddPhotoProfileGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	htmlFile := "./views/addphoto.html"
	http.ServeFile(w, r, htmlFile)
}

func (rt *_router) SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println(err)
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
	json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)
		}
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
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == StatusUnauthorized {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else if err.Error() == StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				r.Header.Add("Token", "Ciao")
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)
		}
		if errJson := json.NewEncoder(w).Encode("/login"); errJson != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

func (s *_router) ServeImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	buf, err := os.ReadFile("./public/images/" + id)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(buf)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusInternalServerError)
		return
	}
}
