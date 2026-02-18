"use client";
import { useState, useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import {
  saveSolveSession,
  loadSolveSession,
  savePreferences,
  loadPreferences,
} from "@/lib/db";
import MazeCanvas from "@/components/maze/MazeCanvas";
import SolveControls from "@/components/maze/SolveControls";
import { solveMaze, getMazeById } from "@/lib/api";
import { MazeData } from "@/types";

const SOLVE_ALGORITHMS = [
  { id: "astar", label: "A*_SEARCH" },
  { id: "bfs", label: "BREADTH_FIRST" },
  { id: "greedy", label: "GREEDY_SEARCH" },
];

export default function SolvePage() {
  return (
    <Suspense
      fallback={
        <div className="p-20 font-black italic animate-pulse">
          BOOTING_SOLVER_SUBSYSTEM...
        </div>
      }
    >
      <SolveCore />
    </Suspense>
  );
}

function SolveCore() {
  const searchParams = useSearchParams();
  const urlId = searchParams.get("id");

  const [activeMaze, setActiveMaze] = useState<MazeData | null>(null);
  const [hasLoaded, setHasLoaded] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [solveType, setSolveType] = useState("astar");
  const [isSolving, setIsSolving] = useState(false);
  const [isAnimating, setIsAnimating] = useState(false);
  const [mazeId, setMazeId] = useState(urlId || "");

  const [startPoint, setStartPoint] = useState<[number, number]>([0, 0]);
  const [endPoint, setEndPoint] = useState<[number, number]>([0, 0]);
  const [solution, setSolution] = useState<{
    visited: [number, number][];
    path: [number, number][];
  } | null>(null);

  useEffect(() => {
    const init = async () => {
      if (urlId) {
        await handleLoadID(urlId); // Load specifically from URL param
      } else {
        const savedMaze = await loadSolveSession();
        if (savedMaze) {
          setActiveMaze(savedMaze);
          setStartPoint(savedMaze.start);
          setEndPoint(savedMaze.end);
        }
      }

      const prefs = await loadPreferences("solve_prefs");
      if (prefs) {
        if (SOLVE_ALGORITHMS.some((a) => a.id === prefs.solveType))
          setSolveType(prefs.solveType);
      }
      setHasLoaded(true);
    };
    init();
  }, [urlId]);

  useEffect(() => {
    if (hasLoaded) {
      savePreferences("solve_prefs", {
        solveType,
        startPoint,
        endPoint,
        mazeId,
      });
      if (activeMaze)
        saveSolveSession({ ...activeMaze, start: startPoint, end: endPoint });
    }
  }, [solveType, startPoint, endPoint, mazeId, activeMaze, hasLoaded]);

  const handleLoadID = async (targetId?: string) => {
    const id = targetId || mazeId;
    if (!id) return;
    setIsSolving(true);
    setError(null);
    try {
      const data = await getMazeById(id);
      setActiveMaze(data);
      setStartPoint(data.start);
      setEndPoint(data.end);
      setSolution(null);
      setIsAnimating(false);
      await saveSolveSession(data);
    } catch (err) {
      setError(`COULD_NOT_FIND_REFERENCE: ${id}`);
    } finally {
      setIsSolving(false);
    }
  };

  const handleAction = async () => {
    if (isAnimating) {
      setSolution(null);
      setIsAnimating(false);
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
      setIsAnimating(false);
    } finally {
      setIsSolving(false);
    }
  };

  return (
    <div className="space-y-8">
      <div className={isSolving ? "opacity-50 pointer-events-none" : ""}>
        <SolveControls
          mazeId={mazeId}
          setMazeId={setMazeId}
          handleLoadID={() => handleLoadID()}
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
          validate={(v, m) => Math.min(Math.max(0, v), m - 1)}
        />
      </div>

      {error && (
        <div className="p-3 bg-red-50 border-2 border-red-600 text-red-600 font-bold uppercase text-[11px]">{`>> ${error}`}</div>
      )}

      <section className="relative border-4 border-black h-[750px] bg-zinc-50 overflow-hidden shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] flex flex-col">
        <div className="h-7 border-b-2 border-black bg-white flex items-center px-3 justify-between z-30 shrink-0 uppercase text-[10px] font-bold">
          <div className="flex items-center gap-3">
            <span>SOLVER_OUTPUT</span>
            {activeMaze?.id && (
              <span className="bg-black text-white px-2 py-0.5 text-[9px] font-black">
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
          {activeMaze ? (
            <MazeCanvas
              maze={activeMaze}
              showSave={false}
              showShare={true}
              highlights={solution?.visited}
              solutionPath={solution?.path}
              overrideStart={startPoint}
              overrideEnd={endPoint}
              isPaused={!isAnimating}
              onComplete={() => setIsAnimating(false)}
            />
          ) : (
            <p className="opacity-20 tracking-[0.5em] font-bold uppercase text-2xl">
              Load a maze to solve
            </p>
          )}
        </div>
      </section>
    </div>
  );
}
