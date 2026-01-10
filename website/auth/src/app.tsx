import { Router } from "@solidjs/router";
import { FileRoutes } from "@solidjs/start/router";
import { Suspense } from "solid-js";
import "./app.css";

import { MetaProvider, Title } from "@solidjs/meta";

import { isServer } from "solid-js/web";

import { ColorModeProvider, ColorModeScript, cookieStorageManagerSSR } from "@kobalte/core";
import { getCookie } from "vinxi/http";
import { Toaster } from "./components/ui/toast";

function getColorMode() {
  "use server"
  const colorMode = getCookie("kb-color-mode")
  return colorMode ? `kb-color-mode=${colorMode}` : "dark"
}

export default function App() {
  const colorModeStore = cookieStorageManagerSSR(isServer ? getColorMode() : document.cookie)

  return (
    <Router
      root={props => (
        <>
          <MetaProvider>
            <ColorModeScript storageType={colorModeStore.type} />
            <ColorModeProvider storageManager={colorModeStore}>
              <Suspense>{props.children}</Suspense>
              <Toaster />
            </ColorModeProvider>
          </MetaProvider>
        </>
      )}
    >
      <FileRoutes />
    </Router>
  );
}
