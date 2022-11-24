package functionalities

import (
	"log"
	"time"
	"wasa-photo/service/api/errors"

	"github.com/google/uuid"
)

var profiles = map[string]Profile{}

type UserFollow struct {
	IdUser string
	Time   string
}

type ProfileUpdate struct {
	NewUsername string `json:"username"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
}

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

type PhotoAdd struct {
	IdImage                string `json:"idimage"`
	ProfilePictureLocation string `json:"profilepicturelocation"`
	Text                   string `json:"text"`
}

type CommentAdd struct {
	Comment string `json:"comment"`
}

type BasicProfile struct {
	Username               string
	Id                     string
	ProfilePictureLocation string
	Followers              int
	Followings             int
	Images                 []BasicImage
}

type BasicImage struct {
	IdImage  string
	Location string
	Time     string
	Likes    int
	Comments int
	Text     string
}

type UltraBasicProfile struct {
	Id                     string
	Username               string
	ProfilePictureLocation string
}

func NewProfile(username string, id string) error {
	profile := Profile{
		Id:       id,
		Username: username,
	}
	profiles[id] = profile
	err := newProfileDB(profile)
	if err != nil {
		return err
	}
	return nil
}

func AddProfile(id string, profile Profile) {
	profiles[id] = profile
}

func newProfileDB(profile Profile) error {
	profile.Followers = []UserFollow{}
	profile.Followings = []UserFollow{}
	profile.Bans = []UserFollow{}
	profile.Images = []Image{}
	profile.AlreadySeen = map[string]string{}
	// profileCollection := database.AppDatabaseMongo.GetProfilesCollection()
	// database.AppDatabaseMongo.InsertOne(profileCollection)
	// _, err := mongodb.CollectionProfiles.InsertOne(mongodb.Ctx, profile)
	//x.followers = new([]UserFollow)
	return nil
}

func GetProfile(id string) Profile {
	return profiles[id]
}

func (p Profile) GetUserProfile() Profile {
	return p
}

func (p Profile) GetUserFollowers() []UserFollow {
	return p.Followers
}

func (p Profile) GetBasicUserFollowers() []UltraBasicProfile {
	followers := p.GetUserFollowers()
	var basicFollowers []UltraBasicProfile
	for _, id := range followers {
		basic := GetUltraBasicProfile(id.IdUser)
		basicFollowers = append(basicFollowers, basic)
	}
	return basicFollowers
}

func (p Profile) GetBasicUserFollowings() []UltraBasicProfile {
	followings := p.GetUserFollowings()
	var basicFollowings []UltraBasicProfile
	for _, id := range followings {
		basic := GetUltraBasicProfile(id.IdUser)
		basicFollowings = append(basicFollowings, basic)
	}
	return basicFollowings
}

func GetUltraBasicProfile(id string) UltraBasicProfile {
	profile := UltraBasicProfile{
		Id:                     profiles[id].Id,
		Username:               profiles[id].Username,
		ProfilePictureLocation: GetPictureLocationById(id),
	}
	return profile
}

func (p Profile) GetUserFollowings() []UserFollow {
	return p.Followings
}

func (p Profile) GetLenUserFollowers() int {
	return len(p.Followers)
}

func (p Profile) GetLenUserFollowings() int {
	return len(p.Followings)
}

func (p Profile) GetLenImages() int {
	return len(p.Images)
}

func (p Profile) GetUserBans() []UserFollow {
	return p.Bans
}

func (p Profile) GetUserId() string {
	return p.Id
}

// Return the profile picture of the user
// by his Id
func GetPictureLocationById(id string) string {
	profile := profiles[id]
	return profile.ProfilePicture.getPictureLocation()
}

// Set the new username of the profile
func (p *Profile) SetMyUsername(newUsername string) {
	p.Username = newUsername
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: newUsername}}}}
	// profileCollection := database.AppDatabaseMongo.GetProfilesCollection()
}

// Set the name of the profile
func (p *Profile) SetMyName(name string) {
	p.Name = name
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: name}}}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, update)
}

// Set the surname of the profile
func (p *Profile) SetMySurname(surname string) {
	p.Surname = surname
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// update := bson.D{{Key: "$set", Value: bson.D{{Key: "surname", Value: surname}}}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, update)
}

// Add a follower in Profile.Followers
func (p *Profile) AddFollowers(id string) {
	newFollow := UserFollow{
		IdUser: id,
		Time:   time.Now().String(),
	}
	p.Followers = append(p.Followers, newFollow)
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// push := bson.M{"$push": bson.M{"followers": newFollow}}
	// result, err := mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, push)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(result)
}

// Add a following in Profile.Followings
func (p *Profile) AddFollowings(id string) {
	user, ok := profiles[id]
	if !ok {
		return
	}

	newFollow := UserFollow{
		IdUser: id,
		Time:   time.Now().String(),
	}
	if p.FindBanUser(id) || p.FindFollowingUser(id) {
		return
	}
	p.Followings = append(p.Followings, newFollow)
	user.AddFollowers(p.Id)
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// push := bson.M{"$push": bson.M{"followings": newFollow}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, push)
}

// Find if the user is banned. It'll return a boolean.
// true if the user is banned, false if not
func (p Profile) FindBanUser(id string) bool {
	for _, i := range p.Bans {
		if i.IdUser == id {
			return true
		}
	}
	return false
}

// Find if the user is followed. It'll return a boolean.
// true if the user is followed, false if not
func (p Profile) FindFollowerUser(id string) bool {
	for _, i := range p.Followers {
		if i.IdUser == id {
			return true
		}
	}
	return false
}

// Find if the user is followd by the user. It'll return
// a boolean. true if is followed by the user, false if not
func (p Profile) FindFollowingUser(id string) bool {
	for _, i := range p.Followings {
		if i.IdUser == id {
			return true
		}
	}
	return false
}

// Unfollow a user by his id
func (p *Profile) UnFollowers(id string) {
	p.UnfollowingUser(id)
	user := profiles[id]
	user.UnFollowerUser(p.Id)
	profiles[p.Id] = *p
}

func (p *Profile) UnfollowingUser(id string) {
	var followListN []UserFollow
	for _, x := range p.Followings {
		if x.IdUser != id {
			follow := UserFollow{
				IdUser: x.IdUser,
				Time:   x.Time,
			}
			followListN = append(followListN, follow)
		}
	}
	p.Followings = followListN
	profiles[p.Id] = *p
}

func (p *Profile) UnFollowerUser(id string) {
	var followList []UserFollow
	for _, x := range p.Followings {
		if x.IdUser != id {
			follow := UserFollow{
				IdUser: x.IdUser,
				Time:   x.Time,
			}
			followList = append(followList, follow)
		}
	}
	p.Followers = followList
	profiles[p.Id] = *p
}

// Add a banned user in Profile.Bans
func (p *Profile) AddBans(id string) {
	user, ok := profiles[id]
	if !ok {
		return
	}
	newBan := UserFollow{
		IdUser: id,
		Time:   time.Now().String(),
	}
	if p.FindBanUser(id) {
		return
	}
	p.DeleteFollower(id)
	p.Bans = append(p.Followings, newBan)
	user.DeleteFollower(p.Id)
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// push := bson.M{"$push": bson.M{"bans": newBan}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, push)
}

// Delete a user from Profile.Followers and
// Profile.Followings
func (p *Profile) DeleteFollower(id string) {
	var followList []UserFollow
	for _, x := range p.Followers {
		if x.IdUser != id {
			follow := UserFollow{
				IdUser: x.IdUser,
				Time:   x.Time,
			}
			followList = append(followList, follow)
		}
	}
	p.Followers = followList
	var followListN []UserFollow
	for _, x := range p.Followings {
		if x.IdUser != id {
			follow := UserFollow{
				IdUser: x.IdUser,
				Time:   x.Time,
			}
			followListN = append(followListN, follow)
		}
	}
	p.Followings = followListN
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// update1 := bson.D{{Key: "$set", Value: bson.D{{Key: "followings", Value: followListN}}}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, update1)
	// update2 := bson.D{{Key: "$set", Value: bson.D{{Key: "followers", Value: followList}}}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, update2)
}

// Delete a user from Profile.Bans
func (p *Profile) UnBans(id string) {
	var bans []UserFollow
	for _, x := range p.Bans {
		if x.IdUser != id {
			ban := UserFollow{
				IdUser: x.IdUser,
				Time:   x.Time,
			}
			bans = append(bans, ban)
		}
	}
	p.Bans = bans
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// update := bson.D{{Key: "$set", Value: bson.D{{Key: "bans", Value: bans}}}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, update)
}

// Update the profile picture
func (p *Profile) UpdateProfilePicture(newImgLocation string) {
	oldPicture := p.ProfilePicture
	oldPicture.updatePicture(newImgLocation)
	p.ProfilePicture = oldPicture
	profiles[p.Id] = *p
}

// Update the profile info, such as name and surname
func (p *Profile) UpdateProfileInfo(newProfileInfo ProfileUpdate) {
	if newProfileInfo.Name != "" {
		p.SetMyName(newProfileInfo.Name)
	}
	if newProfileInfo.Surname != "" {
		p.SetMySurname(newProfileInfo.Surname)
	}
	profiles[p.Id] = *p
}

// Get the basic profile info of a profile.
// It'll return a BasicProfile struct
func GetProfileBasicInfo(id string) BasicProfile {
	profile := profiles[id]
	basic := BasicProfile{
		Username:               profile.Username,
		Id:                     profile.Id,
		ProfilePictureLocation: profile.ProfilePicture.getPictureLocation(),
		Followers:              profile.GetLenUserFollowers(),
		Followings:             profile.GetLenUserFollowings(),
	}
	for i := 0; i < len(profile.Images); i++ {
		image := profile.Images[i].getBasicImage()
		basic.Images = append(basic.Images, image)
	}
	return basic
}

// Photos

// Add a photo in the profile.
// It'll create a new id and set the position
// of the image
func (p *Profile) AddPhoto(text string) string {
	newID := uuid.NewString()
	image := Image{
		IdImage:  newID,
		Location: "/images/" + newID + ".png",
		Text:     text,
		Time:     time.Now().String(),
		Comments: []Comment{},
		Likes:    []Like{},
	}
	p.Images = append(p.Images, image)
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// push := bson.M{"$push": bson.M{"images": image}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, push)
	return image.IdImage
}

// Delete a photo from Profile.Images
func (p *Profile) DeletePhoto(imageId string) {
	var images []Image
	for _, image := range p.Images {
		if image.IdImage != imageId {
			images = append(images, image)
		}
	}
	p.Images = images
	profiles[p.Id] = *p
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// push := bson.M{"$set": bson.M{"images": p.Images}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, push)
}

// Add a comment in the image struct Image.Comments
func (p *Profile) AddPhotoComment(usernameIdComment string, imageId string, content string) {
	images := p.Images
	for i := 0; i < len(images); i++ {
		if images[i].IdImage == imageId {
			images[i].addComment(usernameIdComment, content)
			// filter := bson.M{"id": p.Id, "images.idimage": imageId}
			// push := bson.M{"$set": bson.M{"images.$.comments": images[i].Comments}}
			// go mongodb.CollectionProfiles.FindOneAndUpdate(mongodb.Ctx, filter, push)
		}
	}
	p.Images = images
	profiles[p.Id] = *p
}

// Delete a comment in the image struct Image.Comments
func (p *Profile) DeletePhotoComment(usernameIdComment string, imageId string, index int) {
	images := p.Images
	log.Println(usernameIdComment, imageId, index)
	for i := 0; i < len(images); i++ {
		if images[i].IdImage == imageId {
			image := images[i]
			comments := image.deleteComment(usernameIdComment, index)
			images[i].Comments = comments

			// filter := bson.M{"id": p.Id, "images.idimage": imageId}
			// push := bson.M{"$set": bson.M{"images.$.comments": images[i].Comments}}
			// go mongodb.CollectionProfiles.FindOneAndUpdate(mongodb.Ctx, filter, push)
		}
	}
	p.Images = images
	profiles[p.Id] = *p
}

// Delete a like in the image struct Image.Likes
func (p *Profile) DeletePhotoLike(usernameIdLike string, imageId string) {
	images := p.Images
	for i := 0; i < len(images); i++ {
		if images[i].IdImage == imageId {
			image := images[i]
			likes := image.deleteLike(usernameIdLike)
			images[i].Likes = likes
			// filter := bson.M{"id": p.Id, "images.idimage": imageId}
			// push := bson.M{"$set": bson.M{"images.$.likes": images[i].Likes}}
			// go mongodb.CollectionProfiles.FindOneAndUpdate(mongodb.Ctx, filter, push)
		}
	}
	p.Images = images
	profiles[p.Id] = *p
}

// Add a like in the image struct Image.Likes
func (p *Profile) AddPhotoLike(usernameIdLike string, imageId string) {
	images := p.Images
	for i := 0; i < len(images); i++ {
		if images[i].IdImage == imageId {
			images[i].addLike(usernameIdLike)
			// filter := bson.M{"id": p.Id, "images.idimage": imageId}
			// push := bson.M{"$set": bson.M{"images.$.likes": images[i].Likes}}
			// go mongodb.CollectionProfiles.FindOneAndUpdate(mongodb.Ctx, filter, push)
		}
	}
	p.Images = images
	profiles[p.Id] = *p
}

// Get image information.
// Return an error if the image is not found
func (p Profile) GetImageInfo(imageId string) (Image, error) {
	image := p.getImageById(imageId)
	if image.getLocation() == "" {
		return Image{}, errors.NewErrStatus("Image Not Found")
	}
	return image, nil
}
