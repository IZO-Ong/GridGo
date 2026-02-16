"use client";
import { useState, useEffect, useCallback } from "react";
import {
  loadGenerateSession,
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

  const [solveType, setSolveType] = useState("astar");
  const [isSolving, setIsSolving] = useState(false);
  const [isAnimating, setIsAnimating] = useState(false);
  const [mazeId, setMazeId] = useState("G-7724-X");

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
      alert(`COULD NOT FIND REFERENCE: ${mazeId}. Ensure the ID is valid.`);
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

  const handleLoadLast = async () => {
    try {
      const saved = await loadGenerateSession();
      if (saved) {
        setActiveMaze(saved);
        setStartPoint(saved.start as [number, number]);
        setEndPoint(saved.end as [number, number]);
        setSolution(null);
        setIsAnimating(false);
        await saveSolveSession(saved);
      }
    } catch (e) {
      console.error("LOAD_ERROR:", e);
    }
  };

  return (
    <div className="space-y-8">
      <SolveControls
        mazeId={mazeId}
        setMazeId={setMazeId}
        handleLoadLast={handleLoadLast}
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

      <section className="relative border-4 border-black h-[750px] bg-zinc-50 overflow-hidden shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] flex flex-col">
        <div className="h-7 border-b-2 border-black bg-white flex items-center px-3 justify-between z-30 shrink-0 uppercase text-[10px] font-bold">
          <span>
            SOLVER_OUTPUT{" "}
            {activeMaze && `[${activeMaze.rows}X${activeMaze.cols}]`}
          </span>
          <div className="flex gap-4 opacity-30 font-mono text-[9px]">
            <span>VISITED: {solution?.visited?.length ?? "--"}</span>
            <span>PATH: {solution?.path?.length ?? "--"}</span>
          </div>
        </div>
        <div className="relative flex-1 bg-white overflow-hidden flex items-center justify-center">
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
