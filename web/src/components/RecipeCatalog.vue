<script setup>
defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  recipes: {
    type: Array,
    default: () => [],
  },
})
</script>

<template>
  <section class="panel">
    <div class="panel__header">
      <p class="panel__eyebrow">Reference</p>
      <h2>Recipes</h2>
    </div>

    <div v-if="loading" class="empty-state">Loading recipes...</div>
    <div v-else-if="recipes.length === 0" class="empty-state">No recipes found in the configured recipe directory.</div>
    <ul v-else class="recipe-list">
      <li v-for="recipe in recipes" :key="recipe.name">
        <strong>{{ recipe.name }}</strong>
        <p>{{ recipe.description || 'No description in YAML yet.' }}</p>
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

.recipe-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 12px;
}

.recipe-list li {
  padding: 14px 16px;
  border-radius: var(--radius-lg);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(250, 246, 241, 0.88));
  border: 1px solid var(--color-border);
}

.recipe-list strong {
  display: block;
}

.recipe-list p {
  margin: 6px 0 0;
  color: var(--color-text-muted);
  font-size: 0.92rem;
}

.empty-state {
  border: 1px dashed var(--color-border-strong);
  border-radius: var(--radius-lg);
  padding: 24px;
  color: var(--color-text-muted);
  background: rgba(246, 239, 230, 0.55);
}
</style>