import * as core from '@actions/core'
import Client, {connect} from '@dagger.io/dagger'

export async function run(): Promise<void> {
  try {
    const script: string = core.getInput('script')
    core.debug(`Running ${script} ...`) // debug is only output if you set the secret `ACTIONS_STEP_DEBUG` to true

    core.debug(new Date().toTimeString())
    connect(
      async (client: Client) => {
        // get reference to the local project
        const source = await client.host().workdir(undefined, [script]).id()

        // get Python image
        const py = await client.container().from('python:3.9-slim').id()

        // mount python script and execute
        const runner = client
          .container(py.id)
          .withMountedDirectory('/src', source.id)
          .withWorkdir('/src')
          .exec(['python', script])

        // get stdout
        await runner.stdout().contents() // TODO - update sdk to get string out of stdout
        //const stdout = 'success'
      },
      {LogOutput: process.stdout}
    )
    core.debug(new Date().toTimeString())

    core.setOutput('stdout', new Date().toTimeString())
  } catch (error) {
    if (error instanceof Error) core.setFailed(error.message)
  }
}

run()
