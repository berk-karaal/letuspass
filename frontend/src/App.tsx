import { getCurrentUser } from "@/api/letuspass";
import router from "@/routes";
import { useAppDispatch } from "@/store/hooks";
import { startupComplete, userLoggedIn } from "@/store/slices/user";
import { notifications } from "@mantine/notifications";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useEffect } from "react";
import { RouterProvider } from "react-router-dom";

function App() {
  const dispatch = useAppDispatch();

  const currentUserQuery = useQuery({
    queryKey: ["current-user"],
    queryFn: getCurrentUser,
    refetchOnWindowFocus: false,
    staleTime: Infinity,
    gcTime: Infinity,
    retry: (failureCount: number, error: Error) => {
      if (failureCount > 2) {
        return false;
      }
      if (axios.isAxiosError(error)) {
        if (error.response?.status == 401) {
          return false;
        }
      }
      return true;
    },
  });

  useEffect(() => {
    if (currentUserQuery.isError) {
      dispatch(startupComplete());
      if (
        axios.isAxiosError(currentUserQuery.error) &&
        currentUserQuery.error.response?.status !== 401
      ) {
        notifications.show({
          title: "Failed to fetch current user",
          message: "Failed to fetch current user.",
          color: "red",
        });
      }
    } else if (currentUserQuery.data) {
      dispatch(startupComplete());
      dispatch(
        userLoggedIn({
          email: currentUserQuery.data.email,
          name: currentUserQuery.data.name,
        })
      );
    }
  }, [currentUserQuery, dispatch]);

  return (
    <>
      <RouterProvider router={router} />
    </>
  );
}

export default App;
