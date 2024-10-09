<script>
import { RouterLink, RouterView } from 'vue-router'
export default{
    data:()=>{
        return {
            token: null,
            some_data: null,
            errormsg: null,
			dataBool: false,
        }
    },
    methods:{
        async info(){
			this.token = localStorage.getItem("Token")
			if (this.token === null && window.location.pathname != "/login"){
				location.replace("/login")
				this.dataBool = false
			}
			else if (this.token === null){
				this.dataBool = false
			}
			else{
				try {
					let response = await this.$axios.get("/profile",{headers:{"Token":this.token}});
					this.some_data = response.data;
					this.dataBool = true
				} catch (e) {
					localStorage.removeItem("Token")
					location.replace("/login")
				}
			}
        },
		async rem(){
			this.token = null
			localStorage.removeItem("Token")
			await this.info()
		}
    },
    async mounted(){
        await this.info()
    }

}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6 text-uppercase" href="/">Wasa - Photo</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column" v-if="this.dataBool">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li>
						<li class="nav-item" v-if="this.dataBool">
							<RouterLink to="/profile" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Profile
							</RouterLink>
						</li>
						<li class="nav-item" v-if="this.dataBool">
							<RouterLink to="/add" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
								Add Photo
							</RouterLink>
						</li>
						<li class="nav-item" v-if="this.dataBool">
							<RouterLink to="/search" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
								Search
							</RouterLink>
						</li>
					</ul>
			
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Secondary menu</span>
					</h6>
					<ul class="nav flex-column" v-if="!this.dataBool">
						<li class="nav-item">
							<RouterLink to="/login" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-in"/></svg>
								Log In or Sign In
							</RouterLink>
						</li>
					</ul>
					<ul class="nav flex-column" v-if="this.dataBool">
						<li class="nav-item">
							<div @click="rem()" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
								Log Out
							</div>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
</style>