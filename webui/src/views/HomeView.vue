<script>

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			token: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.token = localStorage.getItem("Token")
				if (this.token === null){
					location.replace("/login")
				}
				let response = await this.$axios.get("/",{headers:{"Token":this.token}});
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			if (this.some_data != null){
				for (let i = 0; i < this.some_data.length; i++){
					console.log(this.some_data[i]);
					let dataPost = {
						idimage: this.some_data[i].IdImage
					};
					try {
						await this.$axios.post("/",dataPost,{headers:{"Token":this.token}})
					} catch (error) {
						this.errormsg = error.toString();
					}
				};
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<div v-if="this.some_data" class="grid-container">
			<div v-for="item in this.some_data" :key="item">
				<div class="grid-item">
					<ImageComponent v-if="item" :imageComp="item" :idUser="item.Id" />
				</div>
			</div>
		</div>
	</div>
</template>

<style>
</style>