import { useNavigate, useSearch } from "@tanstack/react-router";
import { SearchIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { useDebounceValue } from "usehooks-ts";

export default function SearchBox() {
  const searchParams = useSearch({
    from: "/_protected/_layout/(home)/",
  });
  const navigate = useNavigate({ from: "/" });

  const [query, setQuery] = useState(searchParams.query || "");
  const [dbquery] = useDebounceValue(query.trim(), 500);

  useEffect(() => {
    navigate({
      to: "/",
      search: { query: dbquery },
    });
  }, [dbquery, navigate, searchParams]);

  return (
    <label className="input w-full md:max-w-md">
      <span>
        <SearchIcon />
      </span>
      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        type="text"
        className="grow"
        placeholder="Search username..."
      />
    </label>
  );
}
