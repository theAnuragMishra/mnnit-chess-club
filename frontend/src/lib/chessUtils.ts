import { Chess } from 'chess.js';
type PieceKey = 'p' | 'r' | 'n' | 'b' | 'q';
const INITIAL_PIECES: Record<PieceKey, number> = {
	p: 8,
	r: 2,
	n: 2,
	b: 2,
	q: 1
};

const PIECE_VALUES: { [key: string]: number } = {
	p: 1,
	n: 3,
	b: 3,
	r: 5,
	q: 9
};

const TOTAL_MATERIAL = Object.entries(INITIAL_PIECES).reduce(
	(sum, [piece, count]) => sum + PIECE_VALUES[piece] * count,
	0
);

interface CapturedAndMaterial {
	white: Array<'p' | 'q' | 'n' | 'r' | 'b'>;
	black: Array<'p' | 'q' | 'n' | 'r' | 'b'>;
	balance: number;
}

export function getCapturedAndMaterial(chess: Chess): CapturedAndMaterial {
	const countCaptured = (color: 'w' | 'b') => {
		const pieces: Record<PieceKey, number> = { ...INITIAL_PIECES };
		let materialLost = 0;
		const capturedCounts: Record<string, number> = {
			p: 0,
			r: 0,
			n: 0,
			b: 0,
			q: 0,
			k: 0
		};

		const board = chess.board();
		for (const row of board) {
			for (const square of row) {
				if (square && square.color === color && square.type != 'k') {
					pieces[square.type] -= 1;
				}
			}
		}

		for (const [piece, count] of Object.entries(pieces) as [PieceKey, number][]) {
			capturedCounts[piece] = count;
			materialLost += PIECE_VALUES[piece] * count;
		}

		return { capturedCounts, materialLost };
	};

	const whiteLost = countCaptured('w');
	const blackLost = countCaptured('b');

	for (const piece of Object.keys(INITIAL_PIECES)) {
		const min = Math.min(whiteLost.capturedCounts[piece], blackLost.capturedCounts[piece]);
		whiteLost.capturedCounts[piece] -= min;
		blackLost.capturedCounts[piece] -= min;
	}

	// 🔄 CHANGE: helper to build arrays at the end (for UI/output)
	const buildArray = (counts: Record<string, number>) =>
		Object.entries(counts).flatMap(([piece, count]) => Array(count).fill(piece));

	const whiteFinal = buildArray(whiteLost.capturedCounts);
	const blackFinal = buildArray(blackLost.capturedCounts);

	const whiteMaterial = TOTAL_MATERIAL - whiteLost.materialLost;
	const blackMaterial = TOTAL_MATERIAL - blackLost.materialLost;

	return {
		white: whiteFinal,
		black: blackFinal,
		balance: whiteMaterial - blackMaterial
	};
}
