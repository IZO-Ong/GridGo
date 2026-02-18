"use client";
import { useState } from "react";

interface RegisterFormProps {
  onSubmit: (payload: any) => void;
  loading: boolean;
  setError: (msg: string | null) => void;
}

export default function RegisterForm({
  onSubmit,
  loading,
  setError,
}: RegisterFormProps) {
  const [formData, setFormData] = useState({
    email: "",
    username: "",
    password: "",
    verify: "",
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.password !== formData.verify) {
      setError("Passwords do not match"); // Error lifting
      return;
    }
    onSubmit({
      email: formData.email,
      username: formData.username,
      password: formData.password,
    });
  };

  return (
    <form className="space-y-4" onSubmit={handleSubmit}>
      <div className="space-y-1">
        <label className="text-[9px] font-black uppercase opacity-50">
          Email Address
        </label>
        <input
          type="email"
          required
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50"
          placeholder="user@domain.com"
        />
      </div>
      <div className="space-y-1">
        <label className="text-[9px] font-black uppercase opacity-50">
          Public Username
        </label>
        <input
          type="text"
          required
          value={formData.username}
          onChange={(e) =>
            setFormData({ ...formData, username: e.target.value })
          }
          className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50"
          placeholder="e.g. MazeMaster99"
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
          className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50"
          placeholder="••••••••"
        />
      </div>
      <div className="space-y-1">
        <label className="text-[9px] font-black uppercase opacity-50">
          Confirm Password
        </label>
        <input
          type="password"
          required
          value={formData.verify}
          onChange={(e) => setFormData({ ...formData, verify: e.target.value })}
          className="w-full border-2 border-black p-2 outline-none font-bold text-sm focus:bg-zinc-50"
          placeholder="••••••••"
        />
      </div>
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-black text-white py-3 font-black uppercase italic hover:bg-zinc-800 transition-all active:translate-y-1 text-sm border-2 border-black disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer mt-4"
      >
        {loading ? "Creating Account..." : "Register"}
      </button>
    </form>
  );
}
