#!/usr/bin/env node

/**
 * Postinstall script for PromptVault NPM package
 * Downloads the appropriate binary from GitHub Releases based on platform/architecture
 */

const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');

// Platform mappings
const PLATFORMS = {
  darwin: {
    arm64: 'promptvault-darwin-arm64',
    amd64: 'promptvault-darwin-amd64',
  },
  linux: {
    arm64: 'promptvault-linux-arm64',
    amd64: 'promptvault-linux-amd64',
  },
  win32: {
    amd64: 'promptvault-windows-amd64.exe',
  },
};

/**
 * Get the binary name for current platform
 */
function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();
  
  // Map node arch names to go arch names
  const archMap = {
    x64: 'amd64',
    arm64: 'arm64',
  };
  
  const goArch = archMap[arch] || arch;
  
  if (!PLATFORMS[platform] || !PLATFORMS[platform][goArch]) {
    throw new Error(
      `Unsupported platform: ${platform} ${arch}\n` +
      `Supported platforms: ${Object.keys(PLATFORMS).join(', ')}\n` +
      `Supported architectures: x64, arm64`
    );
  }
  
  return PLATFORMS[platform][goArch];
}

/**
 * Download binary from GitHub Releases
 */
function downloadBinary(binaryName) {
  const packageJson = require('./package.json');
  const version = packageJson.version;
  
  // Try GitHub Releases first
  const url = `https://github.com/Bharath-code/promptvault/releases/download/v${version}/${binaryName}`;
  
  return new Promise((resolve, reject) => {
    https.get(url, {
      headers: {
        'User-Agent': 'promptvault-npm-installer'
      }
    }, (response) => {
      // Handle redirects
      if (response.statusCode === 302 || response.statusCode === 301) {
        downloadBinaryFromUrl(response.headers.location)
          .then(resolve)
          .catch(reject);
        return;
      }
      
      if (response.statusCode !== 200) {
        reject(new Error(
          `Failed to download binary (HTTP ${response.statusCode})\n` +
          `URL: ${url}\n` +
          `Please check your internet connection or install manually from:\n` +
          `https://github.com/Bharath-code/promptvault/releases`
        ));
        return;
      }
      
      const binaryPath = path.join(
        __dirname,
        'promptvault' + (os.platform() === 'win32' ? '.exe' : '')
      );
      
      const file = fs.createWriteStream(binaryPath);
      
      response.pipe(file);
      
      file.on('finish', () => {
        file.close();
        
        // Make executable on Unix-like systems
        if (os.platform() !== 'win32') {
          fs.chmodSync(binaryPath, 0o755);
        }
        
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(path.join(__dirname, 'promptvault'), () => {}); // Cleanup
      reject(new Error(
        `Failed to download binary: ${err.message}\n` +
        `Please check your internet connection or install manually from:\n` +
        `https://github.com/Bharath-code/promptvault/releases`
      ));
    });
  });
}

/**
 * Download binary from a direct URL (for redirects)
 */
function downloadBinaryFromUrl(url) {
  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download binary: ${response.statusCode}`));
        return;
      }
      
      const binaryPath = path.join(
        __dirname,
        'promptvault' + (os.platform() === 'win32' ? '.exe' : '')
      );
      
      const file = fs.createWriteStream(binaryPath);
      response.pipe(file);
      
      file.on('finish', () => {
        file.close();
        
        if (os.platform() !== 'win32') {
          fs.chmodSync(binaryPath, 0o755);
        }
        
        resolve();
      });
    }).on('error', reject);
  });
}

/**
 * Main installation function
 */
async function main() {
  try {
    console.log('🔍 Detecting platform...');
    const binaryName = getBinaryName();
    
    console.log(`📦 Downloading PromptVault ${binaryName}...`);
    await downloadBinary(binaryName);
    
    console.log('✅ PromptVault installed successfully!');
    console.log('');
    console.log('Run `promptvault --help` to get started.');
    console.log('');
  } catch (error) {
    console.error('❌ Failed to install PromptVault');
    console.error('');
    console.error(error.message);
    console.error('');
    console.error('Manual installation:');
    console.error('1. Download from https://github.com/Bharath-code/promptvault/releases');
    console.error('2. Extract and add to your PATH');
    process.exit(1);
  }
}

// Run installation
main();
