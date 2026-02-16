"use client";
import { useEffect, useRef, useState } from "react";
import { MazeData } from "@/hooks/useMazeGeneration";
import { useMazeCanvas } from "@/hooks/useMazeCanvas";
import { renderMazeImage } from "@/lib/api";

const PADDING = 800; // Used for panning overflow

interface MazeCanvasProps {
  maze: MazeData;
  showSave?: boolean;
  highlights?: [number, number][];
  solutionPath?: [number, number][];
  // Overrides for real-time start/end marker updates
  overrideStart?: [number, number];
  overrideEnd?: [number, number];
  isPaused?: boolean;
}

export default function MazeCanvas({
  maze,
  showSave = true,
  highlights = [],
  solutionPath = [],
  overrideStart,
  overrideEnd,
  isPaused = false,
}: MazeCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [visibleHighlights, setVisibleHighlights] = useState<number>(0);
  const [visibleSolutionStep, setVisibleSolutionStep] = useState(0);

  const {
    containerRef,
    dynamicCellSize,
    transform,
    cssOffset,
    onMouseDown,
    handleZoom,
    centerMaze,
  } = useMazeCanvas(maze);

  // Dynamic Speed: 12 seconds = 720 frames at 60fps
  const totalNodes = highlights.length + solutionPath.length;
  const stepSize = Math.max(1, Math.ceil(totalNodes / 720));

  const handleSave = async () => {
    try {
      const blob = await renderMazeImage(maze);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `maze_${maze.rows}x${maze.cols}_${Date.now()}.png`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (e: any) {
      alert(`FAILED TO SAVE: ${e.message}`);
    }
  };

  useEffect(() => {
    setVisibleHighlights(0);
    setVisibleSolutionStep(0);
  }, [highlights, solutionPath]);

  useEffect(() => {
    if (!highlights?.length || isPaused) return;

    let frame: number;
    let lastTime = 0;

    const animate = (time: number) => {
      if (!lastTime) lastTime = time;
      if (time - lastTime > 16) {
        setVisibleHighlights((prev) => {
          if (prev < highlights.length)
            return Math.min(prev + stepSize, highlights.length);

          setVisibleSolutionStep((solPrev) => {
            if (solPrev < solutionPath.length)
              return Math.min(solPrev + stepSize, solutionPath.length);
            return solPrev;
          });
          return prev;
        });
        lastTime = time;
      }
      frame = requestAnimationFrame(animate);
    };

    frame = requestAnimationFrame(animate);
    return () => cancelAnimationFrame(frame);
  }, [highlights, solutionPath, isPaused, stepSize]);

  useEffect(() => {
    const canvas = canvasRef.current;
    const container = containerRef.current;
    if (!canvas || !container || !maze || dynamicCellSize === 0) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    // Set internal canvas size to match visible container + panning padding
    canvas.width = container.clientWidth + PADDING * 2;
    canvas.height = container.clientHeight + PADDING * 2;

    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.save();
    ctx.translate(transform.x + PADDING, transform.y + PADDING);
    ctx.scale(transform.s, transform.s);

    const cellSize = dynamicCellSize;

    // 1. Draw Visited Path
    ctx.fillStyle = "rgba(167, 139, 250, 0.4)";
    for (let i = 0; i < visibleHighlights; i++) {
      const point = highlights?.[i];
      if (point)
        ctx.fillRect(
          point[1] * cellSize,
          point[0] * cellSize,
          cellSize,
          cellSize
        );
    }

    // 2. Draw Solution Path
    if (visibleHighlights >= highlights.length && solutionPath.length > 0) {
      ctx.strokeStyle = "#ef4444";
      ctx.lineWidth = cellSize * 0.4;
      ctx.lineCap = "round";
      ctx.lineJoin = "round";
      ctx.beginPath();
      const currentPath = solutionPath.slice(0, visibleSolutionStep);
      currentPath.forEach(([r, c], idx) => {
        const x = c * cellSize + cellSize / 2;
        const y = r * cellSize + cellSize / 2;
        if (idx === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
      });
      ctx.stroke();
    }

    // 3. Draw Maze Walls
    const getWallColor = (w: number) => {
      if (w >= 255) return "black";
      const brightness = Math.floor(230 - w * (230 / 255));
      return `rgb(${brightness}, ${brightness}, ${brightness})`;
    };

    const sPoint = overrideStart || maze.start;
    const ePoint = overrideEnd || maze.end;

    for (let r = 0; r < maze.rows; r++) {
      for (let c = 0; c < maze.cols; c++) {
        const x = c * cellSize;
        const y = r * cellSize;
        if (r === sPoint[0] && c === sPoint[1]) {
          ctx.fillStyle = "#90ee90";
          ctx.fillRect(x, y, cellSize, cellSize);
        } else if (r === ePoint[0] && c === ePoint[1]) {
          ctx.fillStyle = "#ff6347";
          ctx.fillRect(x, y, cellSize, cellSize);
        }

        maze.grid[r][c].walls.forEach((w, i) => {
          if (w) {
            ctx.strokeStyle = getWallColor(maze.grid[r][c].wall_weights[i]);
            ctx.lineWidth = cellSize > 5 ? 1 : 0.5;
            ctx.beginPath();
            if (i === 0) {
              ctx.moveTo(x, y);
              ctx.lineTo(x + cellSize, y);
            }
            if (i === 1) {
              ctx.moveTo(x + cellSize, y);
              ctx.lineTo(x + cellSize, y + cellSize);
            }
            if (i === 2) {
              ctx.moveTo(x, y + cellSize);
              ctx.lineTo(x + cellSize, y + cellSize);
            }
            if (i === 3) {
              ctx.moveTo(x, y);
              ctx.lineTo(x, y + cellSize);
            }
            ctx.stroke();
          }
        });
      }
    }
    ctx.restore();
  }, [
    maze,
    dynamicCellSize,
    transform,
    visibleHighlights,
    visibleSolutionStep,
    overrideStart,
    overrideEnd,
  ]);

  return (
    <div className="relative w-full h-full flex items-center justify-center p-8">
      <div
        ref={containerRef}
        className="w-full h-full relative overflow-hidden border-2 border-black bg-white cursor-grab active:cursor-grabbing"
        onMouseDown={onMouseDown}
      >
        <div
          style={{
            transform: `translate3d(${cssOffset.x - PADDING}px, ${cssOffset.y - PADDING}px, 0)`,
            willChange: "transform",
          }}
        >
          <canvas ref={canvasRef} className="block select-none" />
        </div>
        {/* Brackets */}
        <div className="absolute top-0 left-0 w-8 h-8 border-t-2 border-l-2 border-black z-20 pointer-events-none" />
        <div className="absolute top-0 right-0 w-8 h-8 border-t-2 border-r-2 border-black z-20 pointer-events-none" />
        <div className="absolute bottom-0 left-0 w-8 h-8 border-b-2 border-l-2 border-black z-20 pointer-events-none" />
        <div className="absolute bottom-0 right-0 w-8 h-8 border-b-2 border-r-2 border-black z-20 pointer-events-none" />
      </div>

      {showSave && (
        <div className="absolute bottom-6 left-6 border-2 border-black bg-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] z-30">
          <button
            onClick={handleSave}
            className="p-3 hover:bg-black hover:text-white transition-colors cursor-pointer"
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="3"
            >
              <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
              <polyline points="17 21 17 13 7 13 7 21" />
              <polyline points="7 3 7 8 15 8" />
            </svg>
          </button>
        </div>
      )}

      <div className="absolute bottom-6 right-6 flex flex-col border-2 border-black bg-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] divide-y-2 divide-black z-30">
        <button
          onClick={() => {
            const r = containerRef.current?.getBoundingClientRect();
            if (r) handleZoom(-1, r.width / 2, r.height / 2);
          }}
          className="p-3 hover:bg-black hover:text-white font-bold text-lg"
        >
          +
        </button>
        <button
          onClick={() => {
            const r = containerRef.current?.getBoundingClientRect();
            if (r) handleZoom(1, r.width / 2, r.height / 2);
          }}
          className="p-3 hover:bg-black hover:text-white font-bold text-lg"
        >
          -
        </button>
        <button
          onClick={() => centerMaze(dynamicCellSize)}
          className="p-2 text-[9px] hover:bg-black hover:text-white font-bold uppercase"
        >
          Reset
        </button>
      </div>
    </div>
  );
}
