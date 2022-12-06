package api

import (
	"math/rand"
	"time"
)

type StreamContent struct {
	Username string
	Id       string
	IdImage  string
	Location string
	Time     string
	Comments int
	Likes    int
	Text     string
}

func (p *Profile) GetNewStream() []StreamContent {
	var newStream []StreamContent
	for i := 0; i < p.GetLenUserFollowings(); i++ {
		user := profiles[p.Followings[i].IdUser]
		for j := 0; j < user.GetLenImages(); j++ {
			if _, ok := p.AlreadySeen[user.Images[j].IdImage]; !ok {
				basic := user.Images[j].getBasicImage()
				stream := StreamContent{
					Username: user.Username,
					Id:       user.Id,
					IdImage:  basic.IdImage,
					Location: basic.Location,
					Time:     basic.Time,
					Comments: basic.Comments,
					Likes:    basic.Likes,
					Text:     basic.Text,
				}
				newStream = append(newStream, stream)
			}
		}
	}
	// shuffle
	for i := range newStream {
		j := rand.Intn(len(newStream))
		newStream[i], newStream[j] = newStream[j], newStream[i]
	}
	return newStream
}

func (p *Profile) AddAlreadySeen(imageId string) {
	p.AlreadySeen[imageId] = time.Now().String()
	// filter := bson.D{{Key: "id", Value: p.Id}}
	// update := bson.D{{Key: "$set", Value: bson.D{{Key: "alreadyseen", Value: p.AlreadySeen}}}}
	// go mongodb.CollectionProfiles.UpdateOne(mongodb.Ctx, filter, update)
}
