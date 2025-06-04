import { cn } from "../../lib/utils";

type Props = {
  msg?: string;
  className?: string;
};

export default function SpinnerBtn({ msg = "Loading...", className }: Props) {
  return (
    <>
      <span className={cn("loading loading-xs loading-spinner", className)} />
      <span>{msg}</span>
    </>
  );
}
