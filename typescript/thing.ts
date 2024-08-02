import { useDjangolang } from "./api";

export function Djangolang() {
  const { data, isLoading, error } = useDjangolang(
    "/logical-things/{primaryKey}",
    {
      params: {
        path: {
          primaryKey: "aa-bb-cc",
        },
      },
    },
  );

  if (isLoading) {
    return "Loading...";
  }

  if (error) {
    return `error: ${error}`;
  }

  const objects = data?.objects;

  return data?.objects
    ?.map((object) => {
      return object.Name;
    })
    .join(", ");
}
