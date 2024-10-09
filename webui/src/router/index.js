import {createRouter, createWebHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import AddPhotoView from '../views/AddPhotoView.vue'
import SearchView from '../views/SearchView.vue'

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{path: '/', component: HomeView},
		{path: '/profile', component: ProfileView},
		{path: '/login', component: LoginView},
		{path: '/add', component: AddPhotoView},
		{path: '/search', component: SearchView},
		{path: '/profile/:username', component: ProfileView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
