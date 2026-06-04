const esbuild = require('esbuild');
const { execFileSync } = require('node:child_process');
const fs = require('node:fs');
const path = require('node:path');

const production = process.argv.includes('--production');
const watch = process.argv.includes('--watch');

// Go's WASM runtime support script (an IIFE that defines `globalThis.Go`). It is
// inlined into the worker bundle so the worker is self-contained and does not
// rely on `importScripts` with a relative URL, which fails when the browser
// extension host loads the worker from a blob URL.
const goroot = execFileSync('go', ['env', 'GOROOT']).toString().trim();
const wasmExecSource = fs.readFileSync(path.join(goroot, 'lib', 'wasm', 'wasm_exec.js'), 'utf8');

/**
 * @type {import('esbuild').Plugin}
 */
const esbuildProblemMatcherPlugin = {
	name: 'esbuild-problem-matcher',

	setup(build) {
		build.onStart(() => {
			console.log('[watch] build started');
		});
		build.onEnd((result) => {
			result.errors.forEach(({ text, location }) => {
				console.error(`✘ [ERROR] ${text}`);
				console.error(`    ${location.file}:${location.line}:${location.column}:`);
			});
			console.log('[watch] build finished');
		});
	},
};

const shared = {
	bundle: true,
	minify: production,
	sourcemap: !production,
	sourcesContent: false,
	external: ['vscode'],
	logLevel: 'silent',
	plugins: [
		/* add to the end of plugins array */
		esbuildProblemMatcherPlugin,
	],
};

// Desktop (Node.js) extension host: spawns the native Go server over stdio.
const nodeConfig = {
	...shared,
	entryPoints: ['src/extension.ts'],
	format: 'cjs',
	platform: 'node',
	outfile: 'dist/extension.js',
};

// Browser (Web Worker) extension host: runs the WASM Go server in a worker.
const webConfig = {
	...shared,
	entryPoints: ['src/extension.web.ts'],
	format: 'cjs',
	platform: 'browser',
	outfile: 'dist/extension.web.js',
};

// The worker is loaded as a classic worker via `new Worker(...)`. Go's WASM
// runtime is inlined via the banner so the worker bundle is self-contained.
const workerConfig = {
	...shared,
	entryPoints: ['src/server.worker.ts'],
	format: 'iife',
	platform: 'browser',
	outfile: 'dist/server.worker.js',
	banner: { js: wasmExecSource },
};

async function main() {
	const ctxs = await Promise.all(
		[nodeConfig, webConfig, workerConfig].map((config) => esbuild.context(config))
	);
	if (watch) {
		await Promise.all(ctxs.map((ctx) => ctx.watch()));
	} else {
		await Promise.all(ctxs.map(async (ctx) => {
			await ctx.rebuild();
			await ctx.dispose();
		}));
	}
}

main().catch(e => {
	console.error(e);
	process.exit(1);
});
