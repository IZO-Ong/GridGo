"use client";
import { useState, useEffect } from "react";
import {
  saveGenerateSession,
  loadGenerateSession,
  savePreferences,
  loadPreferences,
} from "@/lib/db";
import { updateThumbnail } from "@/lib/api";
import MazeCanvas from "@/components/maze/MazeCanvas";
import GenerateControls from "@/components/maze/GenerateControls";
import { useImageDimensions } from "@/hooks/useImageDimensions";
import { useMazeGeneration, MazeData } from "@/hooks/useMazeGeneration";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

const ALGORITHMS = [
  { id: "image", label: "IMAGE_KRUSKAL" },
  { id: "kruskal", label: "RANDOM_KRUSKAL" },
  { id: "recursive", label: "DFS_BACKTRACKER" },
];

export default function Home() {
  const [activeMaze, setActiveMaze] = useState<MazeData | null>(null);
  const [hasLoaded, setHasLoaded] = useState(false);

  const { loading, error, executeGeneration } = useMazeGeneration();
  const [genType, setGenType] = useState("image");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const { dims, updateDim, clampDimensions, handleImageChange } =
    useImageDimensions();

  useEffect(() => {
    const init = async () => {
      const savedMaze = await loadGenerateSession();
      if (savedMaze) setActiveMaze(savedMaze);

      const prefs = await loadPreferences("gen_prefs");
      if (prefs) {
        // Validate type against available algorithms
        if (ALGORITHMS.some((a) => a.id === prefs.genType)) {
          setGenType(prefs.genType);
        }
        updateDim("rows", prefs.dims.rows);
        updateDim("cols", prefs.dims.cols);
      }
      setHasLoaded(true);
    };
    init();
  }, [updateDim]);

  // 2. Real-time Auto-save
  useEffect(() => {
    if (hasLoaded) {
      savePreferences("gen_prefs", { genType, dims });
    }
  }, [genType, dims, hasLoaded]);

  const onImageChange = (e: React.ChangeEvent<HTMLInputElement> | File) => {
    let file: File | undefined;
    if (e instanceof File) {
      file = e;
      const fakeEvent = {
        target: { files: [file] },
      } as unknown as React.ChangeEvent<HTMLInputElement>;
      handleImageChange(fakeEvent);
    } else {
      file = e.target.files?.[0];
      handleImageChange(e);
    }
    if (file) setSelectedFile(file);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    clampDimensions();
    const formData = new FormData(e.currentTarget);
    formData.set("type", genType);
    if (genType === "image" && selectedFile)
      formData.set("image", selectedFile);

    const newMaze = await executeGeneration(formData);

    if (newMaze) {
      setActiveMaze(newMaze);
      await saveGenerateSession(newMaze);

      setTimeout(async () => {
        const canvas = document.getElementById(
          "main-maze-canvas"
        ) as HTMLCanvasElement;
        if (canvas) {
          const tempCanvas = document.createElement("canvas");
          const tempCtx = tempCanvas.getContext("2d");
          tempCanvas.width = 400;
          tempCanvas.height = 225;

          if (tempCtx) {
            tempCtx.fillStyle = "white";
            tempCtx.fillRect(0, 0, tempCanvas.width, tempCanvas.height);

            const PADDING = 800;
            const mazeW = canvas.width - PADDING * 2;
            const mazeH = canvas.height - PADDING * 2;
            const scale = Math.max(
              tempCanvas.width / mazeW,
              tempCanvas.height / mazeH
            );
            const xOff = (tempCanvas.width - mazeW * scale) / 2;
            const yOff = (tempCanvas.height - mazeH * scale) / 2;

            tempCtx.drawImage(
              canvas,
              PADDING,
              PADDING,
              mazeW,
              mazeH,
              xOff,
              yOff,
              mazeW * scale,
              mazeH * scale
            );
            const b64 = tempCanvas.toDataURL("image/webp", 0.6);

            await updateThumbnail(newMaze.id, b64).catch((err) =>
              console.error("THUMBNAIL_SYNC_FAILURE", err)
            );
          }
        }
      }, 400);
    }
  };

  return (
    <div className="space-y-8">
      <GenerateControls
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
        <div className="p-3 bg-red-50 border-2 border-red-600 text-red-600 font-bold uppercase text-[11px]">
          {`>> ERROR_SEQUENCE: ${error}`}
        </div>
      )}

      <section className="relative border-4 border-black h-[750px] bg-zinc-50 overflow-hidden shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] flex flex-col">
        <div className="h-7 border-b-2 border-black bg-white flex items-center px-3 justify-between z-30 shrink-0">
          <div className="flex items-center gap-3">
            <span className="text-[10px] font-bold tracking-widest uppercase">
              MAZE_OUTPUT
            </span>
            {activeMaze?.id && (
              <span className="bg-black text-white px-2 py-0.5 text-[9px] font-black tracking-tighter">
                {activeMaze.id}
              </span>
            )}
          </div>

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
            <MazeCanvas maze={activeMaze} showSave={true} showShare={false} />
          ) : (
            <div className="h-full w-full flex items-center justify-center">
              <p className="opacity-20 tracking-[0.5em] font-bold uppercase text-2xl text-center px-12">
                Generate a Maze!
              </p>
            </div>
          )}
        </div>
      </section>
    </div>
  );
}
