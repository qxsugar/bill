const ENV = 'dev'

const configs = {
  dev: {
    baseUrl: 'https://bill.ppapi.cn',
  },
  prod: {
    baseUrl: 'https://bill.ppapi.cn',
  },
}

module.exports = configs[ENV]
