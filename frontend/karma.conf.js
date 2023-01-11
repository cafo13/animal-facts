module.exports = function (config) {
    config.set({
        basePath: '',
        frameworks: ['jasmine'],
        files: ['src/**/*.spec.ts'],
        logLevel: config.LOG_INFO,
        browsers: ['ChromeHeadless', 'Chrome'],
        customLaunchers: {
            ChromeHeadlessCI: {
                base: 'ChromeHeadless',
                flags: ['--no-sandbox']
            }
        }
    })
}
