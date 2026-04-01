/**
 * @file useGameState.js
 * @description Centralized state management for Galaxies: Burn Rate. 
 * Acts as the reactive bridge between the Vue UI and the Go engine backend.
 */

import { reactive, computed } from 'vue'

const state = reactive({
  // Global Player Data
  player: {
    name: 'Unknown Manager',
    companyName: 'Independent Logistics',
    credits: 0
  },
  currentDay: 1,
  
  // Entity Collections
  fleet: [],
  planets: {},
  
  // UI State
  activeShipId: null,
  currentView: 'dashboard' // 'dashboard', 'market', 'shipyard', 'cantina', 'manifest', 'navigation'
})

// --- Getters / Computed ---

const activeShip = computed(() => {
  if (!state.activeShipId || state.fleet.length === 0) return null
  return state.fleet.find(s => s.id === state.activeShipId) || state.fleet[0]
})

const activePlanet = computed(() => {
  const ship = activeShip.value
  if (!ship || ship.status !== 'Idle') return null
  return state.planets[ship.locationId]
})

// --- Actions / Mutators ---

function setActiveShip(shipId) {
  state.activeShipId = shipId
  state.currentView = 'dashboard' // Reset view when switching ships
}

function setView(viewName) {
  state.currentView = viewName
}

/**
 * Mocks the backend call to pass time.
 * @param {number} days 
 */
async function advanceTime(days) {
  // TODO: Call Wails backend (e.g., window.go.main.App.AdvanceTime(days))
  console.log(`Advancing time by ${days} days...`)
  state.currentDay += days
}

export function useGameState() {
  return {
    state,
    activeShip,
    activePlanet,
    setActiveShip,
    setView,
    advanceTime
  }
}