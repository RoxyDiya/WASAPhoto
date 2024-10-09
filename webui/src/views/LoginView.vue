<script>
export default {
    name: "token-login",
    data() {
        return {
            errormsg: null,
            loading: false,
            token_data: null,
            myUsername: ""
        };
    },
    methods: {
        async refresh() {
            this.loading = true;
            this.errormsg = null;
            let dataPost = {
                username: this.myUsername
            };
            try{
                let response = await this.$axios.post("/signin", dataPost);
                this.token_data = response.data;
                localStorage.setItem("Token",this.token_data)
                location.replace("/profile")
            }
            catch(e){
                this.errormsg = "prova "+e.toString();
            }
            this.loading = false;
            
        },
    },
}
</script>

<template>
	<div>
		<div class="login-form">
			<h1>Login Form</h1>
			    <input type="text" placeholder="Username" v-model="myUsername" >
				<input type="submit" value="Log In" @click="refresh">
			<div style="padding:15px">
		</div>
    </div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>