<script>
export default {
    data:()=>{
        return {
            profiles: null,
            errormsg:null,
        }
    },
	props: ["imageUrl", "userNameF"],
}
</script>

<template>
    <div class="profile" v-if="imageUrl">
        <a :href="'/profile/'+userNameF">
            <img v-bind:src="imageUrl" alt="Profile Picture">
            <h4>{{userNameF}}</h4>
        </a>
    </div>
</template>

<style>
.profile img{
    width: 30%;
    height: auto;
    border-radius: 100%;
}
.profile a{
    /* position: fixed; */
    margin-right: 0;
    padding: 0px;
    display:flex;
    color: black;
    text-decoration: none;
    width: 200px;
}
</style>