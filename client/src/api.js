const BASE_URL = 'http://' + import.meta.env.VITE_BACKEND_SERVER_HOST + ':' + import.meta.env.VITE_BACKEND_SERVER_PORT
const NO_RESPONSE_FROM_SERVER_MESSAGE = 'no response body from the server'

export const EXPRESSION_STATUS = {
  WAITING_FOR_PROCESSING: 'waiting for processing',
  PROCESSING: 'processing',
  PROCESSED: 'processed',
  FAILED: 'failed'
}

export const api = {
  async sendExpression(expression) {
    const result = {id: null, error: null}
    let response
    let responseBody

    try {
      response = await fetch(`${BASE_URL}/api/v1/calculate`, {
        method: 'POST',
        body: JSON.stringify({expression})
      })
      responseBody = await response.json()
    } catch (e) {
      result.error = NO_RESPONSE_FROM_SERVER_MESSAGE
      return result
    }

    if (!response.ok) {
      result.error = responseBody.error
      return result
    }

    result.id = responseBody.id
    return result
  },

  async checkExpression(expressionId) {
    const result = {
      id: null,
      status: null,
      result: null,
      error: null
    }
    let response
    let responseBody

    try {
      response = await fetch(`${BASE_URL}/api/v1/expressions/${expressionId}`)
      responseBody = await response.json()
    } catch (e) {
      result.error = NO_RESPONSE_FROM_SERVER_MESSAGE
      return result
    }

    if (!response.ok) {
      result.error = responseBody.error
      return result
    }

    result.id = responseBody.id
    result.result = responseBody.result
    result.status = responseBody.status
    return result
  }
}
