import { ShipWheel } from "lucide-react";
import ThreeDotLoader from "../three-dot-loader";

export default function PageLoader() {
  return (
    <div
      data-theme="forest"
      className="flex h-screen w-full items-center justify-center rounded"
    >
      <div className="space-y-4 text-center">
        <ThreeDotLoader size="lg" />
        <div className="flex items-center gap-4">
          <ShipWheel className="size-9 text-primary" />
          <span className="bg-gradient-to-r from-primary to-secondary bg-clip-text font-mono text-3xl font-bold tracking-wider text-transparent">
            Streamify
          </span>
        </div>
      </div>
    </div>
  );
}
