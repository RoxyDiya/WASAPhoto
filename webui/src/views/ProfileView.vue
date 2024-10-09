<script>
export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            token: null,
            myName: null,
            mySurname: null,
            usernameChange:"",
            myPage: true,
            alreadyFollow: false,
            userId: "",
            profiles: [],
            modal: false,
            isBan: false,
            base_URL: null,
        };
    },
    methods: {
        async refresh() {
            let userName = this.$route.params.username;
            this.loading = true;
            this.errormsg = null;
            this.token = localStorage.getItem("Token");
            if (this.token === null){
                console.log("ASAAS");
                location.replace("/login")
            }
            try {
                if (userName != null) {
                    let users = await this.$axios.get("/search?query=" + userName + "&precise=1", { headers: { "Token": this.token } });
                    let myProfile = await this.$axios.get(`/profile`, { headers: { "Token": this.token } });
                    if (users.data[0].Id === undefined){
                        location.replace("/")
                    }
                    this.userId = "/" + users.data[0].Id;
                    if (myProfile.data.Id === users.data[0].Id) {
                        this.myPage = true;
                        this.some_data = myProfile.data;
                    }
                    else {
                        this.myPage = false;
                        let response = await this.$axios.get("/profile" + this.userId, { headers: { "Token": this.token } });
                        let myFollowers = await this.$axios.get(`/profile/${myProfile.data.Id}/followings`, { headers: { "Token": this.token } });
                        let isBan = await this.$axios.get("/ban" + this.userId, { headers: { "Token": this.token } });
                        if (isBan.data == true){
                            this.isBan = true
                        }
                        myFollowers = myFollowers.data;
                        this.alreadyFollow = false;
                        if (myFollowers != null) {
                            for (let f = 0; f < myFollowers.length; f++) {
                                console.log(myFollowers[f].Id);
                                if (myFollowers[f].Id === users.data[0].Id) {
                                    this.alreadyFollow = true;
                                }
                            }
                        }
                        this.some_data = response.data;
                    }
                }
                else {
                    let myProfile = await this.$axios.get(`/profile`, { headers: { "Token": this.token } });
                    this.myPage = true;
                    this.some_data = myProfile.data;

                    this.userId = "/" + this.some_data.Id;
                }
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.some_data.ProfilePictureLocation = this.some_data.ProfilePictureLocation === "" ? __API_URL__+'/images/icon_standard.png' : __API_URL__+this.some_data.ProfilePictureLocation
        },
        async change() {
            this.loading = true;
            this.errormsg = null;
            try {
                let postData = {
                    name: this.myName,
                    surname: this.mySurname,
                    username: this.usernameChange,
                };
                let response = await this.$axios.post("/profile", postData, { headers: { "Token": this.token } });
                this.some_data = response.data;
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.usernameChange = "";
            location.replace("/profile")
        },
        async follow() {
            this.loading = true;
            this.errormsg = null;
            this.token = localStorage.getItem("Token");
            try {
                let response = null;
                if (!this.alreadyFollow) {
                    console.log(this.userId);
                    //response = await this.$axios.put("/profile"+this.userId,{headers:{"Token":this.token}});
                    response = await this.$axios({
                        method: "put",
                        url: `/profile${this.userId}`,
                        headers: {
                            "Content-Type": "application/json",
                            "Token": this.token,
                        }
                    });
                    this.alreadyFollow = true;
                }
                else {
                    response = await this.$axios({
                        method: "delete",
                        url: `/profile${this.userId}`,
                        headers: {
                            "Content-Type": "application/json",
                            "Token": this.token,
                        }
                    });
                    this.alreadyFollow = false;
                }
                this.some_data = response.data;
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
            this.refresh();
        },
        async banUser(){
            try {
                if (!this.isBan){
                    await this.$axios({
                        method: "put",
                        url: `/ban${this.userId}`,
                        headers: {
                            "Content-Type": "application/json",
                            "Token": this.token,
                        }
                    });
                }
                else{
                    await this.$axios({
                        method: "delete",
                        url: `/ban${this.userId}`,
                        headers: {
                            "Content-Type": "application/json",
                            "Token": this.token,
                        }
                    });
                    this.isBan = false
                }
            } catch (error) {
                this.errormsg = error
            }
            console.log(this.isBan);
            this.refresh()
            
        },
        async openModalWR(typeFollow) {
            try {
                let dataFollow = await this.profilesData(typeFollow);
                this.profiles = dataFollow;
            }
            catch (error) {
                this.errormsg = error.toString();
            }
            document.getElementById(`myFollow-${typeFollow}`).style.display = "flex";
        },
        openModalS(){
            document.getElementById(`settings`).style.display = "flex";
        },
        async profilesData(followCase) {
            try {
                var datap = await this.$axios.get(`/profile${this.userId}/${followCase}`, { headers: { "Token": this.token } });
                let data = datap.data
                if (!data){
                    return []
                }
                for (let i = 0; i < data.length; i++) {
                    data[i].ProfilePictureLocation = data[i].ProfilePictureLocation === "" ? __API_URL__+"/images/icon_standard.png" : __API_URL__+data[i].ProfilePictureLocation
                    
                }
                return data;
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            return null;
        },
        hideModal(id) {
            document.getElementById(id).style.display = "none";
            this.refresh();
        }
    },
    mounted() {
        this.refresh();
    },
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Profile</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<div v-if="this.some_data" style="display: flex;">
			<ProfileImageComponent :userNameF="this.some_data.Username" :imageUrl="this.some_data.ProfilePictureLocation" ></ProfileImageComponent>
            <div v-if="myPage">
                <svg class="feather" @click="openModalS()"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
                <div class="modal-mask" style="display: none;" id="settings">
					<div class="modal-wrapper">
						<div class="modal-container">

							<div class="modal-header">
								<slot name="header">
									<h3>
										Settings
									</h3>
								</slot>
							</div>  

							<div class="modal-body">
								<slot name="body">
									<div style="display: flex; gap: 10%;">
                                        <input type="text" placeholder="Name" v-model="myName" >
                                        <input type="text" placeholder="Surname" v-model="mySurname" >
                                        <input type="submit" value="Change" @click="change">
                                    </div>
                                    <div style="display: flex; gap: 30%;">
                                        <input type="text" placeholder="Username" v-model="usernameChange">
                                        <input type="submit" value="Change Username" @click="change">
                                    </div>
								</slot>
							</div>

							<div class="modal-footer">
								
									<button class="modal-default-button" @click="hideModal('settings')">
										OK
									</button>
								
							</div>
						</div>
					</div>
				</div>
            </div>
        </div>
        <div v-if="!myPage">
            <ul class="wrapper">
                <div v-if="!alreadyFollow && !isBan">
                    <li class="icon followProfile" @click="follow">Follow</li>
                </div>
                <div v-if="alreadyFollow && !isBan">
                    <li class="icon followProfile" @click="follow">Unfollow</li>
                </div>
                <div v-if="!isBan">
                    <li class="icon banProfile" @click="banUser">Ban</li>
                </div>
                <div v-else>
                    <li class="icon banProfile" @click="banUser">Unban</li>
                </div>
            </ul>
        </div>
		<div v-if="this.some_data" class="follow-info">		
			<div>
				<h2 @click="openModalWR('followers')">{{this.some_data.Followers}}</h2>
                <span>Followers</span>
				<div class="modal-mask" style="display: none;" id="myFollow-followers">
					<div class="modal-wrapper">
						<div class="modal-container">

							<div class="modal-header">
								<slot name="header">
									<h3>
										Followers
									</h3>
								</slot>
							</div>  

							<div class="modal-body">
								<slot name="body">
									<FollowComponent v-if="this.profiles" :usersFollow="this.profiles"></FollowComponent>
								</slot>
							</div>

							<div class="modal-footer">
								
									<button class="modal-default-button" @click="hideModal('myFollow-followers')">
										OK
									</button>
								
							</div>
						</div>
					</div>
				</div>

			</div>
			<div>
				<h2 @click="openModalWR('followings')">{{this.some_data.Followings}}</h2>
                <span>Followings</span>
				<div class="modal-mask" style="display: none;" id="myFollow-followings">
					<div class="modal-wrapper">
						<div class="modal-container">

							<div class="modal-header">
								<h3>
									Followings
								</h3> 
								<slot name="header">
								</slot>
							</div>  

							<div class="modal-body">
								<slot name="body">
									<FollowComponent v-if="this.profiles" :usersFollow="this.profiles"></FollowComponent>
								</slot>
							</div>

							<div class="modal-footer">
								
								<button class="modal-default-button" @click="hideModal('myFollow-followings')">
									OK
								</button>
								
							</div>
						</div>
					</div>
				</div>
			</div>
            <div v-if="this.some_data && this.some_data.Images">
				<h2 title="Images">{{this.some_data.Images.length}}</h2>
                <span>Images</span>
			</div>
            <div v-else>
                <h2 title="Images">0</h2>
                <span>Images</span>
            </div>
				
		</div>
		<div v-if="this.some_data" style="display: flex; gap: 20%;">
			<h2>{{this.some_data.Name}}</h2>
			<h2>{{this.some_data.Surname}}</h2>
		</div>
		
		<div v-if="this.some_data && !this.isBan" class="grid-container">
			<div v-for="item in this.some_data.Images" :key="item">
				<div class="grid-item">
					<ImageComponent v-if="item" :imageComp="item" :idUser="this.some_data.Id" />
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.grid-container {
  display: grid;
  grid-template-columns: 33% 33% 33%;
  padding: 10px;
}
.grid-item {
  padding: 20px;
  font-size: 30px;
  text-align: center;
  border: 1px solid #000000;
}
</style>