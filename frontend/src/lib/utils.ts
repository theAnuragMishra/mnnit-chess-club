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

export function formatResultAndReason(result: string, reason: string) {
	if (result === 'aborted') return 'Game Aborted';
	if (result === '1/2-1/2') {
		if (reason === 'ThreefoldRepetition') return 'Draw by threefold repetition';
		if (reason === 'FivefoldRepetition') return 'Draw by fivefold repetition';
		if (reason === 'FiftyMoveRule') return 'Draw by 50 move rule';
		if (reason === 'SeventyFiveMoveRule') return 'Draw by 75 move rule';
		if (reason === 'Stalemate') return 'Draw by stalemate';
		if (reason === 'InsufficientMaterial') return 'Insufficient material | Draw';
		return reason;
	}
	return `${reason} | ${result === '1-0' ? 'White' : 'Black'} is victorious`;
}

export function getTimeControl(baseTime: number, increment: number) {
	const totalTime = baseTime / 60 + (increment * 2) / 3;
	return totalTime < 3
		? 'Bullet'
		: totalTime < 15
			? 'Blitz'
			: totalTime < 60
				? 'Rapid'
				: 'Classical';
}

export function scrollIntoContainerView(container: HTMLDivElement, targetEl: HTMLElement) {
	const containerRect = container.getBoundingClientRect();
	const targetRect = targetEl.getBoundingClientRect();

	// Calculate the scroll difference
	const offsetTop = targetRect.top - containerRect.top + container.scrollTop;
	const offsetBottom = targetRect.bottom - containerRect.bottom + container.scrollTop;

	// If the element is above the container view
	if (targetRect.top < containerRect.top) {
		container.scrollTo({ top: offsetTop, behavior: 'smooth' });
	}
	// If the element is below the container view
	else if (targetRect.bottom > containerRect.bottom) {
		container.scrollTo({ top: offsetBottom, behavior: 'smooth' });
	}
	// If it's already in view — do nothing
}
