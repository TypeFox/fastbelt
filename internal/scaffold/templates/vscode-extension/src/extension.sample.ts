// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Sample extension used only so tsc/VS Code has inputs while editing templates in this repo.
// Scaffolded projects get extension.ts from extension.ts.tmpl instead; this file is not copied.

import type * as vscode from 'vscode';
import type { LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node.js';
import * as path from 'node:path';
import { LanguageClient, TransportKind } from 'vscode-languageclient/node.js';

let client: LanguageClient;

export async function activate(context: vscode.ExtensionContext): Promise<void> {
	client = await startLanguageClient(context);
}

export function deactivate(): Thenable<void> | undefined {
	if (client) {
		return client.stop();
	}
	return undefined;
}

async function startLanguageClient(context: vscode.ExtensionContext): Promise<LanguageClient> {
	const serverOptions: ServerOptions = {
		command: context.asAbsolutePath(path.join('dist', 'server' + (process.platform === 'win32' ? '.exe' : ''))),
		args: [],
		transport: TransportKind.stdio,
	};

	const clientOptions: LanguageClientOptions = {
		documentSelector: [{ scheme: 'file', language: 'placeholder-lang' }],
	};
	const lc = new LanguageClient(
		'placeholder-lang',
		'Placeholder',
		serverOptions,
		clientOptions,
	);
	await lc.start();
	return lc;
}
