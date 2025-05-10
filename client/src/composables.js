import {timeout} from "@/utils.js";
import {api, EXPRESSION_STATUS} from "@/api.js";
import {computed, ref, toValue, watchEffect} from "vue";
import {useLocalStorage} from "@vueuse/core";


const EXPRESSION_FAILED_RESPONSE = 'calculation failed, make sure there is no division by zero'
const EXPRESSION_POLLING_INTERVAL_IN_MS = 200


export const useExpressionServerEvaluation = (expression) => {
  const result = ref(null)
  const status = ref(null)
  const error = ref(null)
  const isLoading = ref(false)

  const send = () => {
    isLoading.value = true
    _sendExpressionAndCheckResult(toValue(expression), result, status, error)
      .then(() => isLoading.value = false)
  }

  const reset = () => {
    toValue(expression)  // Хотим обнулять предыдущий результат после изменения выражения
    result.value = null
    status.value = null
    error.value = null
    isLoading.value = false
  }

  watchEffect(() => {
    reset()
  })

  return {result, status, error, isLoading, send, reset}
}

const _sendExpressionAndCheckResult = async (expression, resultRef, statusRef, errorRef) => {
  const {id: expressionId, error} = await api.sendExpression(expression)
  if (error) {
    errorRef.value = error
    return
  }

  while (true) {
    const {result, status, error} = await api.checkExpression(expressionId)
    if (error) {
      errorRef.value = error
      return
    }

    statusRef.value = status
    if (status === EXPRESSION_STATUS.FAILED || status === EXPRESSION_STATUS.ABORTED) {
      errorRef.value = EXPRESSION_FAILED_RESPONSE
      return
    } else if (status !== EXPRESSION_STATUS.SUCCEED) {
      await timeout(EXPRESSION_POLLING_INTERVAL_IN_MS)
      continue
    }
    resultRef.value = result
    return
  }
}

const accessToken = useLocalStorage('accessToken', null)
const currentUser = ref(null)

export const useAuthentication = () => {
  const isAuthenticated = computed(() => accessToken.value !== null)
  const authHeader = computed(
    () => {
      if (accessToken.value === null) {
        return {}
      }
      return {Authorization: `Bearer ${accessToken.value}`}
    }
  )

  const loginError = ref(null)
  const registerError = ref(null)

  const login = async ({username, password}) => {
    loginError.value = null
    const res = await api.login({
      username: toValue(username),
      password: toValue(password)
    })
    if (res?.error) {
      loginError.value = res.error
      return
    }

    accessToken.value = res.accessToken
    currentUser.value = res.user
  }

  const register = async ({username, password}) => {
    registerError.value = null
    const res = await api.register({
      username: toValue(username),
      password: toValue(password)
    })
    if (res?.error) {
      registerError.value = res.error
      return
    }

    accessToken.value = res.accessToken
    currentUser.value = res.user
  }

  const logout = async () => {
    accessToken.value = null
    currentUser.value = null
  }

  const fetchCurrentUser = async () => {
    currentUser.value = await api.currentUser()
  }

  return {
    currentUser,
    isAuthenticated,
    authHeader,
    login,
    loginError,
    register,
    registerError,
    logout,
    fetchCurrentUser
  }
}
