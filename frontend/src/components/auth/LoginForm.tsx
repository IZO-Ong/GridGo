"use client";
import { useState } from "react";

interface LoginFormProps {
  onSubmit: (payload: any) => void;
  loading: boolean;
}

export default function LoginForm({ onSubmit, loading }: LoginFormProps) {
  const [formData, setFormData] = useState({ identifier: "", password: "" });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ username: formData.identifier, password: formData.password });
  };

  return (
    <form className="space-y-4" onSubmit={handleSubmit}>
      <div className="space-y-1">
        <label className="text-[9px] font-black uppercase opacity-50">
          Username or Email
        </label>
        <input
          type="text"
          required
          value={formData.identifier}
          onChange={(e) =>
            setFormData({ ...formData, identifier: e.target.value })
          }
          className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50 transition-colors"
          placeholder="Enter handle or email"
        />
      </div>
      <div className="space-y-1">
        <label className="text-[9px] font-black uppercase opacity-50">
          Password
        </label>
        <input
          type="password"
          required
          value={formData.password}
          onChange={(e) =>
            setFormData({ ...formData, password: e.target.value })
          }
          className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50 transition-colors"
          placeholder="••••••••"
        />
      </div>
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-black text-white py-3 font-black uppercase italic hover:bg-zinc-800 transition-all active:translate-y-1 text-sm border-2 border-black disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer mt-4"
      >
        {loading ? "Processing..." : "Login"}
      </button>
    </form>
  );
}
