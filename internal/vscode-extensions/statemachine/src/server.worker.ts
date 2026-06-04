// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// This Web Worker hosts the WASM-compiled Go language server. Go's
// `wasm_exec.js` runtime is inlined ahead of this bundle (esbuild banner) and
// defines `Go` on the global scope.
//
// The worker receives the WASM module bytes as its first message (the extension
// reads them and transfers them in), then relays JSON-RPC messages between the
// worker's postMessage channel and the WASM server. The bytes are passed in
// rather than fetched here because the worker is loaded from a blob/file URL it
// cannot resolve resources against.
//
// The browser LanguageClient (vscode-languageclient/browser) exchanges JSON-RPC
// *objects* over postMessage, while the Go server speaks bare JSON *strings*
// (jsonrpc2.RawFramer). This glue translates between the two representations.

declare const Go: {
    new(): {
        importObject: WebAssembly.Imports;
        run(instance: WebAssembly.Instance): Promise<void>
    }
};

interface InitMessage {
    __init: true;
    wasm: Uint8Array<ArrayBuffer>;
}

const glue = globalThis as {
    wasmSend?: (message: string) => void;
    wasmReady?: () => void;
    wasmReceive?: (message: string) => void;
};

// Messages that arrive before the WASM server has registered its receive
// callback are buffered and flushed once it signals readiness.
const pending: string[] = [];
let ready = false;
let started = false;

// Server -> client: the Go side hands us a JSON string; forward it as an object.
glue.wasmSend = (message: string) => {
    self.postMessage(JSON.parse(message));
};

// Called by the Go side once `wasmReceive` is registered.
glue.wasmReady = () => {
    ready = true;
    for (const message of pending) {
        glue.wasmReceive?.(message);
    }
    pending.length = 0;
};

self.onmessage = (event) => {
    const data = event.data as Partial<InitMessage> | unknown;
    // The first message carries the WASM bytes and starts the server.
    if (!started && (data as Partial<InitMessage>)?.__init) {
        started = true;
        startServer((data as InitMessage).wasm);
        return;
    }
    // Client -> server: stringify the JSON-RPC object for the Go side.
    const message = JSON.stringify(data);
    if (ready && glue.wasmReceive) {
        glue.wasmReceive(message);
    } else {
        pending.push(message);
    }
};

async function startServer(wasm: Uint8Array<ArrayBuffer>): Promise<void> {
    const go = new Go();
    const source = await WebAssembly.instantiate(wasm, go.importObject, {});
    const instance = source.instance;
    // go.run resolves only when the Go program exits; the server blocks on its
    // connection, so this keeps the WASM instance alive for the worker's life.
    await go.run(instance);
}
