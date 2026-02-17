import { useState } from "react";

export default function VerifyForm({
  email,
  onVerify,
  loading,
}: {
  email: string;
  onVerify: (otp: string) => void;
  loading: boolean;
}) {
  const [otp, setOtp] = useState("");

  return (
    <form
      className="space-y-4"
      onSubmit={(e) => {
        e.preventDefault();
        onVerify(otp);
      }}
    >
      <div className="space-y-1 text-center">
        <label className="text-[9px] font-black uppercase opacity-50 block">
          Enter 6-Digit Code for {email}
        </label>
        <input
          type="text"
          maxLength={6}
          required
          value={otp}
          onChange={(e) => setOtp(e.target.value)}
          className="w-full border-2 border-black p-4 outline-none font-black text-center text-2xl tracking-[0.5em] focus:bg-zinc-50"
          placeholder="000000"
        />
      </div>
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-black text-white py-3 font-black uppercase italic hover:bg-zinc-800 border-2 border-black"
      >
        {loading ? "Verifying..." : "Verify Identity"}
      </button>
    </form>
  );
}
