/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_LOGO_PATH: string;
  readonly VITE_BACKGROUND_PATH: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
