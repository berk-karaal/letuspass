import axios from "axios";

export const restQueryRetryFunc = (failureCount: number, error: Error) => {
  if (failureCount > 2) {
    return false;
  }
  if (axios.isAxiosError(error)) {
    if (
      (error.response?.status ?? 500 >= 400) &&
      (error.response?.status ?? 500 < 500)
    ) {
      return false;
    }
  }
  return true;
};
