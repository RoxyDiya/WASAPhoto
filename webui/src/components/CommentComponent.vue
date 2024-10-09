<script>
export default {
    data:()=>{
        return {
            profiles: null,
            errormsg:null,
            // idUser:null,
            // follow:null,
        }
    },
	props: ["commentData"],
    methods:{
        async deleteComment(idUser,IdImage,indexComment){
            try {
                await this.$axios({
                    method: "delete",
                    url: `/profile/${idUser}/comment/${IdImage}?index=${indexComment}`,
                    headers: {
                        "Content-Type": "application/json",
                        "Token": localStorage.getItem("Token"),
                    }
                });
                this.$router.go(0);
            } catch (error) {
                this.errormsg = error
            }
        }
    },
}
</script>

<template>
    <table id="table">
        <tr>
            <td>Username</td>
            <td>Comment</td>
            <td>Time</td>
            <td></td>
        </tr>
        <tr v-for="comment in commentData" :key="comment">
            <td>
                <ProfileImageComponent :userNameF="comment.Username" :imageUrl="comment.ProfilePictureLocation" ></ProfileImageComponent>
            </td>
            <td>
                {{comment.Content}}
            </td>
            <td>
                {{comment.Time}}
            </td>
            <td v-if="comment.IsMine" @click="deleteComment(comment.idUser,comment.IdImage,comment.indexComment)">
                <svg class="feather" id="color-like" style="width: 30; height: 30; color: red;"><use href="/feather-sprite-v4.29.0.svg#trash" /></svg>
            </td>
        </tr>
    </table>
    
</template>

<style>
.row{
 display: flex;
}
td{
    font-size: 15px;
}
#table {
 font-family: Arial, Helvetica, sans-serif;
 border-collapse: collapse;
 width: 100%;
}
#table td, #table th {
 border: 1px solid #ddd;
 padding: 20px;
}
#table .title{
  font-size: 25px;
  background-color: #999999;
  text-align: center;
}
#table tr:nth-child(even){background-color: #f2f2f2;}
#table tr:hover {background-color: #ddd;}
#table th {
 padding-top: 12px;
 padding-bottom: 12px;
 text-align: left;
 background-color: #04AA6D;
 color: white;
}

#table a{
  color: #000;
}
</style>