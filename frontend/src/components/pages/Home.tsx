"use client";
import { useState, useEffect } from "react";
import { saveLastMaze, loadLastMaze } from "@/lib/db";
import MazeCanvas from "@/components/maze/MazeCanvas";
import MazeControls from "@/components/maze/MazeControls";
import { useImageDimensions } from "@/hooks/useImageDimensions";
import { useMazeGeneration, MazeData } from "@/hooks/useMazeGeneration";

const ALGORITHMS = [
  { id: "image", label: "IMAGE_KRUSKAL" },
  { id: "kruskal", label: "RANDOM_KRUSKAL" },
  { id: "recursive", label: "DFS_BACKTRACKER" },
];

export default function Home() {
  // We use a local state for the maze to allow manual loading from IndexedDB
  const [activeMaze, setActiveMaze] = useState<MazeData | null>(null);

  const { loading, error, executeGeneration } = useMazeGeneration();
  const [genType, setGenType] = useState("image");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const { dims, updateDim, clampDimensions, handleImageChange } =
    useImageDimensions();

  // Load the last generated maze from disk on component mount
  useEffect(() => {
    const initPersistence = async () => {
      const saved = await loadLastMaze();
      if (saved) {
        setActiveMaze(saved);
      }
    };
    initPersistence();
  }, []);

  const onImageChange = (e: React.ChangeEvent<HTMLInputElement> | File) => {
    let file: File | undefined;

    if (e instanceof File) {
      file = e;
      const fakeEvent = {
        target: {
          files: [file],
        },
      } as unknown as React.ChangeEvent<HTMLInputElement>;
      handleImageChange(fakeEvent);
    } else {
      file = e.target.files?.[0];
      handleImageChange(e);
    }

    if (file) {
      setSelectedFile(file);
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clampDimensions();

    const formData = new FormData(e.currentTarget);

    if (genType === "image" && selectedFile) {
      formData.set("image", selectedFile);
    }

    const newMaze = await executeGeneration(formData);

    if (newMaze) {
      setActiveMaze(newMaze);
      // Persist to IndexedDB (Supports large 300x300 grids)
      await saveLastMaze(newMaze);
    }
  };

  return (
    <div className="space-y-8">
      <MazeControls
        genType={genType}
        setGenType={setGenType}
        dims={dims}
        updateDim={updateDim}
        onImageChange={onImageChange}
        selectedFile={selectedFile}
        onSubmit={handleSubmit}
        loading={loading}
        algorithms={ALGORITHMS}
        isSubmitDisabled={loading || (genType === "image" && !selectedFile)}
      />

      {error && (
        <div className="p-3 bg-red-50 border-2 border-red-600 text-red-600 font-bold uppercase">
          {`>> ERROR: ${error}`}
        </div>
      )}

      <section className="relative border-4 border-black h-[750px] bg-zinc-50 overflow-hidden shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] flex flex-col">
        <div className="h-7 border-b-2 border-black bg-white flex items-center px-3 justify-between z-30 shrink-0">
          <span className="text-[10px] font-bold tracking-widest uppercase">
            MAZE_OUTPUT
          </span>
          <span className="text-[10px] opacity-30 font-bold uppercase">
            DIM:{" "}
            {activeMaze
              ? `${activeMaze.rows}x${activeMaze.cols}`
              : `${dims.rows}x${dims.cols}`}
          </span>
        </div>

        <div className="relative flex-1 bg-white overflow-hidden">
          <div className="absolute inset-0 bg-[radial-gradient(#000000_1px,transparent_1px)] [background-size:32px_32px] opacity-[0.05] pointer-events-none" />

          {activeMaze ? (
            <MazeCanvas maze={activeMaze} />
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
