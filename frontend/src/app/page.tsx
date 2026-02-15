"use client";
import { useState } from "react";
import MazeCanvas from "@/components/maze/MazeCanvas";
import MazeControls from "@/components/maze/MazeControls";
import { useImageDimensions } from "@/hooks/useImageDimensions";
import { useMazeGeneration } from "@/hooks/useMazeGeneration";

const ALGORITHMS = [
  { id: "image", label: "IMAGE_KRUSKAL" },
  { id: "kruskal", label: "RANDOM_KRUSKAL" },
  { id: "recursive", label: "DFS_BACKTRACKER" },
];

export default function Home() {
  const { maze, loading, error, executeGeneration, setError } =
    useMazeGeneration();
  const [genType, setGenType] = useState("image");
  const [hasImage, setHasImage] = useState(false);
  const { dims, updateDim, clampDimensions, handleImageChange } =
    useImageDimensions();

  const onImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleImageChange(e);
    setHasImage(!!e.target.files?.[0]);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clampDimensions();
    const formData = new FormData(e.currentTarget);
    await executeGeneration(formData);
  };

  return (
    <div className="space-y-8">
      <MazeControls
        genType={genType}
        setGenType={setGenType}
        dims={dims}
        updateDim={updateDim}
        onImageChange={onImageChange}
        onSubmit={handleSubmit}
        loading={loading}
        algorithms={ALGORITHMS}
        isSubmitDisabled={loading || (genType === "image" && !hasImage)}
      />

      {error && (
        <div className="p-3 bg-red-50 border-2 border-red-600 text-red-600 font-bold uppercase">
          {`>> ERROR: ${error}`}
        </div>
      )}

      <section className="relative border-4 border-black h-[750px] bg-zinc-50 overflow-hidden shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] flex flex-col">
        {/* Header */}
        <div className="h-7 border-b-2 border-black bg-white flex items-center px-3 justify-between z-30 shrink-0">
          <span className="text-[10px] font-bold tracking-widest uppercase">
            MAZE_OUTPUT
          </span>
          <span className="text-[10px] opacity-30 font-bold uppercase">
            DIM:{" "}
            {maze ? `${maze.rows}x${maze.cols}` : `${dims.rows}x${dims.cols}`}
          </span>
        </div>

        {/* Maze Container */}
        <div className="relative flex-1 bg-white overflow-hidden">
          <div className="absolute inset-0 bg-[radial-gradient(#000000_1px,transparent_1px)] [background-size:32px_32px] opacity-[0.05] pointer-events-none" />
          {maze ? (
            <MazeCanvas maze={maze} />
          ) : (
            <div className="h-full w-full flex items-center justify-center">
              <p className="opacity-20 tracking-[0.5em] font-bold uppercase text-2xl">
                Generate a Maze!
              </p>
            </div>
          )}
        </div>
      </section>
    </div>
  );
}
