export function WhitePromotionDialog({
  handleClick,
}: {
  handleClick: (p: "q" | "r" | "b" | "n") => void;
}) {
  return (
    <div>
      <ul>
        <li>
          <button onClick={() => handleClick("q")}>Q</button>
        </li>
        <li>
          <button onClick={() => handleClick("r")}>R</button>
        </li>
        <li>
          <button onClick={() => handleClick("b")}>B</button>
        </li>
        <li>
          <button onClick={() => handleClick("n")}>N</button>
        </li>
      </ul>
    </div>
  );
}

export function BlackPromotionDialog() {
  return <div></div>;
}
