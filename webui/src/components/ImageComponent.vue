<script>
import $ from 'jquery'
export default {
    data:()=>{
        return{
            loading:false,
            errormsg:null,
            imageContent: null,
            comments:[],
            commentContent: null,
            username:null,
            profilePicture:null,
        }
    },
	props: ["imageComp","idUser"],
    methods:{
        async basicProfile(){
            let profile = await this.$axios.get(`/profile/${this.idUser}/ultra`)
            profile = profile.data;
            this.username = profile.Username;
            this.profilePicture = profile.ProfilePicture;
        },
        async info(){
            this.comments = []
            document.getElementById(`popup-${this.imageComp.IdImage}`).style.display = "flex";
            // let allDrop = document.getElementsByClassName("dropdown");
            // for (var i = 0; i < allDrop.length; i++) {
            //     allDrop[i].style.display = "contents"
            // }
            this.loading = true
            try {
                let comp = await this.$axios.get(`/profile/${this.idUser}/image/${this.imageComp.IdImage}`);
                comp = comp.data;
                this.imageContent = comp;
                
                
                let t = document.getElementById(`title-${this.imageComp.IdImage}`);
                t.innerHTML = "";
                let h1 = document.createElement('h1');
                let ti = document.createTextNode(comp.Text);
                h1.appendChild(ti);
                t.appendChild(h1);
                let d = document.getElementById(`img-${this.imageComp.IdImage}`);
                d.innerHTML = "";
                let img = document.createElement('img');
                img.src = 'http://localhost:3000'+comp.Location;
                d.appendChild(img);
                this.comments = [];
                $(".postId").name = this.imageComp.IdImage
                for (let comment in comp.Comments){
                    comment = comp.Comments[comment]
                    let comm = await this.$axios.get(`/profile/${comment.UserIdComment}/ultra`)
                    comm = comm.data;
                    let date = startTime(comment.Time)
                    //console.log(date);
                    comment.Time = date
                    comment.Username = comm.Username
                    this.comments.push(comment)
                }
            } catch (error) {
                this.errormsg = error;
            }
            this.loading = false
        },
        close(){
            $(".popup").fadeOut();
            let t = document.getElementById(`title-${this.imageComp.IdImage}`);
            let d = document.getElementById(`img-${this.imageComp.IdImage}`);
            t.innerHTML = "";
            d.innerHTML = "";
            this.comments = []

        },
        async commentPost(idUser,imageComp){
            this.loading = true;
            let dataPost = {
                'comment':this.commentContent
            }
            await this.$axios.post(`/profile/${idUser}/comment/${imageComp.IdImage}`,dataPost,{headers:{"Token":localStorage.getItem("Token")}});
            document.getElementById(`comment-${imageComp.IdImage}`).value = "";
            this.commentContent = null
            this.info();
            this.loading = false;
        },
    },
	mounted() {
		this.basicProfile()
	}
}
function checkTime(i) {
    if (i < 10) {
        i = "0" + i;
    }
    return i;
}

function startTime(date) {
    var today = new Date(date*1000);
    var h = today.getHours();
    var m = today.getMinutes();
    var s = today.getSeconds();
    var d = today.getDate();
    var mo = today.getMonth()+1;
    var y = today.getFullYear();
    // add a zero in front of numbers<10
    m = checkTime(m);
    s = checkTime(s);
    return `${h}:${m}:${s}-${d}/${mo}/${y}`;
    
}
</script>

<template>
	<div v-if="imageComp" class="box">
        <h2>
            {{this.username}}
        </h2>
		<img v-bind:src="'http://localhost:3000'+imageComp.Location" @click="info">
		<div style="display: flex;gap: 20%;">
			<h5>{{imageComp.Text}}</h5>
			<h4>{{imageComp.Comments}}</h4>
			<h4>{{imageComp.Likes}}</h4>
		</div>
        <div v-if="imageComp" class="popup" v-bind:id="'popup-'+imageComp.IdImage">
            <div class="popup-content">
                <LoadingSpinner :loading="this.loading"></LoadingSpinner>
                <div style="display: flex;">
                    <h2>{{this.username}} </h2>
                    <div class="title" v-bind:id="'title-'+imageComp.IdImage">
                    </div>
                </div>
                <div class="desc" v-bind:id="'img-'+imageComp.IdImage">
                </div>
                <div>
                    <CommentComponents :commentData="this.comments"></CommentComponents>
                </div>
                <div>
                    <input type="text" placeholder="Comment" v-bind:id="'comment-'+imageComp.IdImage" v-model="commentContent" >
                    <input type="submit" value="Post" id="postId" @click="commentPost(idUser,imageComp)">
                </div>
    
                <button type="button" id="close" @click="close">Close</button>
            </div>
        </div>
	</div>
</template>

<style>
img{
	object-fit: cover;
	width: 280px;
  	height: 280px;
    border-color: black;
    border: 2px solid;
    box-shadow: 5px;
}
.popup{
    background: rgba(0, 0, 0, 0.6);
    position:absolute;
    /* overflow: auto;
    max-height: 100%; */
    top: 10%;
    left: 18%;
    right: 5%;
    /* right: 5%; */
    display: none;
    justify-content: center;
    align-items: center;
}

.popup-content{
    height: 100%;
    width: 100%;
    background: #fff;
    padding: 2px;
    border-radius: 5px;
    position: relative;
}


.popup-content #close{
  position: relative;
  bottom: 0px;
  right: 0px;
}
.popup-content #post{
  position: absolute;
  bottom: 0px;
  left: 0px;
}

.popup-content-large{
    height: 450px;
    width: 800px;
    background: #fff;
    padding: 20px;
    border-radius: 5px;
    position: relative;
}


.popup-content-large #close{
  position: absolute;
  bottom: 0px;
  right: 0px;
}
.popup-content-large #post{
  position: absolute;
  bottom: 0px;
  left: 0px;
}

.popup-close{
    position: absolute;
    top: -15px;
    right: -15px;
    background: #fff;
    height: 20px;
    width: 20px;
    border-radius: 50%;
    box-shadow: 6px 6px 29px -4px rgba(0,0,0,0.75);
    cursor: pointer;
}
</style>
