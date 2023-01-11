process.env['CHROME_BIN'] = require('puppeteer').executablePath()

module.exports = (config: any) => {
    config.set({
        browsers: ['ChromeHeadless']
    })
}
