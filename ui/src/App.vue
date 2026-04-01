<script setup>
import { ref } from 'vue'
import SplashScreen from './components/SplashScreen.vue'
import MainMenu from './components/MainMenu.vue'
import CompanyRegistration from './components/CompanyRegistration.vue'
import GameDashboard from './components/GameDashboard.vue' // Add this import

const currentScreen = ref('splash') 
const activeCharter = ref(null) // State to hold the registered data

const handleBootComplete = () => {
  currentScreen.value = 'menu'
}

const handleNewGame = () => {
  currentScreen.value = 'registration'
}

const handleRegistrationComplete = (charterData) => {
  activeCharter.value = charterData // Store the payload
  currentScreen.value = 'game'      // Mount the dashboard
}
</script>

<template>
  <div id="terminal-root">
    <SplashScreen v-if="currentScreen === 'splash'" @complete="handleBootComplete" />
    
    <MainMenu v-else-if="currentScreen === 'menu'" @new-game="handleNewGame" />
    
    <CompanyRegistration 
      v-else-if="currentScreen === 'registration'" 
      @complete="handleRegistrationComplete"
      @cancel="currentScreen = 'menu'"
    />
    
    <GameDashboard 
      v-else-if="currentScreen === 'game'" 
      :charterData="activeCharter" 
    />
  </div>
</template>

<style>
/* Import VT323 - it needs a larger base font size to be readable */
@import url('https://fonts.googleapis.com/css2?family=VT323&display=swap');

:root {
  /* Base Palette - Using variables allows for easy theme swapping later */
  --bg-color: #050505;
  --text-color: #e0e0e0;
  --highlight-bg: #e0e0e0;
  --highlight-text: #050505;
  --border-style: 2px solid var(--text-color);
  --dim-text: #888888;
}

* {
  box-sizing: border-box;
}

body, html {
  margin: 0;
  padding: 0;
  height: 100vh;
  width: 100%;
  background-color: var(--bg-color);
  color: var(--text-color);
  font-family: 'VT323', monospace;
  font-size: 26px; 
  line-height: 1.2;
  overflow: hidden;
  -webkit-font-smoothing: none;
}

#app {
  width: 100%;
  height: 100%;
  max-width: none;
  margin: 0;
  padding: 0;
  text-align: left;
}

#terminal-root {
  width: 100%; /* Changed from 100vw to prevent horizontal scrollbar bleed */
  height: 100%;
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
  align-items: stretch; 
  overflow: hidden;
}

/* Global Reusable UI Elements */
.btn-block {
  background: transparent;
  color: var(--text-color);
  border: var(--border-style);
  font-family: inherit;
  font-size: 1.2rem;
  padding: 0.5rem 1rem;
  cursor: pointer;
  text-transform: none; 
  /* Crucial for retro feel: 0s transition makes hover states instant */
  transition: background-color 0s, color 0s; 
}

.btn-block:hover, .btn-block:focus {
  background-color: var(--highlight-bg);
  color: var(--highlight-text);
  outline: none;
}
</style>