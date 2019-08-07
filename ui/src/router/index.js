import Vue from 'vue'
import Router from 'vue-router'
import Splash from '@/components/Splash'
import Movie from '@/components/Movie'
import LogIn from '@/components/Login';
import Browse from '@/components/Browse';

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Splash',
      component: Splash
    },
    {
      path: `/movie/:id`,
      name: 'Movie',
      component: Movie
    },
    {
      path: '/login',
      name: 'LogIn',
      component: LogIn
    },
    {
      path: '/browse',
      name: 'Browse',
      component: Browse
    }
  ]
})
