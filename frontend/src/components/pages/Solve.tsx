"use client";
import { useState, useEffect } from "react";
import { loadLastMaze } from "@/lib/db";
import MazeCanvas from "@/components/maze/MazeCanvas";
import AlgorithmSelect from "@/components/ui/AlgorithmSelect";
import { useMazeGeneration, MazeData } from "@/hooks/useMazeGeneration";
import { solveMaze } from "@/lib/api";

const SOLVE_ALGORITHMS = [
  { id: "astar", label: "A*_SEARCH" },
  { id: "bfs", label: "BREADTH_FIRST" },
  { id: "greedy", label: "GREEDY_SEARCH" },
];

export default function Solve() {
  const [activeMaze, setActiveMaze] = useState<MazeData | null>(null);
  const [solveType, setSolveType] = useState("astar");
  const [isSolving, setIsSolving] = useState(false);
  const [mazeId] = useState("G-7724-X");

  // Coordinate states for the input fields
  const [startPoint, setStartPoint] = useState<[number, number]>([0, 0]);
  const [endPoint, setEndPoint] = useState<[number, number]>([0, 0]);

  const [solution, setSolution] = useState<{
    visited: [number, number][];
    path: [number, number][];
  } | null>(null);

  // FIX: Sync state when the component mounts (handles tab switching/refresh)
  useEffect(() => {
    const fetchCurrentSession = async () => {
      const saved = await loadLastMaze();
      if (saved) {
        setActiveMaze(saved);
        // Explicitly set coordinates from the saved data
        setStartPoint(saved.start);
        setEndPoint(saved.end);
      }
    };
    fetchCurrentSession();
  }, []);

  // FIX: Sync state during manual "Load Last" button click
  const handleLoadLast = async () => {
    try {
      const saved = await loadLastMaze();
      if (saved) {
        setActiveMaze(saved);
        // Explicitly set coordinates from the saved data
        setStartPoint(saved.start);
        setEndPoint(saved.end);
        setSolution(null); // Clear previous visual paths
      } else {
        alert("NO_MAZE_FOUND_IN_STORAGE");
      }
    } catch (e) {
      console.error("LOAD_FAILURE:", e);
    }
  };

  const handleRunSolver = async () => {
    if (!activeMaze) return;
    setIsSolving(true);
    setSolution(null);

    // Use current input values as overrides for the solver
    const mazeToSolve = {
      ...activeMaze,
      start: startPoint,
      end: endPoint,
    };

    try {
      const data = await solveMaze(mazeToSolve, solveType);
      setSolution(data);
    } catch (err) {
      console.error("SOLVE_ERROR:", err);
    } finally {
      setIsSolving(false);
    }
  };

  return (
    <div className="space-y-8">
      <form
        className="grid grid-cols-12 gap-6 items-end"
        onSubmit={(e) => e.preventDefault()}
      >
        <div className="col-span-4 space-y-2">
          <label className="block font-bold uppercase tracking-widest text-[10px]">
            Maze_Reference
          </label>
          <div className="flex border-2 border-black h-[38px] divide-x-2 divide-black bg-white">
            <div className="flex-1 flex items-center px-3 font-mono text-[10px] font-bold text-zinc-400">
              ID: {mazeId}
            </div>
            <button
              type="button"
              onClick={handleLoadLast}
              className="px-3 bg-white hover:bg-black hover:text-white transition-colors text-[9px] font-black uppercase tracking-tighter cursor-pointer"
            >
              Load_Last
            </button>
          </div>
        </div>

        <div className="col-span-2 space-y-2">
          <label className="block font-bold uppercase tracking-widest text-[10px]">
            Start [R, C]
          </label>
          <div className="flex border-2 border-black h-[38px] divide-x-2 divide-black bg-white">
            <input
              type="number"
              value={startPoint[0]}
              onChange={(e) =>
                setStartPoint([parseInt(e.target.value) || 0, startPoint[1]])
              }
              className="w-full text-center font-bold text-xs outline-none"
            />
            <input
              type="number"
              value={startPoint[1]}
              onChange={(e) =>
                setStartPoint([startPoint[0], parseInt(e.target.value) || 0])
              }
              className="w-full text-center font-bold text-xs outline-none"
            />
          </div>
        </div>

        <div className="col-span-2 space-y-2">
          <label className="block font-bold uppercase tracking-widest text-[10px]">
            End [R, C]
          </label>
          <div className="flex border-2 border-black h-[38px] divide-x-2 divide-black bg-white">
            <input
              type="number"
              value={endPoint[0]}
              onChange={(e) =>
                setEndPoint([parseInt(e.target.value) || 0, endPoint[1]])
              }
              className="w-full text-center font-bold text-xs outline-none"
            />
            <input
              type="number"
              value={endPoint[1]}
              onChange={(e) =>
                setEndPoint([endPoint[0], parseInt(e.target.value) || 0])
              }
              className="w-full text-center font-bold text-xs outline-none"
            />
          </div>
        </div>

        <div className="col-span-4 space-y-2">
          <label className="block font-bold uppercase tracking-widest text-[10px]">
            Solve_Algorithm
          </label>
          <AlgorithmSelect
            value={solveType}
            onChange={setSolveType}
            options={SOLVE_ALGORITHMS}
          />
        </div>

        <button
          type="button"
          onClick={handleRunSolver}
          disabled={isSolving || !activeMaze}
          className={`col-span-12 border-2 border-black p-4 font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all flex justify-start pl-8 ${
            isSolving || !activeMaze
              ? "bg-zinc-100 text-zinc-400 opacity-50"
              : "bg-white hover:bg-black hover:text-white active:shadow-none active:translate-y-1 cursor-pointer"
          }`}
        >
          {isSolving ? ">>> EXECUTING PATHFINDER..." : ">>> RUN SOLVER"}
        </button>
      </form>

      <section className="relative border-4 border-black h-[750px] bg-zinc-50 overflow-hidden shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] flex flex-col">
        <div className="h-7 border-b-2 border-black bg-white flex items-center px-3 justify-between z-30 shrink-0">
          <span className="text-[10px] font-bold tracking-widest uppercase">
            SOLVER_OUTPUT{" "}
            {activeMaze ? `[${activeMaze.rows}X${activeMaze.cols}]` : ""}
          </span>
          <div className="flex gap-4">
            <span className="text-[10px] opacity-30 font-bold uppercase">
              VISITED: {solution?.visited?.length ?? "--"}
            </span>
            <span className="text-[10px] opacity-30 font-bold uppercase">
              PATH: {solution?.path?.length ?? "--"}
            </span>
          </div>
        </div>

        <div className="relative flex-1 bg-white overflow-hidden">
          <div className="h-full w-full flex items-center justify-center">
            {activeMaze ? (
              <MazeCanvas
                maze={activeMaze}
                showSave={false}
                highlights={solution?.visited}
                solutionPath={solution?.path}
              />
            ) : (
              <p className="opacity-20 tracking-[0.5em] font-bold uppercase text-2xl">
                Load a maze to solve
              </p>
            )}
          </div>
        </div>
      </section>
    </div>
  );
}
