"use client";
import { useState, useRef } from "react";
import AlgorithmSelect from "@/components/ui/AlgorithmSelect";
import GridDimensionsInput from "@/components/ui/GridDimensionsInput";

interface GenerateControlsProps {
  genType: string;
  setGenType: (val: string) => void;
  dims: { rows: number; cols: number };
  updateDim: (dim: "rows" | "cols", val: number) => void;
  onImageChange: (e: React.ChangeEvent<HTMLInputElement> | File) => void;
  selectedFile: File | null;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  loading: boolean;
  isSubmitDisabled: boolean;
  algorithms: { id: string; label: string }[];
}

export default function GenerateControls({
  genType,
  setGenType,
  dims,
  updateDim,
  onImageChange,
  selectedFile,
  onSubmit,
  loading,
  isSubmitDisabled,
  algorithms,
}: GenerateControlsProps) {
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Clamping utility to enforce 2-300 range
  const clamp = (val: number) => Math.min(Math.max(2, val), 300);

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    if (e.type === "dragenter" || e.type === "dragover") setIsDragging(true);
    else if (e.type === "dragleave") setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files?.[0]) onImageChange(e.dataTransfer.files[0]);
  };

  return (
    <form onSubmit={onSubmit} className="grid grid-cols-12 gap-4 items-end">
      <div className="col-span-3 space-y-2">
        <div className="flex justify-between items-end">
          <label className="block font-bold uppercase tracking-widest text-[9px]">
            Grid_Dimensions
          </label>
          <span className="text-[8px] font-mono opacity-40">RANGE: 2-300</span>
        </div>
        <GridDimensionsInput
          rows={dims.rows}
          cols={dims.cols}
          onUpdate={updateDim}
          min={2}
          max={{ rows: 300, cols: 300 }}
        />
      </div>

      <div className="col-span-4 space-y-2">
        <label className="block font-bold uppercase tracking-widest text-[9px]">
          Algorithm
        </label>
        <AlgorithmSelect
          value={genType}
          onChange={setGenType}
          options={algorithms}
        />
      </div>

      <div
        className={`col-span-5 space-y-2 transition-all ${
          genType === "image"
            ? "opacity-100"
            : "opacity-0 pointer-events-none hidden"
        }`}
      >
        <label className="block font-bold text-[9px] uppercase tracking-widest">
          Source_Image
        </label>
        <div
          onDragOver={handleDrag}
          onDrop={handleDrop}
          onClick={() => fileInputRef.current?.click()}
          className={`relative w-full h-[38px] border-2 border-dashed flex items-stretch transition-colors cursor-pointer bg-white ${
            isDragging ? "bg-zinc-100 border-black" : "border-zinc-300"
          }`}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={onImageChange}
            className="hidden"
          />
          <div className="flex-1 flex items-center px-3 min-w-0">
            <span className="truncate text-[10px] font-bold uppercase">
              {selectedFile ? selectedFile.name : "Drop image here"}
            </span>
          </div>
          <div className="border-l-2 border-black bg-black text-white px-3 flex items-center text-[9px] font-black uppercase">
            Browse
          </div>
        </div>
      </div>

      <button
        type="submit"
        disabled={isSubmitDisabled || loading}
        className={`col-span-12 border-2 border-black p-4 font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all flex items-center justify-between px-8 group
        ${
          loading || isSubmitDisabled
            ? "bg-zinc-100 text-zinc-400 opacity-50 cursor-not-allowed" // Greyed out style
            : "bg-white hover:bg-black hover:text-white cursor-pointer active:translate-y-1 active:shadow-none"
        }`}
      >
        <span className="italic tracking-tighter text-lg font-black uppercase">
          {loading ? ">>> PROCESSING_SEQUENCE..." : ">>> GENERATE_MAZE"}
        </span>
        {loading && (
          <div className="flex gap-1">
            <span className="w-1.5 h-1.5 bg-zinc-400 rounded-full animate-bounce [animation-delay:-0.3s]" />
            <span className="w-1.5 h-1.5 bg-zinc-400 rounded-full animate-bounce [animation-delay:-0.15s]" />
            <span className="w-1.5 h-1.5 bg-zinc-400 rounded-full animate-bounce" />
          </div>
        )}
      </button>
    </form>
  );
}
