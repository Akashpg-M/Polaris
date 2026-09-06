const trimTrailingSlash = (value: string) => value.replace(/\/+$/, '');

const pageOrigin = window.location.origin;
const pageWebSocketOrigin = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}`;

// Container deployments use the same browser origin and let nginx route to
// the internal services. Explicit VITE_* values remain available for split
// local development or deployments with separate public endpoints.
export const engineApi = trimTrailingSlash(import.meta.env.VITE_ENGINE_API || `${pageOrigin}/api/engine`);
export const engineReadyz = trimTrailingSlash(import.meta.env.VITE_ENGINE_READYZ || `${pageOrigin}/api/engine/readyz`);
export const gatewayApi = trimTrailingSlash(import.meta.env.VITE_GATEWAY_API || `${pageOrigin}/api/gateway`);
export const gatewayWs = trimTrailingSlash(import.meta.env.VITE_GATEWAY_WS || pageWebSocketOrigin);
