package api

import (
	"log"
	"time"
)

type Like struct {
	UserIdLike string
	Time       int64
}
type Comment struct {
	UserIdComment string
	Time          int64
	Content       string
}

type Image struct {
	IdImage  string
	Location string
	Text     string
	Time     int64
	Likes    []Like
	Comments []Comment
}

type ProfilePicture struct {
	Location string
	Time     int64
}

func (i Image) getLocation() string {
	return i.Location
}

// func (i Image) getLikes() []Like {
// 	return i.Likes
// }

// func (i Image) getComments() []Comment {
// 	return i.Comments
// }

func (p *Profile) getImageById(imageId string) Image {
	images := p.Images
	for i := 0; i < len(images); i++ {
		if images[i].IdImage == imageId {
			return images[i]
		}
	}
	return Image{}
}

func (i Image) getLenLikes() int {
	return len(i.Likes)
}

func (i Image) getLenComments() int {
	return len(i.Comments)
}

func (i *Image) addLike(id string) {
	isin := false
	for _, likes := range i.Likes {
		if likes.UserIdLike == id {
			isin = true
		}
	}
	if !isin {
		like := Like{
			UserIdLike: id,
			Time:       time.Now().Unix(),
		}
		i.Likes = append(i.Likes, like)
	}
}

func (i *Image) deleteComment(usernameIdComment string, index int) []Comment {
	var comments []Comment
	for j := 0; j < len(i.Comments); j++ {
		comment := Comment{
			UserIdComment: i.Comments[j].UserIdComment,
			Time:          i.Comments[j].Time,
			Content:       i.Comments[j].Content,
		}
		if j != index || i.Comments[j].UserIdComment != usernameIdComment {
			log.Println("ENTRA QUA DENTRO")
			comments = append(comments, comment)
		}
	}
	log.Println(len(comments))
	return comments
}

func (i *Image) deleteLike(usernameIdLike string) []Like {
	var likes []Like
	for j := 0; j < len(i.Likes); j++ {
		like := Like{
			UserIdLike: i.Likes[j].UserIdLike,
			Time:       i.Likes[j].Time,
		}
		if i.Likes[j].UserIdLike != usernameIdLike {
			likes = append(likes, like)
		}
	}
	return likes
}

func (i *Image) addComment(id string, content string) {
	comment := Comment{
		UserIdComment: id,
		Time:          time.Now().Unix(),
		Content:       content,
	}
	i.Comments = append(i.Comments, comment)
}

func (pp *ProfilePicture) updatePicture(newLocation string) {
	pp.Location = newLocation
	pp.Time = time.Now().Unix()
}

func (pp *ProfilePicture) getPictureLocation() string {
	return pp.Location
}

func (i *Image) getBasicImage() BasicImage {
	image := BasicImage{
		IdImage:  i.IdImage,
		Location: i.getLocation(),
		Time:     i.Time,
		Likes:    i.getLenLikes(),
		Comments: i.getLenComments(),
		Text:     i.Text,
	}
	return image
}
