import {timeout} from "@/utils.js";
import {api, EXPRESSION_STATUS} from "@/api.js";
import {ref, toValue, watchEffect} from "vue";


const EXPRESSION_FAILED_RESPONSE = 'make sure there is no division by zero'
const EXPRESSION_REQUEST_INTERVAL_IN_MS = 200


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
    if (status === EXPRESSION_STATUS.FAILED) {
      errorRef.value = EXPRESSION_FAILED_RESPONSE
      return
    } else if (status !== EXPRESSION_STATUS.PROCESSED) {
      await timeout(EXPRESSION_REQUEST_INTERVAL_IN_MS)
      continue
    }
    resultRef.value = result
    return
  }
}
