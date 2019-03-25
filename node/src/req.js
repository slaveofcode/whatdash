import axios from 'axios'
import cfg from './cfg'

const request = axios.create({
  baseURL: cfg.API_BASE,
  validateStatus: null,
  timeout: 1000 * 60,
})

export default request
export {
  axios,
}