import { cn } from "../../lib/utils";

interface Props {
  className?: string;
  size?: "sm" | "md" | "lg";
  color?: "primary" | "secondary";
}

export default function ThreeDotLoader({
  className,
  size = "md",
  color = "primary",
}: Props = {}) {
  const sizeClasses = {
    sm: "w-1 h-1",
    md: "w-2 h-2",
    lg: "w-3 h-3",
  };

  const colorClasses = {
    primary: "bg-primary",
    secondary: "bg-secondary",
  };

  const dotClass = cn(
    "rounded-full animate-pulse",
    sizeClasses[size],
    colorClasses[color]
  );

  return (
    <div
      className={cn("flex items-center justify-center space-x-1", className)}
    >
      <div
        className={cn(dotClass)}
        style={{
          animationDelay: "0ms",
          animationDuration: "1.4s",
        }}
      />
      <div
        className={cn(dotClass)}
        style={{
          animationDelay: "160ms",
          animationDuration: "1.4s",
        }}
      />
      <div
        className={cn(dotClass)}
        style={{
          animationDelay: "320ms",
          animationDuration: "1.4s",
        }}
      />
    </div>
  );
}
