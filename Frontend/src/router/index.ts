import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import FixturesView from '@/views/FixturesView.vue'
import StandingsView from '@/views/StandingsView.vue'
import FavoritesView from '@/views/FavoritesView.vue'
import TeamDetailView from '@/views/TeamDetailView.vue'
import MatchDetailView from '@/views/MatchDetailView.vue'
import NewsView from '@/views/NewsView.vue'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import ProfileView from '@/views/ProfileView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/fixtures', name: 'fixtures', component: FixturesView },
    { path: '/standings', name: 'standings', component: StandingsView },
    { path: '/favorites', name: 'favorites', component: FavoritesView },
    { path: '/news', name: 'news', component: NewsView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/register', name: 'register', component: RegisterView },
    { path: '/profile', name: 'profile', component: ProfileView },
    { path: '/teams/:id', name: 'team-detail', component: TeamDetailView },
    { path: '/matches/:id', name: 'match-detail', component: MatchDetailView },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
