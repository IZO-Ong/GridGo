"use client";

interface GridDimensionsInputProps {
  rows: number;
  cols: number;
  labelOverride?: { row: string; col: string };
  onUpdate: (dim: "rows" | "cols", val: number) => void;
  onBlur?: () => void;
  min?: number;
  max?: { rows: number; cols: number };
}

export default function GridDimensionsInput({
  rows,
  cols,
  onUpdate,
  onBlur,
  labelOverride,
  min = 0,
  max = { rows: 999, cols: 999 },
}: GridDimensionsInputProps) {
  const noArrowsClass =
    "[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none";

  const rowLabel = labelOverride?.row || "Rows";
  const colLabel = labelOverride?.col || "Cols";

  const handleBlur = (dim: "rows" | "cols", currentVal: number) => {
    const limit = dim === "rows" ? max.rows : max.cols;
    const clampedVal = Math.min(Math.max(min, currentVal), limit);

    if (clampedVal !== currentVal) {
      onUpdate(dim, clampedVal);
    }

    if (onBlur) onBlur();
  };

  return (
    <div className="flex border-2 border-black bg-white h-[38px] divide-x-2 divide-black">
      <div className="relative flex-1 flex items-center group">
        <input
          name="rows"
          type="number"
          value={rows}
          onChange={(e) => onUpdate("rows", Number(e.target.value))}
          onBlur={(e) => handleBlur("rows", Number(e.target.value))}
          className={`w-full h-full pl-3 pr-10 outline-none font-bold text-sm bg-transparent focus:bg-zinc-50 transition-colors ${noArrowsClass}`}
        />
        <span className="absolute right-2 text-[9px] font-black text-zinc-300 uppercase pointer-events-none tracking-tighter">
          {rowLabel}
        </span>
      </div>

      <div className="relative flex-1 flex items-center group">
        <input
          name="cols"
          type="number"
          value={cols}
          onChange={(e) => onUpdate("cols", Number(e.target.value))}
          onBlur={(e) => handleBlur("cols", Number(e.target.value))}
          className={`w-full h-full pl-3 pr-10 outline-none font-bold text-sm bg-transparent focus:bg-zinc-50 transition-colors ${noArrowsClass}`}
        />
        <span className="absolute right-2 text-[9px] font-black text-zinc-300 uppercase pointer-events-none tracking-tighter">
          {colLabel}
        </span>
      </div>
    </div>
  );
}
