import { readFileSync } from "fs";

function toArtistList(artist) {
  return Array.isArray(artist) ? artist : [artist];
}

function formatArtists(artist) {
  const list = toArtistList(artist);
  if (list.length === 1) return list[0];
  return list.slice(0, -1).join(", ") + " & " + list[list.length - 1];
}

export default function (eleventyConfig) {
  eleventyConfig.addFilter("formatArtists", formatArtists);
  eleventyConfig.addFilter("toArtistList", toArtistList);

  // Create a collection of all songs
  eleventyConfig.addCollection("songs", (collectionApi) => {
    return collectionApi
      .getFilteredByGlob("src/songs/*.md")
      .sort((a, b) => a.data.title.localeCompare(b.data.title));
  });

  // Create a collection grouped by album
  eleventyConfig.addCollection("albums", (collectionApi) => {
    const songs = collectionApi.getFilteredByGlob("src/songs/*.md");
    const albumMap = new Map();
    for (const song of songs) {
      const album = song.data.album;
      if (!album) continue;
      if (!albumMap.has(album)) {
        albumMap.set(album, []);
      }
      albumMap.get(album).push(song);
    }
    const sorted = [...albumMap.entries()]
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([album, songs]) => ({
        name: album,
        slug: eleventyConfig.getFilter("slugify")(album),
        songs: songs.sort((a, b) => a.data.title.localeCompare(b.data.title)),
      }));
    return sorted;
  });

  // Create a collection grouped by artist
  eleventyConfig.addCollection("artists", (collectionApi) => {
    const songs = collectionApi.getFilteredByGlob("src/songs/*.md");
    const artistMap = new Map();
    for (const song of songs) {
      for (const artist of toArtistList(song.data.artist)) {
        if (!artistMap.has(artist)) {
          artistMap.set(artist, []);
        }
        artistMap.get(artist).push(song);
      }
    }
    // Sort artists alphabetically, and songs within each artist
    const sorted = [...artistMap.entries()]
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([artist, songs]) => ({
        name: artist,
        slug: eleventyConfig.getFilter("slugify")(artist),
        songs: songs.sort((a, b) => a.data.title.localeCompare(b.data.title)),
      }));
    return sorted;
  });

  eleventyConfig.addShortcode("inlineCSS", () => {
    const css = readFileSync("src/css/styles.css", "utf-8");
    return `<style>${css}</style>`;
  });

  return {
    pathPrefix: process.env.ELEVENTY_PATH_PREFIX || "/karaokay/",
    dir: {
      input: "src",
      output: "_site",
    },
    templateFormats: ["njk", "md"],
    markdownTemplateEngine: "njk",
  };
}
