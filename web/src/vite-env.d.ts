/// <reference types="svelte" />
/// <reference types="vite/client" />

interface Window {
  hcaptcha?: {
    render: (container: string, options: { sitekey: string }) => string;
    execute: (widgetId: string, options: { async: boolean }) => Promise<{ response: string }>;
  };
  [key: string]: any;
}
