package main

const APP_PORT = 4000
const POSTGRES_CONNECTION_PARAMS = `user=mark dbname=linkshortener sslmode=disable`
const REDIS_HOST = `:6379`
const MAX_IDLE_REDIS_CONNECTIONS = 3
const MAX_IDLE_DB_CONNECTIONS = 10
const REDIS_CACHE_TTL = 10
