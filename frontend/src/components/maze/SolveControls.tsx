"use client";
import AlgorithmSelect from "@/components/ui/AlgorithmSelect";
import GridDimensionsInput from "@/components/ui/GridDimensionsInput";
import { MazeData } from "@/hooks/useMazeGeneration";

interface SolveControlsProps {
  mazeId: string;
  setMazeId: (id: string) => void;
  handleLoadID: () => void; // Removed handleLoadLast
  startPoint: [number, number];
  setStartPoint: (p: [number, number]) => void;
  endPoint: [number, number];
  setEndPoint: (p: [number, number]) => void;
  solveType: string;
  setSolveType: (type: string) => void;
  handleAction: () => void;
  activeMaze: MazeData | null;
  isSolving: boolean;
  isAnimating: boolean;
  algorithms: { id: string; label: string }[];
  validate: (val: number, max: number) => number;
}

export default function SolveControls({
  mazeId,
  setMazeId,
  handleLoadID,
  startPoint,
  setStartPoint,
  endPoint,
  setEndPoint,
  solveType,
  setSolveType,
  handleAction,
  activeMaze,
  isSolving,
  isAnimating,
  algorithms,
  validate,
}: SolveControlsProps) {
  return (
    <form
      className="grid grid-cols-12 gap-4 items-end"
      onSubmit={(e) => {
        e.preventDefault();
        handleLoadID();
      }}
    >
      <div className="col-span-3 space-y-2">
        <label className="block font-bold uppercase tracking-widest text-[9px]">
          Reference_ID
        </label>
        {/* Single-button cloud load interface */}
        <div className="flex border-2 border-black h-[38px] divide-x-2 divide-black bg-white">
          <input
            type="text"
            value={mazeId}
            onChange={(e) => setMazeId(e.target.value.toUpperCase())}
            placeholder="ENTER_ID"
            className="flex-1 min-w-0 px-3 outline-none font-mono font-bold text-[10px] bg-transparent uppercase"
          />
          <button
            type="button"
            onClick={handleLoadID}
            className="px-4 bg-white hover:bg-black hover:text-white transition-colors text-[9px] font-black uppercase cursor-pointer shrink-0"
          >
            LOAD
          </button>
        </div>
      </div>

      <div className="col-span-3 space-y-2">
        <label className="block font-bold uppercase tracking-widest text-[9px]">
          Start [R, C]
        </label>
        <GridDimensionsInput
          rows={startPoint[0]}
          cols={startPoint[1]}
          labelOverride={{ row: "Row", col: "Col" }}
          onUpdate={(dim, val) =>
            setStartPoint([
              dim === "rows" ? val : startPoint[0],
              dim === "cols" ? val : startPoint[1],
            ])
          }
          min={0}
          max={{
            rows: (activeMaze?.rows || 1) - 1,
            cols: (activeMaze?.cols || 1) - 1,
          }}
        />
      </div>

      <div className="col-span-3 space-y-2">
        <label className="block font-bold uppercase tracking-widest text-[9px]">
          End [R, C]
        </label>
        <GridDimensionsInput
          rows={endPoint[0]}
          cols={endPoint[1]}
          labelOverride={{ row: "Row", col: "Col" }}
          onUpdate={(dim, val) =>
            setEndPoint([
              dim === "rows"
                ? validate(val, activeMaze?.rows || 1)
                : endPoint[0],
              dim === "cols"
                ? validate(val, activeMaze?.cols || 1)
                : endPoint[1],
            ])
          }
        />
      </div>

      <div className="col-span-3 space-y-2">
        <label className="block font-bold uppercase tracking-widest text-[9px]">
          Algorithm
        </label>
        <AlgorithmSelect
          value={solveType}
          onChange={setSolveType}
          options={algorithms}
        />
      </div>

      <button
        type="button"
        onClick={handleAction}
        disabled={!activeMaze && !isAnimating}
        className={`col-span-12 border-2 border-black p-4 font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all flex items-center justify-between px-8 group
        ${
          isAnimating
            ? "bg-black text-white cursor-pointer active:translate-y-1 active:shadow-none"
            : !activeMaze
              ? "bg-zinc-100 text-zinc-400 opacity-50"
              : "bg-white hover:bg-black hover:text-white cursor-pointer active:translate-y-1 active:shadow-none"
        }`}
      >
        <span className="italic tracking-tighter text-lg font-black uppercase">
          {isSolving
            ? ">>> INITIALIZING..."
            : isAnimating
              ? ">>> TERMINATE_SOLVE"
              : ">>> RUN_SOLVER"}
        </span>
        {isAnimating && (
          <span className="w-2 h-2 bg-red-500 rounded-full animate-pulse" />
        )}
      </button>
    </form>
  );
}
