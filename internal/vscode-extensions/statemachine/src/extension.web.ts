// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

import * as vscode from 'vscode';
import type { LanguageClientOptions } from 'vscode-languageclient/browser.js';
import { LanguageClient } from 'vscode-languageclient/browser.js';

let client: LanguageClient;

// This function is called when the extension is activated in a browser host
// (e.g. VS Code Web). The language server runs as a WebAssembly module inside a
// Web Worker instead of as a spawned process.
export async function activate(context: vscode.ExtensionContext): Promise<void> {
    client = await startLanguageClient(context);
}

// This function is called when the extension is deactivated.
export function deactivate(): Thenable<void> | undefined {
    if (client) {
        return client.stop();
    }
    return undefined;
}

async function startLanguageClient(context: vscode.ExtensionContext): Promise<LanguageClient> {
    // The worker hosts the WASM-compiled Go language server and bridges
    // JSON-RPC messages to it over postMessage. We read the WASM bytes here (the
    // extension host can access its own bundled resources regardless of scheme)
    // and transfer them to the worker as its first message, because the worker
    // cannot fetch/resolve resources from its own blob/file URL.
    const workerUri = vscode.Uri.joinPath(context.extensionUri, 'dist', 'server.worker.js');
    const wasmUri = vscode.Uri.joinPath(context.extensionUri, 'dist', 'server.wasm');
    const worker = new Worker(workerUri.toString(true));
    const wasm = await loadWasm(wasmUri);
    worker.postMessage({ __init: true, wasm }, [wasm.buffer]);

    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ language: 'statemachine' }]
    };
    const client = new LanguageClient(
        'statemachine',
        'Statemachine',
        clientOptions,
        worker
    );

    await client.start();
    return client;
}

// Reads the bundled WASM module. Uses the workspace file system (handles
// file:// and virtual schemes) and falls back to fetch for http(s) hosts.
async function loadWasm(uri: vscode.Uri): Promise<Uint8Array> {
    try {
        return await vscode.workspace.fs.readFile(uri);
    } catch {
        const response = await fetch(uri.toString(true));
        return new Uint8Array(await response.arrayBuffer());
    }
}
