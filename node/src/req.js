import axios from 'axios'
import cfg from './cfg'

const request = axios.create({
  baseURL: cfg.API_BASE,
  validateStatus: null
})

export default request