import axios from 'axios'
import { BASE_API_GATEWAY } from '../config/env'

const apiClient = axios.create({
  baseURL: BASE_API_GATEWAY,
  timeout: 5000,
})

export { apiClient }
