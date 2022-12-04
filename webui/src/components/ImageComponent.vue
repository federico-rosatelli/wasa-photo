<script>
import $ from 'jquery'
export default {
	props: ["imageComp","idUser"],
    methods:{
        async info(idUser,idImage){
            document.getElementById("popup").style.display = "flex";
            // let allDrop = document.getElementsByClassName("dropdown");
            // for (var i = 0; i < allDrop.length; i++) {
            //     allDrop[i].style.display = "contents"
            // }
            
            let comp = await this.$axios.get("/profile/"+idUser+"/image/"+idImage);
            comp = comp.data
            let t = document.getElementById("title");
            let h1 = document.createElement('h1');
            let ti = document.createTextNode(comp.Text);
            h1.appendChild(ti);
            t.appendChild(h1);
            let d = document.getElementById("img");
            let img = document.createElement('img');
            img.src = 'http://localhost:3000'+comp.Location;
            d.appendChild(img);
        },
        close(){
            $(".popup").fadeOut();
            let t = document.getElementById("title");
            let d = document.getElementById("img");
            t.innerHTML = "";
            d.innerHTML = "";
        }
    }
}
</script>

<template>
	<div v-if="imageComp" class="box">
		<img v-bind:src="'http://localhost:3000'+imageComp.Location" @click="info(idUser,imageComp.IdImage)">
		<div style="display: flex;gap: 20%;">
			<h5>{{imageComp.Text}}</h5>
			<h4>{{imageComp.Comments}}</h4>
			<h4>{{imageComp.Likes}}</h4>
		</div>
	</div>
	<div v-if="!imageComp"><slot /></div>
    <div class="popup" id="popup">
        <div class="popup-content">
            <div class="title" id="title">
            </div>
            <div class="desc" id="img">
            </div>
            <button type="button" id="close" @click="close">Close</button>
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
}
.popup{
    background: rgba(0, 0, 0, 0.6);
    width: 100%;
    height: 100%;
    position: fixed;
    top: 0;
    left: 20%;
    right: 20%;
    display: none;
    justify-content: center;
    align-items: center;
}

.popup-content{
    height: 800px;
    width: 800px;
    background: #fff;
    padding: 2px;
    border-radius: 5px;
    position: relative;
}


.popup-content #close{
  position: absolute;
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
