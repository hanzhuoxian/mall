import { afterEach } from 'vitest';
import { cleanup } from '@testing-library/react';

afterEach(() => cleanup());

const createStorage = () => {
  const values = new Map<string, string>();
  return { getItem: (key: string) => values.get(key) ?? null, setItem: (key: string, value: string) => values.set(key, String(value)), removeItem: (key: string) => values.delete(key), clear: () => values.clear() };
};
if (typeof globalThis.localStorage === 'undefined') Object.defineProperty(globalThis, 'localStorage', { value: createStorage(), configurable: true });
if (typeof globalThis.sessionStorage === 'undefined') Object.defineProperty(globalThis, 'sessionStorage', { value: createStorage(), configurable: true });

if (typeof window !== 'undefined' && !window.matchMedia) {
  window.matchMedia = (query: string) => ({
    matches: false, media: query, onchange: null,
    addListener: () => {}, removeListener: () => {},
    addEventListener: () => {}, removeEventListener: () => {}, dispatchEvent: () => false,
  });
}
if (typeof globalThis.ResizeObserver === 'undefined') {
  globalThis.ResizeObserver = class { observe() {} unobserve() {} disconnect() {} } as unknown as typeof ResizeObserver;
}
