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
            id:null,
            profilePicture:null,
            like: false,
            myProfileId: null,
            NComments: 0,
            NLikes: 0,
            base_URL: null,
            ImageLocation: null,
        }
    },
	props: ["imageComp","idUser"],
    methods:{
        async basicProfile(){
            this.loading = true
            this.token = localStorage.getItem("Token")
            let profile = await this.$axios.get(`/profile/${this.idUser}/ultra`,{headers:{"Token":this.token}})
            profile = profile.data;
            this.username = profile.Username;
            this.id = profile.Id;
            this.profilePicture = profile.ProfilePictureLocation === "" ? __API_URL__+'/images/icon_standard.png' : __API_URL__+profile.ProfilePictureLocation;
            let myProfile = await this.$axios.get(`/profile`,{headers:{"Token":this.token}})
            this.myProfileId = myProfile.data.Id
            this.NComments = this.imageComp.Comments
            this.NLikes = this.imageComp.Likes
            this.ImageLocation = __API_URL__ + this.imageComp.Location
            this.loading = false
        },
        async info(){
            this.token = localStorage.getItem("Token")
            this.comments = []
            document.getElementById(`popup-${this.imageComp.IdImage}`).style.display = "flex";
            // let allDrop = document.getElementsByClassName("dropdown");
            // for (var i = 0; i < allDrop.length; i++) {
            //     allDrop[i].style.display = "contents"
            // }
            this.loading = true
            try {
                let comp = await this.$axios.get(`/profile/${this.idUser}/image/${this.imageComp.IdImage}`,{headers:{"Token":this.token}});
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
                img.src = __API_URL__+comp.Location;
                img.classList.add("box")
                d.appendChild(img);
                this.comments = [];
                $(".postId").name = this.imageComp.IdImage
                let i = 0
                for (let comment in comp.Comments){
                    comment = comp.Comments[comment]
                    
                    let isBan = await this.$axios.get(`/ban/${comment.UserIdComment}`,{headers:{"Token":this.token}})
                    if (!isBan.data){
                        let comm = await this.$axios.get(`/profile/${comment.UserIdComment}/ultra`,{headers:{"Token":this.token}})
                        comm = comm.data;
                        let date = startTime(comment.Time)
                        //console.log(date);
                        comment.Time = date;
                        comment.Username = comm.Username;
                        comment.ProfilePictureLocation = comm.ProfilePictureLocation === "" ? __API_URL__+'/images/icon_standard.png' : __API_URL__+comm.ProfilePictureLocation;
                        let ismine = false
                        if (comment.UserIdComment === this.myProfileId){
                            ismine = true
                        }
                        comment.IsMine = ismine
                        comment.idUser = this.idUser
                        comment.IdImage = this.imageComp.IdImage
                        comment.indexComment = i
                        this.comments.push(comment)
                    }
                    i++;
                }
                for (let like in this.imageContent.Likes){
                    like = this.imageContent.Likes[like]
                    if (like.UserIdLike == this.myProfileId){
                        // document.getElementById("color-like").style.color = "red"
                        document.getElementById(`color-like-inner-${this.imageComp.IdImage}`).style.color = "red"
                        this.like = true
                    }
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
            this.info();

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
            this.NComments += 1
        },
        async likePut(idUser,imageComp){
            this.loading = true;
            if (this.like){
                await this.$axios({
                    method: "delete",
                    url: `/profile/${idUser}/like/${imageComp.IdImage}`,
                    headers: {
                        "Content-Type": "application/json",
                        "Token": this.token,
                    }
                });
                document.getElementById(`color-like-inner-${this.imageComp.IdImage}`).style.color = "black"
                this.like = false;
                this.NLikes -= 1;
            }
            else{
                await this.$axios({
                    method: "put",
                    url: `/profile/${idUser}/like/${imageComp.IdImage}`,
                    headers: {
                        "Content-Type": "application/json",
                        "Token": this.token,
                    }
                });
                this.NLikes += 1;
            }
            this.info();
            this.loading = false;
        },
        async deleteImage(){
            try {
                await this.$axios({
                    method: "delete",
                    url: `/profile/${this.idUser}/image/${this.imageComp.IdImage}`,
                    headers: {
                        "Content-Type": "application/json",
                        "Token": this.token,
                    }
                });
            } catch (error) {
                console.log(error);
            }
            location.reload();
        },
        openDropDown(imageId){
            document.getElementById("myDropdown-"+imageId).classList.toggle("show");
        },
        async profilePictureUpdate(imageLocation){
            var formData = new FormData();
			formData.append('profilePicture',imageLocation)
            await this.$axios.post(`/addphoto`,formData,{headers:{"Token":this.token}});
            this.info();
        }
    },
	async mounted() {
		await this.basicProfile()
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
window.onclick = function(event) {
  if (!event.target.matches('.dropbtn')) {
    var dropdowns = document.getElementsByClassName("dropdown-content");
    var i;
    for (i = 0; i < dropdowns.length; i++) {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('show')) {
        openDropdown.classList.remove('show');
      }
    }
  }
}
</script>

<template>
	<div v-if="imageComp">
        <ProfileImageComponent :userNameF="this.username" :imageUrl="this.profilePicture" ></ProfileImageComponent>
        <LoadingSpinner v-if="this.loading"></LoadingSpinner>
		<img class="box" v-bind:src="this.ImageLocation" @click="info">
		<div style="display: flex;gap: 20%;">
			<h5>{{imageComp.Text}}
            </h5>
			<h4>
                <svg class="feather" id="color-like" style="width: 30; height: 30;"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
                {{this.NLikes}}
            </h4>
			<h4>
                <svg class="feather" style="width: 30; height: 30;"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
                {{this.NComments}}
            </h4>
		</div>
        <div v-if="imageComp" class="popup" v-bind:id="'popup-'+imageComp.IdImage">
            <div class="popup-content">
                <LoadingSpinner :loading="this.loading"></LoadingSpinner>
                <ProfileImageComponent :userNameF="this.username" :imageUrl="this.profilePicture" ></ProfileImageComponent>
                <div style="display: flex; margin-left: 35%;">
                    <div class="desc" v-bind:id="'img-'+imageComp.IdImage">
                    </div>
                    <div class="dropdown">
                        <div style="margin-left: 1%;" class="dropbtn" @click="openDropDown(imageComp.IdImage)">⋮</div>
                        <div v-bind:id="'myDropdown-'+imageComp.IdImage" class="dropdown-content">
                            <a href="#" v-if="this.id === this.myProfileId" @click="deleteImage()">Delete</a>
                            <a href="#" @click="profilePictureUpdate(imageComp.Location)">Update Profile Picture</a>
                        </div>
                    </div>
                </div>
                
                <div style="display: flex; position: relative; left: 35%; gap: 10%;">
                    <h4>
                        <svg class="feather" v-bind:id="'color-like-inner-'+imageComp.IdImage" style="width: 30; height: 30;" @click="likePut(idUser,imageComp)"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
                        {{this.NLikes}}
                    </h4>
                    <h4>
                        <svg class="feather" style="width: 30; height: 30;"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
                        {{this.NComments}}
                    </h4>
                </div>
                <div class="title" v-bind:id="'title-'+imageComp.IdImage"></div>
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
.box{
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
    border: 2px solid black;
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