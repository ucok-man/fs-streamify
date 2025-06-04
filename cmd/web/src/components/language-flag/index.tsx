import { LANGUAGE_TO_FLAG } from "../../constants";

type Props = {
  language: string;
};

export default function LanguageFlag({ language }: Props) {
  if (!language) return null;

  const langLower = language.toLowerCase() as keyof typeof LANGUAGE_TO_FLAG;
  const countryCode = LANGUAGE_TO_FLAG[langLower];

  if (countryCode) {
    return (
      <img
        src={`https://flagcdn.com/24x18/${countryCode}.png`}
        alt={`${langLower} flag`}
        className="mr-1 inline-block h-3"
      />
    );
  }
  return null;
}
