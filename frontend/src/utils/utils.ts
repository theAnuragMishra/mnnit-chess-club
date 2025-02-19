export function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export const chunkArray = (arr: any, size: number) => {
  return arr.reduce(
    (acc, _, i) => (i % size === 0 ? [...acc, arr.slice(i, i + size)] : acc),
    [],
  );
};
