const esbuild = require("esbuild");
const {sassPlugin} = require("esbuild-sass-plugin")

esbuild
    .build({
        entryPoints: ["frontend/Index.tsx", "frontend/styles.scss"],
        outdir: "gallery-server/public/",
        bundle: true,
        minify: true,
        loader: { '.ts': 'tsx', '.scss': 'css' },
        plugins: [sassPlugin()],
    })
    .then(() => console.log("⚡ Build complete! ⚡"))
    .catch(() => process.exit(1));
