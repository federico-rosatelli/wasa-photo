package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	customError "wasa-photo/service/api/errors"
	"wasa-photo/service/api/functionalities"

	"github.com/gorilla/mux"
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
	params := mux.Vars(r)
	id := params["id"]
	imageId := params["imageId"]
	if !userExists(id) {
		return
	}
	profile := functionalities.GetProfile(id)
	image, err := profile.GetImageInfo(imageId)
	if err != nil {
		return
	}
	if errJson := json.NewEncoder(w).Encode(image); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) GetProfileFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
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
	params := mux.Vars(r)
	id := params["id"]
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
	params := mux.Vars(r)
	id := params["id"]
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
	params := mux.Vars(r)
	id := params["id"]
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
	err := validateUserByUsernameID(prof.Username, prof.Id)
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
	profile := functionalities.GetProfile(prof.Id)
	if prof.NewUsername != "" {
		user := users[prof.Id]
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
	prof := functionalities.PhotoAdd{
		Username: r.PostFormValue("Username"),
		Id:       r.PostFormValue("Id"),
		Text:     r.PostFormValue("Text"),
	}

	err := validateUserByUsernameID(prof.Username, prof.Id)
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
	profile := functionalities.GetProfile(prof.Id)
	if prof.ProfilePictureLocation != "" {
		profile.UpdateProfilePicture(prof.ProfilePictureLocation)
	} else {

		profile.AddPhoto(prof.Text)
	}

	lastImageId := profile.Images[0].IdImage
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("myFile")

	if err != nil {
		return
	}
	defer file.Close()

	f, err := os.OpenFile("./public/images/"+lastImageId+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		f.Close()
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

	var prof functionalities.PhotoAdd
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(prof.Id)
	profile.DeletePhoto(prof.IdImage)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddCommentProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.CommentAdd
	params := mux.Vars(r)
	id := params["id"]
	imageId := params["imageId"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(id)
	profile.AddPhotoComment(prof.Id, imageId, prof.Comment)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeleteCommentProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.DeleteElement
	params := mux.Vars(r)
	id := params["id"]
	imageId := params["imageId"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(id)
	profile.DeletePhotoComment(prof.Id, imageId, prof.Index)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) DeleteLikeProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.DeleteElement
	params := mux.Vars(r)
	id := params["id"]
	imageId := params["imageId"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(id)
	profile.DeletePhotoLike(prof.Id, imageId)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddLikeProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.LikeAdd
	params := mux.Vars(r)
	id := params["id"]
	imageId := params["imageId"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(id)
	profile.AddPhotoLike(prof.Id, imageId)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.FollowerAdd
	params := mux.Vars(r)
	id := params["id"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(prof.Id)
	profile.AddFollowings(id)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UnFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.FollowerAdd
	params := mux.Vars(r)
	id := params["id"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(prof.Id)
	profile.DeleteFollower(id)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) BanFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.FollowerAdd
	params := mux.Vars(r)
	id := params["id"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(prof.Id)
	profile.AddBans(id)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) UnBanFollowerProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var prof functionalities.FollowerAdd
	params := mux.Vars(r)
	id := params["id"]
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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

	profile := functionalities.GetProfile(prof.Id)
	profile.UnBans(id)
	// w.Header().Set("Content-type", "application/json")

	if errJson := json.NewEncoder(w).Encode(profile); errJson != nil {
		http.Error(w, errJson.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) ProfileInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	err := validateUserByUsernameID(name, id)
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
	w.Header().Set("content-type", "application/json")
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := creds.returnID()
	if errJson := json.NewEncoder(w).Encode(userID); errJson != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) AddSeen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var prof functionalities.PhotoAdd
	json.NewDecoder(r.Body).Decode(&prof)
	err := validateUserByUsernameID(prof.Username, prof.Id)
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
	profile := functionalities.GetProfile(prof.Id)
	profile.AddAlreadySeen(prof.IdImage)
	if errJson := json.NewEncoder(w).Encode("OK!"); errJson != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (rt *_router) Welcome(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	err := validateUserByUsernameID(name, id)
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
		if errJson := json.NewEncoder(w).Encode("/login"); errJson != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

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
