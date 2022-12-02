import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/profile', component: ProfileView},
		{path: '/login', component: LoginView},
		{path: '/:username', component: HomeView},
		{path: '/:username/:imageid', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router
