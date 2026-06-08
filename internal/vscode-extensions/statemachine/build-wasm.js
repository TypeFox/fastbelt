// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Builds the statemachine language server as a WebAssembly module. Go's
// `wasm_exec.js` runtime is inlined into the worker bundle by esbuild, so it is
// not copied here.

const { execFileSync } = require('node:child_process');
const fs = require('node:fs');
const path = require('node:path');

const dist = path.join(__dirname, 'dist');
fs.mkdirSync(dist, { recursive: true });

// Compile the Go server for the JS/WASM target.
execFileSync('go', ['build', '-o', path.join(dist, 'server.wasm'), '../../../examples/statemachine/server'], {
    cwd: __dirname,
    stdio: 'inherit',
    env: { ...process.env, GOOS: 'js', GOARCH: 'wasm' },
});

console.log('Built dist/server.wasm');
