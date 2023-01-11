module.exports = function (config) {
    config.set({
        basePath: '',
        frameworks: ['jasmine'],
        files: ['src/**/*.spec.ts'],
        browsers: ['ChromeHeadlessNoSandbox', 'ChromeHeadless', 'Chrome'],
        customLaunchers: {
            ChromeHeadlessNoSandbox: {
                base: 'ChromeHeadless',
                flags: ['--no-sandbox']
            }
        }
    })
}
