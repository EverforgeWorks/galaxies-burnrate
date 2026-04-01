<script setup>
import { reactive, onMounted, ref } from 'vue'

const emit = defineEmits(['complete', 'cancel'])

// We use a ref for the first input to auto-focus it on load
const firstInput = ref(null)

const form = reactive({
  managerName: '',
  companyName: '',
  shipName: ''
})

const handleSubmit = () => {
  // Emulate the CLI's fallback logic if fields are left blank
  const payload = {
    managerName: form.managerName.trim() || 'Unknown Manager',
    companyName: form.companyName.trim() || 'Independent Logistics',
    shipName: form.shipName.trim() || 'Vanguard-1'
  }
  
  // Pass the payload up to App.vue to hand off to the game state
  emit('complete', payload)
}

onMounted(() => {
  // Instantly focus the first field for immediate keyboard entry
  if (firstInput.value) {
    firstInput.value.focus()
  }
})
</script>

<template>
  <div class="registration-screen">
    <header class="reg-header">
      <h2>> SECURE CORPORATE REGISTRATION PROTOCOL</h2>
      <p class="sys-msg">STATUS: Awaiting input parameters for new charter...</p>
    </header>

    <form @submit.prevent="handleSubmit" class="reg-form">
      <div class="form-group">
        <label for="managerName">MANAGER DESIGNATION</label>
        <div class="input-wrapper">
          <span class="prompt-char">]</span>
          <input 
            id="managerName" 
            ref="firstInput"
            type="text" 
            v-model="form.managerName" 
            placeholder="Enter Company Owner Name"
            autocomplete="off" 
            spellcheck="false" 
          />
        </div>
      </div>

      <div class="form-group">
        <label for="companyName">PROPOSED COMPANY NAME</label>
        <div class="input-wrapper">
          <span class="prompt-char">]</span>
          <input 
            id="companyName" 
            type="text" 
            v-model="form.companyName" 
            placeholder="Enter Company Name Here"
            autocomplete="off" 
            spellcheck="false" 
          />
        </div>
      </div>

      <div class="form-group">
        <label for="shipName">STARTING VESSEL DESIGNATION</label>
        <div class="input-wrapper">
          <span class="prompt-char">]</span>
          <input 
            id="shipName" 
            type="text" 
            v-model="form.shipName" 
            placeholder="Enter Ship Name Here"
            autocomplete="off" 
            spellcheck="false" 
          />
        </div>
      </div>

      <div class="action-row">
        <button type="button" class="btn-block" @click="emit('cancel')">ABORT</button>
        <button type="submit" class="btn-block primary">AUTHORIZE CHARTER</button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.registration-screen {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 0;
  display: flex;
  flex-direction: column;
  gap: 3rem;
  height: 100%;
}

.reg-header h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  font-weight: normal;
}

.sys-msg {
  color: var(--dim-text);
  margin: 0;
}

.reg-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  font-size: 1.1rem;
  color: var(--dim-text);
}

.input-wrapper {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.prompt-char {
  color: var(--dim-text);
  font-weight: bold;
}

/* The "Reset and Rebuild" for terminal inputs:
  Strips all browser styling, applies our font, and sets up the hard inversion on focus.
*/
input {
  all: unset; /* Obliterates browser defaults */
  flex-grow: 1;
  font-family: inherit;
  font-size: 1.5rem;
  color: var(--text-color);
  background: transparent;
  border-bottom: var(--border-style);
  padding: 0.2rem 0;
  transition: none; /* No smooth fading */
}

input::placeholder {
  color: #333333; /* Very dark gray so it's faintly visible */
}

/* The hard terminal block cursor effect */
input:focus {
  background-color: var(--text-color);
  color: var(--bg-color);
  border-bottom-color: transparent;
}

input:focus::placeholder {
  color: transparent;
}

.action-row {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  border-top: 1px dashed var(--dim-text);
  padding-top: 2rem;
}

/* Specific styling for the primary submit button to make it pop slightly */
.btn-block.primary {
  border-color: var(--text-color);
  background-color: transparent;
  flex-grow: 1;
}

.btn-block.primary:hover, .btn-block.primary:focus {
  background-color: var(--text-color);
  color: var(--bg-color);
}
</style>