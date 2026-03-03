"use client";
import { useState, useRef, useEffect, useCallback } from "react";
import { Maze } from "@/types";

export function useMazeCanvas(maze: Maze | null) {
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
        // Delta > 0 (scrolling down or pinching out) = zoom out
        // Delta < 0 (scrolling up or pinching in) = zoom in
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

  // --- UPDATED: Global drag listeners for Mouse AND Touch ---
  useEffect(() => {
    if (!isDragging) return;

    const handleMove = (x: number, y: number) => {
      setCssOffset({
        x: x - dragStart.current.x,
        y: y - dragStart.current.y,
      });
    };

    const handleUp = (x: number, y: number) => {
      commitDrag({
        x: x - dragStart.current.x,
        y: y - dragStart.current.y,
      });
    };

    // Mouse handlers
    const onMouseMove = (e: MouseEvent) => handleMove(e.clientX, e.clientY);
    const onMouseUp = (e: MouseEvent) => handleUp(e.clientX, e.clientY);

    // Touch handlers (1-finger only)
    const onTouchMove = (e: TouchEvent) => {
      if (e.touches.length === 1)
        handleMove(e.touches[0].clientX, e.touches[0].clientY);
    };
    const onTouchEnd = (e: TouchEvent) => {
      if (e.changedTouches.length === 1)
        handleUp(e.changedTouches[0].clientX, e.changedTouches[0].clientY);
    };

    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
    window.addEventListener("touchmove", onTouchMove, { passive: false });
    window.addEventListener("touchend", onTouchEnd);

    return () => {
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
      window.removeEventListener("touchmove", onTouchMove);
      window.removeEventListener("touchend", onTouchEnd);
    };
  }, [isDragging, commitDrag]);

  // Dynamic cell sizing
  useEffect(() => {
    if (!containerRef.current || !maze) return;
    const updateSize = () => {
      const parent = containerRef.current?.closest("section");
      if (!parent) return;
      // Mobile: reduce the subtraction from 128 to 32 to fill more of the screen
      const sub = window.innerWidth < 768 ? 32 : 128;
      setDynamicCellSize(
        Math.min(
          (parent.clientWidth - sub) / maze.cols,
          (parent.clientHeight - sub) / maze.rows
        )
      );
    };
    updateSize();
    window.addEventListener("resize", updateSize);
    return () => window.removeEventListener("resize", updateSize);
  }, [maze]);

  const onMouseDown = (e: React.MouseEvent | React.TouchEvent) => {
    if ("touches" in e) {
      if (e.touches.length !== 1) return;
      setIsDragging(true);
      dragStart.current = { x: e.touches[0].clientX, y: e.touches[0].clientY };
    } else {
      setIsDragging(true);
      dragStart.current = { x: e.clientX, y: e.clientY };
    }
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
