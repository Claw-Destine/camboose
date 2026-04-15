<script setup>
import { computed, reactive } from 'vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  recipeNames: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['create-project'])

const form = reactive({
  name: '',
  recipe: '',
})

const canSubmit = computed(() => form.name.trim() && form.recipe.trim() && !props.loading)

async function submitForm() {
  if (!canSubmit.value) {
    return
  }

  await emit('create-project', {
    name: form.name.trim(),
    recipe: form.recipe.trim(),
  })

  form.name = ''
  form.recipe = ''
}
</script>

<template>
  <section class="panel">
    <div class="panel__header">
      <p class="panel__eyebrow">Create</p>
      <h2>New project</h2>
    </div>

    <form class="form" @submit.prevent="submitForm">
      <label class="field">
        <span>Project name</span>
        <input v-model="form.name" type="text" name="name" placeholder="inventory-refresh" autocomplete="off" />
      </label>

      <label class="field">
        <span>Recipe</span>
        <input v-model="form.recipe" list="recipe-options" type="text" name="recipe" placeholder="static_html_with_blog" />
        <datalist id="recipe-options">
          <option v-for="recipeName in recipeNames" :key="recipeName" :value="recipeName" />
        </datalist>
      </label>

      <button class="button" type="submit" :disabled="!canSubmit">
        {{ loading ? 'Creating...' : 'Create project' }}
      </button>
    </form>
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
  margin-bottom: 20px;
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

.form {
  display: grid;
  gap: 16px;
}

.field {
  display: grid;
  gap: 8px;
}

.field span {
  color: var(--color-text-muted);
  font-size: 0.92rem;
}

.field input {
  width: 100%;
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.9);
  padding: 14px 16px;
  font: inherit;
  color: var(--color-text);
  transition: border-color 140ms ease, box-shadow 140ms ease;
}

.field input:focus {
  outline: none;
  border-color: var(--color-accent-strong);
  box-shadow: 0 0 0 4px rgba(212, 138, 88, 0.18);
}

.button {
  min-height: 48px;
  border: 0;
  border-radius: 999px;
  background: linear-gradient(135deg, var(--color-accent-strong), var(--color-accent));
  color: white;
  font: inherit;
  font-weight: 700;
  cursor: pointer;
  transition: transform 140ms ease, opacity 140ms ease;
}

.button:hover:enabled {
  transform: translateY(-1px);
}

.button:disabled {
  opacity: 0.58;
  cursor: not-allowed;
}
</style>