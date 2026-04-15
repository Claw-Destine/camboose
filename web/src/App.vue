<script setup>
import { computed, onMounted, reactive, ref } from 'vue'

import ProjectForm from './components/ProjectForm.vue'
import ProjectList from './components/ProjectList.vue'
import RecipeCatalog from './components/RecipeCatalog.vue'
import { createProject, deleteProject, listProjects, listRecipies } from './api'

const recipes = ref([])
const projects = ref([])
const loading = reactive({
  recipes: false,
  projects: false,
  create: false,
  deleteName: '',
})
const errorMessage = ref('')
const successMessage = ref('')

const availableRecipeNames = computed(() => recipes.value.map((recipe) => recipe.name))

function setError(message) {
  errorMessage.value = message
  if (message) {
    successMessage.value = ''
  }
}

function setSuccess(message) {
  successMessage.value = message
  if (message) {
    errorMessage.value = ''
  }
}

async function loadRecipes() {
  loading.recipes = true

  try {
    recipes.value = await listRecipies()
  } catch (error) {
    setError(error instanceof Error ? error.message : 'Failed to load recipes.')
  } finally {
    loading.recipes = false
  }
}

async function loadProjects() {
  loading.projects = true

  try {
    projects.value = await listProjects()
  } catch (error) {
    setError(error instanceof Error ? error.message : 'Failed to load projects.')
  } finally {
    loading.projects = false
  }
}

async function handleCreateProject(payload) {
  loading.create = true

  try {
    const project = await createProject(payload)
    projects.value = [...projects.value, project].sort((left, right) => left.name.localeCompare(right.name))
    setSuccess(`Created project ${project.name}.`)
  } catch (error) {
    setError(error instanceof Error ? error.message : 'Failed to create project.')
    throw error
  } finally {
    loading.create = false
  }
}

async function handleDeleteProject(name) {
  loading.deleteName = name

  try {
    await deleteProject(name)
    projects.value = projects.value.filter((project) => project.name !== name)
    setSuccess(`Deleted project ${name}.`)
  } catch (error) {
    setError(error instanceof Error ? error.message : 'Failed to delete project.')
  } finally {
    loading.deleteName = ''
  }
}

onMounted(async () => {
  await Promise.all([loadRecipes(), loadProjects()])
})
</script>

<template>
  <div class="app-shell">
    <header class="hero">
      <div class="hero__copy">
        <p class="eyebrow">Recipe-driven delivery orchestration</p>
        <h1>camboose control room</h1>
        <p class="hero__lede">
          Bootstrap projects from recipes, inspect what is available, and keep the first frontend layer close to the
          current Go API.
        </p>
      </div>
      <div class="hero__stats">
        <div class="stat-card">
          <span class="stat-card__label">Recipes</span>
          <strong>{{ recipes.length }}</strong>
        </div>
        <div class="stat-card">
          <span class="stat-card__label">Projects</span>
          <strong>{{ projects.length }}</strong>
        </div>
      </div>
    </header>

    <section v-if="errorMessage || successMessage" class="status-strip" :class="{ 'status-strip--error': errorMessage }">
      {{ errorMessage || successMessage }}
    </section>

    <main class="dashboard">
      <ProjectForm
        :loading="loading.create"
        :recipe-names="availableRecipeNames"
        @create-project="handleCreateProject"
      />

      <ProjectList
        :loading="loading.projects"
        :deleting-name="loading.deleteName"
        :projects="projects"
        @delete-project="handleDeleteProject"
      />

      <RecipeCatalog :loading="loading.recipes" :recipes="recipes" />
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  max-width: 1180px;
  margin: 0 auto;
  padding: 48px 20px 72px;
}

.hero {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(280px, 0.9fr);
  gap: 24px;
  align-items: stretch;
  margin-bottom: 28px;
}

.hero__copy,
.hero__stats {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  background: var(--color-surface);
  box-shadow: var(--shadow-soft);
}

.hero__copy {
  padding: 32px;
  position: relative;
  overflow: hidden;
}

.hero__copy::after {
  content: '';
  position: absolute;
  inset: auto -60px -80px auto;
  width: 220px;
  height: 220px;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(212, 138, 88, 0.28) 0%, rgba(212, 138, 88, 0) 72%);
}

.hero__stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  padding: 20px;
}

.eyebrow {
  margin: 0 0 14px;
  color: var(--color-accent-strong);
  text-transform: uppercase;
  letter-spacing: 0.16em;
  font-size: 0.72rem;
  font-weight: 700;
}

h1 {
  margin: 0;
  max-width: 12ch;
  font-size: clamp(2.8rem, 6vw, 5.2rem);
  line-height: 0.94;
}

.hero__lede {
  position: relative;
  z-index: 1;
  max-width: 58ch;
  margin: 18px 0 0;
  color: var(--color-text-muted);
  font-size: 1.02rem;
}

.stat-card {
  padding: 20px;
  border-radius: var(--radius-lg);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.92), rgba(246, 239, 230, 0.88));
  border: 1px solid rgba(100, 76, 57, 0.12);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-height: 160px;
}

.stat-card strong {
  font-size: clamp(2rem, 4vw, 3rem);
}

.stat-card__label {
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 0.78rem;
}

.status-strip {
  margin-bottom: 24px;
  padding: 14px 16px;
  border-radius: var(--radius-md);
  background: rgba(52, 114, 88, 0.12);
  border: 1px solid rgba(52, 114, 88, 0.25);
  color: var(--color-success);
}

.status-strip--error {
  background: rgba(171, 73, 48, 0.12);
  border-color: rgba(171, 73, 48, 0.22);
  color: var(--color-danger);
}

.dashboard {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 20px;
}

.dashboard :deep(.panel:nth-child(1)) {
  grid-column: span 4;
}

.dashboard :deep(.panel:nth-child(2)) {
  grid-column: span 5;
}

.dashboard :deep(.panel:nth-child(3)) {
  grid-column: span 3;
}

@media (max-width: 960px) {
  .hero {
    grid-template-columns: 1fr;
  }

  .dashboard :deep(.panel) {
    grid-column: 1 / -1;
  }
}

@media (max-width: 640px) {
  .app-shell {
    padding: 28px 14px 56px;
  }

  .hero__copy,
  .hero__stats {
    padding: 18px;
  }

  .hero__stats {
    grid-template-columns: 1fr 1fr;
  }
}
</style>