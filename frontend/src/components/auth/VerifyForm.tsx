"use client";
import { useState, useEffect } from "react";

export default function VerifyForm({
  email,
  onVerify,
  onResend,
  loading,
}: {
  email: string;
  onVerify: (otp: string) => void;
  onResend: () => Promise<void>;
  loading: boolean;
}) {
  const [otp, setOtp] = useState("");
  const [cooldown, setCooldown] = useState(60);
  const [resending, setResending] = useState(false);

  // Timer logic
  useEffect(() => {
    let timer: NodeJS.Timeout;
    if (cooldown > 0) {
      timer = setTimeout(() => setCooldown(cooldown - 1), 1000);
    }
    return () => clearTimeout(timer);
  }, [cooldown]);

  const handleResend = async () => {
    if (cooldown > 0 || resending) return;

    setResending(true);
    try {
      await onResend();
      setCooldown(60);
    } catch (err) {
      alert("RESEND_FAILED");
    } finally {
      setResending(false);
    }
  };

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
          onChange={(e) => setOtp(e.target.value.replace(/\D/g, ""))}
          className="w-full border-2 border-black p-4 outline-none font-black text-center text-2xl tracking-[0.5em] focus:bg-zinc-50"
          placeholder="000000"
        />
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-black text-white py-3 font-black uppercase italic hover:bg-zinc-800 border-2 border-black disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer transition-all"
      >
        {loading ? "Verifying..." : "Verify Identity"}
      </button>

      {/* Resend Action */}
      <div className="text-center">
        <button
          type="button"
          onClick={handleResend}
          disabled={cooldown > 0 || resending}
          className="text-[10px] font-black uppercase tracking-tighter disabled:opacity-30 cursor-pointer hover:underline disabled:no-underline"
        >
          {resending
            ? "Transmitting..."
            : cooldown > 0
              ? `Resend Code in ${cooldown}s`
              : "[ Resend New Code ]"}
        </button>
      </div>
    </form>
  );
}
