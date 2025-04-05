import type { Key } from 'chessground/types';
import { Chess, type Color, type Piece } from 'chess.js';

export function getBaseURL(): string {
	return 'http://localhost:8080';
}

export function colorToCgColor(chessjsColor: Color): 'white' | 'black' {
	return chessjsColor === 'w' ? 'white' : 'black';
}
export function cgColorToColor(chessgroundColor: 'white' | 'black'): Color {
	return chessgroundColor === 'white' ? 'w' : 'b';
}

export const getValidMoves = (chess: Chess) => {
	const moves = new Map<Key, Key[]>();

	chess.board().forEach((row) => {
		row.forEach((square) => {
			if (square) {
				const from = square.square as Key;
				const legalMoves = chess.moves({ square: square.square, verbose: true }).map((m) => m.to);
				if (legalMoves.length) moves.set(from, legalMoves);
			}
		});
	});

	return moves;
};

export function isPromoting(dest: Key, piece: Piece) {
	return (
		piece.type == 'p' &&
		((piece.color == 'w' && dest[1] == '8') || (piece.color == 'b' && dest[1] == '1'))
	);
}

export function formatPostgresTimestamp(dateObj: Date): string {
	const now = new Date();
	const diffInSeconds = Math.floor((now.getTime() - dateObj.getTime()) / 1000);

	if (diffInSeconds < 60) {
		return `${diffInSeconds} second${diffInSeconds > 1 ? 's' : ''} ago`;
	}

	const diffInMinutes = Math.floor(diffInSeconds / 60);
	if (diffInMinutes < 60) {
		return `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`;
	}

	const diffInHours = Math.floor(diffInMinutes / 60);
	if (diffInHours < 24) {
		return `${diffInHours} hour${diffInHours > 1 ? 's' : ''} ago`;
	}

	// Format date for timestamps older than 1 day
	const options: Intl.DateTimeFormatOptions = {
		day: 'numeric',
		month: 'long',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit',
		hour12: false
	};
	return dateObj.toLocaleDateString('en-GB', options).replace(',', '');
}
