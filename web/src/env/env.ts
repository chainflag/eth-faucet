/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_LOGO_PATH: string;
  readonly VITE_BACKGROUND_PATH: string;
  readonly VITE_FAVICON_PATH: string;
  readonly VITE_SYMBOL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
