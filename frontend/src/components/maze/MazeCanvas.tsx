"use client";
import { useEffect, useRef, useState, useMemo } from "react";
import { MazeData } from "@/types";
import { useMazeCanvas } from "@/hooks/useMazeCanvas";
import { renderMazeImage } from "@/lib/api";

const PADDING = 800;

interface MazeCanvasProps {
  maze: MazeData;
  showSave?: boolean;
  showShare?: boolean;
  highlights?: [number, number][];
  solutionPath?: [number, number][];
  overrideStart?: [number, number];
  overrideEnd?: [number, number];
  isPaused?: boolean;
  onComplete?: () => void;
}

export default function MazeCanvas({
  maze,
  showSave = false,
  showShare = false,
  highlights = [],
  solutionPath = [],
  overrideStart,
  overrideEnd,
  isPaused = false,
  onComplete,
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

  const totalNodes = (highlights?.length || 0) + (solutionPath?.length || 0);
  const stepSize = useMemo(() => Math.max(1, totalNodes / 540), [totalNodes]);

  useEffect(() => {
    setVisibleHighlights(0);
    setVisibleSolutionStep(0);
  }, [highlights, solutionPath]);

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

  const handleShare = () => {
    if (!maze.id) return;
    const url = `${window.location.origin}/solve?id=${maze.id}`;
    navigator.clipboard.writeText(url);
    alert("ACCESS_LINK_COPIED_TO_CLIPBOARD");
  };

  useEffect(() => {
    if (!highlights?.length || isPaused) return;

    let frame: number;
    let lastTime = 0;

    const animate = (time: number) => {
      if (!lastTime) lastTime = time;
      const delta = time - lastTime;

      if (delta > 16) {
        setVisibleHighlights((prevH) => {
          if (prevH < highlights.length) {
            return Math.min(prevH + stepSize, highlights.length);
          }

          setVisibleSolutionStep((prevS) => {
            const nextS = Math.min(prevS + stepSize, solutionPath.length);
            return nextS;
          });

          return prevH;
        });
        lastTime = time;
      }

      // Trigger completion after the UI has processed the update
      if (
        visibleHighlights >= highlights.length &&
        visibleSolutionStep >= solutionPath.length &&
        onComplete
      ) {
        onComplete();
        return; // Stop animation loop
      }

      frame = requestAnimationFrame(animate);
    };

    frame = requestAnimationFrame(animate);
    return () => cancelAnimationFrame(frame);
  }, [
    highlights,
    solutionPath,
    isPaused,
    stepSize,
    onComplete,
    visibleHighlights,
    visibleSolutionStep,
  ]);

  const highlightPath = useMemo(() => {
    const path = new Path2D();
    if (!dynamicCellSize || !highlights) return path;
    const currentBatch = highlights.slice(0, Math.floor(visibleHighlights));
    currentBatch.forEach(([r, c]) => {
      path.rect(
        c * dynamicCellSize,
        r * dynamicCellSize,
        dynamicCellSize,
        dynamicCellSize
      );
    });
    return path;
  }, [visibleHighlights, highlights, dynamicCellSize]);

  useEffect(() => {
    const canvas = canvasRef.current;
    const container = containerRef.current;
    if (!canvas || !container || !maze || dynamicCellSize === 0) return;
    const ctx = canvas.getContext("2d", { alpha: false });
    if (!ctx) return;

    canvas.width = container.clientWidth + PADDING * 2;
    canvas.height = container.clientHeight + PADDING * 2;
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    ctx.save();
    ctx.translate(transform.x + PADDING, transform.y + PADDING);
    ctx.scale(transform.s, transform.s);

    const cellSize = dynamicCellSize;
    const sPoint = overrideStart || maze.start;
    const ePoint = overrideEnd || maze.end;

    ctx.fillStyle = "rgba(113, 113, 122, 0.4)";
    ctx.fill(highlightPath);

    if (visibleHighlights >= highlights.length && solutionPath.length > 0) {
      ctx.strokeStyle = "#ef4444";
      ctx.lineWidth = cellSize * 0.4;
      ctx.lineCap = "round";
      ctx.lineJoin = "round";
      ctx.beginPath();
      const currentPath = solutionPath.slice(
        0,
        Math.floor(visibleSolutionStep)
      );
      currentPath.forEach(([r, c], idx) => {
        const x = c * cellSize + cellSize / 2;
        const y = r * cellSize + cellSize / 2;
        if (idx === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
      });
      ctx.stroke();
    }

    const getWallColor = (w: number) =>
      w >= 255
        ? "black"
        : `rgb(${Math.floor(230 - w * (230 / 255))},${Math.floor(230 - w * (230 / 255))},${Math.floor(230 - w * (230 / 255))})`;
    ctx.fillStyle = "#90ee90";
    ctx.fillRect(
      sPoint[1] * cellSize,
      sPoint[0] * cellSize,
      cellSize,
      cellSize
    );
    ctx.fillStyle = "#ff6347";
    ctx.fillRect(
      ePoint[1] * cellSize,
      ePoint[0] * cellSize,
      cellSize,
      cellSize
    );

    const wallBatches: Record<string, Path2D> = {};
    for (let r = 0; r < maze.rows; r++) {
      for (let c = 0; c < maze.cols; c++) {
        const x = c * cellSize,
          y = r * cellSize;
        maze.grid[r][c].walls.forEach((w, i) => {
          if (w) {
            const color = getWallColor(maze.grid[r][c].wall_weights[i]);
            if (!wallBatches[color]) wallBatches[color] = new Path2D();
            const p = wallBatches[color];
            if (i === 0) {
              p.moveTo(x, y);
              p.lineTo(x + cellSize, y);
            }
            if (i === 1) {
              p.moveTo(x + cellSize, y);
              p.lineTo(x + cellSize, y + cellSize);
            }
            if (i === 2) {
              p.moveTo(x, y + cellSize);
              p.lineTo(x + cellSize, y + cellSize);
            }
            if (i === 3) {
              p.moveTo(x, y);
              p.lineTo(x, y + cellSize);
            }
          }
        });
      }
    }
    ctx.lineWidth = cellSize > 5 ? 1 : 0.5;
    Object.entries(wallBatches).forEach(([color, path]) => {
      ctx.strokeStyle = color;
      ctx.stroke(path);
    });
    ctx.restore();
  }, [
    maze,
    dynamicCellSize,
    transform,
    highlightPath,
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
          <canvas
            id="main-maze-canvas"
            ref={canvasRef}
            className="block select-none"
          />
        </div>
      </div>

      {(showSave || showShare) && (
        <div className="absolute bottom-6 left-6 border-2 border-black bg-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] z-30">
          {showSave ? (
            <button
              onClick={handleSave}
              title="Save to PNG"
              className="p-3 hover:bg-black hover:text-white transition-colors cursor-pointer"
            >
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="3"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
                <polyline points="17 21 17 13 7 13 7 21" />
                <polyline points="7 3 7 8 15 8" />
              </svg>
            </button>
          ) : (
            <button
              onClick={handleShare}
              title="Share Matrix Link"
              className="p-3 hover:bg-black hover:text-white transition-colors cursor-pointer"
            >
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="currentColor"
                stroke="currentColor"
                strokeWidth="1"
              >
                <path d="M15,5 L15,9 C7,9 4,14 4,20 C7,15 11,13 15,13 L15,17 L22,10 L15,5 Z" />
              </svg>
            </button>
          )}
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
