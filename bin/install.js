#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const http = require('http');
const os = require('os');

const BINARY_DIR = path.join(os.homedir(), '.wlt', 'bin');
const GITHUB_REPO = 'weiliantong/cli';
const VERSION = require('../package.json').version;

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();
  const map = {
    'darwin-x64': 'Darwin_x86_64',
    'darwin-arm64': 'Darwin_arm64',
    'linux-x64': 'Linux_x86_64',
    'linux-arm64': 'Linux_arm64',
    'win32-x64': 'Windows_x86_64',
  };
  const key = `${platform}-${arch}`;
  const suffix = map[key];
  if (!suffix) {
    console.error(`Unsupported platform: ${key}`);
    process.exit(1);
  }
  return { suffix, ext: platform === 'win32' ? '.exe' : '' };
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    const follow = (url) => {
      const mod = url.startsWith('https') ? https : http;
      mod.get(url, { headers: { 'User-Agent': 'node' } }, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          follow(res.headers.location);
          return;
        }
        if (res.statusCode !== 200) {
          reject(new Error(`HTTP ${res.statusCode}`));
          return;
        }
        res.pipe(file);
        file.on('finish', () => { file.close(); resolve(); });
      }).on('error', (err) => { fs.unlink(dest, () => {}); reject(err); });
    };
    follow(url);
  });
}

async function main() {
  const { suffix, ext } = getPlatform();
  const binaryName = `wlt_${suffix}${ext}`;
  const downloadUrl = `https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${binaryName}`;

  if (!fs.existsSync(BINARY_DIR)) {
    fs.mkdirSync(BINARY_DIR, { recursive: true });
  }

  const dest = path.join(BINARY_DIR, `wlt${ext}`);

  console.log(`Downloading wlt v${VERSION} for ${suffix}...`);
  try {
    await download(downloadUrl, dest);
    fs.chmodSync(dest, 0o755);
    console.log(`Installed to ${dest}`);
  } catch (err) {
    console.error(`Download failed: ${err.message}`);
    console.error('You can build from source: go build -o wlt .');
    process.exit(1);
  }
}

main();
