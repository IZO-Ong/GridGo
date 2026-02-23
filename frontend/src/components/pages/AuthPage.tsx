"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import AuthCard from "@/components/auth/AuthCard";
import LoginForm from "@/components/auth/LoginForm";
import RegisterForm from "@/components/auth/RegisterForm";
import VerifyForm from "@/components/auth/VerifyForm";
import OAuthButton from "@/components/auth/OAuthButton";
import {
  login as apiLogin,
  register as apiRegister,
  verifyAccount,
} from "@/lib/api";

export default function AuthPage() {
  const { login } = useAuth();
  const router = useRouter();

  const [view, setView] = useState<"login" | "register" | "verify">("login");
  // Store full registration data in case we need to resend
  const [formData, setFormData] = useState({
    email: "",
    username: "",
    password: "",
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAuthSubmit = async (payload: any) => {
    setError(null);
    setLoading(true);
    try {
      if (view === "login") {
        const data = await apiLogin(payload);
        login(data.token, data.username);
        router.push("/");
      } else if (view === "register") {
        // Save the full payload so resending actually works
        // (your backend expects username/password in HandleRegister)
        setFormData({
          email: payload.email,
          username: payload.username,
          password: payload.password,
        });
        await apiRegister(payload);
        setView("verify");
        setError("SUCCESS: Code sent to email.");
      } else if (view === "verify") {
        // payload here is just the 6-digit string from VerifyForm
        await verifyAccount(formData.email, payload);
        setView("login");
        setError("SUCCESS: Verified. Please login.");
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleResend = async () => {
    setError(null);
    try {
      await apiRegister(formData);
      setError("SUCCESS: New code transmitted.");
    } catch (err: any) {
      setError(err.message || "RESEND_FAILURE");
      throw err;
    }
  };

  return (
    <AuthCard
      title={
        view === "login"
          ? "Login"
          : view === "register"
            ? "Register"
            : "Verify Code"
      }
      error={error}
      footerLabel={view === "login" ? "Create an account" : "Back to Login"}
      footerAction={() => {
        setError(null);
        setView(view === "login" ? "register" : "login");
      }}
    >
      {view === "login" && (
        <LoginForm onSubmit={handleAuthSubmit} loading={loading} />
      )}
      {view === "register" && (
        <RegisterForm
          onSubmit={handleAuthSubmit}
          loading={loading}
          setError={setError}
        />
      )}
      {view === "verify" && (
        <VerifyForm
          email={formData.email}
          onVerify={handleAuthSubmit}
          onResend={handleResend}
          loading={loading}
        />
      )}

      {view !== "verify" && <OAuthButton />}
    </AuthCard>
  );
}
