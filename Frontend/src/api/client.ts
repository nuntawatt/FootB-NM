import axios from 'axios'

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api/v1'

function getApiHost(): string | null {
  try {
    const parsed = new URL(apiBaseUrl)
    if (parsed.hostname && parsed.hostname !== 'localhost' && parsed.hostname !== '127.0.0.1') {
      return parsed.hostname
    }
  } catch {
    // ignore invalid URL
  }
  return null
}

// Replace localhost/127.0.0.1 in returned media URLs (e.g. MinIO :9000) with the API host when using LAN/remote backend
function rewriteLocalhostUrls(data: unknown, targetHost: string): unknown {
  if (typeof data === 'string') {
    return data.replace(/^(https?:\/\/)(localhost|127\.0\.0\.1)(:\d+)?/i, `$1${targetHost}$3`)
  }
  if (Array.isArray(data)) {
    return data.map((item) => rewriteLocalhostUrls(item, targetHost))
  }
  if (data !== null && typeof data === 'object') {
    const obj = data as Record<string, unknown>
    for (const key of Object.keys(obj)) {
      obj[key] = rewriteLocalhostUrls(obj[key], targetHost)
    }
    return obj
  }
  return data
}

const apiClient = axios.create({
  baseURL: apiBaseUrl,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
})

// Attach JWT Bearer token automatically if present
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor for URL normalization & error handling
apiClient.interceptors.response.use(
  (response) => {
    const apiHost = getApiHost()
    if (apiHost && response.data) {
      response.data = rewriteLocalhostUrls(response.data, apiHost)
    }
    return response
  },
  (error) => {
    const message = error.response?.data?.error || error.message || 'เกิดข้อผิดพลาด'
    console.error('[API Error]', message)
    return Promise.reject(error)
  }
)

export default apiClient
