import {snakeToCamel} from "@/utils.js";
import {useAuthentication} from "@/composables.js";

const BASE_URL = import.meta.env.VITE_BACKEND_BASE_URL ?? ''
const NO_RESPONSE_FROM_SERVER_MESSAGE = 'no response body from the server'
const UNKNOWN_ERROR_MESSAGE = 'something went wrong...'

export const EXPRESSION_STATUS = {
  NEW: 'new',
  PROCESSING: 'processing',
  SUCCEED: 'succeed',
  ABORTED: 'processed',
  FAILED: 'failed'
}

const apiFetch = async (url, options = {}) => {
  const result = {}
  let response
  let responseBody = {}
  const {authHeader, logout} = useAuthentication()

  try {
    options = {...options}
    options.headers = {...options.headers, ...authHeader.value}
    response = await fetch(url, options)
    responseBody = await response.json()
  } catch (e) {
    responseBody.error = NO_RESPONSE_FROM_SERVER_MESSAGE
  }

  if (response.status === 401) {
    await logout()
  }

  if (!response.ok && !responseBody?.error) {
    responseBody.error = UNKNOWN_ERROR_MESSAGE
  }

  Object.entries(responseBody).forEach(([key, value]) => {
    result[snakeToCamel(key)] = value
  })

  return result
}

export const api = {
  async sendExpression(expression) {
    const schema = {id: null, error: null}

    const response = await apiFetch(`${BASE_URL}/api/v1/calculate`, {
      method: 'POST',
      body: JSON.stringify({expression})
    })

    return {...schema, ...response}
  },

  async checkExpression(expressionId) {
    const schema = {
      id: null,
      status: null,
      result: null,
      error: null
    }

    const response = await apiFetch(`${BASE_URL}/api/v1/expressions/${expressionId}`)

    return {...schema, ...response}
  },

  async login({username, password}) {
    return this._auth({username, password, type: 'login'})
  },

  async register({username, password}) {
    return this._auth({username, password, type: 'register'})
  },

  async _auth({username, password, type}) {
    const schema = {
      error: null,
      accessToken: null,
      user: null
    }

    const response = await apiFetch(`${BASE_URL}/api/v1/auth/${type}`, {
      method: 'POST',
      body: JSON.stringify({username, password})
    })

    return {...schema, ...response}
  },

  async currentUser() {
    const schema = {
      id: null,
      username: null,
    }

    const response = await apiFetch(`${BASE_URL}/api/v1/users/me`)

    return {...schema, ...response}
  }
}
