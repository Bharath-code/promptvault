#!/usr/bin/env node

/**
 * PromptVault CLI wrapper
 * Spawns the Go binary with all command-line arguments
 */

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

// Determine binary path
const binaryName = process.platform === 'win32' ? 'promptvault.exe' : 'promptvault';
const binaryPath = path.join(__dirname, '..', binaryName);

// Check if binary exists
if (!fs.existsSync(binaryPath)) {
  console.error('❌ PromptVault binary not found');
  console.error('');
  console.error('The postinstall script may have failed. Try reinstalling:');
  console.error('  npm uninstall -g promptvault && npm install -g promptvault');
  console.error('');
  console.error('Or download manually from:');
  console.error('  https://github.com/Bharath-code/promptvault/releases');
  process.exit(1);
}

// Spawn the Go binary with all arguments
const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: 'inherit',
  windowsHide: false,
});

// Handle errors
child.on('error', (error) => {
  if (error.code === 'ENOENT') {
    console.error('❌ Failed to execute PromptVault binary');
    console.error('');
    console.error('The binary may be corrupted or incompatible with your system.');
    console.error('Try reinstalling:');
    console.error('  npm uninstall -g promptvault && npm install -g promptvault');
  } else {
    console.error('❌ Error running PromptVault:', error.message);
  }
  process.exit(1);
});

// Exit with the same code as the child process
child.on('exit', (code) => {
  process.exit(code);
});

// Handle process signals
process.on('SIGINT', () => {
  child.kill('SIGINT');
});

process.on('SIGTERM', () => {
  child.kill('SIGTERM');
});
