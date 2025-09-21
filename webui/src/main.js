
// 1. IMPORT PACKAGES *********************************************************************************************

/* 1. CSS File for Web UI Styling */
import './assets/main.css'
/* 2. CreateApp Function from Vue to Initialize the Vue Application */
import { createApp } from 'vue'
/* 3. The App Component - The Root Vue.js component that will be rendered first */
import App from './App.vue'
/* 4. The Vue Router Configuration - Allows the app to handle navigation through pages/views */
import router from './router'
/* 5. Axios - Library allowing to make HTTP requests to servers or APIs */
import axios from "axios";


// 2. CREATE VUE APP **********************************************************************************************

/* 1. Create a new Vue App Istance - using the main App component */
const app = createApp(App)


// 3. SETUP AXIOS INSTANCE ****************************************************************************************

/* 2. Add a Custom Axios instance to the app's global properties */
app.config.globalProperties.$axios = axios.create({
    /* Set the Back-End (Golang) default server address for all the HTTP Requests */
    baseURL: "http://localhost:8080",
    /* Set Max Wait Time for HTTP requests before failure */
    timeout:1000*5
});


// 4. MOUNT ROUTER and VUE APP ************************************************************************************

/* Let the Vue App use the Router - Allow navigation between different pages/view */
app.use(router)

/* Attach the Vue App to the HTML Element having the id 'app' - Make App appear on the web page */
app.mount('#app')
