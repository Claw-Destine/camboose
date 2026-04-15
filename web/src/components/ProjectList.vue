<script setup>
defineProps({
  deletingName: {
    type: String,
    default: '',
  },
  loading: {
    type: Boolean,
    default: false,
  },
  projects: {
    type: Array,
    default: () => [],
  },
})

defineEmits(['delete-project'])
</script>

<template>
  <section class="panel">
    <div class="panel__header">
      <div>
        <p class="panel__eyebrow">Observe</p>
        <h2>Projects</h2>
      </div>
      <span class="badge">{{ projects.length }}</span>
    </div>

    <div v-if="loading" class="empty-state">Loading projects...</div>
    <div v-else-if="projects.length === 0" class="empty-state">No projects yet. Create the first one from a recipe.</div>
    <ul v-else class="project-list">
      <li v-for="project in projects" :key="project.name" class="project-list__item">
        <div>
          <strong>{{ project.name }}</strong>
          <p>{{ project.recipe }}</p>
        </div>
        <button
          class="ghost-button"
          type="button"
          :disabled="deletingName === project.name"
          @click="$emit('delete-project', project.name)"
        >
          {{ deletingName === project.name ? 'Removing...' : 'Delete' }}
        </button>
      </li>
    </ul>
  </section>
</template>

<style scoped>
.panel {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  background: var(--color-surface);
  box-shadow: var(--shadow-soft);
  padding: 24px;
}

.panel__header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  margin-bottom: 18px;
}

.panel__eyebrow {
  margin: 0 0 8px;
  color: var(--color-accent-strong);
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-size: 0.72rem;
  font-weight: 700;
}

h2 {
  margin: 0;
  font-size: 1.5rem;
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 2.5rem;
  min-height: 2.5rem;
  padding: 0 0.75rem;
  border-radius: 999px;
  background: rgba(212, 138, 88, 0.16);
  color: var(--color-accent-strong);
  font-weight: 700;
}

.project-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 12px;
}

.project-list__item {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: rgba(255, 255, 255, 0.76);
}

.project-list__item strong {
  display: block;
  font-size: 1rem;
}

.project-list__item p {
  margin: 4px 0 0;
  color: var(--color-text-muted);
}

.ghost-button {
  border: 1px solid rgba(171, 73, 48, 0.24);
  border-radius: 999px;
  background: rgba(171, 73, 48, 0.08);
  color: var(--color-danger);
  padding: 10px 14px;
  font: inherit;
  font-weight: 700;
  cursor: pointer;
}

.ghost-button:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.empty-state {
  border: 1px dashed var(--color-border-strong);
  border-radius: var(--radius-lg);
  padding: 24px;
  color: var(--color-text-muted);
  background: rgba(246, 239, 230, 0.55);
}

@media (max-width: 640px) {
  .project-list__item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>