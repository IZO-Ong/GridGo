"use client";
import { useState, useRef, useEffect } from "react";

interface AlgorithmOption {
  id: string;
  label: string;
}

interface AlgorithmSelectProps {
  value: string;
  onChange: (id: string) => void;
  options: AlgorithmOption[];
}

export default function AlgorithmSelect({
  value,
  onChange,
  options,
}: AlgorithmSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const clickOutside = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", clickOutside);
    return () => document.removeEventListener("mousedown", clickOutside);
  }, []);

  return (
    <div
      className="relative border-2 border-black bg-white font-mono h-[38px] flex items-center"
      ref={ref}
    >
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="w-full h-full px-3 text-left flex justify-between items-center focus:bg-zinc-50 transition-colors cursor-pointer group"
      >
        <span className="truncate text-xs font-bold uppercase tracking-tight">
          {options.find((o) => o.id === value)?.label}
        </span>

        <div className="ml-2 flex items-center pointer-events-none">
          <svg
            className={`w-3 h-3 transition-transform duration-200 ${isOpen ? "rotate-180" : ""}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            strokeWidth="4"
          >
            <path
              strokeLinecap="square"
              strokeLinejoin="miter"
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </div>
      </button>

      {isOpen && (
        <div className="absolute top-[calc(100%-2px)] left-[-2px] w-[calc(100%+4px)] border-2 border-black bg-white z-50 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] divide-y-2 divide-black">
          {options.map((opt) => (
            <button
              key={opt.id}
              type="button"
              onClick={() => {
                onChange(opt.id);
                setIsOpen(false);
              }}
              className={`w-full p-2 text-left text-xs font-bold uppercase transition-colors cursor-pointer ${
                value === opt.id
                  ? "bg-zinc-100"
                  : "hover:bg-black hover:text-white"
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
