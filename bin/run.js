#!/usr/bin/env node

const { execFileSync } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const ext = os.platform() === 'win32' ? '.exe' : '';
const binaryPath = path.join(os.homedir(), '.wlt', 'bin', `wlt${ext}`);

if (!fs.existsSync(binaryPath)) {
  console.error('wlt binary not found. Run: npm install to download it.');
  process.exit(1);
}

try {
  execFileSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
    env: process.env,
  });
} catch (err) {
  process.exit(err.status || 1);
}
