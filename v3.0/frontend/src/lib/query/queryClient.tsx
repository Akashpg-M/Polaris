import { createContext, useContext, useEffect, useMemo, useRef, useSyncExternalStore, type ReactNode } from 'react';

export type QueryKey = readonly unknown[];
type QueryStatus = 'idle' | 'loading' | 'success' | 'error';

interface QueryRecord<T = unknown> {
  key: QueryKey;
  data?: T;
  error?: unknown;
  status: QueryStatus;
  updatedAt: number;
  promise?: Promise<T>;
  fetcher?: (signal: AbortSignal) => Promise<T>;
  controller?: AbortController;
  listeners: Set<() => void>;
}

const emptyRecord: QueryRecord = { key: [], status: 'idle', updatedAt: 0, listeners: new Set() };

function encoded(key: QueryKey) { return JSON.stringify(key); }
function startsWith(key: QueryKey, prefix: QueryKey) {
  return prefix.every((value, index) => JSON.stringify(key[index]) === JSON.stringify(value));
}

export class QueryClient {
  private records = new Map<string, QueryRecord>();

  private record<T>(key: QueryKey): QueryRecord<T> {
    const id = encoded(key);
    let record = this.records.get(id);
    if (!record) {
      record = { key, status: 'idle', updatedAt: 0, listeners: new Set() };
      this.records.set(id, record);
    }
    return record as QueryRecord<T>;
  }

  snapshot<T>(key: QueryKey) { return (this.records.get(encoded(key)) || emptyRecord) as QueryRecord<T>; }

  subscribe(key: QueryKey, listener: () => void) {
    const record = this.record(key);
    record.listeners.add(listener);
    return () => record.listeners.delete(listener);
  }

  private notify(record: QueryRecord) {
    this.records.set(encoded(record.key), { ...record });
    this.records.get(encoded(record.key))?.listeners.forEach(listener => listener());
  }

  fetch<T>(key: QueryKey, fetcher: (signal: AbortSignal) => Promise<T>, staleTime = 30_000, force = false) {
    const record = this.record<T>(key);
    record.fetcher = fetcher;
    if (record.promise) return record.promise;
    if (!force && record.status === 'success' && Date.now() - record.updatedAt < staleTime) return Promise.resolve(record.data as T);
    record.controller?.abort();
    const controller = new AbortController();
    record.controller = controller;
    record.status = record.data === undefined ? 'loading' : record.status;
    record.error = undefined;
    this.notify(record);
    const promise = fetcher(controller.signal)
      .then(data => {
        const current = this.record<T>(key);
        current.data = data;
        current.error = undefined;
        current.status = 'success';
        current.updatedAt = Date.now();
        current.promise = undefined;
        this.notify(current);
        return data;
      })
      .catch(error => {
        const current = this.record<T>(key);
        current.promise = undefined;
        if (error instanceof DOMException && error.name === 'AbortError') return Promise.reject(error);
        current.error = error;
        current.status = 'error';
        this.notify(current);
        throw error;
      });
    const live = this.record<T>(key);
    live.promise = promise;
    this.notify(live);
    return promise;
  }

  set<T>(key: QueryKey, updater: T | ((current: T | undefined) => T | undefined)) {
    const record = this.record<T>(key);
    const next = typeof updater === 'function' ? (updater as (current: T | undefined) => T | undefined)(record.data) : updater;
    if (next === undefined) return;
    record.data = next;
    record.status = 'success';
    record.updatedAt = Date.now();
    this.notify(record);
  }

  updateWhere<T>(predicate: (key: QueryKey) => boolean, updater: (current: T, key: QueryKey) => T) {
    this.records.forEach(record => {
      if (record.data !== undefined && predicate(record.key)) this.set<T>(record.key, updater(record.data as T, record.key));
    });
  }

  refresh(prefix: QueryKey) {
    this.records.forEach(record => {
      if (startsWith(record.key, prefix) && record.fetcher && record.listeners.size > 0) {
        void this.fetch(record.key, record.fetcher, 0, true).catch(() => undefined);
      }
    });
  }

  clear() {
    this.records.forEach(record => record.controller?.abort());
    this.records.clear();
  }
}

const QueryContext = createContext<QueryClient | null>(null);

export function QueryProvider({ children }: { children: ReactNode }) {
  const client = useMemo(() => new QueryClient(), []);
  return <QueryContext.Provider value={client}>{children}</QueryContext.Provider>;
}

export function useQueryClient() {
  const value = useContext(QueryContext);
  if (!value) throw new Error('QueryProvider is missing');
  return value;
}

export function useQuery<T>(input: {
  key: QueryKey;
  query: (signal: AbortSignal) => Promise<T>;
  enabled?: boolean;
  staleTime?: number;
  refetchInterval?: number | false | ((data: T | undefined) => number | false);
}) {
  const client = useQueryClient();
  const keyId = JSON.stringify(input.key);
  const key = useMemo(() => JSON.parse(keyId) as QueryKey, [keyId]);
  const pollingEnabled = Boolean(input.refetchInterval);
  const query = useRef(input.query);
  const intervalConfig = useRef(input.refetchInterval);
  useEffect(() => { query.current = input.query; }, [input.query]);
  useEffect(() => { intervalConfig.current = input.refetchInterval; }, [input.refetchInterval]);
  const state = useSyncExternalStore(
    listener => client.subscribe(key, listener),
    () => client.snapshot<T>(key),
    () => client.snapshot<T>(key),
  );
  useEffect(() => {
    if (input.enabled === false) return;
    void client.fetch(key, signal => query.current(signal), input.staleTime).catch(() => undefined);
  }, [client, key, keyId, input.enabled, input.staleTime]);
  useEffect(() => {
    if (input.enabled === false || !input.refetchInterval) return;
    let timer: number | undefined;
    let disposed = false;
    const schedule = () => {
      const configured = intervalConfig.current;
      const delay = typeof configured === 'function' ? configured(client.snapshot<T>(key).data) : configured;
      if (delay === false || disposed) return;
      timer = window.setTimeout(async () => {
        await client.fetch(key, signal => query.current(signal), 0, true).catch(() => undefined);
        schedule();
      }, delay);
    };
    schedule();
    return () => { disposed = true; if (timer !== undefined) window.clearTimeout(timer); };
  }, [client, key, keyId, input.enabled, input.refetchInterval, pollingEnabled]);
  return {
    data: state.data,
    error: state.error,
    isLoading: state.status === 'loading' || state.status === 'idle',
    isRefreshing: state.status === 'success' && Boolean(state.promise),
    refetch: () => client.fetch(key, signal => query.current(signal), 0, true),
  };
}
