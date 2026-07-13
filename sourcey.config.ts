export default {
  name: "txeh Reference",
  siteUrl: "https://txeh.txn2.com",
  baseUrl: "/reference",
  repo: "https://github.com/txn2/txeh",
  editBranch: "master",
  navigation: {
    tabs: [
      {
        tab: "Guides",
        slug: "",
        mkdocs: "./mkdocs.yml",
      },
      {
        tab: "Go API",
        slug: "api",
        godoc: {
          module: ".",
          packages: ["./..."],
          mode: "live",
        },
      },
    ],
  },
};
