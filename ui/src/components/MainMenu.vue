<script setup>
const emit = defineEmits(['new-game', 'load-game', 'settings'])

// Placeholder for actual engine version
const version = "v0.1.0-alpha"

// Depending on how Wails is configured, window.runtime.Quit() will handle this eventually
const handleQuit = () => {
  console.log('Terminate Connection Triggered')
}
</script>

<template>
  <div class="main-menu">
    <header class="title-banner">
      <pre class="ascii-art">
  _____       _           _            
 / ____|     | |         (_)           
| |  __  __ _| | __ ___  ___  ___  ___ 
| | |_ |/ _' | |/ _' \ \/ / |/ _ \/ __|
| |__| | (_| | | (_| |>  <| |  __/\__ \
 \_____|\__,_|_|\__,_/_/\_\_|\___||___/
                               
     B U R N   R A T E
      </pre>
      <div class="version-tag">{{ version }}</div>
    </header>

    <nav class="menu-options">
      <button class="btn-block" @click="emit('new-game')">New Company Registration</button>
      <button class="btn-block disabled" disabled>Load Corporate Data</button>
      <button class="btn-block disabled" disabled>Terminal Settings</button>
      <button class="btn-block" @click="handleQuit">Terminate Connection</button>
    </nav>
    
    <footer class="footer-info">
      System Ready. Awaiting input.
    </footer>
  </div>
</template>

<style scoped>
.main-menu {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;       /* NEW: Forces the menu to span the entire browser width */
  flex-grow: 1;      /* NEW: Forces it to consume all available vertical space */
  gap: 4rem;
}

.title-banner {
  text-align: center;
}

.ascii-art {
  /* Enforce standard monospace for ASCII to maintain shape */
  font-family: monospace; 
  font-size: 16px; 
  white-space: pre;
  line-height: 1.1;
  margin: 0;
  color: var(--text-color);
}

.version-tag {
  margin-top: 1rem;
  font-size: 1rem;
  color: var(--dim-text);
}

.menu-options {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
  width: 100%;
  max-width: 400px;
}

/* Styling for non-implemented features */
.btn-block.disabled {
  opacity: 0.4;
  cursor: not-allowed;
  border-style: dashed;
}

.btn-block.disabled:hover {
  background-color: transparent;
  color: var(--text-color);
}

.footer-info {
  position: absolute;
  bottom: 2rem;
  font-size: 1.1rem;
  color: var(--dim-text);
}
</style>