import { Performer_findPerformer_measurements as Measurements } from "src/graphql/definitions/Performer";

export const formatCareer = (
  start?: number | null,
  end?: number | null
): string | undefined =>
  start || end ? `Active ${start ?? "????"}\u2013${end ?? ""}` : undefined;

export const formatMeasurements = (val?: Measurements): string | undefined => {
  if ((val?.cup_size && val.band_size) || val?.hip || val?.waist) {
    const bust =
      val.cup_size && val.band_size ? `${val.band_size}${val.cup_size}` : "??";
    return `${bust}-${val.waist ?? ""}-${val.hip ?? ""}`;
  }
  return undefined;
};

export const getBraSize = (measurements: Measurements): string | undefined =>
  (measurements.cup_size &&
    measurements.cup_size &&
    `${measurements.band_size}${measurements.cup_size}`) ??
  undefined;

export interface URL {
  url: string;
  type: string;
}

export interface Image {
  url: string;
  id: string;
  width: number;
  height: number;
}

export const sortImageURLs = (
  urls: Image[],
  orientation: "portrait" | "landscape"
) =>
  urls
    .map((u) => ({
      ...u,
      aspect:
        orientation === "portrait"
          ? u.height / u.width > 1
          : u.width / u.height > 1,
    }))
    .sort((a, b) => {
      if (a.aspect > b.aspect) return -1;
      if (a.aspect < b.aspect) return 1;
      if (orientation === "portrait" && a.height > b.height) return -1;
      if (orientation === "portrait" && a.height < b.height) return 1;
      if (orientation === "landscape" && a.width > b.width) return -1;
      if (orientation === "landscape" && a.width < b.width) return 1;
      return 0;
    });

export const getImage = (
  urls: Image[],
  orientation: "portrait" | "landscape"
) => {
  const images = sortImageURLs(urls, orientation);
  return images?.[0]?.url ?? "";
};

export const getUrlByType = (urls: (URL | null)[], type: string) =>
  (urls && (urls.find((url) => url?.type === type) || {}).url) || "";

export const formatBodyModification = (
  bodyMod?: { location: string; description?: string | null } | null
) =>
  bodyMod
    ? bodyMod.location +
      (bodyMod.description ? ` (${bodyMod.description})` : "")
    : null;

export const formatBodyModifications = (
  bodyMod?: { location: string; description?: string | null }[] | null
) => (bodyMod ?? []).map(formatBodyModification).join(", ");

export const formatPendingEdits = (count?: number) =>
  count ? ` (${count} Pending)` : "";

export const formatDuration = (dur?: number | null) => {
  if (!dur) return "";
  let value = dur;
  let hour = 0;
  let minute = 0;
  let seconds = 0;
  if (value >= 3600) {
    hour = Math.floor(value / 3600);
    value -= hour * 3600;
  }
  minute = Math.floor(value / 60);
  value -= minute * 60;
  seconds = value;

  const res = [
    minute.toString().padStart(2, "0"),
    seconds.toString().padStart(2, "0"),
  ];
  if (hour) res.unshift(hour.toString());
  return res.join(":");
};

export const parseDuration = (
  dur: string | null | undefined
): number | null => {
  if (!dur) return null;

  const regex = /^((?<hours>\d+:)?(?<minutes>[0-5]?\d):)?(?<seconds>[0-5]?\d)$/;
  const matches = regex.exec(dur);
  const hours = matches?.groups?.hours ?? "0";
  const minutes = matches?.groups?.minutes ?? "0";
  const seconds = matches?.groups?.seconds ?? "0";

  const duration =
    Number.parseInt(seconds, 10) +
    Number.parseInt(minutes, 10) * 60 +
    Number.parseInt(hours, 10) * 3600;
  return duration > 0 ? duration : null;
};
