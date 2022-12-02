package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	customError "wasa-photo/service/api/errors"
	"wasa-photo/service/api/functionalities"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) SearchProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	query := r.URL.Query().Get("query")
	users := searchUsername(query)
	if errJson := json.NewEncoder(w).Encode(users); errJson != nil {
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
	profile := functionalities.GetProfile(id)
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
	if !userExists(id) {
		return
	}
	profile := functionalities.GetProfile(id)
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
	profile := functionalities.GetProfile(id)
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
	data := functionalities.GetProfileBasicInfo(id)
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
	data := functionalities.GetUltraBasicProfile(id)
	if errJson := json.NewEncoder(w).Encode(data); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UpdateProfileInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var prof functionalities.ProfileUpdate
	json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(idSession)
	if prof.NewUsername != "" {
		user := users[idSession]
		user.updateUsername(prof.NewUsername)
	} else {
		profile.UpdateProfileInfo(prof)
	}
	// w.Header().Set("Content-type", "application/json")
	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddPhotoProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	//r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//var prof functionalities.PhotoAdd

	//json.NewDecoder(r.Body).Decode(&prof)
	// TODO

	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	prof := functionalities.PhotoAdd{
		Text:                   r.PostFormValue("text"),
		ProfilePictureLocation: r.PostFormValue("profilePicture"),
	}
	profile := functionalities.GetProfile(id)
	if prof.ProfilePictureLocation != "" {
		profile.UpdateProfilePicture(prof.ProfilePictureLocation)
		if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
			http.Error(w, errJson.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	idImage := profile.AddPhoto(prof.Text)

	lastImageId := idImage
	r.ParseMultipartForm(10 << 20)
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

	io.Copy(f, file)

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeletePhotoProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// var prof functionalities.PhotoAdd
	// json.NewDecoder(r.Body).Decode(&prof)
	prof := r.URL.Query().Get("imageid")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(id)
	profile.DeletePhoto(prof)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddCommentProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.CommentAdd
	//params := mux.Vars(r)
	id := ps.ByName("id")
	imageId := ps.ByName("imageid")
	json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(id)
	profile.AddPhotoComment(idSession, imageId, prof.Comment)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeleteCommentProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//var prof functionalities.DeleteElement
	//params := mux.Vars(r)
	id := ps.ByName("id")
	imageId := ps.ByName("imageid")
	//json.NewDecoder(r.Body).Decode(&prof)
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
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(id)
	profile.DeletePhotoComment(sessionId, imageId, index)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeleteLikeProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//var prof functionalities.DeleteElement
	//params := mux.Vars(r)
	id := ps.ByName("id")
	imageId := ps.ByName("imageid")
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(id)
	profile.DeletePhotoLike(sessionId, imageId)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddLikeProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//var prof functionalities.LikeAdd
	//params := mux.Vars(r)
	id := ps.ByName("id")
	imageId := ps.ByName("imageid")
	//json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(id)
	profile.AddPhotoLike(idSession, imageId)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//var prof functionalities.FollowerAdd
	id := ps.ByName("id")
	//json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(idSession)
	profile.AddFollowings(id)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UnFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	//var prof functionalities.FollowerAdd
	//params := mux.Vars(r)
	id := ps.ByName("id")
	//json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(idSession)
	profile.UnFollowers(id)

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) BanFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// var prof functionalities.FollowerAdd
	// params := mux.Vars(r)
	id := ps.ByName("id")
	//json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(idSession)
	profile.AddBans(id)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UnBanFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// var prof functionalities.FollowerAdd
	// params := mux.Vars(r)
	id := ps.ByName("id")
	//json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(idSession)
	profile.UnBans(id)
	// w.Header().Set("Content-type", "application/json")

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
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)

			} else if err.Error() == "StatusInternalServerError" {
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
	profile := functionalities.GetProfile(id)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	id := r.URL.Query().Get("id")
	if errJson := json.NewEncoder(w).Encode(users[id]); errJson != nil {
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
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST")
	// //w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// w.Header().Set("content-type", "application/json")
	// w.WriteHeader(200)
	log.Println("E QUA ENTRA PERO")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Endcoding, Content-Type, Content-Length, Authorization, X-CSRF-token")
	//setupCorsResponse(&w, r)
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	//log.Printf(creds.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := creds.returnID()
	log.Println(userID)
	r.Header.Add("Token", userID)
	type returnId struct{ id string }
	if errJson := json.NewEncoder(w).Encode(userID); errJson != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func (rt *_router) AddSeen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var prof functionalities.PhotoAdd
	json.NewDecoder(r.Body).Decode(&prof)
	ua := r.Header.Get("Token")
	session, err := returnSessionFromId(ua)
	if err != nil {
		switch err.(type) {
		case *customError.ErrStatus:
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else if err.Error() == "StatusInternalServerError" {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		default:
			http.Error(w, "Enable to get Error Type", http.DefaultMaxHeaderBytes)
		}
	}
	id := session.Id
	profile := functionalities.GetProfile(id)
	profile.AddAlreadySeen(prof.IdImage)
	if errJson := json.NewEncoder(w).Encode("OK!"); errJson != nil {
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
			if err.Error() == "StatusUnauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else if err.Error() == "StatusInternalServerError" {
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
	entry := users[id]
	userData := users[id].Data
	userData.updateInfo(r, "/")
	entry.Data = userData
	users[id] = entry
	// filter := bson.D{{Key: "id", Value: id}}
	// update := bson.D{{Key: "$set", Value: bson.D{{Key: "data", Value: entry.Data}}}}
	//go mongodb.CollectionUsers.UpdateOne(mongodb.Ctx, filter, update)
	// w.Header().Set("Content-type", "application/json")
	profile := functionalities.GetProfile(id)
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
	buf, err := ioutil.ReadFile("./public/images/" + id)
	if err != nil {
		http.Error(w, "File Not Found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	//w.Header().Set("Content-Disposition", `attachment;filename="${id}"`)
	w.Write(buf)
	return
}
