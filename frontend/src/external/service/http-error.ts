import type { ResponseError } from "@/external/client/api/generated/runtime";

export function isResponseError(error: unknown): error is ResponseError {
  return (
    typeof error === "object" &&
    error !== null &&
    "response" in error &&
    typeof (error as { response?: Response }).response?.status === "number"
  );
}

export function isNotFoundError(error: unknown): error is ResponseError {
  return isResponseError(error) && error.response.status === 404;
}

export async function readErrorMessage(
  error: ResponseError,
): Promise<string | undefined> {
  try {
    const cloned = error.response.clone();
    const body = await cloned.json();
    if (typeof body?.message === "string") {
      return body.message;
    }
  } catch {
    // no-op
  }
  return undefined;
}
