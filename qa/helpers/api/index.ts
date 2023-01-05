import supertest from 'supertest'
import dotenv from 'dotenv'

dotenv.config()

const apiST = supertest(process.env.API_BASE_URL)
const env = process.env.ENV || 'development'

const contentTypeHandler = (api: any, contentType: string, data: any) => {
  const payload = data === undefined ? {} : data[env]
  return contentType === 'multipart/form-data'
    ? api.field(payload)
    : api.send(payload)
}

const get = (endpoint: string, {authorization, params = {}}: any) => {
  const api = apiST.get(endpoint)

  return authorization !== undefined
    ? api.set('Authorization', `Bearer ${authorization}`).query(params[env])
    : api.query(params[env])
}

const patch = (
  endpoint: string,
  {data, authorization, contentType = 'application/json'}: any,
) => {
  let api = apiST
    .patch(endpoint)
    .set('Content-Type', contentType)
    .set('Authorization', `Bearer ${authorization}`)
  api = contentTypeHandler(api, contentType, data)

  return api
}

const put = (
  endpoint: string,
  {data, authorization, contentType = 'application/json'}: any,
) => {
  let api = apiST
    .put(endpoint)
    .set('Content-Type', contentType)
    .set('Authorization', `Bearer ${authorization}`)
  api = contentTypeHandler(api, contentType, data)

  return api
}

const post = (
  endpoint: string,
  {data, cookie, authorization, contentType = 'application/json'}: any,
) => {
  let api = apiST
    .post(endpoint)
    .set('Content-Type', contentType)
    .set('Authorization', `Bearer ${authorization}`)
  if (cookie !== undefined) {
    api = api.set('Cookie', cookie)
  }
  api = contentTypeHandler(api, contentType, data)

  return api
}

const del = (
  endpoint: string,
  {data, authorization, contentType = 'application/json'}: any,
) =>
  apiST
    .delete(endpoint)
    .set('Content-Type', contentType)
    .set('Authorization', `Bearer ${authorization}`)
    .send(data[env])

const api = {
  get,
  post,
  patch,
  delete: del,
  put,
}

export {api}
