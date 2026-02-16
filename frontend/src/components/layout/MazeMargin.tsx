"use client";
import { useEffect, useRef, useState } from "react";

const seededRandom = (r: number, c: number) => {
  let h = (r * 0x45d9f3b) ^ c;
  h = ((h >>> 16) ^ h) * 0x45d9f3b;
  h = ((h >>> 16) ^ h) * 0x45d9f3b;
  h = (h >>> 16) ^ h;
  return h % 100 > 50;
};

export default function MazeMargin({ side }: { side: "left" | "right" }) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [opacity, setOpacity] = useState(1);

  // Map the side to specific Tailwind classes so the compiler can find them
  const sideClasses = {
    left: "left-0",
    right: "right-0",
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

    const updateCanvas = () => {
      const canvas = canvasRef.current;
      if (!canvas) return;
      const ctx = canvas.getContext("2d");
      if (!ctx) return;

      const marginWidth = 70;
      const cols = 7;
      const cellSize = marginWidth / cols;

      const rawHeight = document.documentElement.scrollHeight;
      const fixedHeight = Math.ceil(rawHeight / cellSize) * cellSize;

      if (canvas.height !== fixedHeight) {
        canvas.width = marginWidth;
        canvas.height = fixedHeight;
      }

      ctx.strokeStyle = "black";
      ctx.lineWidth = 1.5;
      ctx.lineCap = "square";
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      const rows = fixedHeight / cellSize;

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

      // Vertical Sidebar Border
      ctx.beginPath();
      const borderX = side === "left" ? canvas.width - 1 : 1;
      ctx.moveTo(borderX, 0);
      ctx.lineTo(borderX, canvas.height);
      ctx.stroke();
    };

    handleZoomOpacity();
    updateCanvas();

    window.addEventListener("resize", handleZoomOpacity);
    const resizeObserver = new ResizeObserver(() => updateCanvas());
    resizeObserver.observe(document.body);

    return () => {
      window.removeEventListener("resize", handleZoomOpacity);
      resizeObserver.disconnect();
    };
  }, [side]);

  return (
    <canvas
      ref={canvasRef}
      className={`absolute top-0 ${sideClasses[side]} pointer-events-none z-0 transition-opacity duration-150 ease-out`}
      style={{
        width: "70px",
        opacity: opacity,
      }}
    />
  );
}
