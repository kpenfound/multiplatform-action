import {run} from '../github-action.mjs'
import {expect, test} from '@jest/globals'

test('run', async () => {
  process.env.INPUT_SCRIPT = 'test.py'
  await run()
})
