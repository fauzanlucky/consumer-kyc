import dotenv from 'dotenv'

dotenv.config()

const host = process.env.CDN_HOST

const REGEX_EMAIL = /^$|^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/

const REGEX_PASSPORT = /^[0-9a-zA-Z]+$/g

const REGEX_KTP = /^[0-9]+$/g

const REGEX_IMAGES_URL = new RegExp(
  `^${host}([\\w-:./])+(.jpg|.png|.webp|.svg)$`,
)

export {REGEX_EMAIL, REGEX_PASSPORT, REGEX_KTP, REGEX_IMAGES_URL}
