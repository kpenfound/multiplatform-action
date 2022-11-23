const core = require('@actions/core')
const os = require('os')
const cp = require('child_process')

async function run() {
  try {
    let script = core.getInput('script')
    core.debug(`Running ${script} ...`) // debug is only output if you set the secret `ACTIONS_STEP_DEBUG` to true
    core.debug(new Date().toTimeString())

    let goos = os.platform()
    let goarch = os.arch()

    if (goos == 'win32') {
      goos = 'windows'
    }
    if (goarch == 'x86') {
      goarch = 'amd64'
    }

    const goBinary = `${__dirname}/build/action_${goos}_${goarch}`
    const out = cp.spawnSync(goBinary, { stdio: 'inherit' })
    core.setOutput('stdout', new Date().toTimeString())
  } catch (error) {
    if (error instanceof Error) core.setFailed(error.message)
  }
}

run()