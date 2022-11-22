import {run} from '../src/main'
import {expect, test} from '@jest/globals'

test('run', async () => {
  process.env.INPUT_SCRIPT = 'test.py'
  await run()
})
