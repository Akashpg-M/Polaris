export class PolarisApiError extends Error {
  code?: string;
  requestId?: string;
  status: number;
  details?: unknown;

  constructor(input: { message: string; status: number; code?: string; requestId?: string; details?: unknown }) {
    super(input.message);
    this.name = 'PolarisApiError';
    this.status = input.status;
    this.code = input.code;
    this.requestId = input.requestId;
    this.details = input.details;
  }
}

export function readableError(error: unknown): PolarisApiError {
  if (error instanceof PolarisApiError) return error;
  if (error instanceof DOMException && error.name === 'AbortError') {
    return new PolarisApiError({ message: 'The request was cancelled.', status: 0, code: 'CANCELLED' });
  }
  return new PolarisApiError({
    message: error instanceof Error ? error.message : 'An unexpected request error occurred.',
    status: 0,
  });
}

