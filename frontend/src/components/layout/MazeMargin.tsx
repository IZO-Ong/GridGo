"use client";
import { useEffect, useRef, useState } from "react";
import { usePathname } from "next/navigation";

const seededRandom = (r: number, c: number) => {
  let h = (r * 0x45d9f3b) ^ c;
  h = ((h >>> 16) ^ h) * 0x45d9f3b;
  h = ((h >>> 16) ^ h) * 0x45d9f3b;
  h = (h >>> 16) ^ h;
  return h % 100 > 50;
};

export default function MazeMargin({ side }: { side: "left" | "right" }) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const pathname = usePathname();
  const [opacity, setOpacity] = useState(1);

  const sideClasses = {
    left: "left-0 border-r-2",
    right: "right-0 border-l-2",
  };

  useEffect(() => {
    const handleZoomOpacity = () => {
      const width = window.innerWidth;
      const startFade = 1200;
      const endFade = 1100;

      if (width < endFade) {
        setOpacity(0);
      } else if (width > startFade) {
        setOpacity(1);
      } else {
        setOpacity((width - endFade) / (startFade - endFade));
      }
    };

    window.addEventListener("resize", handleZoomOpacity);
    handleZoomOpacity();
    return () => window.removeEventListener("resize", handleZoomOpacity);
  }, []);

  useEffect(() => {
    const drawMaze = () => {
      const canvas = canvasRef.current;
      if (!canvas) return;
      const ctx = canvas.getContext("2d");
      if (!ctx) return;

      const marginWidth = 70;
      const cols = 7;
      const cellSize = marginWidth / cols;
      const BUFFER_ROWS = 3; // Added extra rows to ensure full coverage

      const content = document.getElementById("main-content-wrapper");
      const pageHeight = content
        ? Math.max(content.offsetHeight, window.innerHeight)
        : window.innerHeight;

      // Calculate rows and add the buffer
      const rows = Math.ceil(pageHeight / cellSize) + BUFFER_ROWS;
      const fixedHeight = rows * cellSize;

      if (Math.abs(canvas.height - fixedHeight) < cellSize) return;

      canvas.width = marginWidth;
      canvas.height = fixedHeight;

      ctx.strokeStyle = "black";
      ctx.lineWidth = 1.5;
      ctx.lineCap = "square";
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
          const x = c * cellSize;
          const y = r * cellSize;
          const carveNorth = seededRandom(r, c);

          ctx.beginPath();
          if (r === 0 && c === cols - 1) continue;

          if (r === 0) {
            ctx.moveTo(x + cellSize, y);
            ctx.lineTo(x + cellSize, y + cellSize);
          } else if (c === cols - 1) {
            ctx.moveTo(x, y);
            ctx.lineTo(x + cellSize, y);
          } else {
            if (carveNorth) {
              ctx.moveTo(x, y);
              ctx.lineTo(x + cellSize, y);
            } else {
              ctx.moveTo(x + cellSize, y);
              ctx.lineTo(x + cellSize, y + cellSize);
            }
          }
          ctx.stroke();
        }
      }

      ctx.beginPath();
      const borderX = side === "left" ? canvas.width - 1 : 1;
      ctx.moveTo(borderX, 0);
      ctx.lineTo(borderX, canvas.height);
      ctx.stroke();
    };

    const observer = new ResizeObserver(() => drawMaze());
    const content = document.getElementById("main-content-wrapper");
    if (content) observer.observe(content);

    drawMaze();
    return () => observer.disconnect();
  }, [side, pathname]);

  return (
    <canvas
      ref={canvasRef}
      className={`absolute top-0 ${sideClasses[side]} pointer-events-none z-0 transition-opacity duration-150 border-black`}
      style={{ width: "70px", opacity: opacity }}
    />
  );
}
