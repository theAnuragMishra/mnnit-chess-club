import type { Key } from '@lichess-org/chessground/types';
import { Chess, type Color, type Piece } from 'chess.js';
import { PUBLIC_BASE_URL } from '$env/static/public';

export function getBaseURL(): string {
	return PUBLIC_BASE_URL;
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

const methods: { [index: number]: string } = {
	0: 'Game is live',
	1: 'Checkmate',
	2: 'Resignation',
	3: 'DrawOffer',
	4: 'Draw by stalemate',
	5: 'Draw by threefold repetition',
	6: 'Draw by fivefold repetition',
	7: 'Draw by 50 move rule',
	8: 'Draw by 75 move rule',
	9: 'Insufficient material | Draw',
	10: 'Draw by mutual agreement',
	11: 'White resigned',
	12: 'Black resigned',
	13: 'Game Aborted',
	14: 'White timeout',
	15: 'Black timeout',
	16: "White didn't play",
	17: "Black didn't play"
};

export function formatResultAndMethod(result: number, method: number) {
	return `${methods[method]}${result === 1 ? ' | White wins' : result === 2 ? ' | Black wins' : ''}`;
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

export function dateTimeToDate(date: string | unknown, time: string | unknown) {
	if (!date || !time) return null;
	const combined = new Date(`${date}T${time}`);
	return combined;
}

export function getTimeLeft(startTimeStr: string, seconds: number) {
	const startTime = new Date(startTimeStr);
	const currentTime = new Date();
	// console.log(startTime, currentTime);
	if (isNaN(startTime.getTime())) {
		return -1;
	}
	return seconds * 1000 - Math.max(0, currentTime.getTime() - startTime.getTime());
}
