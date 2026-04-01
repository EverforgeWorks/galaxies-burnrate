<script setup>
import { ref } from 'vue'

const props = defineProps({
  charterData: {
    type: Object,
    required: true
  }
})

// Local mock state to build the UI before wiring up the Go backend
const currentDay = ref(1)
const credits = ref(150000)
const activeView = ref('starmap') // 'starmap', 'market', 'shipyard', etc.

const fleet = ref([
  { id: 's_1', name: props.charterData.shipName, class: 'Scout', status: 'IDLE', fuel: 100, maxFuel: 100 },
])

const activeShipId = ref('s_1')

const passTime = () => {
  currentDay.value += 1
  // Trigger Go engine advance time later
}
</script>

<template>
  <div class="game-dashboard">
    
    <header class="top-bar">
      <div class="corp-identity">
        <span class="company-name">{{ charterData.companyName }}</span>
        <span class="manager-name">MGR: {{ charterData.managerName }}</span>
      </div>
      
      <div class="global-metrics">
        <div class="metric">
          <span class="label">FUNDS:</span>
          <span class="value">{{ credits.toLocaleString() }} C</span>
        </div>
        <div class="metric">
          <span class="label">DAY:</span>
          <span class="value">{{ currentDay }}</span>
        </div>
      </div>

      <div class="global-actions">
        <button class="btn-block alt" @click="passTime">PASS TIME [+1]</button>
      </div>
    </header>

    <aside class="side-panel">
      <div class="panel-header">
        <h3>> ACTIVE FLEET</h3>
      </div>
      
      <div class="roster">
        <button 
          v-for="ship in fleet" 
          :key="ship.id"
          class="ship-card btn-block"
          :class="{ 'active-card': activeShipId === ship.id }"
          @click="activeShipId = ship.id"
        >
          <div class="ship-header">
            <span class="ship-name">{{ ship.name }}</span>
            <span class="ship-status">[{{ ship.status }}]</span>
          </div>
          <div class="ship-details">
            <span class="dim-text">CLS: {{ ship.class }}</span>
            <span class="dim-text">FUEL: {{ ship.fuel }}/{{ ship.maxFuel }}</span>
          </div>
        </button>
      </div>

      <div class="panel-footer">
        <span class="dim-text">FLEET CAP: 1/5</span>
      </div>
    </aside>

    <main class="main-view">
      <nav class="sub-nav">
        <button class="nav-tab" :class="{ active: activeView === 'starmap' }" @click="activeView = 'starmap'">[ STARMAP ]</button>
        <button class="nav-tab" :class="{ active: activeView === 'market' }" @click="activeView = 'market'">[ MARKET ]</button>
        <button class="nav-tab" :class="{ active: activeView === 'shipyard' }" @click="activeView = 'shipyard'">[ SHIPYARD ]</button>
        <button class="nav-tab" :class="{ active: activeView === 'manifest' }" @click="activeView = 'manifest'">[ MANIFEST ]</button>
      </nav>

      <section class="view-content">
        <div class="mock-starmap" v-if="activeView === 'starmap'">
          <pre class="ascii-map">
       .            * .
             * [SOL]
  .                       *
          * .                  .
                      [CENTAURI]
      .          * *
          </pre>
          <div class="map-overlay">
            <p>> AWAITING NAV-COMPUTER INPUT...</p>
          </div>
        </div>
        
        <div class="generic-view" v-else>
          <h2>> {{ activeView.toUpperCase() }} MODULE OFFLINE</h2>
          <p class="dim-text">Awaiting implementation.</p>
        </div>
      </section>
    </main>

  </div>
</template>

<style scoped>
/* Master Grid Layout */
.game-dashboard {
  display: grid;
  grid-template-columns: clamp(260px, 20vw, 400px) 1fr;
  grid-template-rows: auto 1fr;
  grid-template-areas: 
    "topbar topbar"
    "sidebar main";
  width: 100%;
  flex-grow: 1; /* Forces the dashboard to stretch and fill the App.vue container */
  border: var(--border-style);
  overflow: hidden; /* Locks the outer frame so only internal panels scroll */
}

/* --- Top Toolbar --- */
.top-bar {
  grid-area: topbar;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: var(--border-style);
  padding: 0.5rem 1rem;
  background-color: var(--highlight-bg); /* Inverted for header */
  color: var(--highlight-text);
}

.corp-identity {
  display: flex;
  flex-direction: column;
}

.company-name {
  font-size: 1.4rem;
  font-weight: bold;
}

.manager-name {
  font-size: 1rem;
}

.global-metrics {
  display: flex;
  gap: 2rem;
}

.metric {
  display: flex;
  gap: 0.5rem;
  align-items: baseline;
}

.metric .label {
  font-size: 1rem;
}

.metric .value {
  font-size: 1.4rem;
  font-weight: bold;
}

/* Specific button style for the inverted header */
.btn-block.alt {
  border-color: var(--highlight-text);
  color: var(--highlight-text);
}
.btn-block.alt:hover, .btn-block.alt:focus {
  background-color: var(--highlight-text);
  color: var(--highlight-bg);
}

/* --- Left Sidebar --- */
.side-panel {
  grid-area: sidebar;
  border-right: var(--border-style);
  display: flex;
  flex-direction: column;
}

.panel-header, .panel-footer {
  padding: 1rem;
  background-color: transparent;
}

.panel-header {
  border-bottom: 1px dashed var(--dim-text);
}

.panel-header h3 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: normal;
}

.roster {
  flex-grow: 1;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  overflow-y: auto;
}

.ship-card {
  display: flex;
  flex-direction: column;
  text-align: left;
  padding: 0.8rem;
}

.ship-card.active-card {
  background-color: var(--text-color);
  color: var(--bg-color);
}

.ship-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.ship-name {
  font-weight: bold;
}

.ship-details {
  display: flex;
  justify-content: space-between;
  font-size: 1rem;
}

.active-card .dim-text {
  color: #555555; /* Darker dim text when inverted */
}

/* --- Main Central Area --- */
.main-view {
  grid-area: main;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-color);
}

.sub-nav {
  display: flex;
  border-bottom: var(--border-style);
}

.nav-tab {
  all: unset;
  padding: 0.5rem 1.5rem;
  font-size: 1.2rem;
  color: var(--dim-text);
  cursor: pointer;
  border-right: var(--border-style);
}

.nav-tab:hover {
  color: var(--text-color);
}

.nav-tab.active {
  background-color: var(--text-color);
  color: var(--bg-color);
}

.view-content {
  flex-grow: 1;
  padding: 2rem;
  position: relative;
  overflow-y: auto; /* Ensures long manifests scroll inside this panel, not the whole page */
  min-height: 0; /* Critical flex/grid quirk: allows children to shrink and trigger internal overflow */
}

/* Temporary Map Styles */
.mock-starmap {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 1px dashed var(--dim-text);
  position: relative;
  overflow: hidden;
}

.ascii-map {
  color: var(--text-color);
}

.map-overlay {
  position: absolute;
  bottom: 1rem;
  left: 1rem;
  background-color: var(--bg-color);
  padding: 0.5rem;
  border: 1px solid var(--text-color);
}

.dim-text {
  color: var(--dim-text);
}
</style>