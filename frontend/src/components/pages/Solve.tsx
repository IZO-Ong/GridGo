"use client";
import { useState, useEffect, useCallback } from "react";
import {
  saveSolveSession,
  loadSolveSession,
  savePreferences,
  loadPreferences,
} from "@/lib/db";
import MazeCanvas from "@/components/maze/MazeCanvas";
import SolveControls from "@/components/maze/SolveControls";
import { solveMaze, getMazeById } from "@/lib/api";
import { MazeData } from "@/hooks/useMazeGeneration";

const SOLVE_ALGORITHMS = [
  { id: "astar", label: "A*_SEARCH" },
  { id: "bfs", label: "BREADTH_FIRST" },
  { id: "greedy", label: "GREEDY_SEARCH" },
];

export default function Solve() {
  const [activeMaze, setActiveMaze] = useState<MazeData | null>(null);
  const [hasLoaded, setHasLoaded] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [solveType, setSolveType] = useState("astar");
  const [isSolving, setIsSolving] = useState(false);
  const [isAnimating, setIsAnimating] = useState(false);
  const [mazeId, setMazeId] = useState("");

  const [startPoint, setStartPoint] = useState<[number, number]>([0, 0]);
  const [endPoint, setEndPoint] = useState<[number, number]>([0, 0]);
  const [solution, setSolution] = useState<{
    visited: [number, number][];
    path: [number, number][];
  } | null>(null);

  useEffect(() => {
    const init = async () => {
      const savedMaze = await loadSolveSession();
      if (savedMaze) setActiveMaze(savedMaze);

      const prefs = await loadPreferences("solve_prefs");
      if (prefs) {
        if (SOLVE_ALGORITHMS.some((a) => a.id === prefs.solveType)) {
          setSolveType(prefs.solveType);
        }
        setStartPoint(prefs.startPoint);
        setEndPoint(prefs.endPoint);
        if (prefs.mazeId) setMazeId(prefs.mazeId);
      }
      setHasLoaded(true);
    };
    init();
  }, []);

  useEffect(() => {
    if (hasLoaded) {
      savePreferences("solve_prefs", {
        solveType,
        startPoint,
        endPoint,
        mazeId,
      });
      if (activeMaze) {
        saveSolveSession({ ...activeMaze, start: startPoint, end: endPoint });
      }
    }
  }, [solveType, startPoint, endPoint, mazeId, activeMaze, hasLoaded]);

  const validate = (val: number, max: number) =>
    Math.min(Math.max(0, val), max - 1);

  const handleAnimationComplete = useCallback(() => {
    setIsAnimating(false);
  }, []);

  const handleLoadID = async () => {
    if (!mazeId) return;
    setIsSolving(true);
    setError(null);
    try {
      const data = await getMazeById(mazeId);
      setActiveMaze(data);
      setStartPoint(data.start);
      setEndPoint(data.end);
      setSolution(null);
      setIsAnimating(false);
      await saveSolveSession(data);
    } catch (err) {
      console.error("CLOUD_FETCH_ERROR:", err);
      setError(`COULD_NOT_FIND_REFERENCE: ${mazeId}`);
    } finally {
      setIsSolving(false);
    }
  };

  const handleAction = async () => {
    if (isAnimating) {
      setIsAnimating(false);
      setSolution(null);
      return;
    }
    if (!activeMaze) return;
    setIsSolving(true);
    setIsAnimating(true);
    setSolution(null);

    try {
      const data = await solveMaze(
        { ...activeMaze, start: startPoint, end: endPoint },
        solveType
      );
      setSolution(data);
    } catch (err) {
      console.error("SOLVE_ERROR:", err);
      setIsAnimating(false);
    } finally {
      setIsSolving(false);
    }
  };

  return (
    <div className="space-y-8">
      {/* Greyed out logic handled by passing isSolving to controls */}
      <div
        className={
          isSolving ? "opacity-50 pointer-events-none transition-opacity" : ""
        }
      >
        <SolveControls
          mazeId={mazeId}
          setMazeId={setMazeId}
          handleLoadID={handleLoadID}
          startPoint={startPoint}
          setStartPoint={setStartPoint}
          endPoint={endPoint}
          setEndPoint={setEndPoint}
          solveType={solveType}
          setSolveType={setSolveType}
          handleAction={handleAction}
          activeMaze={activeMaze}
          isSolving={isSolving}
          isAnimating={isAnimating}
          algorithms={SOLVE_ALGORITHMS}
          validate={validate}
        />
      </div>

      {error && (
        <div className="p-3 bg-red-50 border-2 border-red-600 text-red-600 font-bold uppercase text-[11px]">
          {`>> ERROR_SEQUENCE: ${error}`}
        </div>
      )}

      <section className="relative border-4 border-black h-[750px] bg-zinc-50 overflow-hidden shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] flex flex-col">
        {/* Updated Header with Black ID Chip */}
        <div className="h-7 border-b-2 border-black bg-white flex items-center px-3 justify-between z-30 shrink-0 uppercase text-[10px] font-bold">
          <div className="flex items-center gap-3">
            <span>SOLVER_OUTPUT</span>
            {activeMaze?.id && (
              <span className="bg-black text-white px-2 py-0.5 text-[9px] font-black tracking-tighter">
                {activeMaze.id}
              </span>
            )}
          </div>
          <div className="flex gap-4 opacity-30 font-mono text-[9px]">
            <span>
              DIM: {activeMaze ? `${activeMaze.rows}X${activeMaze.cols}` : "--"}
            </span>
            <span>VISITED: {solution?.visited?.length ?? "--"}</span>
            <span>PATH: {solution?.path?.length ?? "--"}</span>
          </div>
        </div>

        <div className="relative flex-1 bg-white overflow-hidden flex items-center justify-center">
          {isSolving && (
            <div className="absolute inset-0 bg-white/50 z-20 flex items-center justify-center backdrop-blur-[1px]">
              <span className="font-black italic animate-pulse">
                INITIALIZING_DATA_STREAM...
              </span>
            </div>
          )}

          {activeMaze ? (
            <MazeCanvas
              maze={activeMaze}
              showSave={false}
              highlights={solution?.visited}
              solutionPath={solution?.path}
              overrideStart={startPoint}
              overrideEnd={endPoint}
              isPaused={!isAnimating}
              onComplete={handleAnimationComplete}
            />
          ) : (
            <div className="h-full w-full flex items-center justify-center">
              <p className="opacity-20 tracking-[0.5em] font-bold uppercase text-2xl text-center px-12">
                Load a maze to solve
              </p>
            </div>
          )}
        </div>
      </section>
    </div>
  );
}
