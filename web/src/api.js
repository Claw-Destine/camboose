const jsonHeaders = {
  'Content-Type': 'application/json',
}

async function request(path, options = {}) {
  const response = await fetch(path, options)

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || `Request failed with status ${response.status}`)
  }

  if (response.status === 204) {
    return null
  }

  return response.json()
}

export function listRecipies() {
  return request('/api/recipies')
}

export function listProjects() {
  return request('/api/project')
}

export function createProject(payload) {
  return request('/api/project', {
    method: 'POST',
    headers: jsonHeaders,
    body: JSON.stringify(payload),
  })
}

export function deleteProject(name) {
  return request(`/api/project/${encodeURIComponent(name)}`, {
    method: 'DELETE',
  })
}