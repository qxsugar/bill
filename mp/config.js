const ENV = 'dev'

const configs = {
  dev: {
    baseUrl: 'https://dev-api.example.com',
  },
  prod: {
    baseUrl: 'https://api.example.com',
  },
}

module.exports = configs[ENV]
