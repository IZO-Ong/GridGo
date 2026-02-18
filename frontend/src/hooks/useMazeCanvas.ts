"use client";
import { useState, useRef, useEffect, useCallback } from "react";
import { MazeData } from "@/types";

export function useMazeCanvas(maze: MazeData | null) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [dynamicCellSize, setDynamicCellSize] = useState(0);
  const [transform, setTransform] = useState({ s: 1, x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const dragStart = useRef({ x: 0, y: 0 });
  const [cssOffset, setCssOffset] = useState({ x: 0, y: 0 });

  const centerMaze = useCallback(
    (cellSize: number) => {
      const container = containerRef.current;
      if (!container || !maze || cellSize === 0) return;
      const viewW = container.clientWidth;
      const viewH = container.clientHeight;
      const mazeW = maze.cols * cellSize;
      const mazeH = maze.rows * cellSize;

      setTransform({
        s: 1,
        x: (viewW - mazeW) / 2,
        y: (viewH - mazeH) / 2,
      });
      setCssOffset({ x: 0, y: 0 });
    },
    [maze]
  );

  const commitDrag = useCallback((offset: { x: number; y: number }) => {
    setTransform((p) => ({ ...p, x: p.x + offset.x, y: p.y + offset.y }));
    setCssOffset({ x: 0, y: 0 });
    setIsDragging(false);
  }, []);

  const handleZoom = useCallback(
    (delta: number, mouseX?: number, mouseY?: number) => {
      setTransform((prev) => {
        const scaleFactor = delta > 0 ? 0.9 : 1.1;
        const newScale = Math.min(Math.max(prev.s * scaleFactor, 1), 20);
        if (newScale === prev.s) return prev;
        return mouseX !== undefined && mouseY !== undefined
          ? {
              s: newScale,
              x: mouseX - (mouseX - prev.x) * (newScale / prev.s),
              y: mouseY - (mouseY - prev.y) * (newScale / prev.s),
            }
          : { ...prev, s: newScale };
      });
      setCssOffset({ x: 0, y: 0 });
    },
    []
  );

  // Re-center
  useEffect(() => {
    if (dynamicCellSize > 0) centerMaze(dynamicCellSize);
  }, [maze, dynamicCellSize, centerMaze]);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;
    const handleNativeWheel = (e: WheelEvent) => {
      e.preventDefault();
      const rect = container.getBoundingClientRect();
      handleZoom(e.deltaY, e.clientX - rect.left, e.clientY - rect.top);
    };
    container.addEventListener("wheel", handleNativeWheel, { passive: false });
    return () => container.removeEventListener("wheel", handleNativeWheel);
  }, [handleZoom]);

  // Global drag listeners
  useEffect(() => {
    if (!isDragging) return;
    const handleGlobalMove = (e: MouseEvent) => {
      setCssOffset({
        x: e.clientX - dragStart.current.x,
        y: e.clientY - dragStart.current.y,
      });
    };
    const handleGlobalUp = (e: MouseEvent) => {
      commitDrag({
        x: e.clientX - dragStart.current.x,
        y: e.clientY - dragStart.current.y,
      });
    };
    window.addEventListener("mousemove", handleGlobalMove);
    window.addEventListener("mouseup", handleGlobalUp);
    return () => {
      window.removeEventListener("mousemove", handleGlobalMove);
      window.removeEventListener("mouseup", handleGlobalUp);
    };
  }, [isDragging, commitDrag]);

  // Handle dynamic sizing based on parent container
  useEffect(() => {
    if (!containerRef.current || !maze) return;
    const updateSize = () => {
      const parent = containerRef.current?.closest("section");
      if (!parent) return;
      setDynamicCellSize(
        Math.min(
          (parent.clientWidth - 128) / maze.cols,
          (parent.clientHeight - 128) / maze.rows
        )
      );
    };
    updateSize();
    window.addEventListener("resize", updateSize);
    return () => window.removeEventListener("resize", updateSize);
  }, [maze]);

  const onMouseDown = (e: React.MouseEvent) => {
    setIsDragging(true);
    dragStart.current = { x: e.clientX, y: e.clientY };
  };

  return {
    containerRef,
    dynamicCellSize,
    transform,
    isDragging,
    cssOffset,
    onMouseDown,
    handleZoom,
    centerMaze,
  };
}
